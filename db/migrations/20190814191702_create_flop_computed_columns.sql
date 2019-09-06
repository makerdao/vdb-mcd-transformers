-- +goose Up

--Extend flop_state with bid events
CREATE FUNCTION api.flop_state_bid_events(flop api.flop_state)
    RETURNS SETOF api.flop_bid_event AS
$$
SELECT *
FROM api.all_flop_bid_events()
WHERE bid_id = flop.bid_id
$$
    LANGUAGE sql
    STABLE;

-- +goose Down

DROP FUNCTION api.flop_state_bid_events(api.flop_state);