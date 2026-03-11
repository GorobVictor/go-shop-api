package repositories

import (
	"context"
	"log"
	"shop-api/generated/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product db.CreateProductParams) (db.Product, error) {
	q := db.New(r.db)
	return q.CreateProduct(ctx, product)
}

func (r *ProductRepository) GetProducts(ctx context.Context, limit int32, offset int32) ([]db.Product, error) {
	q := db.New(r.db)
	products, err := q.GetProducts(ctx, db.GetProductsParams{Limit: limit, Offset: offset})

	if err != nil {
		log.Println(err.Error())
		if err.Error() == "no rows in result set" {
			return products, nil
		}
	}

	return products, err
}

func (r *ProductRepository) CountProducts(ctx context.Context) (int64, error) {
	q := db.New(r.db)
	return q.CountProducts(ctx)
}

func (r *ProductRepository) GetProductsByName(ctx context.Context, name string, limit int32, offset int32) ([]db.Product, error) {
	q := db.New(r.db)
	products, err := q.GetProductsByName(ctx, db.GetProductsByNameParams{Name: name, Limit: limit, Offset: offset})

	if err != nil {
		log.Println(err.Error())
		if err.Error() == "no rows in result set" {
			return products, nil
		}
	}

	return products, err
}

func (r *ProductRepository) CountProductsByName(ctx context.Context, name string) (int64, error) {
	q := db.New(r.db)
	return q.CountProductsByName(ctx, name)
}

func (r *ProductRepository) GetProductByIds(ctx context.Context, ids []int64) ([]db.Product, error) {
	q := db.New(r.db)
	return q.GetProductByIds(ctx, ids)
}

func (r *ProductRepository) UpdateReceiptStatus(ctx context.Context, receipt db.UpdateReceiptStatusParams) (db.Receipt, error) {
	q := db.New(r.db)
	return q.UpdateReceiptStatus(ctx, receipt)
}
