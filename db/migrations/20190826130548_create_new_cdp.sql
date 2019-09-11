-- +goose Up
CREATE TABLE maker.new_cdp
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    usr       TEXT,
    own       TEXT,
    cdp       NUMERIC,
    tx_idx    INTEGER NOT NULL,
    log_idx   INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

COMMENT ON COLUMN maker.new_cdp.id
    IS E'@omit';

ALTER TABLE public.checked_headers
    ADD COLUMN new_cdp INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP TABLE maker.new_cdp;
ALTER TABLE public.checked_headers
    DROP COLUMN new_cdp;
