alter table if exists integrator_config
    add column if not exists extra json default '{}';