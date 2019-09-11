-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.all_flaps(max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flap_state AS
$BODY$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT bid_id
            FROM maker.flap
            ORDER BY bid_id DESC
            LIMIT all_flaps.max_results OFFSET all_flaps.result_offset
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
-- +goose Down
DROP FUNCTION api.all_flaps(INTEGER, INTEGER);
