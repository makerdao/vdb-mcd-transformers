-- +goose Up
CREATE TABLE maker.deal
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

CREATE INDEX deal_header_index
    ON maker.deal (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN deal INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.deal_header_index;

DROP TABLE maker.deal;

ALTER TABLE public.checked_headers
    DROP COLUMN deal;
