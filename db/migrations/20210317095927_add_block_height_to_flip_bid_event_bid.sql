-- +goose Up
CREATE OR REPLACE FUNCTION api.flip_bid_event_bid(event api.flip_bid_event)
    RETURNS api.flip_bid_snapshot AS
$$
SELECT * FROM api.get_flip_with_address(event.bid_id, event.contract_address, event.block_height)
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
CREATE OR REPLACE FUNCTION api.flip_bid_event_bid(event api.flip_bid_event)
    RETURNS api.flip_bid_snapshot AS
$$
SELECT * FROM api.get_flip_with_address(event.bid_id, event.contract_address)
$$
    LANGUAGE sql
    STABLE;
