-- +goose Up

--Extend flop with bid events
CREATE FUNCTION api.flop_bid_events(flop api.flop)
    RETURNS SETOF api.flop_bid_event AS
$$
SELECT *
FROM api.all_flop_bid_events()
WHERE bid_id = flop.bid_id
$$
    LANGUAGE sql
    STABLE;

-- +goose Down

DROP FUNCTION api.flop_bid_events(api.flop);