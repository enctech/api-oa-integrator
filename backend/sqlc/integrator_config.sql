-- name: GetIntegratorConfigs :many
select *
from integrator_config;

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
insert into integrator_config (client_id, provider_id, name, sp_id, plaza_id_map, url, insecure_skip_verify,
                               integrator_name, extra, tax_rate, surcharge, surchange_type)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
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
    surchange_type       = coalesce($13, surchange_type)
where id = $1
returning *;

-- name: DeleteIntegratorConfig :execresult
delete
from integrator_config
where id = $1;