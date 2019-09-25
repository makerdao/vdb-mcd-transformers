-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.all_flops(max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flop_state
AS
$BODY$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT bid_id
            FROM maker.flop
            ORDER BY bid_id DESC
            LIMIT all_flops.max_results OFFSET all_flops.result_offset
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
-- +goose Down
DROP FUNCTION api.all_flops(INTEGER, INTEGER);
