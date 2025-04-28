ALTER TABLE submissions
  ADD COLUMN compile_error TEXT,
  ADD COLUMN runtime_error TEXT,
  ADD COLUMN output TEXT;