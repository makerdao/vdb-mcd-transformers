-- +goose Up
ALTER TABLE maker.dog_file_hole
    ALTER COLUMN data TYPE NUMERIC USING (data::NUMERIC);
-- +goose Down
ALTER TABLE maker.dog_file_hole
    DROP COLUMN data;
ALTER TABLE maker.dog_file_hole
    ADD COLUMN data BIGINT;
