-- Drop existing foreign key constraint
ALTER TABLE auth_tokens DROP CONSTRAINT IF EXISTS fk_user;

-- Drop existing unique constraint if exists
ALTER TABLE auth_tokens DROP CONSTRAINT IF EXISTS auth_tokens_user_id_key;

-- Add unique constraint
ALTER TABLE auth_tokens ADD CONSTRAINT auth_tokens_user_id_key UNIQUE (user_id);

-- Re-add foreign key constraint
ALTER TABLE auth_tokens ADD CONSTRAINT fk_user 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE; 