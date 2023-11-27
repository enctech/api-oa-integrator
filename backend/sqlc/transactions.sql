-- name: CreateOATransaction :one
insert into oa_transactions (businesstransactionid, lpn, customerid, jobid, facility, device, extra, entry_lane,
                             exit_lane)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning *;

-- name: GetOATransaction :one
select *
from oa_transactions
where businesstransactionid = $1;

-- name: GetOATransactions :many
select *
from oa_transactions
where lpn like concat('%', sqlc.arg(lpn)::text, '%')
  and jobid::text like concat('%', sqlc.arg(jobid)::text, '%')
  and facility::text like concat('%', sqlc.arg(facility)::text, '%')
  and entry_lane::text like concat('%', sqlc.arg(entry_lane)::text, '%')
  and (exit_lane::text LIKE concat('%', sqlc.arg(exit_lane)::text, '%') or (exit_lane is null and (sqlc.arg(exit_lane)::text) = ''))
  and created_at >= sqlc.arg(after)
  and created_at <= sqlc.arg(before);

-- name: GetOATransactionsCount :one
select count(*)
from oa_transactions
where lpn like concat('%', sqlc.arg(lpn)::text, '%')
  and jobid::text like concat('%', sqlc.arg(jobid)::text, '%')
  and facility::text like concat('%', sqlc.arg(facility)::text, '%')
  and entry_lane::text like concat('%', sqlc.arg(entry_lane)::text, '%')
  and (exit_lane::text LIKE concat('%', sqlc.arg(exit_lane)::text, '%') or (exit_lane is null and (sqlc.arg(exit_lane)::text) = ''))
  and created_at >= sqlc.arg(after)
  and created_at <= sqlc.arg(before);

-- name: UpdateOATransaction :one
update oa_transactions
set lpn        = coalesce($2, lpn),
    customerid = coalesce($3, customerid),
    jobid      = coalesce($4, jobid),
    facility   = coalesce($5, facility),
    device     = coalesce($6, device),
    extra      = coalesce($7, extra),
    exit_lane  = coalesce($8, exit_lane)
where businesstransactionid = $1
returning *;