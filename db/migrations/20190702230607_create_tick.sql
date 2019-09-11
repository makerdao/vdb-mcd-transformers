-- +goose Up
CREATE TABLE maker.tick
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    log_idx    INTEGER NOT NULL,
    tx_idx     INTEGER NOT NULL,
    raw_log    JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX tick_header_index
    ON maker.tick (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN tick INTEGER NOT NULL DEFAULT 0;

CREATE INDEX tick_bid_id_index
    ON maker.tick (bid_id);

-- +goose Down
ALTER TABLE public.checked_headers
    DROP COLUMN tick;

DROP INDEX maker.tick_header_index;
DROP INDEX maker.tick_bid_id_index;
DROP TABLE maker.tick;
