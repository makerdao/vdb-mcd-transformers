-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.yank
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

CREATE INDEX yank_header_index
    ON maker.yank (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN yank INTEGER NOT NULL DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.yank_header_index;

DROP TABLE maker.yank;

ALTER TABLE public.checked_headers
    DROP COLUMN yank;