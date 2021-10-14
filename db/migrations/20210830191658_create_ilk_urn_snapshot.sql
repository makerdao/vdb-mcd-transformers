-- +goose Up

CREATE FUNCTION api.clip_sale_snapshot_ilk(clip_sale_snapshot api.clip_sale_snapshot) RETURNS api.ilk_snapshot
    LANGUAGE sql
    STABLE
AS $$
SELECT i.*
FROM api.ilk_snapshot i
         LEFT JOIN maker.ilks ON ilks.identifier = i.ilk_identifier
WHERE ilks.id = clip_sale_snapshot.ilk_id
  AND i.block_number <= clip_sale_snapshot.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$;


CREATE FUNCTION api.clip_sale_snapshot_urn(clip api.clip_sale_snapshot) RETURNS api.urn_snapshot
    LANGUAGE sql
    STABLE
AS
$$
SELECT *
FROM api.get_urn(
        (SELECT identifier FROM maker.ilks WHERE ilks.id = clip.ilk_id),
        (SELECT identifier FROM maker.urns WHERE urns.id = clip.urn_id),
        clip.block_height)
$$;

-- +goose Down
DROP FUNCTION api.clip_sale_snapshot_ilk(api.clip_sale_snapshot);
DROP FUNCTION api.clip_sale_snapshot_urn(api.clip_sale_snapshot);