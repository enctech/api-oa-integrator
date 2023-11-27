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