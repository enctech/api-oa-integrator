create table logs
(
    id      uuid primary key,
    level   varchar,
    message varchar,
    fields  jsonb,
    created_at timestamp not null default NOW()
);

create table snb_config
(
    id       uuid primary key,
    endpoint varchar,
    facility varchar,
    device   varchar
);

create table users
(
    id         uuid primary key,
    username   varchar,
    password   varchar,
    permission varchar
);

create table integrator_config
(
    id                   uuid primary key,
    client_id            varchar,
    sp_id                varchar,
    plaza_id             varchar,
    url                  varchar,
    insecure_skip_verify boolean            default false,
    created_at           timestamp not null default NOW(),
    updated_at           timestamp not null default NOW()
);