create table logs
(
    id         uuid primary key,
    level      varchar,
    message    varchar,
    fields     jsonb,
    created_at timestamp not null default NOW()
);

create table snb_config
(
    id       uuid primary key,
    name     varchar,
    username varchar,
    password varchar,
    endpoint varchar,
    facility varchar[],
    device   varchar[]
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
    provider_id          int,
    name                 varchar,
    sp_id                varchar,
    plaza_id_map         json,
    url                  varchar,
    insecure_skip_verify boolean            default false,
    created_at           timestamp not null default NOW(),
    updated_at           timestamp not null default NOW()
);

create table oa_transactions
(
    id                    uuid primary key,
    businessTransactionId varchar   not null unique,
    lpn                   varchar,
    customerId            varchar,
    jobId                 varchar,
    facility              varchar,
    device                varchar,
    extra                 jsonb,
    entry_lane            varchar,
    exit_lane             varchar,
    created_at            timestamp not null default NOW(),
    updated_at            timestamp not null default NOW()
);