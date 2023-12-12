alter table integrator_transactions
    alter column id drop default;
alter table integrator_transactions
    drop column business_transaction_id;
alter table integrator_transactions
    rename column id to business_transaction_id;
