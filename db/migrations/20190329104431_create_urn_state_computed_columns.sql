-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend urn_state with ilk_state
CREATE FUNCTION api.urn_state_ilk(state api.urn_state)
    RETURNS api.historical_ilk_state AS
$$
SELECT *
FROM api.historical_ilk_state i
WHERE i.ilk_identifier = state.ilk_identifier
  AND i.block_number <= state.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$
    LANGUAGE sql
    STABLE;

-- Extend urn_state with frob_events
CREATE FUNCTION api.urn_state_frobs(state api.urn_state, max_results INTEGER DEFAULT NULL,
                                    result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.frob_event AS
$$
SELECT *
FROM api.urn_frobs(state.ilk_identifier, state.urn_identifier)
WHERE block_height <= state.block_height
LIMIT urn_state_frobs.max_results
OFFSET
urn_state_frobs.result_offset
$$
    LANGUAGE sql
    STABLE;


-- Extend urn_state with bite events
CREATE FUNCTION api.urn_state_bites(state api.urn_state, max_results INTEGER DEFAULT NULL,
                                    result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
SELECT *
FROM api.urn_bites(state.ilk_identifier, state.urn_identifier)
WHERE block_height <= state.block_height
LIMIT urn_state_bites.max_results
OFFSET
urn_state_bites.result_offset
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.urn_state_bites(api.urn_state, INTEGER, INTEGER);
DROP FUNCTION api.urn_state_frobs(api.urn_state, INTEGER, INTEGER);
DROP FUNCTION api.urn_state_ilk(api.urn_state);
