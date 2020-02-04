-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend file_event with ilk_snapshot
CREATE FUNCTION api.ilk_file_event_ilk(event api.ilk_file_event)
    RETURNS SETOF api.ilk_snapshot AS
$$
SELECT *
FROM api.ilk_snapshot i
WHERE i.ilk_identifier = event.ilk_identifier
  AND i.block_number <= event.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$
    LANGUAGE sql
    STABLE;

-- Extend file_event with txs
CREATE FUNCTION api.ilk_file_event_tx(event api.ilk_file_event)
    RETURNS api.tx AS
$$
SELECT *
FROM get_tx_data(event.block_height, event.log_id)
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.ilk_file_event_ilk(api.ilk_file_event);
DROP FUNCTION api.ilk_file_event_tx(api.ilk_file_event);
