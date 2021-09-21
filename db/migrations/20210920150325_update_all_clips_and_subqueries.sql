-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.get_clip_with_address(sale_id numeric, clip_address text, block_height bigint DEFAULT api.max_block()) RETURNS api.clip_sale_snapshot
    LANGUAGE sql STABLE STRICT
AS $$
WITH address_id AS (SELECT id FROM public.addresses WHERE address = get_clip_with_address.clip_address),
     ilk_id AS (SELECT DISTINCT ilk_id FROM maker.dog_bark WHERE clip = (SELECT id FROM address_id)),
     urn_id AS (SELECT urn_id FROM maker.dog_bark WHERE clip = (SELECT id FROM address_id) AND maker.dog_bark.sale_id = get_clip_with_address.sale_id),
     storage_values AS (
         SELECT pos,
                tab,
                lot,
                usr,
                tic,
                top,
                created,
                updated,
                block_number
         FROM maker.clip
         WHERE clip.sale_id = get_clip_with_address.sale_id
           AND clip.address_id = (SELECT id FROM address_id)
           AND block_number <= get_clip_with_address.block_height
         ORDER BY block_number DESC
         LIMIT 1
     )
SELECT storage_values.block_number,
       get_clip_with_address.sale_id,
       (SELECT ilk_id FROM ilk_id),
       (SELECT urn_id FROM urn_id),
       storage_values.pos,
       storage_values.tab,
       storage_values.lot,
       storage_values.usr,
       storage_values.tic,
       storage_values.top,
       storage_values.created,
       storage_values.updated,
       get_clip_with_address.clip_address
FROM storage_values
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.all_clips(ilk text, max_results integer DEFAULT '-1'::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.clip_sale_snapshot
    LANGUAGE plpgsql STABLE STRICT
AS $$
BEGIN
    RETURN QUERY (
        WITH ilk_ids AS (SELECT id
                         FROM maker.ilks
                         WHERE identifier = all_clips.ilk),
             clip_ids AS (
                 SELECT DISTINCT clip
                 FROM maker.dog_bark
                 WHERE dog_bark.ilk_id = (SELECT id from ilk_ids)
             ),
             sales AS (
                 SELECT DISTINCT sale_id, address
                 FROM maker.clip
                          JOIN addresses on clip.address_id = addresses.id
                 WHERE maker.clip.address_id IN (SELECT clip FROM clip_ids)
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
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.get_clip_with_address(sale_id NUMERIC, clip_address TEXT,
                                          block_height BIGINT DEFAULT api.max_block()) RETURNS api.clip_sale_snapshot
    LANGUAGE sql
    STABLE STRICT
AS
$$
WITH address_id AS (SELECT id FROM public.addresses WHERE address = get_clip_with_address.clip_address),
     ilk_id AS (SELECT ilk_id FROM maker.dog_bark WHERE clip = (SELECT id FROM address_id)),
     urn_id AS (SELECT urn_id FROM maker.dog_bark WHERE clip = (SELECT id FROM address_id)),
     storage_values AS (
         SELECT pos,
                tab,
                lot,
                usr,
                tic,
                top,
                created,
                updated,
                block_number
         FROM maker.clip
         WHERE clip.sale_id = get_clip_with_address.sale_id
           AND clip.address_id = (SELECT id FROM address_id)
           AND block_number <= get_clip_with_address.block_height
         ORDER BY block_number DESC
         LIMIT 1
     )
SELECT storage_values.block_number,
       get_clip_with_address.sale_id,
       (SELECT ilk_id FROM ilk_id),
       (SELECT urn_id FROM urn_id),
       storage_values.pos,
       storage_values.tab,
       storage_values.lot,
       storage_values.usr,
       storage_values.tic,
       storage_values.top,
       storage_values.created,
       storage_values.updated,
       get_clip_with_address.clip_address
FROM storage_values
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.all_clips(ilk text, max_results integer DEFAULT '-1'::integer,
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
             clip_ids AS (
                 SELECT DISTINCT clip
                 FROM maker.dog_bark
                 WHERE dog_bark.ilk_id = (SELECT id from ilk_ids)
             ),
             sales AS (
                 SELECT DISTINCT sale_id, address
                 FROM maker.clip
                          JOIN addresses on clip.address_id = addresses.id
                 WHERE maker.clip.address_id IN (SELECT id FROM clip_ids)
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
