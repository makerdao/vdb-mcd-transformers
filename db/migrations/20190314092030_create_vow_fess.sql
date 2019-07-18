-- +goose Up
CREATE TABLE maker.vow_fess
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    tab       NUMERIC NOT NULL,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX vow_fess_header_index
    ON maker.vow_fess (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN vow_fess INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.vow_fess_header_index;
DROP TABLE maker.vow_fess;
ALTER TABLE public.checked_headers
    DROP COLUMN vow_fess;

