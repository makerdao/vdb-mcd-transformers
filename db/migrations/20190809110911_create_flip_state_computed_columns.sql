-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend flip_state with ilk_state
CREATE FUNCTION api.flip_state_ilk(flip api.flip_state)
    RETURNS api.ilk_state AS
$$
SELECT *
FROM api.get_ilk((SELECT identifier FROM maker.ilks WHERE ilks.id = flip.ilk_id), flip.block_height)
$$
    LANGUAGE sql
    STABLE;


-- Extend flip_state with urn_state
CREATE FUNCTION api.flip_state_urn(flip api.flip_state)
    RETURNS SETOF api.urn_state AS
$$
SELECT *
FROM api.get_urn((SELECT identifier FROM maker.ilks WHERE ilks.id = flip.ilk_id),
                 (SELECT identifier FROM maker.urns WHERE urns.id = flip.urn_id), flip.block_height)
$$
    LANGUAGE sql
    STABLE;

-- Extend flip_state with bid events
-- there are different kinds of flippers, so we need to filter on contract address
CREATE FUNCTION api.flip_state_bid_events(flip api.flip_state, max_results INTEGER DEFAULT NULL)
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
SELECT bid_id, lot, bid_amount, act, block_height, tx_idx, events.contract_address
FROM api.all_flip_bid_events() AS events
WHERE bid_id = flip.bid_id
  AND contract_address = (SELECT address FROM addresses)
ORDER BY block_height DESC
LIMIT flip_state_bid_events.max_results
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.flip_state_bid_events(api.flip_state, INTEGER);
DROP FUNCTION api.flip_state_ilk(api.flip_state);
DROP FUNCTION api.flip_state_urn(api.flip_state);