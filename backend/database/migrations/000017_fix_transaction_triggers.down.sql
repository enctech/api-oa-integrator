-- Drop the indexes
DROP INDEX IF EXISTS idx_oa_transactions_created_at;
DROP INDEX IF EXISTS idx_integrator_transactions_created_at;

-- Recreate the functions and triggers (not recommended, but needed for rollback)
CREATE OR REPLACE FUNCTION delete_old_oa_transaction_data() RETURNS TRIGGER AS
$$
BEGIN
    DELETE FROM oa_transactions WHERE created_at < NOW() - interval '100 days';
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER oa_old_data_cleanup_trigger
    AFTER INSERT ON oa_transactions
    FOR EACH ROW
EXECUTE FUNCTION delete_old_oa_transaction_data();

CREATE OR REPLACE FUNCTION delete_old_integrator_transaction_data() RETURNS TRIGGER AS
$$
BEGIN
    DELETE FROM integrator_transactions WHERE created_at < NOW() - interval '100 days';
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER integrator_transaction_cleanup_trigger
    AFTER INSERT ON integrator_transactions
    FOR EACH ROW
EXECUTE FUNCTION delete_old_integrator_transaction_data();
