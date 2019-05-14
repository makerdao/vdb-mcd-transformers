-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE maker.log_value AS (
  val           NUMERIC,
  block_number  BIGINT,
  tx_idx        INTEGER
  -- tx
  );


CREATE OR REPLACE FUNCTION maker.log_values(beginTime INT, endTime INT)
  RETURNS SETOF maker.log_value AS
$body$
  SELECT val, pip_log_value.block_number, tx_idx FROM maker.pip_log_value
    LEFT JOIN public.headers ON pip_log_value.header_id = headers.id
    WHERE block_timestamp BETWEEN $1 AND $2
$body$
  LANGUAGE sql STABLE;


CREATE OR REPLACE FUNCTION maker.log_value_tx(priceUpdate maker.log_value)
  RETURNS maker.tx AS
$body$
  SELECT txs.hash, txs.tx_index, headers.block_number, headers.hash, txs.tx_from, txs.tx_to
    FROM maker.pip_log_value plv
    LEFT JOIN public.header_sync_transactions txs ON plv.header_id = txs.header_id
    LEFT JOIN headers ON plv.header_id = headers.id
    WHERE headers.block_number = priceUpdate.block_number AND priceUpdate.tx_idx = txs.tx_index
    ORDER BY headers.block_number DESC
$body$
  LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP FUNCTION IF EXISTS maker.log_values(INT, INT);
DROP FUNCTION maker.log_value_tx(maker.log_value);
DROP TYPE maker.log_value CASCADE;