package receipt

import (
	"context"
	"shop-api/generated/db"
	"shop-api/internal/config"
	"shop-api/internal/database/repositories"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stripe/stripe-go/v84"
)

var CURRENCY = "usd"
var MODE = "payment"

type ReceiptService struct {
	receiptRepo  *repositories.ReceiptRepository
	productRepo  *repositories.ProductRepository
	stripeClient *stripe.Client
	config       *config.Config
}

func NewReceiptService(receiptRepo *repositories.ReceiptRepository, productRepo *repositories.ProductRepository, stripeClient *stripe.Client, config *config.Config) *ReceiptService {
	return &ReceiptService{receiptRepo: receiptRepo, productRepo: productRepo, stripeClient: stripeClient, config: config}
}

func (s *ReceiptService) GetReceipts(ctx context.Context, userId int64, limit int32, offset int32) (result ReceiptsPaginationDto, err error) {
	result.Limit = limit
	result.Offset = offset
	count, err := s.receiptRepo.CountReceipts(ctx, userId)
	if err != nil {
		return ReceiptsPaginationDto{}, err
	}
	result.Total = count
	res, err := s.receiptRepo.GetReceipts(ctx, userId, limit, offset)
	if err != nil {
		return ReceiptsPaginationDto{}, err
	}
	for _, r := range res {
		idx := slices.IndexFunc(result.Receipts, func(temp ReceiptDto) bool { return temp.ID == r.ID })
		if idx == -1 {
			result.Receipts = append(result.Receipts, NewReceiptDto(r))
		} else {
			result.Receipts[idx].Products = append(result.Receipts[idx].Products, NewReceiptProductDto(r))
		}
	}
	return result, nil
}

func (s *ReceiptService) CreateReceipt(ctx context.Context, userId int64, model CreateReceiptDto) (LinkDto, error) {

	ids := make([]int64, 0)
	for id := range model.Products {
		ids = append(ids, id)
	}

	products, err := s.productRepo.GetProductByIds(ctx, ids)
	if err != nil {
		return LinkDto{}, err
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

	successUrl := s.config.BackendUrl + "/api/payment/success?sessionId={{CHECKOUT_SESSION_ID}}"
	cancelUrl := s.config.BackendUrl + "/api/payment/cancel?sessionId={{CHECKOUT_SESSION_ID}}"

	session, err := s.stripeClient.V1CheckoutSessions.Create(ctx, &stripe.CheckoutSessionCreateParams{
		LineItems:  stripeParams,
		Mode:       &MODE,
		SuccessURL: &successUrl,
		CancelURL:  &cancelUrl,
	})

	if err != nil {
		return LinkDto{}, err
	}

	_, _, err = s.receiptRepo.CreateReceipt(ctx, db.CreateReceiptParams{
		UserID:       userId,
		SumPrice:     sumPrice,
		SumDiscount:  sumDiscount,
		StripeID:     pgtype.Text{String: session.ID, Valid: true},
		StripeStatus: db.StripeStatusPending,
	}, productsParams)

	if err != nil {
		return LinkDto{}, err
	}

	return LinkDto{Link: session.URL}, nil
}

func (s *ReceiptService) CancelReceipt(ctx context.Context, sessionId string) (LinkDto, error) {
	session, err := s.stripeClient.V1CheckoutSessions.Retrieve(ctx, sessionId, nil)
	if err != nil {
		return LinkDto{}, err
	}

	status := parseStripeStatus(session.Status)

	_, err = s.productRepo.UpdateReceiptStatus(ctx, db.UpdateReceiptStatusParams{
		StripeID:     pgtype.Text{String: sessionId, Valid: true},
		StripeStatus: status,
	})

	return LinkDto{Link: s.getRedirectUrl(status)}, err
}

func (s *ReceiptService) SuccessReceipt(ctx context.Context, sessionId string) (LinkDto, error) {
	session, err := s.stripeClient.V1CheckoutSessions.Retrieve(ctx, sessionId, nil)
	if err != nil {
		return LinkDto{}, err
	}

	status := parseStripeStatus(session.Status)

	_, err = s.productRepo.UpdateReceiptStatus(ctx, db.UpdateReceiptStatusParams{
		StripeID:     pgtype.Text{String: sessionId, Valid: true},
		StripeStatus: status,
	})

	return LinkDto{Link: s.getRedirectUrl(status)}, err
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

type LinkDto struct {
	Link string `json:"link"`
}

type ReceiptsPaginationDto struct {
	Receipts []ReceiptDto `json:"receipts"`
	Total    int64        `json:"total"`
	Limit    int32        `json:"limit"`
	Offset   int32        `json:"offset"`
}

type ReceiptDto struct {
	ID           int64               `json:"id"`
	SumPrice     int64               `json:"sumPrice"`
	SumDiscount  int64               `json:"sumDiscount"`
	CreatedAt    time.Time           `json:"createdAt"`
	StripeID     string              `json:"stripeId"`
	StripeStatus db.StripeStatus     `json:"stripeStatus"`
	Products     []ReceiptProductDto `json:"products"`
}

func NewReceiptDto(r db.GetReceiptsRow) ReceiptDto {
	return ReceiptDto{ID: r.ID, SumPrice: r.SumPrice, SumDiscount: r.SumDiscount, CreatedAt: r.CreatedAt.Time, StripeID: r.StripeID.String, StripeStatus: r.StripeStatus, Products: []ReceiptProductDto{NewReceiptProductDto(r)}}
}

type ReceiptProductDto struct {
	ProductID int64  `json:"productId"`
	Quantity  int32  `json:"quantity"`
	Price     int64  `json:"price"`
	Discount  int64  `json:"discount"`
	Name      string `json:"name"`
}

func NewReceiptProductDto(r db.GetReceiptsRow) ReceiptProductDto {
	return ReceiptProductDto{ProductID: r.ProductID.Int64, Quantity: r.Quantity.Int32, Price: r.Price.Int64, Discount: r.Discount.Int64, Name: r.Name.String}
}
