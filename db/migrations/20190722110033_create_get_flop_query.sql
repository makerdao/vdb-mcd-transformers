-- +goose Up
CREATE TYPE api.flop_state AS (
    bid_id NUMERIC,
    guy TEXT,
    tic BIGINT,
    "end" BIGINT,
    lot NUMERIC,
    bid NUMERIC,
    dealt BOOLEAN,
    created TIMESTAMP,
    updated TIMESTAMP
    );

CREATE OR REPLACE FUNCTION api.get_flop(bid_id NUMERIC, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flop_state
AS
$$
WITH address_id AS (
    SELECT address_id
    FROM maker.flop
    WHERE flop.bid_id = get_flop.bid_id
      AND block_number <= block_height
    LIMIT 1
),
     storage_values AS (
         SELECT bid_id,
                guy,
                tic,
                "end",
                lot,
                bid,
                created,
                updated
         FROM maker.flop
         WHERE bid_id = get_flop.bid_id
           AND block_number <= block_height
         ORDER BY block_number DESC
         LIMIT 1
     ),
     deal AS (
         SELECT deal.bid_id
         FROM maker.deal
                  LEFT JOIN public.headers ON deal.header_id = headers.id
         WHERE deal.bid_id = get_flop.bid_id
           AND deal.address_id = (SELECT address_id FROM address_id)
           AND headers.block_number <= block_height
         ORDER BY bid_id, block_number DESC
         LIMIT 1
     )

SELECT get_flop.bid_id,
       storage_values.guy,
       storage_values.tic,
       storage_values."end",
       storage_values.lot,
       storage_values.bid,
       CASE (SELECT COUNT(*) FROM deal)
           WHEN 0 THEN FALSE
           ELSE TRUE
           END AS dealt,
       storage_values.created,
       storage_values.updated
FROM storage_values
$$
    LANGUAGE SQL
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;
-- +goose Down
DROP FUNCTION api.get_flop(NUMERIC, BIGINT);
DROP TYPE api.flop_state CASCADE;