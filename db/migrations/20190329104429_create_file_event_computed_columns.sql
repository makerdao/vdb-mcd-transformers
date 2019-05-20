-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend file_event with ilk_state
CREATE FUNCTION api.file_event_ilk(event api.file_event)
  RETURNS SETOF api.ilk_state AS
$$
  SELECT * FROM api.get_ilk(
    event.block_height,
    event.ilk_name
  )
$$ LANGUAGE sql STABLE;

-- Extend file_event with txs
CREATE FUNCTION api.file_event_tx(event api.file_event)
  RETURNS api.tx AS
$$
  SELECT txs.hash, txs.tx_index, headers.block_number AS block_height, headers.hash, tx_from, tx_to
  FROM public.header_sync_transactions txs
  LEFT JOIN headers ON txs.header_id = headers.id
  WHERE block_number <= event.block_height AND txs.tx_index = event.tx_idx
  ORDER BY block_height DESC
  LIMIT 1
$$ LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.file_event_ilk(event api.file_event);
DROP FUNCTION api.file_event_tx(api.file_event);
