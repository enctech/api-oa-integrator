create extension if not exists "uuid-ossp";

create or replace function trigger_set_timestamp()
    returns trigger as
$$
begin
    new.updated_at = now();
    return new;
end;
$$ language plpgsql;