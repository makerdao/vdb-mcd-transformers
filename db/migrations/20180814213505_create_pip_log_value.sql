-- +goose Up
CREATE TABLE maker.pip_log_value (
  id               SERIAL PRIMARY KEY,
  block_number     BIGINT  NOT NULL,
  header_id        INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  contract_address TEXT,
  val              NUMERIC,
  log_idx          INTEGER NOT NULL,
  tx_idx           INTEGER NOT NULL,
  raw_log          JSONB,
  UNIQUE (header_id, contract_address, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN pip_log_value_checked BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
DROP TABLE maker.pip_log_value;

ALTER TABLE public.checked_headers
  DROP COLUMN pip_log_value_checked;
