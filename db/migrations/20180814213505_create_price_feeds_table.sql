-- +goose Up
CREATE TABLE maker.price_feeds (
  id                  SERIAL PRIMARY KEY,
  block_number        BIGINT  NOT NULL,
  header_id           INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  medianizer_address  TEXT,
  usd_value           NUMERIC,
  log_idx             INTEGER NOT NULL,
  tx_idx              INTEGER NOT NULL,
  raw_log             JSONB,
  UNIQUE (header_id, medianizer_address, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN price_feeds_checked BOOLEAN NOT NULL DEFAULT FALSE;


-- +goose Down
DROP TABLE maker.price_feeds;

ALTER TABLE public.checked_headers
  DROP COLUMN price_feeds_checked;
