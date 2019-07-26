-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE OR REPLACE FUNCTION api.all_flips(ilk TEXT) RETURNS SETOF api.flip_state AS
-- +goose StatementBegin
$BODY$
DECLARE
    bid_id NUMERIC;
BEGIN
    FOR bid_id IN
        WITH address AS (SELECT DISTINCT contract_address
                         FROM maker.flip_ilk
                         WHERE flip_ilk.ilk = all_flips.ilk
                         LIMIT 1)
        SELECT DISTINCT flip_bid_guy.bid_id
        FROM maker.flip_bid_guy
        WHERE contract_address = (SELECT contract_address FROM address)
        ORDER BY flip_bid_guy.bid_id
        LOOP
            RETURN NEXT api.get_flip(bid_id, ilk);
        END LOOP;
    RETURN;
END
$BODY$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.all_flips(ilk TEXT);
