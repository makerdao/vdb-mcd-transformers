-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION api.all_clips(clip_address text, max_results integer DEFAULT '-1'::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.clip_sale_snapshot
    LANGUAGE plpgsql STABLE STRICT
AS $$
BEGIN
    RETURN QUERY (
        WITH clip_address AS (
            SELECT DISTINCT addresses.id
            FROM addresses
            WHERE all_clips.clip_address = addresses.address),
             sales AS (
                 SELECT DISTINCT sale_id, address
                 FROM maker.clip
                          JOIN addresses on addresses.id = maker.clip.address_id
                 WHERE maker.clip.address_id = (SELECT id from clip_address)
                 ORDER BY sale_id DESC
                 LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
                     OFFSET
                     all_clips.result_offset
             )
        SELECT f.*
        FROM sales,
             LATERAL api.get_clip_with_address(sales.sale_id, sales.address) f
    );
END
$$;
-- +goose StatementEnd
-- +goose Down
DROP FUNCTION api.all_clips(TEXT, INTEGER, INTEGER);
