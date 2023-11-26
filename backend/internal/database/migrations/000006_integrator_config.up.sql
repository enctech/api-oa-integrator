create table if not exists integrator_config
(
    id                   uuid                    default uuid_generate_v4() primary key,
    client_id            varchar,
    provider_id          int unique     not null,
    name                 varchar unique not null,
    sp_id                varchar,
    plaza_id_map         json,
    url                  varchar,
    insecure_skip_verify boolean                 default false,
    created_at           timestamp      not null default NOW(),
    updated_at           timestamp      not null default NOW()
);

insert into integrator_config (client_id, name, provider_id, sp_id, plaza_id_map, url, insecure_skip_verify)
VALUES ('CETA0109', 'TNG', 2, 'ET', '{
  "1230": "A01"
}', 'http://47.254.241.45:8081', false);

create trigger set_integrator_config_timestamp
    before update
    on integrator_config
    for each row
execute procedure trigger_set_timestamp();