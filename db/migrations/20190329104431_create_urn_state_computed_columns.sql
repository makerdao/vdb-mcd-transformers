-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend urn_state with ilk_state
CREATE OR REPLACE FUNCTION maker.urn_state_ilk(state maker.urn_state)
  RETURNS maker.ilk_state AS
$$
  SELECT * FROM maker.get_ilk(
    state.block_height,
    (SELECT id FROM maker.ilks WHERE name = state.ilk_name)
  )
$$ LANGUAGE sql STABLE;

-- Extend urn_state with frob_events
CREATE OR REPLACE FUNCTION maker.urn_state_frobs(state maker.urn_state)
  RETURNS SETOF maker.frob_event AS
$$
  SELECT * FROM maker.urn_frobs(state.ilk_name, state.urn_id)
  WHERE block_height <= state.block_height
$$ LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION maker.urn_state_frobs(maker.urn_state);
DROP FUNCTION maker.urn_state_ilk(maker.ilk_state);
