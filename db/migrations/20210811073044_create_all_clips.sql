-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION api.all_clips(ilk text, max_results integer DEFAULT '-1'::integer,
                              result_offset integer DEFAULT 0) RETURNS SETOF api.clip_sale_snapshot
    LANGUAGE plpgsql
    STABLE STRICT
AS
$$
BEGIN
    RETURN QUERY (
        WITH ilk_ids AS (SELECT id
                         FROM maker.ilks
                         WHERE identifier = all_clips.ilk),
             clip_id AS (
                 SELECT DISTINCT clip
                 FROM maker.dog_bark
                 WHERE dog_bark.ilk_id = (SELECT id from ilk_ids)
                 LIMIT 1),
             sales AS (
                 SELECT DISTINCT sale_id, address
                 FROM maker.clip
                          JOIN addresses on clip.address_id = addresses.id
                 WHERE maker.clip.address_id = (SELECT id FROM clip_id)
                 ORDER BY sale_id DESC
                 LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END OFFSET all_clips.result_offset
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
