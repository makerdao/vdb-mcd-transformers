-- +goose Up
CREATE TYPE api.relevant_flop_block AS (
    block_height BIGINT,
    block_hash TEXT,
    bid_id NUMERIC
    );

CREATE FUNCTION api.get_flop_blocks_before(bid_id NUMERIC, contract_address TEXT, block_height BIGINT)
    RETURNS SETOF api.relevant_flop_block AS
$$
SELECT block_number AS block_height, block_hash, kicks AS bid_id
FROM maker.flop_kicks
WHERE block_number <= get_flop_blocks_before.block_height
  AND kicks = get_flop_blocks_before.bid_id
  AND flop_kicks.contract_address = get_flop_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flop_bid_bid.bid_id
FROM maker.flop_bid_bid
WHERE block_number <= get_flop_blocks_before.block_height
  AND flop_bid_bid.bid_id = get_flop_blocks_before.bid_id
  AND flop_bid_bid.contract_address = get_flop_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flop_bid_lot.bid_id
FROM maker.flop_bid_lot
WHERE block_number <= get_flop_blocks_before.block_height
  AND flop_bid_lot.bid_id = get_flop_blocks_before.bid_id
  AND flop_bid_lot.contract_address = get_flop_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flop_bid_guy.bid_id
FROM maker.flop_bid_guy
WHERE block_number <= get_flop_blocks_before.block_height
  AND flop_bid_guy.bid_id = get_flop_blocks_before.bid_id
  AND flop_bid_guy.contract_address = get_flop_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flop_bid_tic.bid_id
FROM maker.flop_bid_tic
WHERE block_number <= get_flop_blocks_before.block_height
  AND flop_bid_tic.bid_id = get_flop_blocks_before.bid_id
  AND flop_bid_tic.contract_address = get_flop_blocks_before.contract_address
UNION
SELECT block_number AS block_height, block_hash, flop_bid_end.bid_id
FROM maker.flop_bid_end
WHERE block_number <= get_flop_blocks_before.block_height
  AND flop_bid_end.bid_id = get_flop_blocks_before.bid_id
  AND flop_bid_end.contract_address = get_flop_blocks_before.contract_address
ORDER BY block_height DESC
$$
    LANGUAGE sql
    STABLE;


CREATE TYPE api.flop AS (
    bid_id NUMERIC,
    guy TEXT,
    tic BIGINT,
    "end" BIGINT,
    lot NUMERIC,
    bid NUMERIC,
    dealt BOOLEAN,
    created TIMESTAMP,
    updated TIMESTAMP
    );

CREATE OR REPLACE FUNCTION api.get_flop(bid_id NUMERIC, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.flop
AS
$$
WITH address AS (
    SELECT contract_address
    FROM maker.flop_bid_guy
    WHERE flop_bid_guy.bid_id = get_flop.bid_id
      AND block_number <= block_height
    LIMIT 1
),
     guy AS (
         SELECT guy, bid_id
         FROM maker.flop_bid_guy
         WHERE flop_bid_guy.bid_id = get_flop.bid_id
           AND block_number <= block_height
         ORDER BY flop_bid_guy.bid_id, block_number DESC
         LIMIT 1
     ),
     tic AS (
         SELECT tic, bid_id
         FROM maker.flop_bid_tic
         WHERE flop_bid_tic.bid_id = get_flop.bid_id
           AND block_number <= block_height
         ORDER BY flop_bid_tic.bid_id, block_number DESC
         LIMIT 1
     ),
     "end" AS (
         SELECT "end", bid_id
         FROM maker.flop_bid_end
         WHERE flop_bid_end.bid_id = get_flop.bid_id
           AND block_number <= block_height
         ORDER BY flop_bid_end.bid_id, block_number DESC
         LIMIT 1
     ),
     lot AS (
         SELECT lot, bid_id
         FROM maker.flop_bid_lot
         WHERE flop_bid_lot.bid_id = get_flop.bid_id
           AND block_number <= block_height
         ORDER BY flop_bid_lot.bid_id, block_number DESC
         LIMIT 1
     ),
     bid AS (
         SELECT bid, bid_id
         FROM maker.flop_bid_bid
         WHERE flop_bid_bid.bid_id = get_flop.bid_id
           AND block_number <= block_height
         ORDER BY flop_bid_bid.bid_id, block_number DESC
         LIMIT 1
     ),
     deal AS (
         SELECT deal.bid_id
         FROM maker.deal
                  LEFT JOIN public.headers ON deal.header_id = headers.id
         WHERE deal.bid_id = get_flop.bid_id
           AND deal.contract_address IN (SELECT * FROM address)
           AND headers.block_number <= block_height
         ORDER BY bid_id, block_number DESC
         LIMIT 1
     ),
     relevant_blocks AS (
         SELECT *
         FROM api.get_flop_blocks_before(bid_id, (SELECT * FROM address), get_flop.block_height)
     ),
     created AS (
         SELECT DISTINCT ON (relevant_blocks.bid_id, relevant_blocks.block_height) relevant_blocks.block_height,
                                                                                   relevant_blocks.block_hash,
                                                                                   relevant_blocks.bid_id,
                                                                                   api.epoch_to_datetime(headers.block_timestamp) AS datetime
         FROM relevant_blocks
                  LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
         ORDER BY relevant_blocks.block_height ASC
         LIMIT 1
     ),
     updated AS (
         SELECT DISTINCT ON (relevant_blocks.bid_id, relevant_blocks.block_height) relevant_blocks.block_height,
                                                                                   relevant_blocks.block_hash,
                                                                                   relevant_blocks.bid_id,
                                                                                   api.epoch_to_datetime(headers.block_timestamp) AS datetime
         FROM relevant_blocks
                  LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
         ORDER BY relevant_blocks.block_height DESC
         LIMIT 1
     )

SELECT get_flop.bid_id,
       guy.guy,
       tic.tic,
       "end"."end",
       lot.lot,
       bid.bid,
       CASE (SELECT COUNT(*) FROM deal)
           WHEN 0 THEN FALSE
           ELSE TRUE
           END AS dealt,
       created.datetime,
       updated.datetime
FROM lot
         LEFT JOIN guy ON guy.bid_id = lot.bid_id
         LEFT JOIN tic ON tic.bid_id = lot.bid_id
         LEFT JOIN "end" ON "end".bid_id = lot.bid_id
         LEFT JOIN bid ON bid.bid_id = lot.bid_id
         LEFT JOIN created on created.bid_id = lot.bid_id
         LEFT JOIN updated on updated.bid_id = lot.bid_id
$$
    LANGUAGE SQL
    STABLE;
-- +goose Down
DROP FUNCTION api.get_flop_blocks_before(NUMERIC, TEXT, BIGINT);
DROP TYPE api.relevant_flop_block CASCADE;
DROP FUNCTION api.get_flop(NUMERIC, BIGINT);
DROP TYPE api.flop CASCADE;