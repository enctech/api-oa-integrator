create table logs
(
    id     uuid primary key,
    module varchar not null,
    info   varchar not null,
    extra  jsonb
);