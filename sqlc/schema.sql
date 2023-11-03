create table logs
(
    id     uuid primary key,
    module varchar not null,
    info   varchar not null,
    extra  jsonb
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