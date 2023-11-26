create table if not exists integrator_transactions
(
    business_transaction_id uuid primary key,
    lpn                     varchar,
    status                  varchar,
    amount                  varchar
);

create trigger set_integrator_transactions_timestamp
    before update
    on integrator_transactions
    for each row
execute procedure trigger_set_timestamp();