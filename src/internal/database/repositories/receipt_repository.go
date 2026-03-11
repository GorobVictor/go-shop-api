package repositories

import (
	"context"
	"shop-api/generated/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReceiptRepository struct {
	db *pgxpool.Pool
}

func NewReceiptRepository(db *pgxpool.Pool) *ReceiptRepository {
	return &ReceiptRepository{db: db}
}

func (r *ReceiptRepository) CreateReceipt(ctx context.Context, receipt db.CreateReceiptParams, products []db.CreateReceiptProductParams) (db.Receipt, []db.ReceiptProduct, error) {
	q := db.New(r.db)
	tx, err := r.db.Begin(ctx)

	if err != nil {
		return db.Receipt{}, []db.ReceiptProduct{}, err
	}

	defer tx.Rollback(ctx)

	q = q.WithTx(tx)

	resReceipt, err := q.CreateReceipt(ctx, receipt)

	if err != nil {
		return db.Receipt{}, []db.ReceiptProduct{}, err
	}

	resProducts := make([]db.ReceiptProduct, 0)

	for _, product := range products {
		resProductReceipt, err := q.CreateReceiptProduct(ctx, product)
		resProducts = append(resProducts, resProductReceipt)

		if err != nil {
			return db.Receipt{}, []db.ReceiptProduct{}, err
		}
	}

	err = tx.Commit(ctx)

	if err != nil {
		return db.Receipt{}, []db.ReceiptProduct{}, err
	}

	return resReceipt, resProducts, nil
}
