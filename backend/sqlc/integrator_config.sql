-- name: GetIntegratorConfigs :many
select *
from integrator_config
where deleted_at is null;

-- name: GetIntegratorConfig :one
select *
from integrator_config
where id = $1
  and deleted_at is null;

-- name: GetIntegratorConfigByClient :one
select *
from integrator_config
where client_id = $1
  and deleted_at is null;

-- name: GetIntegratorConfigByName :one
select *
from integrator_config
where name = $1
  and deleted_at is null;

-- name: CreateIntegratorConfig :one
insert into integrator_config (client_id, provider_id, name, sp_id, plaza_id_map, url, insecure_skip_verify,
                               integrator_name, extra, tax_rate, surcharge, surchange_type, display_name)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
returning *;

-- name: UpdateIntegratorConfig :one
update integrator_config
set provider_id          = coalesce($2, provider_id),
    client_id            = coalesce($3, client_id),
    name                 = coalesce($4, name),
    sp_id                = coalesce($5, sp_id),
    plaza_id_map         = coalesce($6, plaza_id_map),
    url                  = coalesce($7, url),
    insecure_skip_verify = coalesce($8, insecure_skip_verify),
    integrator_name      = coalesce($9, integrator_name),
    extra                = coalesce($10, extra),
    tax_rate             = coalesce($11, tax_rate),
    surcharge            = coalesce($12, surcharge),
    surchange_type       = coalesce($13, surchange_type),
    display_name         = coalesce($14, display_name)
where id = $1
  and deleted_at is null
returning *;

-- name: DeleteIntegratorConfig :execresult
-- Soft delete: transactions reference this row and hold financial data.
update integrator_config
set deleted_at = now()
where id = $1
  and deleted_at is null;