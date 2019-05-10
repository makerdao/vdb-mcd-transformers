-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend ilk_state with frob_events
CREATE OR REPLACE FUNCTION maker.ilk_state_frobs(state maker.ilk_state)
  RETURNS SETOF maker.frob_event AS
$$
  SELECT * FROM maker.all_frobs(state.ilk_name)
  WHERE block_height <= state.block_height
$$ LANGUAGE sql STABLE;


-- Extend ilk_state with file events
CREATE OR REPLACE FUNCTION maker.ilk_state_files(state maker.ilk_state)
  RETURNS SETOF maker.file_event AS
$$
  SELECT * FROM maker.ilk_files(state.ilk_name)
  WHERE block_height <= state.block_height
$$ LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION maker.ilk_state_frobs(maker.ilk_state);
DROP FUNCTION maker.ilk_state_files(maker.ilk_state);
