package product

import (
	"context"
	"shop-api/generated/db"
	"shop-api/internal/database/repositories"

	"github.com/jackc/pgx/v5/pgtype"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService(productRepo *repositories.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product CreateProductDto) (ProductDto, error) {
	result, err := s.productRepo.CreateProduct(ctx, db.CreateProductParams{
		Name:        product.Name,
		Price:       product.Price,
		Discount:    product.Discount,
		Description: pgtype.Text{String: product.Description, Valid: true},
		Image:       pgtype.Text{String: product.Image, Valid: true},
	})

	if err != nil {
		return ProductDto{}, err
	}

	return NewProductDto(result), nil
}

type CreateProductDto struct {
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Discount    int64  `json:"discount"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type ProductDto struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Discount    int64  `json:"discount"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func NewProductDto(product db.Product) ProductDto {
	return ProductDto{ID: product.ID, Name: product.Name, Price: product.Price, Discount: product.Discount, Description: product.Description.String, Image: product.Image.String}
}
