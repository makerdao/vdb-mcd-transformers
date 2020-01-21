-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE FUNCTION api.historical_urn_state_ilk(state api.historical_urn_state)
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

CREATE FUNCTION api.historical_urn_state_frobs(state api.historical_urn_state, max_results INTEGER DEFAULT NULL,
                                               result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.frob_event AS
$$
SELECT *
FROM api.urn_frobs(state.ilk_identifier, state.urn_identifier)
WHERE block_height <= state.block_height
LIMIT historical_urn_state_frobs.max_results
OFFSET
historical_urn_state_frobs.result_offset
$$
    LANGUAGE sql
    STABLE;


CREATE FUNCTION api.historical_urn_state_bites(state api.historical_urn_state, max_results INTEGER DEFAULT NULL,
                                               result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
SELECT *
FROM api.urn_bites(state.ilk_identifier, state.urn_identifier)
WHERE block_height <= state.block_height
LIMIT historical_urn_state_bites.max_results
OFFSET
historical_urn_state_bites.result_offset
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.historical_urn_state_bites(api.historical_urn_state, INTEGER, INTEGER);
DROP FUNCTION api.historical_urn_state_frobs(api.historical_urn_state, INTEGER, INTEGER);
DROP FUNCTION api.historical_urn_state_ilk(api.historical_urn_state);
