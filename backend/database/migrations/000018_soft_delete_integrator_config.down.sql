DROP INDEX IF EXISTS integrator_config_name_key;
DROP INDEX IF EXISTS integrator_config_provider_id_key;

ALTER TABLE integrator_config
    DROP COLUMN IF EXISTS deleted_at;

-- These will fail if a retired config shares a name or provider_id with a
-- live one, which is possible once a site has been deleted and re-added.
-- That is deliberate: failing loudly is better than deleting config rows
-- that transactions still reference. Resolve the duplicates by hand first.
ALTER TABLE integrator_config
    ADD CONSTRAINT integrator_config_name_key UNIQUE (name);
ALTER TABLE integrator_config
    ADD CONSTRAINT integrator_config_provider_id_key UNIQUE (provider_id);
