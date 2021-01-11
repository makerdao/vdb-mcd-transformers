-- +goose Up
DROP FUNCTION api.get_urns_by_ilk(TEXT, BIGINT, INTEGER, INTEGER);

CREATE FUNCTION api.get_urns_by_ilk(ilk_identifier TEXT,
                                    block_height bigint DEFAULT api.max_block(),
                                    min_block_height bigint DEFAULT 0,
                                    max_results integer DEFAULT NULL::integer,
                                    result_offset integer DEFAULT 0) RETURNS SETOF api.urn_snapshot
    LANGUAGE sql
    STABLE
AS
$$
SELECT *
FROM (SELECT DISTINCT ON (urn_identifier, urn_snapshot.ilk_identifier) urn_identifier,
                                                                       urn_snapshot.ilk_identifier,
                                                                       urn_snapshot.block_height,
                                                                       ink,
                                                                       coalesce(art, 0),
                                                                       created,
                                                                       updated
      FROM api.urn_snapshot
      WHERE urn_snapshot.block_height <= get_urns_by_ilk.block_height
        AND urn_snapshot.block_height > get_urns_by_ilk.min_block_height
        AND urn_snapshot.ilk_identifier = get_urns_by_ilk.ilk_identifier
      ORDER BY urn_identifier, ilk_identifier, updated DESC) AS latest_urns
ORDER BY updated DESC
LIMIT get_urns_by_ilk.max_results OFFSET get_urns_by_ilk.result_offset
$$;

-- +goose Down
DROP FUNCTION api.get_urns_by_ilk(TEXT, BIGINT, BIGINT, INTEGER, INTEGER);

CREATE FUNCTION api.get_urns_by_ilk(ilk_identifier TEXT, block_height bigint DEFAULT api.max_block(),
                                    max_results integer DEFAULT NULL::integer,
                                    result_offset integer DEFAULT 0) RETURNS SETOF api.urn_snapshot
    LANGUAGE sql
    STABLE
AS
$$
SELECT *
FROM (SELECT DISTINCT ON (urn_identifier, urn_snapshot.ilk_identifier) urn_identifier,
                                                                       urn_snapshot.ilk_identifier,
                                                                       urn_snapshot.block_height,
                                                                       ink,
                                                                       coalesce(art, 0),
                                                                       created,
                                                                       updated
      FROM api.urn_snapshot
      WHERE urn_snapshot.block_height <= get_urns_by_ilk.block_height
        AND urn_snapshot.ilk_identifier = get_urns_by_ilk.ilk_identifier
      ORDER BY urn_identifier, ilk_identifier, updated DESC) AS latest_urns
ORDER BY updated DESC
LIMIT get_urns_by_ilk.max_results OFFSET get_urns_by_ilk.result_offset
$$;
