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

-------------Region S&B Config start-------------
-- name: CreateSnbConfig :one
insert into snb_config (endpoint, facility, device, name, username, password)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: UpdateSnbConfig :one
update snb_config
set endpoint = coalesce($2, endpoint),
    facility = coalesce($3, facility),
    device   = coalesce($4, device),
    name     = coalesce($5, name),
    username = coalesce($6, username),
    password = coalesce($7, password)
where id = $1
returning *;

-- name: GetAllSnbConfig :many
select *
from snb_config;

-- name: GetSnbConfig :one
select *
from snb_config
where id = $1;

-- name: GetSnbConfigByFacilityAndDevice :one
select *
from snb_config
where sqlc.arg(facility)::text = any (facility)
  and sqlc.arg(device)::text = any (device);

-- name: DeleteSnbConfig :execresult
delete
from snb_config
where id = $1;
-------------Region S&B Config end---------------

-------------Region OA Transaction start-------------
-- name: CreateOATransaction :one
insert into oa_transactions (businesstransactionid, lpn, customerid, jobid, facility, device, extra, entry_lane,
                             exit_lane)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning *;

-- name: GetOATransaction :one
select *
from oa_transactions
where businesstransactionid = $1;

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

-------------Region OA Transaction end-------------

-------------Region Integrator Config start-------------
-- name: GetIntegratorConfig :one
select *
from integrator_config
where id = $1;

-- name: GetIntegratorConfigByClient :one
select *
from integrator_config
where client_id = $1;

-- name: GetIntegratorConfigByName :one
select *
from integrator_config
where name = $1;

-- name: CreateIntegratorConfig :one
insert into integrator_config (client_id, provider_id, name, sp_id, plaza_id_map, url, insecure_skip_verify)
values ($1, $2, $3, $4, $5, $6, $7)
returning *;

-- name: UpdateIntegratorConfig :one
update integrator_config
set provider_id          = coalesce($2, provider_id),
    client_id            = coalesce($3, client_id),
    name                 = coalesce($4, name),
    sp_id                = coalesce($5, sp_id),
    plaza_id_map         = coalesce($6, plaza_id_map),
    url                  = coalesce($7, url),
    insecure_skip_verify = coalesce($8, insecure_skip_verify)
where id = $1
returning *;

-- name: DeleteIntegratorConfig :execresult
delete
from integrator_config
where id = $1;
-------------Region Integrator Config end---------------

-------------Region Authentication start-------------
-- name: CreateUser :one
insert into users (username, password, permission)
values ($1, $2, $3)
returning *;

-- name: GetUser :one
select *
from users
where username = $1;

-- name: DeleteUser :execresult
delete
from users
where id = $1;
-------------Region Authentication end---------------