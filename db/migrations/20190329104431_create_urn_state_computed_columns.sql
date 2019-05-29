-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend urn_state with ilk_state
CREATE FUNCTION api.urn_state_ilk(state api.urn_state)
  RETURNS api.ilk_state AS
$$
SELECT *
FROM api.get_ilk(state.ilk_identifier, state.block_height)
$$
LANGUAGE sql
STABLE;

-- Extend urn_state with frob_events
CREATE FUNCTION api.urn_state_frobs(state api.urn_state)
  RETURNS SETOF api.frob_event AS
$$
SELECT *
FROM api.urn_frobs(state.ilk_identifier, state.urn_guy)
WHERE block_height <= state.block_height
$$
LANGUAGE sql
STABLE;


-- Extend urn_state with bite events
CREATE FUNCTION api.urn_state_bites(state api.urn_state)
  RETURNS SETOF api.bite_event AS
$$
SELECT *
FROM api.urn_bites(state.ilk_identifier, state.urn_guy)
WHERE block_height <= state.block_height
$$
LANGUAGE sql
STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.urn_state_bites(api .urn_state );
DROP FUNCTION api.urn_state_frobs(api .urn_state );
DROP FUNCTION api.urn_state_ilk(api .urn_state );
