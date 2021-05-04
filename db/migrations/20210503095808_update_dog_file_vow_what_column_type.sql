-- +goose Up
ALTER TABLE maker.dog_file_vow
ALTER COLUMN what TYPE TEXT USING (what::TEXT);

DROP INDEX maker.dog_file_vow_what_index;
-- +goose Down
ALTER TABLE maker.dog_file_vow
    ALTER COLUMN what TYPE BIGINT USING (what::BIGINT);

CREATE INDEX dog_file_vow_what_index
    ON maker.dog_file_vow (what);