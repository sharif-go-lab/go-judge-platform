ALTER TABLE submissions
  DROP COLUMN IF EXISTS compile_error,
  DROP COLUMN IF EXISTS runtime_error,
  DROP COLUMN IF EXISTS output;