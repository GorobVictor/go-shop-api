-- name: CreateUser :one
INSERT INTO public.users(first_name, last_name, email, password_hash, user_role)
	VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 limit 1;

-- name: GetUserProfile :one
SELECT ID, first_name, last_name, email, user_role, created_at FROM users WHERE id = $1 limit 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 limit 1;

-- name: AnyEmail :one
SELECT email FROM users WHERE email = $1 limit 1;

-- name: GetUsers :many
Select ID, first_name, last_name, email, user_role, created_at from users order by id limit $1 offset $2;

-- name: CountUsers :one
Select count(*) from users;