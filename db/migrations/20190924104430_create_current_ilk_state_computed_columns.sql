-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend ilk_state with frob_events
CREATE FUNCTION api.current_ilk_state_frobs(state api.current_ilk_state, max_results INTEGER DEFAULT NULL,
                                            result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.frob_event AS
$$
SELECT *
FROM api.all_frobs(state.ilk_identifier)
ORDER BY block_height DESC
LIMIT max_results
OFFSET
result_offset
$$
    LANGUAGE sql
    STABLE;


-- Extend ilk_state with file events
CREATE FUNCTION api.current_ilk_state_ilk_file_events(state api.current_ilk_state, max_results INTEGER DEFAULT NULL,
                                                      result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.ilk_file_event AS
$$
SELECT *
FROM api.all_ilk_file_events(state.ilk_identifier)
LIMIT max_results
OFFSET
result_offset
$$
    LANGUAGE sql
    STABLE;


-- Extend ilk_state with bite events
CREATE FUNCTION api.current_ilk_state_bites(state api.current_ilk_state, max_results INTEGER DEFAULT NULL,
                                            result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
SELECT *
FROM api.all_bites(state.ilk_identifier)
ORDER BY block_height DESC
LIMIT max_results
OFFSET
result_offset
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.current_ilk_state_bites(api.current_ilk_state, INTEGER, INTEGER);
DROP FUNCTION api.current_ilk_state_frobs(api.current_ilk_state, INTEGER, INTEGER);
DROP FUNCTION api.current_ilk_state_ilk_file_events(api.current_ilk_state, INTEGER, INTEGER);