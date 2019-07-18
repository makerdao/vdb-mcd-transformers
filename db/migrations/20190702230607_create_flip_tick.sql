-- +goose Up
CREATE TABLE maker.flip_tick
(
    id               SERIAL PRIMARY KEY,
    header_id        INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    bid_id           NUMERIC NOT NULL,
    contract_address TEXT,
    log_idx          INTEGER NOT NULL,
    tx_idx           INTEGER NOT NULL,
    raw_log          JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX flip_tick_header_index
    ON maker.flip_tick (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN flip_tick INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE public.checked_headers
    DROP COLUMN flip_tick;

DROP INDEX maker.flip_tick_header_index;
DROP TABLE maker.flip_tick;
