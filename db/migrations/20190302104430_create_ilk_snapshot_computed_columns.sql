-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend ilk_snapshot with frob_events
CREATE FUNCTION api.ilk_snapshot_frobs(state api.ilk_snapshot, max_results INTEGER DEFAULT NULL,
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


-- Extend ilk_snapshot with file events
CREATE FUNCTION api.ilk_snapshot_ilk_file_events(state api.ilk_snapshot, max_results INTEGER DEFAULT NULL,
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


-- Extend ilk_snapshot with bite events
CREATE FUNCTION api.ilk_snapshot_bites(state api.ilk_snapshot, max_results INTEGER DEFAULT NULL,
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
DROP FUNCTION api.ilk_snapshot_bites(api.ilk_snapshot, INTEGER, INTEGER);
DROP FUNCTION api.ilk_snapshot_ilk_file_events(api.ilk_snapshot, INTEGER, INTEGER);
DROP FUNCTION api.ilk_snapshot_frobs(api.ilk_snapshot, INTEGER, INTEGER);
