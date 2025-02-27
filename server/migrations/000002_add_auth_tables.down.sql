DROP TABLE IF EXISTS login_attempts;
DROP TABLE IF EXISTS auth_tokens;
DROP TABLE IF EXISTS registration_verifications;
ALTER TABLE users DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS user_status; 