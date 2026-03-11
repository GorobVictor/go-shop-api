-- name: CreateReceipt :one
INSERT INTO public.receipts(user_id, sum_price, sum_discount, stripe_id, stripe_status)
	VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateReceiptStatus :one
UPDATE public.receipts
SET stripe_status = $2
WHERE stripe_id = $1
RETURNING *;