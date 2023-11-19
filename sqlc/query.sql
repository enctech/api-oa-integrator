-- name: CreateLog :one
INSERT INTO logs (level, message, fields, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-------------Region S&B Config start-------------
-- name: CreateSnbConfig :one
insert into snb_config (endpoint, facility, device, name)
values ($1, $2, $3, $4)
returning *;

-- name: UpdateSnbConfig :one
update snb_config
set endpoint = coalesce($2, endpoint),
    facility = coalesce($3, facility),
    device   = coalesce($4, device),
    name   = coalesce($5, name)
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
where sqlc.arg(facility)::text = any(facility)
  and sqlc.arg(device)::text = any(device);

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
where client_id = $1;
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