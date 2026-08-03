-- Integrator configs are referenced by integrator_transactions and
-- oa_transactions, which hold amounts and tax data. A hard delete either
-- fails on those foreign keys or destroys financial history, so configs are
-- retired by setting deleted_at instead.
ALTER TABLE integrator_config
    ADD COLUMN IF NOT EXISTS deleted_at timestamp;

-- The unique constraints have to ignore retired rows, otherwise a deleted
-- config keeps its name and provider_id reserved forever and the same site
-- can never be re-added.
--
-- Recreated as partial unique indexes under the original constraint names so
-- that the duplicate-detection in createIntegratorConfig, which matches on
-- "integrator_config_name_key" and "integrator_config_provider_id_key" in the
-- error text, keeps working.
ALTER TABLE integrator_config
    DROP CONSTRAINT IF EXISTS integrator_config_name_key;
ALTER TABLE integrator_config
    DROP CONSTRAINT IF EXISTS integrator_config_provider_id_key;

CREATE UNIQUE INDEX IF NOT EXISTS integrator_config_name_key
    ON integrator_config (name)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS integrator_config_provider_id_key
    ON integrator_config (provider_id)
    WHERE deleted_at IS NULL;
