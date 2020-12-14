-- +goose Up

-- add a new function that works the same as get_flip but takes into account a specific flip_address
-- this is required because there are now more than one flip contract per ilk
CREATE OR REPLACE FUNCTION api.get_flip_with_address(bid_id NUMERIC, flip_address TEXT, ilk TEXT, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flip_bid_snapshot
AS
$$
WITH ilk_ids AS (SELECT id FROM maker.ilks WHERE ilks.identifier = get_flip_with_address.ilk),
     address_id AS (SELECT id FROM public.addresses WHERE address = get_flip_with_address.flip_address),
     kicks AS (SELECT usr
               FROM maker.flip_kick
               WHERE flip_kick.bid_id = get_flip_with_address.bid_id
                 AND address_id = (SELECT * FROM address_id)
               LIMIT 1),
     urn_id AS (SELECT id
                FROM maker.urns
                WHERE urns.ilk_id = (SELECT id FROM ilk_ids)
                  AND urns.identifier = (SELECT usr FROM kicks)),

     storage_values AS (
         SELECT guy,
                tic,
                "end",
                lot,
                bid,
                gal,
                tab,
                created,
                updated
         FROM maker.flip
         WHERE flip.bid_id = get_flip_with_address.bid_id
           AND flip.address_id = (SELECT id FROM address_id)
           AND block_number <= block_height
         ORDER BY block_number DESC
         LIMIT 1
     ),
     deals AS (SELECT deal.bid_id
               FROM maker.deal
                        LEFT JOIN public.headers ON deal.header_id = headers.id
               WHERE deal.bid_id = get_flip_with_address.bid_id
                 AND deal.address_id = (SELECT * FROM address_id)
                 AND headers.block_number <= block_height)

SELECT get_flip_with_address.block_height,
       get_flip_with_address.bid_id,
       (SELECT id FROM ilk_ids),
       (SELECT id FROM urn_id),
       storage_values.guy,
       storage_values.tic,
       storage_values."end",
       storage_values.lot,
       storage_values.bid,
       storage_values.gal,
       CASE (SELECT COUNT(*) FROM deals)
           WHEN 0 THEN FALSE
           ELSE TRUE END,
       storage_values.tab,
       storage_values.created,
       storage_values.updated,
       flip_address AS flip_address
FROM storage_values
$$
    LANGUAGE SQL
    STABLE
    STRICT;


-- +goose Down
DROP FUNCTION api.get_flip_with_address(bid_id NUMERIC, flip_address TEXT, ilk TEXT, block_height BIGINT);