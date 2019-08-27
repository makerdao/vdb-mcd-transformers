-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.all_flops()
    RETURNS SETOF api.flop
AS
$BODY$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT flop_bid_guy.bid_id
            FROM maker.flop_bid_guy
            UNION
            SELECT DISTINCT flop_bid_tic.bid_id
            FROM maker.flop_bid_tic
            UNION
            SELECT DISTINCT flop_bid_bid.bid_id
            FROM maker.flop_bid_bid
            UNION
            SELECT DISTINCT flop_bid_lot.bid_id
            FROM maker.flop_bid_lot
            UNION
            SELECT DISTINCT flop_bid_end.bid_id
            FROM maker.flop_bid_end
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

