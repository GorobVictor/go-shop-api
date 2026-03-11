package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"shop-api/internal/usecase/receipt"
	"strconv"

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
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	result, err := h.receiptSvc.CreateReceipt(context.Background(), GetUserId(w, r), model)
	if err != nil {
		writeError(w, err, 500)
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
	sessionId := r.URL.Query().Get("sessionId")
	if sessionId == "" {
		http.Error(w, "Session ID is required", 400)
		return
	}

	result, err := h.receiptSvc.SuccessReceipt(context.Background(), sessionId)
	if err != nil {
		writeError(w, err, 500)
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
	sessionId := r.URL.Query().Get("sessionId")
	if sessionId == "" {
		http.Error(w, "Session ID is required", 400)
		return
	}

	result, err := h.receiptSvc.CancelReceipt(context.Background(), sessionId)
	if err != nil {
		writeError(w, err, 500)
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
	queries := r.URL.Query()
	limit, err := strconv.ParseInt(queries.Get("limit"), 10, 32)
	if err != nil {
		writeBadRequest(w, err)
	}
	offset, err := strconv.ParseInt(queries.Get("offset"), 10, 64)
	if err != nil {
		writeBadRequest(w, err)
	}

	result, err := h.receiptSvc.GetReceipts(context.Background(), GetUserId(w, r), int32(limit), int32(offset))

	if err != nil {
		writeError(w, err, 500)
		return
	}

	json.NewEncoder(w).Encode(result)
}
