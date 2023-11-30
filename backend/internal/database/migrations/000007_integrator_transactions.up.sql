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

    FOREIGN KEY (integrator_id) REFERENCES integrator_config (id)
);

create trigger set_integrator_transactions_timestamp
    before update
    on integrator_transactions
    for each row
execute procedure trigger_set_timestamp();