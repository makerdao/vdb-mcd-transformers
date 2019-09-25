-- +goose Up

-- Extend type flop_bid_event with bid field
CREATE FUNCTION api.flop_bid_event_bid(event api.flop_bid_event)
    RETURNS api.flop_state AS
$$
SELECT *
FROM api.get_flop(event.bid_id, event.block_height)
$$
    LANGUAGE sql
    STABLE;

-- Extend type flop_bid_event with txs field
CREATE FUNCTION api.flop_bid_event_tx(event api.flop_bid_event)
    RETURNS SETOF api.tx AS
$$
SELECT *
FROM get_tx_data(event.block_height, event.log_id)
$$
    LANGUAGE SQL
    STABLE;

-- +goose Down
DROP FUNCTION api.flop_bid_event_tx(api.flop_bid_event);
DROP FUNCTION api.flop_bid_event_bid(api.flop_bid_event);
