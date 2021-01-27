-- +goose NO TRANSACTION
-- +goose Up
CREATE INDEX CONCURRENTLY median_slot_index
    ON maker.median_slot (slot);

CREATE INDEX CONCURRENTLY flip_file_cat_data_index
    ON maker.flip_file_cat (data);

-- +goose Down
DROP INDEX maker.median_slot_index;
DROP INDEX maker.flip_file_cat_data_index;
