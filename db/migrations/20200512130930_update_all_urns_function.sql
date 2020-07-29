-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE FUNCTION api.all_urns(block_height bigint DEFAULT api.max_block(), max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.urn_snapshot
    LANGUAGE sql STABLE
    AS $$
    SELECT *
    FROM ( SELECT DISTINCT ON (urn_identifier, ilk_identifier) urn_identifier, ilk_identifier,
                              all_urns.block_height, ink, coalesce(art, 0), created, updated
           FROM api.urn_snapshot
           WHERE block_height <= all_urns.block_height
           ORDER BY urn_identifier, ilk_identifier, updated DESC) AS latest_urns
    ORDER BY updated DESC
    LIMIT all_urns.max_results
    OFFSET all_urns.result_offset
$$;

-- +goose Down

DROP FUNCTION api.all_urns;
