drop type if exists SURCHARGE_TYPE;

alter table integrator_config
    drop column if exists tax_rate,
    drop column if exists surcharge,
    drop column if exists surchange_type;