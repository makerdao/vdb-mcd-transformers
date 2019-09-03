-- +goose Up

-- Extend type flip_bid_event with bid field
CREATE FUNCTION api.flip_bid_event_bid(event api.flip_bid_event)
    RETURNS api.flip_state AS
$$
WITH ilks AS (
    SELECT ilks.identifier
    FROM maker.flip_ilk
       LEFT JOIN maker.ilks ON ilks.id = flip_ilk.ilk_id
    WHERE contract_address = event.contract_address
    LIMIT 1
)
SELECT *
FROM api.get_flip(event.bid_id, (SELECT identifier FROM ilks))
$$
    LANGUAGE sql
    STABLE;

-- Extend type flip_bid_event with txs field
CREATE FUNCTION api.flip_bid_event_tx(event api.flip_bid_event)
    RETURNS SETOF api.tx AS
$$
    SELECT * FROM get_tx_data(event.block_height, event.tx_idx)
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
DROP FUNCTION api.flip_bid_event_tx(api.flip_bid_event);
DROP FUNCTION api.flip_bid_event_bid(api.flip_bid_event);
