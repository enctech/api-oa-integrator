create table if not exists logs
(
    id         uuid               default uuid_generate_v4() primary key,
    level      varchar,
    message    varchar,
    fields     jsonb,
    created_at timestamp not null default NOW(),
    updated_at timestamp not null default NOW()
);

create trigger set_logs_updated_timestamp
    before update
    on logs
    for each row
execute procedure trigger_set_timestamp();

CREATE OR REPLACE FUNCTION delete_old_data() RETURNS TRIGGER AS
$$
BEGIN
    DELETE
    FROM logs
    WHERE created_at < NOW() - interval '90 days';

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER data_cleanup_trigger
    AFTER INSERT
    ON logs
    FOR EACH ROW
EXECUTE FUNCTION delete_old_data();