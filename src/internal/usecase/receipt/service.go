package receipt

import (
	"context"
	"shop-api/generated/db"
	"shop-api/internal/config"
	"shop-api/internal/database/repositories"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stripe/stripe-go/v84"
)

var CURRENCY = "usd"
var MODE = "payment"
var SUCCESS_URL = "http://localhost:3000"
var CANCEL_URL = "http://localhost:3000/api/receipts/cancel"

type ReceiptService struct {
	receiptRepo  *repositories.ReceiptRepository
	productRepo  *repositories.ProductRepository
	stripeClient *stripe.Client
	config       *config.Config
}

func NewReceiptService(receiptRepo *repositories.ReceiptRepository, productRepo *repositories.ProductRepository, stripeClient *stripe.Client, config *config.Config) *ReceiptService {
	return &ReceiptService{receiptRepo: receiptRepo, productRepo: productRepo, stripeClient: stripeClient, config: config}
}

func (s *ReceiptService) CreateReceipt(ctx context.Context, userId int64, model CreateReceiptDto) (ReceiptDto, error) {

	ids := make([]int64, 0)
	for id := range model.Products {
		ids = append(ids, id)
	}

	products, err := s.productRepo.GetProductByIds(ctx, ids)
	if err != nil {
		return ReceiptDto{}, err
	}

	productsParams := make([]db.CreateReceiptProductParams, len(products))
	stripeParams := make([]*stripe.CheckoutSessionCreateLineItemParams, len(products))
	sumPrice, sumDiscount := int64(0), int64(0)

	for i, product := range products {
		quantity := model.Products[product.ID]
		sumPrice += product.Price * int64(quantity)
		sumDiscount += product.Discount * int64(quantity)
		productsParams[i] = db.CreateReceiptProductParams{
			ReceiptID: 0,
			ProductID: product.ID,
			Quantity:  quantity,
			Price:     product.Price,
			Discount:  product.Discount,
		}
		p := product.Price - product.Discount
		q := int64(quantity)
		stripeParams[i] = &stripe.CheckoutSessionCreateLineItemParams{
			PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
				Currency: &CURRENCY,
				ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
					Name: &product.Name,
				},
				UnitAmount: &p,
			},
			Quantity: &q,
		}
	}

	successUrl := s.config.BackendUrl + "/api/receipts/success?sessionId={{CHECKOUT_SESSION_ID}}"
	cancelUrl := s.config.BackendUrl + "/api/receipts/cancel?sessionId={{CHECKOUT_SESSION_ID}}"

	session, err := s.stripeClient.V1CheckoutSessions.Create(ctx, &stripe.CheckoutSessionCreateParams{
		LineItems:  stripeParams,
		Mode:       &MODE,
		SuccessURL: &successUrl,
		CancelURL:  &cancelUrl,
	})

	if err != nil {
		return ReceiptDto{}, err
	}

	_, _, err = s.receiptRepo.CreateReceipt(ctx, db.CreateReceiptParams{
		UserID:       userId,
		SumPrice:     sumPrice,
		SumDiscount:  sumDiscount,
		StripeID:     pgtype.Text{String: session.ID, Valid: true},
		StripeStatus: db.StripeStatusPending,
	}, productsParams)

	if err != nil {
		return ReceiptDto{}, err
	}

	return ReceiptDto{Link: session.URL}, nil
}

func (s *ReceiptService) CancelReceipt(ctx context.Context, sessionId string) (ReceiptDto, error) {
	session, err := s.stripeClient.V1CheckoutSessions.Retrieve(ctx, sessionId, nil)
	if err != nil {
		return ReceiptDto{}, err
	}

	status := parseStripeStatus(session.Status)

	_, err = s.productRepo.UpdateReceiptStatus(ctx, db.UpdateReceiptStatusParams{
		StripeID:     pgtype.Text{String: sessionId, Valid: true},
		StripeStatus: status,
	})

	return ReceiptDto{Link: s.getRedirectUrl(status)}, err
}

func (s *ReceiptService) SuccessReceipt(ctx context.Context, sessionId string) (ReceiptDto, error) {
	session, err := s.stripeClient.V1CheckoutSessions.Retrieve(ctx, sessionId, nil)
	if err != nil {
		return ReceiptDto{}, err
	}

	status := parseStripeStatus(session.Status)

	_, err = s.productRepo.UpdateReceiptStatus(ctx, db.UpdateReceiptStatusParams{
		StripeID:     pgtype.Text{String: sessionId, Valid: true},
		StripeStatus: status,
	})

	return ReceiptDto{Link: s.getRedirectUrl(status)}, err
}

func (s *ReceiptService) getRedirectUrl(status db.StripeStatus) string {
	switch status {
	case db.StripeStatusFailed:
		return s.config.FrontendUrl + "/payment/failed"
	case db.StripeStatusSucceeded:
		return s.config.FrontendUrl + "/payment/success"
	case db.StripeStatusCanceled:
		return s.config.FrontendUrl + "/payment/canceled"
	}
	return s.config.FrontendUrl + "/payment/something-went-wrong"
}

func parseStripeStatus(status stripe.CheckoutSessionStatus) db.StripeStatus {
	switch status {
	case stripe.CheckoutSessionStatusComplete:
		return db.StripeStatusSucceeded
	case stripe.CheckoutSessionStatusExpired:
		return db.StripeStatusCanceled
	case stripe.CheckoutSessionStatusOpen:
		return db.StripeStatusPending
	default:
		return db.StripeStatusPending
	}
}

type CreateReceiptDto struct {
	Products map[int64]int32 `json:"products"`
}

type ReceiptDto struct {
	Link string `json:"link"`
}
