-- name: CreateUser :one
INSERT INTO public.users(first_name, last_name, email, password_hash)
	VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 limit 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 limit 1;

-- name: AnyEmail :one
SELECT email FROM users WHERE email = $1 limit 1;