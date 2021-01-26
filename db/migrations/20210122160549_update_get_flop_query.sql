-- +goose Up
-- add a flop_address attribute to the flop_bid_snapshot type
ALTER TYPE api.flop_bid_snapshot ADD ATTRIBUTE flop_address text;

-- add a new function that works the same as get_flop but takes into account a specific flop_address
-- this is required because there are now more than one flop contract
CREATE OR REPLACE FUNCTION api.get_flop_with_address(bid_id NUMERIC, flop_address TEXT,
                                                     block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flop_bid_snapshot
AS
$$
WITH address_id AS (SELECT id FROM public.addresses WHERE address = get_flop_with_address.flop_address),
     storage_values AS (
         SELECT flop.bid_id,
                guy,
                tic,
                "end",
                lot,
                bid,
                created,
                updated
         FROM maker.flop
         WHERE flop.bid_id = get_flop_with_address.bid_id
           AND flop.address_id = (SELECT * FROM address_id)
           AND block_number <= block_height
         ORDER BY block_number DESC
         LIMIT 1
     ),
     deal AS (
         SELECT deal.bid_id
         FROM maker.deal
                  LEFT JOIN public.headers ON deal.header_id = headers.id
         WHERE deal.bid_id = get_flop_with_address.bid_id
           AND deal.address_id = (SELECT * FROM address_id)
           AND headers.block_number <= block_height
         ORDER BY bid_id, block_number DESC
         LIMIT 1
     )

SELECT get_flop_with_address.bid_id,
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
       storage_values.updated,
       get_flop_with_address.flop_address
FROM storage_values
$$
    LANGUAGE SQL
    STABLE
    STRICT;

-- update all_flops function to use get_flop_with_address
CREATE OR REPLACE FUNCTION api.all_flops(max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flop_bid_snapshot AS
-- +goose StatementBegin
$BODY$
BEGIN
    RETURN QUERY (
        WITH bids AS (
            SELECT DISTINCT bid_id, address
            FROM maker.flop
                     JOIN addresses on addresses.id = flop.address_id
            ORDER BY bid_id DESC
            LIMIT all_flops.max_results OFFSET all_flops.result_offset
        )
        SELECT f.*
        FROM bids,
             LATERAL api.get_flop_with_address(bids.bid_id, bids.address) f
    );
END
$BODY$
    LANGUAGE plpgsql
    STABLE;
-- +goose StatementEnd

-- update computed column to use get_flop_with_address
CREATE OR REPLACE FUNCTION api.flop_bid_event_bid(event api.flop_bid_event)
    RETURNS api.flop_bid_snapshot AS
$$
SELECT *
FROM api.get_flop_with_address(event.bid_id, event.contract_address, event.block_height)
$$
    LANGUAGE sql
    STABLE;

DROP FUNCTION api.get_flop(NUMERIC, BIGINT);


-- +goose Down
-- remove new attribute
ALTER TYPE api.flop_bid_snapshot DROP ATTRIBUTE flop_address;

CREATE OR REPLACE FUNCTION api.get_flop(bid_id NUMERIC, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flop_bid_snapshot
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
         SELECT flop.bid_id,
                guy,
                tic,
                "end",
                lot,
                bid,
                created,
                updated
         FROM maker.flop
         WHERE flop.bid_id = get_flop.bid_id
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

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.all_flops(max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flop_bid_snapshot
AS
$BODY$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT bid_id
            FROM maker.flop
            ORDER BY bid_id DESC
            LIMIT all_flops.max_results OFFSET all_flops.result_offset
        )
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flop(bid_ids.bid_id) f
    );
END
$BODY$
    LANGUAGE plpgsql
    STABLE;
-- +goose StatementEnd

CREATE OR REPLACE FUNCTION api.flop_bid_event_bid(event api.flop_bid_event)
    RETURNS api.flop_bid_snapshot AS
$$
SELECT *
FROM api.get_flop(event.bid_id, event.block_height)
$$
    LANGUAGE sql
    STABLE;

DROP FUNCTION api.get_flop_with_address(NUMERIC, TEXT, BIGINT);
