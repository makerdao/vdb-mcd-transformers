-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.relevant_flip_block AS (
    block_height BIGINT,
    block_hash TEXT,
    bid_id NUMERIC
    );

CREATE FUNCTION api.get_flip_blocks_before(bid_id NUMERIC, contract_address TEXT, block_height BIGINT)
    RETURNS SETOF api.relevant_flip_block AS
$$
SELECT block_number AS block_height, block_hash, kicks AS bid_id
FROM maker.flip_kicks
WHERE block_number <= get_flip_blocks_before.block_height
  AND kicks = get_flip_blocks_before.bid_id
  AND flip_kicks.contract_address = get_flip_blocks_before.contract_address
UNION
SELECT headers.block_number AS block_height, headers.hash AS block_hash, flip_kick.bid_id
FROM maker.flip_kick
         LEFT JOIN public.headers ON flip_kick.header_id = headers.id
WHERE headers.block_number <= get_flip_blocks_before.block_height
  AND flip_kick.bid_id = get_flip_blocks_before.bid_id
  AND flip_kick.contract_address = get_flip_blocks_before.contract_address
UNION
SELECT headers.block_number AS block_height, headers.hash AS block_hash, deal.bid_id
FROM maker.deal
         LEFT JOIN public.headers ON deal.header_id = headers.id
WHERE headers.block_number <= get_flip_blocks_before.block_height
  AND deal.bid_id = get_flip_blocks_before.bid_id
  AND deal.contract_address = get_flip_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flip_bid_bid.bid_id
FROM maker.flip_bid_bid
WHERE block_number <= get_flip_blocks_before.block_height
  AND flip_bid_bid.bid_id = get_flip_blocks_before.bid_id
  AND flip_bid_bid.contract_address = get_flip_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flip_bid_lot.bid_id
FROM maker.flip_bid_lot
WHERE block_number <= get_flip_blocks_before.block_height
  AND flip_bid_lot.bid_id = get_flip_blocks_before.bid_id
  AND flip_bid_lot.contract_address = get_flip_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flip_bid_guy.bid_id
FROM maker.flip_bid_guy
WHERE block_number <= get_flip_blocks_before.block_height
  AND flip_bid_guy.bid_id = get_flip_blocks_before.bid_id
  AND flip_bid_guy.contract_address = get_flip_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flip_bid_tic.bid_id
FROM maker.flip_bid_tic
WHERE block_number <= get_flip_blocks_before.block_height
  AND flip_bid_tic.bid_id = get_flip_blocks_before.bid_id
  AND flip_bid_tic.contract_address = get_flip_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flip_bid_end.bid_id
FROM maker.flip_bid_end
WHERE block_number <= get_flip_blocks_before.block_height
  AND flip_bid_end.bid_id = get_flip_blocks_before.bid_id
  AND flip_bid_end.contract_address = get_flip_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flip_bid_usr.bid_id
FROM maker.flip_bid_usr
WHERE block_number <= get_flip_blocks_before.block_height
  AND flip_bid_usr.bid_id = get_flip_blocks_before.bid_id
  AND flip_bid_usr.contract_address = get_flip_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flip_bid_gal.bid_id
FROM maker.flip_bid_gal
WHERE block_number <= get_flip_blocks_before.block_height
  AND flip_bid_gal.bid_id = get_flip_blocks_before.bid_id
  AND flip_bid_gal.contract_address = get_flip_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flip_bid_tab.bid_id
FROM maker.flip_bid_tab
WHERE block_number <= get_flip_blocks_before.block_height
  AND flip_bid_tab.bid_id = get_flip_blocks_before.bid_id
  AND flip_bid_tab.contract_address = get_flip_blocks_before.contract_address
ORDER BY block_height DESC
$$
    LANGUAGE sql
    STABLE;

COMMENT ON FUNCTION api.get_flip_blocks_before(NUMERIC, TEXT, BIGINT)
    IS E'@omit';

-- Function returning state for a single flip as of given block and ilk

CREATE TYPE api.flip_state AS (
    block_height BIGINT,
    bid_id NUMERIC,
    ilk_id INTEGER,
    urn_id INTEGER,
    guy TEXT,
    tic BIGINT,
    "end" BIGINT,
    lot NUMERIC,
    bid NUMERIC,
    gal TEXT,
    dealt BOOLEAN,
    tab NUMERIC,
    created TIMESTAMP,
    updated TIMESTAMP
    );


-- Function returning the state for a single flip as of the given block height and ilk
CREATE FUNCTION api.get_flip(bid_id NUMERIC, ilk TEXT, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flip_state
AS
$$
WITH ilk_ids AS (SELECT id FROM maker.ilks WHERE ilks.ilk = get_flip.ilk),
     address AS (SELECT contract_address
                 FROM maker.flip_ilk
                 WHERE flip_ilk.ilk = get_flip.ilk
                   AND block_number <= block_height
                 LIMIT 1),
     kicks AS (SELECT usr
               FROM maker.flip_kick
               WHERE flip_kick.bid_id = get_flip.bid_id
                 AND contract_address IN (SELECT * FROM address)),
     urn_ids AS (SELECT id
                 FROM maker.urns
                 WHERE urns.ilk_id IN (SELECT * FROM ilk_ids)
                   AND urns.identifier IN (SELECT usr FROM kicks)),
     guys AS (SELECT flip_bid_guy.bid_id, guy
              FROM maker.flip_bid_guy
              WHERE flip_bid_guy.bid_id = get_flip.bid_id
                AND contract_address IN (SELECT * FROM address)
                AND block_number <= block_height
              ORDER BY flip_bid_guy.bid_id, block_number DESC
              LIMIT 1),
     tics AS (SELECT flip_bid_tic.bid_id, tic
              FROM maker.flip_bid_tic
              WHERE flip_bid_tic.bid_id = get_flip.bid_id
                AND contract_address IN (SELECT * FROM address)
                AND block_number <= block_height
              ORDER BY block_number DESC
              LIMIT 1),
     ends AS (SELECT flip_bid_end.bid_id, "end"
              FROM maker.flip_bid_end
              WHERE flip_bid_end.bid_id = get_flip.bid_id
                AND contract_address IN (SELECT * FROM address)
                AND block_number <= block_height
              ORDER BY block_number DESC
              LIMIT 1),
     lots AS (SELECT flip_bid_lot.bid_id, lot
              FROM maker.flip_bid_lot
              WHERE flip_bid_lot.bid_id = get_flip.bid_id
                AND contract_address IN (SELECT * FROM address)
                AND block_number <= block_height
              ORDER BY block_number DESC
              LIMIT 1),
     bids AS (SELECT flip_bid_bid.bid_id, bid
              FROM maker.flip_bid_bid
              WHERE flip_bid_bid.bid_id = get_flip.bid_id
                AND contract_address IN (SELECT * FROM address)
                AND block_number <= block_height
              ORDER BY block_number DESC
              LIMIT 1),
     gals AS (SELECT flip_bid_gal.bid_id, gal
              FROM maker.flip_bid_gal
              WHERE flip_bid_gal.bid_id = get_flip.bid_id
                AND contract_address IN (SELECT * FROM address)
                AND block_number <= block_height
              ORDER BY block_number DESC
              LIMIT 1),
     tabs AS (SELECT flip_bid_tab.bid_id, tab
              FROM maker.flip_bid_tab
              WHERE flip_bid_tab.bid_id = get_flip.bid_id
                AND contract_address IN (SELECT * FROM address)
                AND block_number <= block_height
              ORDER BY block_number DESC
              LIMIT 1),
     deals AS (SELECT deal.bid_id
               FROM maker.deal
                        LEFT JOIN public.headers ON deal.header_id = headers.id
               WHERE deal.bid_id = get_flip.bid_id
                 AND deal.contract_address IN (SELECT * FROM address)
                 AND headers.block_number <= block_height),
     relevant_blocks AS (SELECT *
                         FROM api.get_flip_blocks_before(bid_id, (SELECT * FROM address), get_flip.block_height)),
     created AS (SELECT DISTINCT ON (relevant_blocks.bid_id, relevant_blocks.block_height) relevant_blocks.block_height,
                                                                                           relevant_blocks.block_hash,
                                                                                           relevant_blocks.bid_id,
                                                                                           api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM relevant_blocks
                          LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
                 ORDER BY relevant_blocks.block_height ASC
                 LIMIT 1),
     updated AS (SELECT DISTINCT ON (relevant_blocks.bid_id, relevant_blocks.block_height) relevant_blocks.block_height,
                                                                                           relevant_blocks.block_hash,
                                                                                           relevant_blocks.bid_id,
                                                                                           api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM relevant_blocks
                          LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
                 ORDER BY relevant_blocks.block_height DESC
                 LIMIT 1),
     latest_bid_state AS (
         SELECT get_flip.bid_id,
                guys.guy,
                tics.tic,
                ends.end,
                lots.lot,
                bids.bid,
                gals.gal,
                tabs.tab,
                updated.block_height
         FROM guys
                  LEFT JOIN tics ON tics.bid_id = guys.bid_id
                  LEFT JOIN ends ON ends.bid_id = guys.bid_id
                  LEFT JOIN lots ON lots.bid_id = guys.bid_id
                  LEFT JOIN bids ON bids.bid_id = guys.bid_id
                  LEFT JOIN gals ON gals.bid_id = guys.bid_id
                  LEFT JOIN tabs ON tabs.bid_id = guys.bid_id
                  LEFT JOIN updated ON updated.bid_id = guys.bid_id)
SELECT (CASE
            WHEN (latest_bid_state.guy <> '0x0000000000000000000000000000000000000000'
                OR latest_bid_state.tic <> 0
                OR latest_bid_state."end" <> 0
                OR latest_bid_state.lot <> 0
                OR latest_bid_state.bid <> 0
                OR latest_bid_state.gal <> '0x0000000000000000000000000000000000000000'
                OR latest_bid_state.tab <> 0)
                THEN (SELECT (get_flip.block_height,
                              get_flip.bid_id,
                              (SELECT id from ilk_ids),
                              (SELECT id from urn_ids),
                              guy,
                              tic,
                              "end",
                              lot,
                              bid,
                              gal,
                              CASE (SELECT COUNT(*) FROM deals)
                                  WHEN 0 THEN FALSE
                                  ELSE TRUE
                                  END,
                              tab,
                              created.datetime,
                              updated.datetime)::api.flip_state
                      FROM latest_bid_state
                               LEFT JOIN created ON created.bid_id = latest_bid_state.bid_id
                               LEFT JOIN updated ON updated.bid_id = latest_bid_state.bid_id)
            ELSE (SELECT (get_flip.block_height,
                          get_flip.bid_id,
                          prev_state.ilk_id,
                          prev_state.urn_id,
                          prev_state.guy,
                          prev_state.tic,
                          prev_state."end",
                          prev_state.lot,
                          prev_state.bid,
                          prev_state.gal,
                          CASE (SELECT COUNT(*) FROM deals)
                              WHEN 0 THEN FALSE
                              ELSE TRUE
                              END,
                          prev_state.tab,
                          created.datetime,
                          updated.datetime)::api.flip_state
                  FROM api.get_flip(get_flip.bid_id, get_flip.ilk,
                                    (SELECT block_height FROM latest_bid_state) - 1) prev_state
                           LEFT JOIN created ON created.bid_id = get_flip.bid_id
                           LEFT JOIN updated ON updated.bid_id = get_flip.bid_id) END)
FROM latest_bid_state
$$
    LANGUAGE SQL
    STABLE
    STRICT;

-- +goose Down
DROP FUNCTION api.get_flip_blocks_before(NUMERIC, TEXT, BIGINT);
DROP TYPE api.relevant_flip_block CASCADE;
DROP FUNCTION api.get_flip(NUMERIC, TEXT, BIGINT);
DROP TYPE api.flip_state CASCADE;
