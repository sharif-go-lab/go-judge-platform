-- migrations/0004_create_sessions.up.sql
CREATE TABLE sessions (
                          token       VARCHAR(128) PRIMARY KEY,
                          user_id     BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          expires_at  TIMESTAMP WITH TIME ZONE NOT NULL,
                          created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user ON sessions(user_id);
