create table if not exists integrator_transactions
(
    business_transaction_id uuid primary key,
    lpn                     varchar,
    integrator_id           uuid,
    status                  varchar,
    amount                  numeric,
    error                   varchar,
    tax_data                jsonb,
    extra                   jsonb,
    created_at              timestamp not null default NOW(),
    updated_at              timestamp not null default NOW(),

    FOREIGN KEY (integrator_id) REFERENCES integrator_config (id)
);

create trigger set_integrator_transactions_timestamp
    before update
    on integrator_transactions
    for each row
execute procedure trigger_set_timestamp();

CREATE OR REPLACE FUNCTION delete_old_integrator_transaction_data() RETURNS TRIGGER AS
$$
BEGIN
    DELETE
    FROM integrator_transactions
    WHERE created_at < NOW() - interval '100 days';

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER integrator_transaction_cleanup_trigger
    AFTER INSERT
    ON integrator_transactions
    FOR EACH ROW
EXECUTE FUNCTION delete_old_integrator_transaction_data();