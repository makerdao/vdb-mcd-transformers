-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.all_flops()
    RETURNS SETOF api.flop_state
AS
$BODY$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT bid_id
            FROM maker.flop
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
DROP FUNCTION api.all_flops();

