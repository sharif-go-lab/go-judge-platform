-- migrations/0006_add_password_hash_to_users.up.sql

-- 1) Add the column only if missing
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS password_hash TEXT NOT NULL DEFAULT '';

-- 2) Drop that default so future inserts must supply a real hash
ALTER TABLE users
    ALTER COLUMN password_hash DROP DEFAULT;
