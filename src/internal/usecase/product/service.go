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

func (s *ProductService) GetProducts(ctx context.Context, limit int32, offset int32) (result ProductsPaginationDto, err error) {
	result.Limit = limit
	result.Offset = offset

	result.Total, err = s.productRepo.CountProducts(ctx)
	if err != nil {
		return result, err
	}

	products, err := s.productRepo.GetProducts(ctx, limit, offset)
	if err != nil {
		return result, err
	}

	for _, product := range products {
		result.Products = append(result.Products, NewProductDto(product))
	}

	return result, nil
}

type ProductsPaginationDto struct {
	Products []ProductDto `json:"products"`
	Total    int64        `json:"total"`
	Limit    int32        `json:"limit"`
	Offset   int32        `json:"offset"`
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
