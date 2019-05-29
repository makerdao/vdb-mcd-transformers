-- +goose Up
CREATE TABLE maker.dent (
  id        SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  bid_id    NUMERIC NOT NULL,
  lot       NUMERIC,
  bid       NUMERIC,
  guy       BYTEA,
  tic       NUMERIC,
  log_idx   INTEGER NOT NULL,
  tx_idx    INTEGER NOT NULL,
  raw_log   JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN dent_checked BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
DROP TABLE maker.dent;

ALTER TABLE public.checked_headers
  DROP COLUMN dent_checked;