-- +goose Up
ALTER TABLE maker.clip_kick RENAME COLUMN bid_id TO sale_id;

-- +goose Down
ALTER TABLE maker.clip_kick RENAME COLUMN sale_id TO bid_id;
