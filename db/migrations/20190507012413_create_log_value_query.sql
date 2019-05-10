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
    FROM public.header_sync_transactions txs
    LEFT JOIN headers ON txs.header_id = headers.id
    LEFT JOIN maker.pip_log_value plv ON txs.header_id = plv.header_id
    WHERE headers.block_number = priceUpdate.block_number
    ORDER BY headers.block_number DESC
$body$
  LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP FUNCTION IF EXISTS maker.log_values(INT, INT);
DROP FUNCTION maker.log_value_tx(maker.log_value);
DROP TYPE maker.log_value CASCADE;