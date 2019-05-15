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
SELECT block_timestamp::BIGINT AS "epoch", (SELECT TIMESTAMP 'epoch' + block_timestamp * INTERVAL '1 second') AS iso
  FROM headers WHERE block_number = tx.block_height
$$ LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.tx_era(api.tx);
DROP TYPE api.tx CASCADE;
DROP TYPE api.era CASCADE;