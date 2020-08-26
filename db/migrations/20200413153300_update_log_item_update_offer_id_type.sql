-- +goose Up
ALTER TABLE maker.log_item_update
    ALTER COLUMN offer_id TYPE NUMERIC;

-- +goose Down

ALTER TABLE maker.log_item_update
    ALTER COLUMN offer_id TYPE INT;

