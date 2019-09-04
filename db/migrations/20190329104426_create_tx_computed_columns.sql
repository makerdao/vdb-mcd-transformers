-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TYPE api.tx AS (
    transaction_hash TEXT,
    transaction_index INTEGER,
    block_height BIGINT,
    block_hash TEXT,
    -- Era object
    tx_from TEXT,
    tx_to TEXT
    );

CREATE TYPE api.era AS (
    "epoch" BIGINT,
    iso TIMESTAMP
    );

-- Extend tx type with era
CREATE FUNCTION api.tx_era(tx api.tx)
    RETURNS api.era AS
$$
SELECT block_timestamp :: BIGINT AS "epoch", api.epoch_to_datetime(block_timestamp) AS iso
FROM headers
WHERE block_number = tx.block_height
$$
    LANGUAGE sql
    STABLE;

CREATE FUNCTION get_tx_data(block_height bigint, tx_idx integer)
    RETURNS SETOF api.tx AS
$$
SELECT txs.hash, txs.tx_index, headers.block_number, headers.hash, tx_from, tx_to
FROM header_sync_transactions txs
         LEFT JOIN headers ON txs.header_id = headers.id
WHERE block_number = block_height
  AND txs.tx_index = tx_idx
ORDER BY block_number DESC
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION get_tx_data(bigint, integer);
DROP FUNCTION api.tx_era(api.tx);
DROP TYPE api.tx CASCADE;
DROP TYPE api.era CASCADE;