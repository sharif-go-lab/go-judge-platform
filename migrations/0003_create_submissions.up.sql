-- migrations/0003_create_submissions.up.sql
CREATE TABLE submissions (
                             id             BIGSERIAL PRIMARY KEY,
                             user_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                             problem_id     BIGINT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
                             code           TEXT    NOT NULL,
                             language       VARCHAR(20) NOT NULL,
                             status         VARCHAR(20) NOT NULL DEFAULT 'pending',
                             compile_error  TEXT,
                             runtime_error  TEXT,
                             output         TEXT,
                             created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                             updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_submissions_user        ON submissions(user_id, created_at DESC);
CREATE INDEX idx_submissions_problem     ON submissions(problem_id, status);
