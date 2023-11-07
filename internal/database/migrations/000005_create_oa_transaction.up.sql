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
    created_at            timestamp not null default NOW(),
    updated_at            timestamp not null default NOW()
);

create trigger set_oa_transaction_updated_timestamp
    before update
    on users
    for each row
execute procedure trigger_set_timestamp();