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
		product.ReceiptID = resReceipt.ID
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

func (r *ReceiptRepository) GetReceipts(ctx context.Context, userId int64, limit int32, offset int32) ([]db.GetReceiptsRow, error) {
	q := db.New(r.db)
	return q.GetReceipts(ctx, db.GetReceiptsParams{UserID: userId, Limit: limit, Offset: offset})
}

func (r *ReceiptRepository) CountReceipts(ctx context.Context, userId int64) (int64, error) {
	q := db.New(r.db)
	return q.CountReceipts(ctx, userId)
}
