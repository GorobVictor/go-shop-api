-- name: CreateReceipt :one
INSERT INTO public.receipts(user_id, sum_price, sum_discount, stripe_id, stripe_status)
	VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateReceiptStatus :one
UPDATE public.receipts
SET stripe_status = $2
WHERE stripe_id = $1
RETURNING *;

-- name: CountReceipts :one
select count(*) from public.receipts where user_id = $1;

-- name: GetReceipts :many
select r.*, rp.*, p.name from 
(select * from public.receipts where user_id = $1 order by id desc limit $2 offset $3) r
left join public.receipt_products rp on r.id = rp.receipt_id
left join public.products p on rp.product_id = p.id
order by r.id desc;