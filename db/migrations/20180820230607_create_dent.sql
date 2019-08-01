-- +goose Up
CREATE TABLE maker.dent
(
    id               SERIAL PRIMARY KEY,
    header_id        INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    bid_id           NUMERIC NOT NULL,
    lot              NUMERIC,
    bid              NUMERIC,
    contract_address TEXT,
    log_idx          INTEGER NOT NULL,
    tx_idx           INTEGER NOT NULL,
    raw_log          JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX dent_header_index
    ON maker.dent (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN dent INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.dent_header_index;

DROP TABLE maker.dent;

ALTER TABLE public.checked_headers
    DROP COLUMN dent;