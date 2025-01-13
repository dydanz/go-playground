-- Drop unique constraint
ALTER TABLE auth_tokens DROP CONSTRAINT IF EXISTS auth_tokens_user_id_key; 