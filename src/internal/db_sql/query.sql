-- name: GetUser :one
select * from users where id = $1 limit 1;