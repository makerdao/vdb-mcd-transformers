-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE OR REPLACE FUNCTION api.all_flips(ilk TEXT, max_results INTEGER DEFAULT NULL) RETURNS SETOF api.flip_state AS
-- +goose StatementBegin
$BODY$
BEGIN
    RETURN QUERY (
        WITH ilk_ids AS (SELECT id
                         FROM maker.ilks
                         WHERE identifier = all_flips.ilk),
             address AS (
                 SELECT DISTINCT address_id
                 FROM maker.flip_ilk
                 WHERE flip_ilk.ilk_id = (SELECT id FROM ilk_ids)
                 LIMIT 1),
             bid_ids AS (
                 SELECT DISTINCT flip_kicks.kicks
                 FROM maker.flip_kicks
                 WHERE address_id = (SELECT * FROM address)
                 ORDER BY flip_kicks.kicks DESC
                 LIMIT all_flips.max_results)
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flip(bid_ids.kicks, all_flips.ilk) f
    );
END
$BODY$
    LANGUAGE plpgsql
    STABLE;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.all_flips(TEXT, INTEGER);
