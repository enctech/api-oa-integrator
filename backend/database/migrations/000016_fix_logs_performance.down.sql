-- Drop the index
DROP INDEX IF EXISTS idx_logs_created_at;

-- Recreate the cleanup function
CREATE OR REPLACE FUNCTION delete_old_data() RETURNS TRIGGER AS
$$
BEGIN
    DELETE FROM logs WHERE created_at < NOW() - interval '100 days';
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Recreate the trigger
CREATE TRIGGER data_cleanup_trigger
    AFTER INSERT
    ON logs
    FOR EACH ROW
EXECUTE FUNCTION delete_old_data();
