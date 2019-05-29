-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend frob_event with ilk_state
CREATE FUNCTION api.frob_event_ilk(event api.frob_event)
  RETURNS api.ilk_state AS
$$
  SELECT * FROM api.get_ilk(event.ilk_identifier, event.block_height)
$$ LANGUAGE sql STABLE;


-- Extend frob_event with urn_state
CREATE FUNCTION api.frob_event_urn(event api.frob_event)
  RETURNS SETOF api.urn_state AS
$$
  SELECT * FROM api.get_urn(event.ilk_identifier, event.urn_guy, event.block_height)
$$ LANGUAGE sql STABLE;


-- Extend frob_event with txs
CREATE FUNCTION api.frob_event_tx(event api.frob_event)
  RETURNS api.tx AS
$$
  SELECT txs.hash, txs.tx_index, headers.block_number, headers.hash, tx_from, tx_to
  FROM public.header_sync_transactions txs
  LEFT JOIN headers ON txs.header_id = headers.id
  WHERE block_number <= event.block_height AND txs.tx_index = event.tx_idx
  ORDER BY block_number DESC
  LIMIT 1 -- Should always be true anyway?
$$ LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.frob_event_ilk(api.frob_event);
DROP FUNCTION api.frob_event_urn(api.frob_event);
DROP FUNCTION api.frob_event_tx(api.frob_event);
