-- name: CountLogs :one
select count(*)
from logs
where created_at >= sqlc.arg(after)
  and created_at <= sqlc.arg(before);

-- name: CreateLog :execresult
insert into logs (level, message, fields, created_at)
values ($1, $2, $3, $4)
returning *;

-- name: CreateLogs :execresult
insert into logs (level, message, fields, created_at)
select l, m, f, t
from unnest($1::text[]) with ordinality as a(l, i)
         join unnest($2::text[]) with ordinality as b(m, i2) on a.i = b.i2
         join unnest($3::jsonb[]) with ordinality as c(f, i3) on a.i = i3
         join unnest($4::timestamptz[]) with ordinality as d(t, i4) on a.i = i4
returning *;

-- name: GetLogs :many
select *
from logs
where message like concat('%', sqlc.arg(message)::text, '%')
  and fields::text like concat('%', sqlc.arg(fields)::text, '%')
  and created_at >= sqlc.arg(after)
  and created_at <= sqlc.arg(before)
order by created_at desc
limit $1 offset $2;