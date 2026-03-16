package routes

import (
	"context"
	"encoding/json"
	"net/http"
	customerrors "shop-api/internal/custom_errors"
	"shop-api/internal/usecase/receipt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

type PaymentHandler struct {
	receiptSvc *receipt.ReceiptService
	tokenAuth  *jwtauth.JWTAuth
}

func NewPaymentHandler(receiptSvc *receipt.ReceiptService, tokenAuth *jwtauth.JWTAuth) *PaymentHandler {
	return &PaymentHandler{receiptSvc: receiptSvc, tokenAuth: tokenAuth}
}

func (h *PaymentHandler) Payment(r *chi.Mux) {
	r.Route("/api/payment", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/success", h.successPayment)
			r.Get("/cancel", h.cancelPayment)
		})
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(h.tokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Post("/create", h.createPayment)
			r.Get("/get", h.getReceipts)
		})
	})
}

// Create Payment
// @Summary Create Payment
// @Tags payment
// @Security ApiKeyAuth
// @Param payment body receipt.CreateReceiptDto true "Payment details"
// @Success 200 {object} receipt.ReceiptDto
// @Router /payment/create [post]
func (h *PaymentHandler) createPayment(w http.ResponseWriter, r *http.Request) {
	var model receipt.CreateReceiptDto
	ReadBody(w, r, &model)
	userId, err := GetUserId(w, r)
	if err != nil {
		CheckError(w, err)
		return
	}

	result, err := h.receiptSvc.CreateReceipt(context.Background(), userId, model)
	if err != nil {
		CheckError(w, err)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// Success Payment
// @Summary Success Payment
// @Tags payment
// @Security ApiKeyAuth
// @Param sessionId query string true "Session ID"
// @Success 200 {object} receipt.ReceiptDto
// @Router /payment/success [get]
func (h *PaymentHandler) successPayment(w http.ResponseWriter, r *http.Request) {
	sessionId, err := parseSessionId(r)
	if err != nil {
		CheckError(w, err)
		return
	}
	result, err := h.receiptSvc.SuccessReceipt(context.Background(), sessionId)
	if err != nil {
		CheckError(w, err)
		return
	}

	http.Redirect(w, r, result.Link, http.StatusSeeOther)
}

// Cancel Payment
// @Summary Cancel Payment
// @Tags payment
// @Security ApiKeyAuth
// @Param sessionId query string true "Session ID"
// @Success 200 {object} receipt.ReceiptDto
// @Router /payment/cancel [get]
func (h *PaymentHandler) cancelPayment(w http.ResponseWriter, r *http.Request) {
	sessionId, err := parseSessionId(r)
	if err != nil {
		CheckError(w, err)
		return
	}
	result, err := h.receiptSvc.CancelReceipt(context.Background(), sessionId)
	if err != nil {
		CheckError(w, err)
		return
	}
	http.Redirect(w, r, result.Link, http.StatusSeeOther)
}

// Get Receipts
// @Summary Get Receipts
// @Tags payment
// @Security ApiKeyAuth
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Success 200 {object} receipt.ReceiptsPaginationDto
// @Router /payment/get [get]
func (h *PaymentHandler) getReceipts(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserId(w, r)
	if err != nil {
		CheckError(w, err)
		return
	}
	result, err := h.receiptSvc.GetReceipts(context.Background(), userId, GetQueryInt32(r, "limit"), GetQueryInt32(r, "offset"))

	if err != nil {
		CheckError(w, err)
		return
	}

	WriteOkResponse(w, result)
}

func parseSessionId(r *http.Request) (string, error) {
	sessionId := r.URL.Query().Get("sessionId")
	if sessionId == "" {
		return "", &customerrors.BadRequestError{Message: "Session ID is required"}
	}

	if len(sessionId) > 0 && sessionId[0] == '{' {
		sessionId = sessionId[1:]
	}
	if len(sessionId) > 0 && sessionId[len(sessionId)-1] == '}' {
		sessionId = sessionId[:len(sessionId)-1]
	}
	return sessionId, nil
}
