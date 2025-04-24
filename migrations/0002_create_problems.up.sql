-- migrations/0002_create_problems.up.sql
CREATE TABLE problems (
                          id              BIGSERIAL PRIMARY KEY,
                          owner_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          title           VARCHAR(200) NOT NULL,
                          statement       TEXT       NOT NULL,
                          time_limit_ms   INT        NOT NULL,
                          memory_limit_mb INT        NOT NULL,
                          status          VARCHAR(20) NOT NULL DEFAULT 'draft',
                          publish_date    TIMESTAMP WITH TIME ZONE,
                          created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                          updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_problems_status_publish ON problems(status, publish_date DESC);
