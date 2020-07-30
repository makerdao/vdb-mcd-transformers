-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend flip_bid_snapshot with ilk_snapshot
CREATE FUNCTION api.flip_bid_snapshot_ilk(flip_bid_snapshot api.flip_bid_snapshot)
    RETURNS api.ilk_snapshot AS
$$
SELECT i.*
FROM api.ilk_snapshot i
         LEFT JOIN maker.ilks ON ilks.identifier = i.ilk_identifier
WHERE ilks.id = flip_bid_snapshot.ilk_id
  AND i.block_number <= flip_bid_snapshot.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$
    LANGUAGE sql
    STABLE;


-- Extend flip_bid_snapshot with bid events
-- there are different kinds of flippers, so we need to filter on contract address
CREATE FUNCTION api.flip_bid_snapshot_bid_events(flip api.flip_bid_snapshot, max_results INTEGER DEFAULT NULL,
                                                 result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flip_bid_event AS
$$
WITH address_ids AS ( -- get the contract address from flip_ilk table using the ilk_id from flip
    SELECT address_id
    FROM maker.flip_ilk
    WHERE ilk_id = flip.ilk_id
    LIMIT 1
),
     addresses AS (
         SELECT address
         FROM public.addresses
         WHERE id = (SELECT address_id FROM address_ids)
     )
SELECT bid_id, lot, bid_amount, act, block_height, events.log_id, events.contract_address
FROM api.all_flip_bid_events() AS events
WHERE bid_id = flip.bid_id
  AND contract_address = (SELECT address FROM addresses)
ORDER BY block_height DESC
LIMIT flip_bid_snapshot_bid_events.max_results
OFFSET
flip_bid_snapshot_bid_events.result_offset
$$
    LANGUAGE sql
    STABLE;

--- Extend managed_cdp with urn_snapshot
CREATE FUNCTION api.flip_bid_snapshot_urn(flip api.flip_bid_snapshot) RETURNS api.urn_snapshot
       LANGUAGE sql STABLE
       AS $$
SELECT *
FROM api.get_urn(
     (SELECT identifier FROM maker.ilks WHERE ilks.id = flip.ilk_id),
     (SELECT identifier FROM maker.urns WHERE urns.id = flip.urn_id),
     flip.block_height)
$$;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.flip_bid_snapshot_bid_events(api.flip_bid_snapshot, INTEGER, INTEGER);
DROP FUNCTION api.flip_bid_snapshot_ilk(api.flip_bid_snapshot);
DROP FUNCTION api.flip_bid_snapshot_urn(api.flip_bid_snapshot);
