-- migrations/0001_create_users.up.sql
CREATE TABLE users (
                       id             BIGSERIAL PRIMARY KEY,
                       username       VARCHAR(50) NOT NULL UNIQUE,
                       password_hash  TEXT       NOT NULL,
                       email            VARCHAR(100) NOT NULL UNIQUE,
                       is_admin         BOOLEAN     NOT NULL DEFAULT FALSE,
                       created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                       updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()

);

CREATE INDEX idx_users_admin ON users(is_admin);
