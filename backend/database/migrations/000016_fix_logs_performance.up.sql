-- Drop the problematic trigger that runs DELETE on every insert
DROP TRIGGER IF EXISTS data_cleanup_trigger ON logs;

-- Drop the function as it's no longer needed
DROP FUNCTION IF EXISTS delete_old_data();

-- Add index on created_at for efficient time-range queries and cleanup
CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs (created_at DESC);
