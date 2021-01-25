-- +goose Up

ALTER TYPE api.flap_bid_snapshot ADD ATTRIBUTE flap_address TEXT;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.get_flap_with_address(bid_id NUMERIC, flap_address TEXT, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flap_bid_snapshot
AS
$$
WITH address_id AS (SELECT id FROM public.addresses WHERE address = get_flap_with_address.flap_address),
     storage_values AS (
         SELECT bid_id,
                guy,
                tic,
                "end",
                lot,
                bid,
                created,
                updated
         FROM maker.flap
         WHERE bid_id = get_flap_with_address.bid_id
		   AND flap.address_id = (SELECT * from address_id)
           AND block_number <= block_height
         ORDER BY block_number DESC
         LIMIT 1
     ),
     deal AS (
         SELECT deal, bid_id
         FROM maker.deal
                  LEFT JOIN public.headers ON deal.header_id = headers.id
         WHERE deal.bid_id = get_flap_with_address.bid_id
           AND deal.address_id = (SELECT * FROM address_id)
           AND headers.block_number <= block_height
         ORDER BY bid_id, block_number DESC
         LIMIT 1
     )

SELECT get_flap_with_address.bid_id,
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
       get_flap_with_address.flap_address
FROM storage_values
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.all_flaps(max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flap_bid_snapshot AS
$BODY$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT bid_id, address
            FROM maker.flap
				JOIN addresses on addresses.id = flap.address_id
            ORDER BY bid_id DESC
            LIMIT all_flaps.max_results
            OFFSET
            all_flaps.result_offset
        )
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flap_with_address(bid_ids.bid_id, bid_ids.address) f
    );
END
$BODY$
    LANGUAGE plpgsql
    STABLE;
-- +goose StatementEnd

CREATE OR REPLACE FUNCTION api.flap_bid_event_bid(event api.flap_bid_event)
    RETURNS api.flap_bid_snapshot AS
$$
SELECT *
FROM api.get_flap_with_address(event.bid_id, event.contract_address, event.block_height)
$$
    LANGUAGE sql
    STABLE;

DROP FUNCTION api.get_flap;

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.all_flaps(max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flap_bid_snapshot AS
$BODY$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT bid_id
            FROM maker.flap
            ORDER BY bid_id DESC
            LIMIT all_flaps.max_results
            OFFSET
            all_flaps.result_offset
        )
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flap(bid_ids.bid_id) f
    );
END
$BODY$
    LANGUAGE plpgsql
    STABLE;
-- +goose StatementEnd

ALTER TYPE api.flap_bid_snapshot DROP ATTRIBUTE flap_address;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.get_flap(bid_id NUMERIC, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flap_bid_snapshot
AS
$$
WITH address_id AS (
    SELECT address_id
    FROM maker.flap
    WHERE flap.bid_id = get_flap.bid_id
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
           AND deal.address_id = (SELECT * FROM address_id)
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
       CASE (SELECT COUNT(*) FROM deal)
           WHEN 0 THEN FALSE
           ELSE TRUE
           END AS dealt,
       storage_values.created,
       storage_values.updated
FROM storage_values
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;
-- +goose StatementEnd

CREATE OR REPLACE FUNCTION api.flap_bid_event_bid(event api.flap_bid_event)
    RETURNS api.flap_bid_snapshot AS
$$
SELECT *
FROM api.get_flap(event.bid_id, event.block_height)
$$
    LANGUAGE sql
    STABLE;

DROP FUNCTION api.get_flap_with_address;
