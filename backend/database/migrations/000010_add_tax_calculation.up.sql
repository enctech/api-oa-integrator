CREATE TYPE SURCHARGE_TYPE AS ENUM ('percentage', 'exact');

alter table integrator_config
    add column if not exists tax_rate       varchar(255)   default '0',
    add column if not exists surcharge      varchar(255)   default '0',
    add column if not exists surchange_type SURCHARGE_TYPE default 'exact';