package repositories

import (
	"context"
	"errors"
	"log"
	"shop-api/generated/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
	q  *db.Queries
}

func NewProductRepository(db *pgxpool.Pool, q *db.Queries) *ProductRepository {
	return &ProductRepository{db: db, q: q}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product db.CreateProductParams) (db.Product, error) {
	return r.q.CreateProduct(ctx, product)
}

func (r *ProductRepository) GetProducts(ctx context.Context, limit int32, offset int32) ([]db.Product, error) {
	products, err := r.q.GetProducts(ctx, db.GetProductsParams{Limit: limit, Offset: offset})

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return products, nil
		}
	}

	return products, err
}

func (r *ProductRepository) CountProducts(ctx context.Context) (int64, error) {
	return r.q.CountProducts(ctx)
}

func (r *ProductRepository) GetProductsByName(ctx context.Context, name string, limit int32, offset int32) ([]db.Product, error) {
	products, err := r.q.GetProductsByName(ctx, db.GetProductsByNameParams{Name: name, Limit: limit, Offset: offset})

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return products, nil
		}
	}

	return products, err
}

func (r *ProductRepository) CountProductsByName(ctx context.Context, name string) (int64, error) {
	return r.q.CountProductsByName(ctx, name)
}

func (r *ProductRepository) GetProductByIds(ctx context.Context, ids []int64) ([]db.Product, error) {
	return r.q.GetProductByIds(ctx, ids)
}

func (r *ProductRepository) UpdateReceiptStatus(ctx context.Context, receipt db.UpdateReceiptStatusParams) (db.Receipt, error) {
	return r.q.UpdateReceiptStatus(ctx, receipt)
}
