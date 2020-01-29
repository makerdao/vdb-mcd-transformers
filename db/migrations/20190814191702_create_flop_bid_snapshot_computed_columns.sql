-- +goose Up

--Extend flop_bid_snapshot with bid events
CREATE FUNCTION api.flop_bid_snapshot_bid_events(flop api.flop_bid_snapshot, max_results INTEGER DEFAULT NULL,
                                                 result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flop_bid_event AS
$$
SELECT *
FROM api.all_flop_bid_events() bids
WHERE bid_id = flop.bid_id
ORDER BY bids.block_height DESC
LIMIT flop_bid_snapshot_bid_events.max_results
OFFSET
flop_bid_snapshot_bid_events.result_offset
$$
    LANGUAGE sql
    STABLE;

-- +goose Down

DROP FUNCTION api.flop_bid_snapshot_bid_events(api.flop_bid_snapshot, INTEGER, INTEGER);
