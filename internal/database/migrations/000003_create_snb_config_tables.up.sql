create table if not exists snb_config
(
    id         uuid               default uuid_generate_v4() primary key,
    endpoint   varchar,
    facility   varchar[],
    device     varchar[],
    created_at timestamp not null default NOW(),
    updated_at timestamp not null default NOW()
);

insert into snb_config (endpoint, facility, device)
VALUES ('http://localhost:8080', ARRAY ['fid'], ARRAY ['did']);

create trigger set_snb_config_updated_timestamp
    before update
    on snb_config
    for each row
execute procedure trigger_set_timestamp();