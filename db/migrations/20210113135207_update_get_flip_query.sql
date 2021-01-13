-- +goose Up
-- add a flip_address attribute to the flip_bid_snapshot type
ALTER TYPE api.flip_bid_snapshot ADD ATTRIBUTE flip_address text;

-- update the get_flip function to return a flip_address in it's flip_bid_snapshot
CREATE OR REPLACE FUNCTION api.get_flip(bid_id NUMERIC, ilk TEXT, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flip_bid_snapshot
AS
$$
WITH ilk_ids AS (SELECT id FROM maker.ilks WHERE ilks.identifier = get_flip.ilk),
     -- there should only ever be 1 address for a given ilk, which is why there's a LIMIT with no ORDER BY
     address_id AS (SELECT address_id
                    FROM maker.flip_ilk
                             LEFT JOIN public.headers ON flip_ilk.header_id = headers.id
                    WHERE flip_ilk.ilk_id = (SELECT id FROM ilk_ids)
                      AND block_number <= block_height
                    LIMIT 1),
     kicks AS (SELECT usr
               FROM maker.flip_kick
               WHERE flip_kick.bid_id = get_flip.bid_id
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
         WHERE flip.bid_id = get_flip.bid_id
           AND flip.address_id = (SELECT address_id FROM address_id)
           AND block_number <= block_height
         ORDER BY block_number DESC
         LIMIT 1
     ),
     deals AS (SELECT deal.bid_id
               FROM maker.deal
                        LEFT JOIN public.headers ON deal.header_id = headers.id
               WHERE deal.bid_id = get_flip.bid_id
                 AND deal.address_id = (SELECT * FROM address_id)
                 AND headers.block_number <= block_height)

SELECT get_flip.block_height,
       get_flip.bid_id,
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
       (SELECT address from addresses where id = (SELECT * FROM address_id)) AS flip_address
FROM storage_values
$$
    LANGUAGE SQL
    STABLE
    STRICT;

-- add a new function that works the same as get_flip but takes into account a specific flip_address
-- this is required because there are now more than one flip contract per ilk
CREATE OR REPLACE FUNCTION api.get_flip_with_address(bid_id NUMERIC, flip_address TEXT, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flip_bid_snapshot
AS
$$
WITH address_id AS (SELECT id FROM public.addresses WHERE address = get_flip_with_address.flip_address),
     ilk_id as (SELECT DISTINCT ilk_id FROM maker.flip_ilk WHERE flip_ilk.address_id = (SELECT id FROM address_id)),
     kick AS (SELECT usr
              FROM maker.flip_kick
              WHERE flip_kick.bid_id = get_flip_with_address.bid_id
                AND address_id = (SELECT * FROM address_id)
              LIMIT 1),
     urn_id AS (SELECT id
                FROM maker.urns
                WHERE urns.ilk_id = (SELECT ilk_id FROM ilk_id)
                  AND urns.identifier = (SELECT usr FROM kick)),

     storage_values AS (
         SELECT guy,
                tic,
                "end",
                lot,
                bid,
                gal,
                tab,
                created,
                updated,
                block_number
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
SELECT storage_values.block_number,
       get_flip_with_address.bid_id,
       (SELECT ilk_id FROM ilk_id),
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
       get_flip_with_address.flip_address
FROM storage_values
$$
    LANGUAGE SQL
    STABLE
    STRICT;


-- update all_flips function to use get_flip_with_address
CREATE OR REPLACE FUNCTION api.all_flips(ilk TEXT, max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flip_bid_snapshot AS
-- +goose StatementBegin
$BODY$
BEGIN
    RETURN QUERY (
        WITH ilk_id AS (SELECT id
                        FROM maker.ilks
                        WHERE identifier = all_flips.ilk),
             address_ids AS (
                 SELECT DISTINCT address_id as id
                 FROM maker.flip_ilk
                 WHERE flip_ilk.ilk_id = (SELECT id FROM ilk_id)
             ),
             bids AS (
                 SELECT DISTINCT bid_id, address
                 FROM maker.flip
                          JOIN addresses on addresses.id = maker.flip.address_id
                 WHERE maker.flip.address_id IN (SELECT * FROM address_ids)
                 ORDER BY bid_id DESC
                 LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
                     OFFSET
                     all_flips.result_offset
             )
        SELECT f.*
        FROM bids,
             LATERAL api.get_flip_with_address(bids.bid_id, bids.address) f
    );
END
$BODY$
    LANGUAGE plpgsql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;
-- +goose StatementEnd

--add flip_address attribute to bite_event type
ALTER TYPE api.bite_event ADD ATTRIBUTE flip_address text;

-- update all_bites to return a bite_event with a flip_address attribute
CREATE OR REPLACE FUNCTION api.all_bites(ilk_identifier TEXT, max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier)

SELECT ilk_identifier,
       identifier AS urn_identifier,
       bid_id,
       ink,
       art,
       tab,
       block_number,
       log_id,
       flip
FROM maker.bite
         LEFT JOIN maker.urns ON bite.urn_id = urns.id
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE urns.ilk_id = (SELECT id FROM ilk)
ORDER BY urn_identifier, block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
    OFFSET
    all_bites.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;

-- update urn_bites to return a bite_event with a flip_address attribute
CREATE OR REPLACE FUNCTION api.urn_bites(ilk_identifier TEXT, urn_identifier TEXT, max_results INTEGER DEFAULT -1,
                                         result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier),
     urn AS (SELECT id
             FROM maker.urns
             WHERE ilk_id = (SELECT id FROM ilk)
               AND identifier = urn_bites.urn_identifier)

SELECT ilk_identifier,
       urn_bites.urn_identifier,
       bid_id,
       ink,
       art,
       tab,
       block_number,
       log_id,
       flip
FROM maker.bite
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE bite.urn_id = (SELECT id FROM urn)
ORDER BY block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
    OFFSET
    urn_bites.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;

-- update bite_event_bid function to make sure to return the flip from the correct flip contract
CREATE OR REPLACE FUNCTION api.bite_event_bid(event api.bite_event)
    RETURNS api.flip_bid_snapshot AS
$$
SELECT *
FROM api.get_flip_with_address(event.bid_id, event.flip_address, event.block_height)
$$
    LANGUAGE sql
    STABLE;

CREATE OR REPLACE FUNCTION api.flip_bid_event_bid(event api.flip_bid_event)
    RETURNS api.flip_bid_snapshot AS
$$
SELECT * FROM api.get_flip_with_address(event.bid_id, event.contract_address)
$$
    LANGUAGE sql
    STABLE;

-- drop api.get_flip to use api.get_flip_with_address instead
DROP FUNCTION api.get_flip(NUMERIC, TEXT, BIGINT);

-- +goose Down

-- remove new attribute
ALTER TYPE api.flip_bid_snapshot DROP ATTRIBUTE flip_address;

-- put get_flip back to how it was before this migration
CREATE OR REPLACE FUNCTION api.get_flip(bid_id NUMERIC, ilk TEXT, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flip_bid_snapshot
AS
$$
WITH ilk_ids AS (SELECT id FROM maker.ilks WHERE ilks.identifier = get_flip.ilk),
     -- there should only ever be 1 address for a given ilk, which is why there's a LIMIT with no ORDER BY
     address_id AS (SELECT address_id
                    FROM maker.flip_ilk
                             LEFT JOIN public.headers ON flip_ilk.header_id = headers.id
                    WHERE flip_ilk.ilk_id = (SELECT id FROM ilk_ids)
                      AND block_number <= block_height
                    LIMIT 1),
     kicks AS (SELECT usr
               FROM maker.flip_kick
               WHERE flip_kick.bid_id = get_flip.bid_id
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
         WHERE flip.bid_id = get_flip.bid_id
           AND flip.address_id = (SELECT address_id FROM address_id)
           AND block_number <= block_height
         ORDER BY block_number DESC
         LIMIT 1
     ),
     deals AS (SELECT deal.bid_id
               FROM maker.deal
                        LEFT JOIN public.headers ON deal.header_id = headers.id
               WHERE deal.bid_id = get_flip.bid_id
                 AND deal.address_id = (SELECT * FROM address_id)
                 AND headers.block_number <= block_height)

SELECT get_flip.block_height,
       get_flip.bid_id,
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
       storage_values.updated
FROM storage_values
$$
    LANGUAGE SQL
    STABLE
    STRICT;

DROP FUNCTION api.get_flip_with_address(bid_id NUMERIC, flip_address TEXT, block_height BIGINT);

-- put all_flips query back to where it was before this migration
CREATE OR REPLACE FUNCTION api.all_flips(ilk TEXT, max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flip_bid_snapshot AS
-- +goose StatementBegin
$BODY$
BEGIN
    RETURN QUERY (
        WITH ilk_ids AS (SELECT id
                         FROM maker.ilks
                         WHERE identifier = all_flips.ilk),
             address AS (
                 SELECT DISTINCT address_id
                 FROM maker.flip_ilk
                 WHERE flip_ilk.ilk_id = (SELECT id FROM ilk_ids)
                 LIMIT 1),
             bid_ids AS (
                 SELECT DISTINCT bid_id
                 FROM maker.flip
                 WHERE address_id = (SELECT * FROM address)
                 ORDER BY bid_id DESC
                 LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
                     OFFSET
                     all_flips.result_offset)
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flip(bid_ids.bid_id, all_flips.ilk) f
    );
END
$BODY$
    LANGUAGE plpgsql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;
-- +goose StatementEnd

-- put bite event back to where it was before
ALTER TYPE api.bite_event DROP ATTRIBUTE flip_address;

-- put all_bites query back to where it was before this migration
CREATE OR REPLACE FUNCTION api.all_bites(ilk_identifier TEXT, max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier)

SELECT ilk_identifier,
       identifier AS urn_identifier,
       bid_id,
       ink,
       art,
       tab,
       block_number,
       log_id
FROM maker.bite
         LEFT JOIN maker.urns ON bite.urn_id = urns.id
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE urns.ilk_id = (SELECT id FROM ilk)
ORDER BY urn_identifier, block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
    OFFSET
    all_bites.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;

-- put urn_bites query back to where it was before this migration
CREATE OR REPLACE FUNCTION api.urn_bites(ilk_identifier TEXT, urn_identifier TEXT, max_results INTEGER DEFAULT -1,
                                         result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier),
     urn AS (SELECT id
             FROM maker.urns
             WHERE ilk_id = (SELECT id FROM ilk)
               AND identifier = urn_bites.urn_identifier)

SELECT ilk_identifier,
       urn_bites.urn_identifier,
       bid_id,
       ink,
       art,
       tab,
       block_number,
       log_id
FROM maker.bite
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE bite.urn_id = (SELECT id FROM urn)
ORDER BY block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
    OFFSET
    urn_bites.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;

-- put bite_event_bit query back to where it was before this migration
CREATE OR REPLACE FUNCTION api.bite_event_bid(event api.bite_event)
    RETURNS api.flip_bid_snapshot AS
$$
SELECT *
FROM api.get_flip(event.bid_id, event.ilk_identifier, event.block_height)
$$
    LANGUAGE sql
    STABLE;

-- put flip_bid_event_bid back as it was
CREATE OR REPLACE FUNCTION api.flip_bid_event_bid(event api.flip_bid_event)
    RETURNS api.flip_bid_snapshot AS
$$
WITH ilk AS (
    SELECT ilks.identifier
    FROM maker.flip_ilk
             LEFT JOIN maker.ilks ON ilks.id = flip_ilk.ilk_id
    WHERE flip_ilk.address_id = (SELECT id FROM addresses WHERE address = event.contract_address)
    LIMIT 1
)
SELECT *
FROM api.get_flip(event.bid_id, (SELECT identifier FROM ilk))
$$
    LANGUAGE sql
    STABLE;

