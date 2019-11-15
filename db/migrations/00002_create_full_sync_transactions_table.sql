-- +goose Up
CREATE TABLE public.full_sync_transactions (
  id          SERIAL PRIMARY KEY,
  block_id    INTEGER NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
  gas_limit    NUMERIC,
  gas_price    NUMERIC,
  hash        VARCHAR(66),
  input_data  BYTEA,
  nonce       NUMERIC,
  raw         BYTEA,
  tx_from     VARCHAR(66),
  tx_index    INTEGER,
  tx_to       VARCHAR(66),
  "value"     NUMERIC
);

COMMENT ON TABLE public.full_sync_transactions
    IS E'@omit';

-- +goose Down
DROP TABLE full_sync_transactions;