-- +goose Up
CREATE TABLE maker.flip_kick
(
    id               SERIAL PRIMARY KEY,
    header_id        INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    bid_id           NUMERIC NOT NULL,
    lot              NUMERIC,
    bid              NUMERIC,
    tab              NUMERIC,
    usr              TEXT,
    gal              TEXT,
    contract_address TEXT,
    tx_idx           INTEGER NOT NULL,
    log_idx          INTEGER NOT NULL,
    raw_log          JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX flip_kick_header_index
    ON maker.flip_kick (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN flip_kick_checked INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.flip_kick_header_index;

DROP TABLE maker.flip_kick;

ALTER TABLE public.checked_headers
    DROP COLUMN flip_kick_checked;
