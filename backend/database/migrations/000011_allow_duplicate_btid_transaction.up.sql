alter table integrator_transactions
    rename column business_transaction_id to id;
alter table integrator_transactions
    add column if not exists business_transaction_id uuid;
UPDATE integrator_transactions
SET business_transaction_id = id
where true;

alter table integrator_transactions
    alter column id set default uuid_generate_v4();
