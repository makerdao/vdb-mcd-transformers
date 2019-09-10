-- +goose Up

-- Extend flap_state with bid events
CREATE FUNCTION api.flap_state_bid_events(flap api.flap_state, max_results INTEGER DEFAULT NULL)
    RETURNS SETOF api.flap_bid_event AS
$$
SELECT *
FROM api.all_flap_bid_events() bids
WHERE bid_id = flap.bid_id
ORDER BY bids.block_height DESC
LIMIT flap_state_bid_events.max_results
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
DROP FUNCTION api.flap_state_bid_events(api.flap_state, INTEGER);


