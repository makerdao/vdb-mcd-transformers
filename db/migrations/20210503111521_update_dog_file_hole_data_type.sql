-- +goose Up
ALTER TABLE maker.dog_file_hole
    ALTER COLUMN data TYPE NUMERIC USING (data::NUMERIC);
-- +goose Down
ALTER TABLE maker.dog_file_hole
    ALTER COLUMN what TYPE BIGINT USING (what::BIGINT);