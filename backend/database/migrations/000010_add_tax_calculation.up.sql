drop type if exists SURCHARGE_TYPE;
create type SURCHARGE_TYPE as enum ('percentage', 'exact');

alter table integrator_config
    add column if not exists tax_rate       numeric        default 0,
    add column if not exists surcharge      numeric        default 0,
    add column if not exists surchange_type SURCHARGE_TYPE default 'exact';