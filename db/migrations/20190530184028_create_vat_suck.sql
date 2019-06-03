-- +goose Up
CREATE TABLE maker.vat_suck
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    u         TEXT,
    v         TEXT,
    rad       NUMERIC,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
    ADD COLUMN vat_suck_checked INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP TABLE maker.vat_suck;
ALTER TABLE public.checked_headers
    DROP COLUMN vat_suck_checked;
