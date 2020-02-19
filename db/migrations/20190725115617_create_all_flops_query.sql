-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.all_flops(max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flop_bid_snapshot
AS
$BODY$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT bid_id
            FROM maker.flop
            ORDER BY bid_id DESC
            LIMIT all_flops.max_results
            OFFSET
            all_flops.result_offset
        )
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flop(bid_ids.bid_id) f
    );
END
$BODY$
    LANGUAGE plpgsql
    STABLE;
-- +goose StatementEnd

COMMENT ON FUNCTION api.all_flops(max_results INTEGER, result_offset INTEGER)
    IS E'Get the state of all Flop auctions as of the most recent block. maxResults and resultOffset arguments are optional and default to no max/offset.';

-- +goose Down
DROP FUNCTION api.all_flops(INTEGER, INTEGER);
