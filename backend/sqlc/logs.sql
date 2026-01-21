-- name: CountLogs :one
select count(*)
from logs
where created_at >= sqlc.arg(after)
  and created_at <= sqlc.arg(before);

-- name: CreateLog :execresult
insert into logs (level, message, fields, created_at)
values ($1, $2, $3, $4)
returning *;

-- name: GetLogs :many
select *
from logs
where created_at >= sqlc.arg(after)
  and created_at <= sqlc.arg(before)
  and (sqlc.arg(message)::text = '' or message ilike concat('%', sqlc.arg(message)::text, '%'))
  and (sqlc.arg(fields)::text = '' or fields::text ilike concat('%', sqlc.arg(fields)::text, '%'))
order by created_at desc
limit $1 offset $2;