-- +goose Up

-- Extend type flap_bid_event with bid field
CREATE FUNCTION api.flap_bid_event_bid(event api.flap_bid_event)
    RETURNS api.flap_state AS
$$
SELECT *
FROM api.get_flap(event.bid_id, event.block_height)
$$
    LANGUAGE sql
    STABLE;

-- Extend type flap_bid_event with txs field
CREATE FUNCTION api.flap_bid_event_tx(event api.flap_bid_event)
    RETURNS SETOF api.tx AS
$$
SELECT *
FROM get_tx_data(event.block_height, event.tx_idx)
$$
    LANGUAGE SQL
    STABLE;

-- +goose Down
DROP FUNCTION api.flap_bid_event_tx(api.flap_bid_event);
DROP FUNCTION api.flap_bid_event_bid(api.flap_bid_event);
