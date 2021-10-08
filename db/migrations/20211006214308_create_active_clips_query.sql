-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.active_clips(ilk text, block_height bigint DEFAULT api.max_block(),
                                            max_results integer DEFAULT '-1'::integer,
                                            result_offset integer DEFAULT 0) RETURNS SETOF api.clip_sale_snapshot
    LANGUAGE plpgsql
    STABLE STRICT
AS
$$
BEGIN
    RETURN QUERY (
        WITH ilk_ids AS (SELECT id
                         FROM maker.ilks
                         WHERE identifier = active_clips.ilk),
             clip_ids AS (
                 SELECT DISTINCT clip
                 FROM maker.dog_bark
                 WHERE dog_bark.ilk_id = (SELECT id from ilk_ids)
             ),
             active_sales AS (
                 SELECT DISTINCT sale_id, address
                 FROM maker.clip_active_sales
                          JOIN addresses ON clip_active_sales.address_id = addresses.id
                          JOIN headers ON clip_active_sales.header_id = headers.id
                 WHERE maker.clip_active_sales.address_id IN (SELECT clip FROM clip_ids)
                   AND headers.block_number = active_clips.block_height
                 ORDER BY sale_id DESC
                 LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END OFFSET active_clips.result_offset
             )
        SELECT f.*
        FROM active_sales,
             LATERAL api.get_clip_with_address(active_sales.sale_id, active_sales.address) f
    );
END
$$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.active_clips(text, bigint, integer, integer);
