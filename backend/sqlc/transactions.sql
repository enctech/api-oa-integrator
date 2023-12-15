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
  and (exit_lane::text LIKE concat('%', sqlc.arg(exit_lane)::text, '%') or
       (exit_lane is null and (sqlc.arg(exit_lane)::text) = ''))
  and created_at >= sqlc.arg(start_at)
  and created_at <= sqlc.arg(end_at)
order by created_at desc
limit $1 offset $2;

-- name: GetLatestOATransactions :many
select *
from oa_transactions
where updated_at >= sqlc.arg(start_at)
  and updated_at <= sqlc.arg(end_at)
order by updated_at desc
limit $1 offset $2;

-- name: GetOAEntryTransactions :one
select count(*)
from oa_transactions
where entry_lane is not null
  and (extra ->> 'steps' = 'leave_loop_entry_done'
    or exit_lane is not null)
  and created_at >= sqlc.arg(start_at)
  and created_at <= sqlc.arg(end_at);

-- name: GetOAExitTransactions :one
select count(*)
from oa_transactions
where exit_lane is not null
  and created_at >= sqlc.arg(start_at)
  and created_at <= sqlc.arg(end_at);

-- name: GetOATransactionsCount :one
select count(*)
from oa_transactions
where lpn like concat('%', sqlc.arg(lpn)::text, '%')
  and jobid::text like concat('%', sqlc.arg(jobid)::text, '%')
  and facility::text like concat('%', sqlc.arg(facility)::text, '%')
  and entry_lane::text like concat('%', sqlc.arg(entry_lane)::text, '%')
  and (exit_lane::text LIKE concat('%', sqlc.arg(exit_lane)::text, '%') or
       (exit_lane is null and (sqlc.arg(exit_lane)::text) = ''))
  and created_at >= sqlc.arg(start_at)
  and created_at <= sqlc.arg(end_at);

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

-- name: CreateIntegratorTransaction :one
with inserted_transaction as (
    insert into integrator_transactions (business_transaction_id, lpn, integrator_id, status, amount, error, tax_data,
                                         extra)
        values ($1, $2, $3, $4, $5, $6, $7, $8)
        returning *)
select *
from inserted_transaction
         inner join integrator_config on integrator_config.id = $3;

-- name: GetIntegratorTransactions :many
select it.*, ic.name as integrator_name
from integrator_transactions it
         inner join public.integrator_config ic on ic.id = it.integrator_id
where lpn like concat('%', sqlc.arg(lpn)::text, '%')
  and name::text like concat('%', sqlc.arg(name)::text, '%')
  and status::text like concat('%', sqlc.arg(status)::text, '%')
  and it.created_at >= sqlc.arg(start_at)
  and it.created_at <= sqlc.arg(end_at)
order by it.created_at desc
limit $1 offset $2;

-- name: GetIntegratorTransactionsCount :one
select count(*)
from integrator_transactions it
         inner join public.integrator_config ic on ic.id = it.integrator_id
where lpn like concat('%', sqlc.arg(lpn)::text, '%')
  and name::text like concat('%', sqlc.arg(name)::text, '%')
  and status::text like concat('%', sqlc.arg(status)::text, '%')
  and it.created_at >= sqlc.arg(start_at)
  and it.created_at <= sqlc.arg(end_at);

-- name: GetTotalTransactionAmount :one
select sum(amount)::numeric
from integrator_transactions
where status::text like concat('%', sqlc.arg(status)::text, '%')
  and created_at >= sqlc.arg(start_at)
  and created_at <= sqlc.arg(end_at);