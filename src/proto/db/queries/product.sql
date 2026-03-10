-- name: CreateProduct :one
INSERT INTO public.products(name, price, discount, description, image)
	VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetProducts :many
Select id, name, price, discount, description, image, created_at from products order by id limit $1 offset $2;

-- name: CountProducts :one
Select count(*) from products;

-- name: GetProductsByName :many
Select id, name, price, discount, description, image, created_at from products
where name ilike $1
order by id limit $2 offset $3;

-- name: CountProductsByName :one
Select count(*) from products
where name ilike $1;