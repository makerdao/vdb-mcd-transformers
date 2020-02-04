-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend frob_event with ilk_snapshot
CREATE FUNCTION api.frob_event_ilk(event api.frob_event)
    RETURNS api.ilk_snapshot AS
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


-- Extend frob_event with urn_state
CREATE FUNCTION api.frob_event_urn(event api.frob_event)
    RETURNS SETOF api.urn_state AS
$$
SELECT *
FROM api.get_urn(event.ilk_identifier, event.urn_identifier, event.block_height)
$$
    LANGUAGE sql
    STABLE;


-- Extend frob_event with txs
CREATE FUNCTION api.frob_event_tx(event api.frob_event)
    RETURNS api.tx AS
$$
SELECT * FROM get_tx_data(event.block_height, event.log_id)
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.frob_event_ilk(api.frob_event);
DROP FUNCTION api.frob_event_urn(api.frob_event);
DROP FUNCTION api.frob_event_tx(api.frob_event);
