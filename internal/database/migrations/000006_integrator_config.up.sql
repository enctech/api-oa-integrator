create table if not exists integrator_config
(
    id                   uuid               default uuid_generate_v4() primary key,
    client_id            varchar,
    name                 varchar,
    sp_id                varchar,
    plaza_id             varchar,
    url                  varchar,
    insecure_skip_verify boolean            default false,
    created_at           timestamp not null default NOW(),
    updated_at           timestamp not null default NOW()
);

insert into integrator_config (client_id, sp_id, plaza_id, url, insecure_skip_verify)
VALUES ('CETA0109', 'ET', 'A01', 'http://47.254.241.45:8081', false);

create trigger set_integrator_config_timestamp
    before update
    on integrator_config
    for each row
execute procedure trigger_set_timestamp();