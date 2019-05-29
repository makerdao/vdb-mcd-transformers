-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend type sin_queue_event with txs field
CREATE FUNCTION api.sin_queue_event_tx(event api.sin_queue_event)
  RETURNS api.tx AS
$$
SELECT
  txs.hash,
  txs.tx_index,
  headers.block_number AS block_height,
  headers.hash,
  tx_from,
  tx_to
FROM public.header_sync_transactions txs
  LEFT JOIN headers ON txs.header_id = headers.id
WHERE block_number <= event.block_height AND txs.tx_index = event.tx_idx
ORDER BY block_height DESC
$$
LANGUAGE sql
STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.sin_queue_event_tx(api .sin_queue_event )
