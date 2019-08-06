-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend frob_event with ilk_state
CREATE FUNCTION api.frob_event_ilk(event api.frob_event)
    RETURNS api.ilk_state AS
$$
SELECT *
FROM api.get_ilk(event.ilk_identifier, event.block_height)
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
