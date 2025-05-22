-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON contracts;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_contracts_status;
DROP INDEX IF EXISTS idx_contracts_created_by;
DROP INDEX IF EXISTS idx_contracts_dates;

-- Drop tables
DROP TABLE IF EXISTS contracts;
DROP TABLE IF EXISTS users;

-- Drop enum
DROP TYPE IF EXISTS user_role; 