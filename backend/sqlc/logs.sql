-- name: CountLogs :one
select count(*)
from logs
where created_at >= sqlc.arg(after)
  and created_at <= sqlc.arg(before);

-- name: CreateLog :one
insert into logs (level, message, fields, created_at)
values ($1, $2, $3, $4)
returning *;

-- name: GetLogs :many
select *
from logs
where message like sqlc.arg(message)::text
  and fields::text like sqlc.arg(fields)::text
  and created_at >= sqlc.arg(after)
  and created_at <= sqlc.arg(before)
order by created_at desc
limit $1 offset $2;