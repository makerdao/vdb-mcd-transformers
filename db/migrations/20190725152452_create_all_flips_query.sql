-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE OR REPLACE FUNCTION api.all_flips(ilk TEXT) RETURNS SETOF api.flip_state AS
-- +goose StatementBegin
$BODY$
BEGIN
    RETURN QUERY (
        WITH address AS (
            SELECT DISTINCT contract_address
            FROM maker.flip_ilk
            WHERE flip_ilk.ilk = all_flips.ilk
            LIMIT 1),
             bid_ids AS (
                 SELECT DISTINCT flip_kicks.kicks
                 FROM maker.flip_kicks
                 WHERE contract_address = (SELECT * FROM address)
                 ORDER BY flip_kicks.kicks)
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flip(bid_ids.kicks, ilk) f
    );
END
$BODY$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.all_flips(ilk TEXT);
