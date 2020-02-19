-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE OR REPLACE FUNCTION api.all_flips(ilk TEXT, max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flip_bid_snapshot AS
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
                 SELECT DISTINCT bid_id
                 FROM maker.flip
                 WHERE address_id = (SELECT * FROM address)
                 ORDER BY bid_id DESC
                 LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
                 OFFSET
                 all_flips.result_offset)
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flip(bid_ids.bid_id, all_flips.ilk) f
    );
END
$BODY$
    LANGUAGE plpgsql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;
-- +goose StatementEnd

COMMENT ON FUNCTION api.all_flips(ilk TEXT, max_results INTEGER, result_offset INTEGER)
    IS E'Get the state of all Flip auctions for a given Ilk as of the most recent block. ilk (e.g. "ETH-A") argument is required. maxResults and resultOffset arguments are optional and default to no max/offset.';

-- +goose Down
DROP FUNCTION api.all_flips(TEXT, INTEGER, INTEGER);
