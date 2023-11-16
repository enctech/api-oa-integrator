create table if not exists oa_transactions
(
    id                    uuid               default uuid_generate_v4() primary key,
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

create trigger set_oa_transaction_updated_timestamp
    before update
    on users
    for each row
execute procedure trigger_set_timestamp();

CREATE OR REPLACE FUNCTION delete_old_oa_transaction_data() RETURNS TRIGGER AS
$$
BEGIN
    DELETE
    FROM oa_transactions
    WHERE created_at < NOW() - interval '90 days';

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER oa_old_data_cleanup_trigger
    AFTER INSERT
    ON oa_transactions
    FOR EACH ROW
EXECUTE FUNCTION delete_old_oa_transaction_data();