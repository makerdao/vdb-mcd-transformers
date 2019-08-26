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
CREATE FUNCTION api.flip_state_bid_events(flip api.flip_state)
    RETURNS SETOF api.flip_bid_event AS
$$
WITH addresses AS ( -- get the contract address from flip_ilk table using the ilk_id from flip
    SELECT contract_address
    FROM maker.flip_ilk
             LEFT JOIN maker.ilks ON ilks.id = flip_ilk.ilk_id
    WHERE ilks.id = flip.ilk_id
    ORDER BY block_number DESC
    LIMIT 1
)
SELECT bid_id, lot, bid_amount, act, block_height, tx_idx, events.contract_address
FROM api.all_flip_bid_events() AS events
         LEFT JOIN addresses ON events.contract_address = addresses.contract_address
WHERE bid_id = flip.bid_id
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.flip_state_bid_events(api.flip_state);
DROP FUNCTION api.flip_state_ilk(api.flip_state);
DROP FUNCTION api.flip_state_urn(api.flip_state);