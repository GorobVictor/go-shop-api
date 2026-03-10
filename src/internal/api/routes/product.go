package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"shop-api/internal/usecase/product"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

type ProductHandler struct {
	productSvc *product.ProductService
	tokenAuth  *jwtauth.JWTAuth
}

func NewProductHandler(productSvc *product.ProductService, tokenAuth *jwtauth.JWTAuth) *ProductHandler {
	return &ProductHandler{productSvc: productSvc, tokenAuth: tokenAuth}
}

func (h *ProductHandler) Products(r *chi.Mux) {
	r.Route("/api/products", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Get("/get", h.getProducts)
		})
		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(h.tokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Use(GetAdminMiddleware)

			r.Post("/create", h.createProduct)
		})
	})
}

// Create Product
// @Summary Create Product
// @Tags products
// @Security ApiKeyAuth
// @Param product body product.CreateProductDto true "Product details"
// @Success 200 {object} product.ProductDto
// @Router /products/create [post]
func (h *ProductHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	var product product.CreateProductDto
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	result, err := h.productSvc.CreateProduct(context.Background(), product)
	if err != nil {
		writeError(w, err, 500)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// Get Products
// @Summary Get Products
// @Tags products
// @Security ApiKeyAuth
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Success 200 {object} product.ProductsPaginationDto
// @Router /products/get [get]
func (h *ProductHandler) getProducts(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	limit, err := strconv.ParseInt(queries.Get("limit"), 10, 32)
	if err != nil {
		writeBadRequest(w, err)
	}
	offset, err := strconv.ParseInt(queries.Get("offset"), 10, 32)
	if err != nil {
		writeBadRequest(w, err)
	}

	result, err := h.productSvc.GetProducts(context.Background(), int32(limit), int32(offset))
	if err != nil {
		writeError(w, err, 500)
		return
	}

	json.NewEncoder(w).Encode(result)
}
