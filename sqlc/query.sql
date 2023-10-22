-- name: CreateLog :one
INSERT INTO logs (module, info, extra)
VALUES ($1, $2, $3)
RETURNING *;