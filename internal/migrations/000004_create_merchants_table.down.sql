DROP TRIGGER IF EXISTS update_merchants_updated_at ON merchants;
DROP FUNCTION IF EXISTS update_merchants_updated_at();
DROP TABLE IF EXISTS merchants;
DROP TYPE IF EXISTS merchant_type;

DROP TABLE IF EXISTS merchant_customers;
DROP INDEX IF EXISTS idx_users_phone;
