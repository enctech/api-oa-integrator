create table if not exists integrator_config
(
    id                   uuid               default uuid_generate_v4() primary key,
    identifier           varchar   not null unique,
    customerGroup        varchar,
    url                  varchar,
    insecure_skip_verify boolean            default false,
    created_at           timestamp not null default NOW(),
    updated_at           timestamp not null default NOW()
);

create trigger set_integrator_config_timestamp
    before update
    on users
    for each row
execute procedure trigger_set_timestamp();