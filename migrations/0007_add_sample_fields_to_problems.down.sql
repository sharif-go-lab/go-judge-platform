ALTER TABLE problems
DROP COLUMN IF EXISTS sample_input,
  DROP COLUMN IF EXISTS sample_output;