-- +goose Up
CREATE TYPE api.flap_state AS (
    bid_id NUMERIC,
    guy TEXT,
    tic BIGINT,
    "end" BIGINT,
    lot NUMERIC,
    bid NUMERIC,
    gal TEXT,
    dealt BOOLEAN,
    created TIMESTAMP,
    updated TIMESTAMP
    );

CREATE FUNCTION api.get_flap(bid_id NUMERIC, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flap_state
AS
$$
WITH address AS (
    SELECT contract_address
    FROM maker.flap
    WHERE flap.bid_id = get_flap.bid_id
      AND block_number <= block_height
    LIMIT 1
),
     storage_values AS (
         SELECT bid_id, guy, tic, "end", lot, bid, gal, created, updated
         FROM maker.flap
         WHERE bid_id = get_flap.bid_id
           AND block_number <= block_height
         ORDER BY block_number DESC
         LIMIT 1
     ),
     deal AS (
         SELECT deal, bid_id
         FROM maker.deal
                  LEFT JOIN public.headers ON deal.header_id = headers.id
         WHERE deal.bid_id = get_flap.bid_id
           AND deal.contract_address IN (SELECT * FROM address)
           AND headers.block_number <= block_height
         ORDER BY bid_id, block_number DESC
         LIMIT 1
     )

SELECT get_flap.bid_id,
       storage_values.guy,
       storage_values.tic,
       storage_values."end",
       storage_values.lot,
       storage_values.bid,
       storage_values.gal,
       CASE (SELECT COUNT(*) FROM deal)
           WHEN 0 THEN FALSE
           ELSE TRUE
           END AS dealt,
       storage_values.created,
       storage_values.updated
FROM storage_values
$$
    LANGUAGE sql
    STABLE;
-- +goose Down
DROP FUNCTION api.get_flap(NUMERIC, BIGINT);
DROP TYPE api.flap_state CASCADE;