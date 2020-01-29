-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE FUNCTION api.urn_snapshot_ilk(state api.urn_snapshot)
    RETURNS api.ilk_snapshot AS
$$
SELECT *
FROM api.ilk_snapshot i
WHERE i.ilk_identifier = state.ilk_identifier
  AND i.block_number <= state.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$
    LANGUAGE sql
    STABLE;

CREATE FUNCTION api.urn_snapshot_frobs(state api.urn_snapshot, max_results INTEGER DEFAULT NULL,
                                               result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.frob_event AS
$$
SELECT *
FROM api.urn_frobs(state.ilk_identifier, state.urn_identifier)
WHERE block_height <= state.block_height
LIMIT urn_snapshot_frobs.max_results
OFFSET
urn_snapshot_frobs.result_offset
$$
    LANGUAGE sql
    STABLE;


CREATE FUNCTION api.urn_snapshot_bites(state api.urn_snapshot, max_results INTEGER DEFAULT NULL,
                                               result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
SELECT *
FROM api.urn_bites(state.ilk_identifier, state.urn_identifier)
WHERE block_height <= state.block_height
LIMIT urn_snapshot_bites.max_results
OFFSET
urn_snapshot_bites.result_offset
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.urn_snapshot_bites(api.urn_snapshot, INTEGER, INTEGER);
DROP FUNCTION api.urn_snapshot_frobs(api.urn_snapshot, INTEGER, INTEGER);
DROP FUNCTION api.urn_snapshot_ilk(api.urn_snapshot);
