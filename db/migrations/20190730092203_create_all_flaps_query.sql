-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION api.all_flaps()
    RETURNS SETOF api.flap AS
$BODY$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT flap_bid_guy.bid_id
            FROM maker.flap_bid_guy
            UNION
            SELECT DISTINCT flap_bid_tic.bid_id
            FROM maker.flap_bid_tic
            UNION
            SELECT DISTINCT flap_bid_bid.bid_id
            FROM maker.flap_bid_bid
            UNION
            SELECT DISTINCT flap_bid_lot.bid_id
            FROM maker.flap_bid_lot
            UNION
            SELECT DISTINCT flap_bid_end.bid_id
            FROM maker.flap_bid_end
            UNION
            SELECT DISTINCT flap_bid_gal.bid_id
            FROM maker.flap_bid_gal
        )
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flap(bid_ids.bid_id) f
    );
END
$BODY$
    LANGUAGE plpgsql;
-- +goose StatementEnd
-- +goose Down
DROP FUNCTION api.all_flaps();
