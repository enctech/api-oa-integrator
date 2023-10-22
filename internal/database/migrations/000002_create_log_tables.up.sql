create table if not exists logs
(
    id     uuid default uuid_generate_v4() primary key,
    module varchar,
    info   varchar,
    extra  jsonb,
    created_at timestamp not null default NOW(),
    updated_at timestamp not null default NOW()
);

create trigger set_logs_updated_timestamp
    before update on logs
    for each row
execute procedure trigger_set_timestamp();