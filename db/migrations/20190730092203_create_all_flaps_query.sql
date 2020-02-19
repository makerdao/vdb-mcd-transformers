-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.all_flaps(max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flap_bid_snapshot AS
$BODY$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT bid_id
            FROM maker.flap
            ORDER BY bid_id DESC
            LIMIT all_flaps.max_results
            OFFSET
            all_flaps.result_offset
        )
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flap(bid_ids.bid_id) f
    );
END
$BODY$
    LANGUAGE plpgsql
    STABLE;
-- +goose StatementEnd

COMMENT ON FUNCTION api.all_flaps(max_results integer, result_offset integer)
    IS E'Get the state of all Flap auctions as of the most recent block. maxResults and resultOffset arguments are optional and default to no max/offset.';

-- +goose Down
DROP FUNCTION api.all_flaps(INTEGER, INTEGER);
