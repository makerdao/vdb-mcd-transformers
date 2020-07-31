-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE OR REPLACE FUNCTION api.all_urns(block_height bigint DEFAULT api.max_block(), max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.urn_snapshot
    LANGUAGE sql STABLE
    AS $$
WITH distinct_urn_snapshots AS (SELECT urn_identifier, ilk_identifier, MAX(block_height) AS block_height
                                FROM api.urn_snapshot
                                WHERE block_height <= all_urns.block_height
                                GROUP BY urn_identifier, ilk_identifier)
SELECT us.urn_identifier, us.ilk_identifier, us.block_height, us.ink, coalesce(us.art, 0), us.created, us.updated
    FROM api.urn_snapshot AS us, distinct_urn_snapshots AS dus
    WHERE us.urn_identifier = dus.urn_identifier
    AND us.ilk_identifier = dus.ilk_identifier
    AND us.block_height = dus.block_height
    LIMIT all_urns.max_results
    OFFSET all_urns.result_offset
$$;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.all_urns;
