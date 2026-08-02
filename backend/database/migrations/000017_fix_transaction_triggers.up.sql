-- Drop problematic triggers that run DELETE on every insert
DROP TRIGGER IF EXISTS oa_old_data_cleanup_trigger ON oa_transactions;
DROP FUNCTION IF EXISTS delete_old_oa_transaction_data();

DROP TRIGGER IF EXISTS integrator_transaction_cleanup_trigger ON integrator_transactions;
DROP FUNCTION IF EXISTS delete_old_integrator_transaction_data();

-- Add indexes on created_at for efficient time-range queries and cleanup
CREATE INDEX IF NOT EXISTS idx_oa_transactions_created_at ON oa_transactions (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_integrator_transactions_created_at ON integrator_transactions (created_at DESC);
