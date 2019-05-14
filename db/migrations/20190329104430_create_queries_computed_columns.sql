-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Missing: ilk files/bites, urn bites
-- + anything not rooted in urn/ilk type

-- Extend type frob_event with ilk field
CREATE FUNCTION api.frob_event_ilk(event api.frob_event)
  RETURNS SETOF api.ilk AS
$$
  SELECT * FROM api.get_ilk(
    event.block_height,
    (SELECT id FROM maker.ilks WHERE name = event.ilk_name))
$$ LANGUAGE sql STABLE;


-- Extend type frob_event with urn field
CREATE FUNCTION api.frob_event_urn(event api.frob_event)
  RETURNS SETOF api.urn AS
$$
  SELECT * FROM api.get_urn(event.ilk_name, event.urn_id, event.block_height)
$$ LANGUAGE sql STABLE;


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


-- Extend tx type with era object
CREATE FUNCTION api.tx_era(tx api.tx)
  RETURNS api.era AS
$$
SELECT block_timestamp::BIGINT AS "epoch", (SELECT TIMESTAMP 'epoch' + block_timestamp * INTERVAL '1 second') AS iso
  FROM headers WHERE block_number = tx.block_height
$$ LANGUAGE sql STABLE;


-- Extend type frob_event with txs field
CREATE FUNCTION api.frob_event_tx(event api.frob_event)
  RETURNS api.tx AS
$$
  SELECT txs.hash, txs.tx_index, headers.block_number AS block_height, headers.hash, tx_from, tx_to
  FROM public.header_sync_transactions txs
  LEFT JOIN headers ON txs.header_id = headers.id
  WHERE block_number <= event.block_height
  ORDER BY block_height DESC
  LIMIT 1 -- Should always be true anyway?
$$ LANGUAGE sql STABLE;


-- Extend ilk with frob events
CREATE FUNCTION api.ilk_frobs(state api.ilk)
  RETURNS SETOF api.frob_event AS
$$
  SELECT * FROM api.all_frobs(state.ilk_name)
  WHERE block_height <= state.block_height
$$ LANGUAGE sql STABLE;


-- Extend urn with frob_events
CREATE FUNCTION api.urn_frobs(state api.urn)
  RETURNS SETOF api.frob_event AS
$$
  SELECT * FROM api.frobs_for_urn(state.ilk_name, state.urn_id)
  WHERE block_height <= state.block_height
$$ LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.frob_event_ilk(api.frob_event);
DROP FUNCTION api.frob_event_urn(api.frob_event);
DROP FUNCTION api.tx_era(api.tx);
DROP FUNCTION api.frob_event_tx(api.frob_event);
DROP FUNCTION api.ilk_frobs(api.ilk);
DROP FUNCTION api.urn_frobs(api.urn);
DROP TYPE api.tx CASCADE;
DROP TYPE api.era CASCADE;