-- +goose Up
CREATE TABLE maker.vow_flog
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    era       INTEGER NOT NULL,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX vow_flog_era_index
    ON maker.vow_flog (era);

CREATE INDEX vow_flog_header_index
    ON maker.vow_flog (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN vow_flog INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.vow_flog_era_index;
DROP INDEX maker.vow_flog_header_index;

DROP TABLE maker.vow_flog;

ALTER TABLE public.checked_headers
    DROP COLUMN vow_flog;
