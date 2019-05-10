-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TYPE maker.tx AS (
  transaction_hash TEXT,
  transaction_index INTEGER,
  block_height BIGINT,
  block_hash TEXT,
  -- Era object
  tx_from TEXT,
  tx_to TEXT
);

CREATE TYPE maker.era AS (
  "epoch" BIGINT,
  iso TIMESTAMP
);

-- Extend tx type with era
CREATE OR REPLACE FUNCTION maker.tx_era(tx maker.tx)
  RETURNS maker.era AS
$$
SELECT block_timestamp::BIGINT AS "epoch", (SELECT TIMESTAMP 'epoch' + block_timestamp * INTERVAL '1 second') AS iso
  FROM headers WHERE block_number = tx.block_height
$$ LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION maker.tx_era(maker.tx);
DROP TYPE maker.tx CASCADE;
DROP TYPE maker.era CASCADE;