-- migrations/0005_create_test_cases.up.sql
CREATE TABLE test_cases (
                            id              BIGSERIAL PRIMARY KEY,
                            problem_id      BIGINT    NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
                            input_data      TEXT      NOT NULL,
                            expected_output TEXT      NOT NULL,
                            created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_testcases_problem ON test_cases(problem_id);
