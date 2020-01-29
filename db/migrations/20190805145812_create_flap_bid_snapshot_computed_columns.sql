-- +goose Up

-- Extend flap_bid_snapshot with bid events
CREATE FUNCTION api.flap_bid_snapshot_bid_events(flap api.flap_bid_snapshot, max_results INTEGER DEFAULT NULL,
                                                 result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flap_bid_event AS
$$
SELECT *
FROM api.all_flap_bid_events() bids
WHERE bid_id = flap.bid_id
ORDER BY bids.block_height DESC
LIMIT flap_bid_snapshot_bid_events.max_results
OFFSET
flap_bid_snapshot_bid_events.result_offset
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
DROP FUNCTION api.flap_bid_snapshot_bid_events(api.flap_bid_snapshot, INTEGER, INTEGER);


