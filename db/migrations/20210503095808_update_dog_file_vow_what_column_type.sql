-- +goose Up
ALTER TABLE maker.dog_file_vow
ALTER COLUMN what TYPE TEXT USING (what::TEXT);

DROP INDEX maker.dog_file_vow_what_index;
-- +goose Down
ALTER TABLE maker.dog_file_vow
  DROP COLUMN what;

ALTER TABLE maker.dog_file_vow
  ADD COLUMN what BIGINT;
