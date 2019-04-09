-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Missing: ilk files/bites, urn bites
-- + anything not rooted in urn/ilk type

-- Extend type frob_event with ilk field
CREATE OR REPLACE FUNCTION maker.frob_event_ilk(event maker.frob_event)
  RETURNS SETOF maker.ilk_state AS
$$
  SELECT * FROM maker.get_ilk_at_block_number(
    event.block_number,
    (SELECT id FROM maker.ilks WHERE ilk = event.ilkid))
$$ LANGUAGE sql STABLE;


-- Extend type frob_event with urn field
CREATE OR REPLACE FUNCTION maker.frob_event_urn(event maker.frob_event)
  RETURNS SETOF maker.urn_state AS
$$
  SELECT * FROM maker.get_urn_state_at_block(event.ilkid, event.urnid, event.block_number)
$$ LANGUAGE sql STABLE;


CREATE TYPE maker.tx AS (
  transaction_hash TEXT,
  transaction_index INTEGER,
  block_number BIGINT,
  block_hash TEXT,
  -- Era object
  tx_from TEXT,
  tx_to TEXT
);

CREATE TYPE maker.era AS (
  "epoch" BIGINT,
  iso TIMESTAMP
);


-- Extend tx type with era object
CREATE OR REPLACE FUNCTION maker.tx_era(tx maker.tx)
  RETURNS maker.era AS
$$
SELECT block_timestamp::BIGINT AS "epoch", (SELECT TIMESTAMP 'epoch' + block_timestamp * INTERVAL '1 second') AS iso
  FROM headers WHERE block_number = tx.block_number
$$ LANGUAGE sql STABLE;


-- Extend type frob_event with txs field
CREATE OR REPLACE FUNCTION maker.frob_event_tx(event maker.frob_event)
  RETURNS maker.tx AS
$$
  SELECT txs.hash, txs.tx_index, block_number, headers.hash, tx_from, tx_to
  FROM public.light_sync_transactions txs
  LEFT JOIN headers ON txs.header_id = headers.id
  WHERE block_number <= event.block_number
  ORDER BY block_number DESC
  LIMIT 1 -- Should always be true anyway?
$$ LANGUAGE sql STABLE;


-- Extend ilk_state with frob events
CREATE OR REPLACE FUNCTION maker.ilk_state_frobs(state maker.ilk_state)
  RETURNS SETOF maker.frob_event AS
$$
  SELECT * FROM maker.all_frobs(state.ilk)
  WHERE block_number <= state.block_number
$$ LANGUAGE sql STABLE;


-- Extend urn_state with frob_events
CREATE OR REPLACE FUNCTION maker.urn_state_frobs(state maker.urn_state)
  RETURNS SETOF maker.frob_event AS
$$
  SELECT * FROM maker.urn_frobs(state.ilkid, state.urnid)
  WHERE block_number <= state.blockheight
$$ LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION maker.frob_event_ilk(maker.frob_event);
DROP FUNCTION maker.frob_event_urn(maker.frob_event);
DROP FUNCTION maker.tx_era(maker.tx);
DROP FUNCTION maker.frob_event_tx(maker.frob_event);
DROP FUNCTION maker.ilk_state_frobs(maker.ilk_state);
DROP FUNCTION maker.urn_state_frobs(maker.urn_state);
DROP TYPE maker.tx CASCADE;
DROP TYPE maker.era CASCADE;