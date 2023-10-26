-- name: CreateLog :one
INSERT INTO logs (module, info, extra)
VALUES ($1, $2, $3)
RETURNING *;


-- name: GetConfig :one
select *
from snb_config
where facility = $1
  and device = $2;