-- name: GetSiteByConfigAndFacility :one
-- Which site of this config owns the facility. Unique per config, so at most
-- one row can match.
select s.*
from integrator_site s
         join integrator_site_facility f on f.site_id = s.id
where s.integrator_config_id = $1
  and f.facility = $2;

-- name: GetSitesByConfig :many
select s.*,
       coalesce(
                       array_agg(f.facility order by f.facility)
                       filter (where f.facility is not null),
                       '{}'
       )::varchar[] as facilities
from integrator_site s
         left join integrator_site_facility f on f.site_id = s.id
where s.integrator_config_id = $1
group by s.id
order by s.created_at;

-- name: CreateSite :one
insert into integrator_site (id, integrator_config_id, provider_id, client_id, vendor_location_id)
values (gen_random_uuid(), $1, $2, $3, $4)
returning *;

-- name: AddSiteFacility :exec
insert into integrator_site_facility (site_id, integrator_config_id, facility)
values ($1, $2, $3);

-- name: DeleteSitesByConfig :exec
-- The form submits every site at once, so a save replaces the whole set.
delete
from integrator_site
where integrator_config_id = $1;
