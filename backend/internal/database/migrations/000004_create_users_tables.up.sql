create table if not exists users
(
    id         uuid               default uuid_generate_v4() primary key,
    username   varchar   not null unique,
    password   varchar,
    permission varchar,
    created_at timestamp not null default NOW(),
    updated_at timestamp not null default NOW()
);

create trigger set_users_updated_timestamp
    before update
    on users
    for each row
execute procedure trigger_set_timestamp();