package routes

import (
	"context"
	"encoding/json"
	"net/http"
	customerrors "shop-api/internal/custom_errors"
	"shop-api/internal/usecase/product"

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
	var model product.CreateProductDto
	ReadBody(w, r, &model)

	result, err := h.productSvc.CreateProduct(context.Background(), model)
	if err != nil {
		panic(customerrors.NewInternalServerError())
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
	result, err := h.productSvc.GetProducts(context.Background(), GetQueryInt32(r, "limit"), GetQueryInt32(r, "offset"))
	if err != nil {
		panic(customerrors.NewInternalServerError())
	}

	json.NewEncoder(w).Encode(result)
}
