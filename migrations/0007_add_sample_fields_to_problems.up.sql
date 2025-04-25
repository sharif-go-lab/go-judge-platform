-- 1) add the two new columns with a benign default so existing rows are valid
ALTER TABLE problems
    ADD COLUMN IF NOT EXISTS sample_input  TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS sample_output TEXT NOT NULL DEFAULT '';

-- 2) drop the defaults so future INSERTs must supply real values
ALTER TABLE problems
    ALTER COLUMN sample_input  DROP DEFAULT,
ALTER COLUMN sample_output DROP DEFAULT;