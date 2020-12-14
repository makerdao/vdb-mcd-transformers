-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- update all_flips function to use get_flip_with_address
CREATE OR REPLACE FUNCTION api.all_flips(ilk TEXT, max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flip_bid_snapshot AS
-- +goose StatementBegin
$BODY$
BEGIN
    RETURN QUERY (
        WITH ilk_ids AS (SELECT id
                         FROM maker.ilks
                         WHERE identifier = all_flips.ilk),
             address_ids AS (
                 SELECT DISTINCT address_id as id
                 FROM maker.flip_ilk
                 WHERE flip_ilk.ilk_id = (SELECT id FROM ilk_ids)
             ),
             bids AS (
                 SELECT DISTINCT bid_id, address
                 FROM maker.flip
                 JOIN addresses on addresses.id = maker.flip.address_id
                 WHERE maker.flip.address_id IN (SELECT * FROM address_ids)
                 ORDER BY bid_id DESC
                 LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
                     OFFSET
                     all_flips.result_offset
        )
        SELECT f.*
        FROM bids,
             LATERAL api.get_flip_with_address(bids.bid_id, bids.address, all_flips.ilk) f
    );
END
$BODY$
    LANGUAGE plpgsql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;
-- +goose StatementEnd

-- +goose Down
-- put all_flips query back to where it was before this migration
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
