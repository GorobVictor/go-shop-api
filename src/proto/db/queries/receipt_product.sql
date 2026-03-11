-- name: CreateReceiptProduct :one
INSERT INTO public.receipt_products(receipt_id, product_id, quantity, price, discount)
	VALUES ($1, $2, $3, $4, $5) RETURNING *;