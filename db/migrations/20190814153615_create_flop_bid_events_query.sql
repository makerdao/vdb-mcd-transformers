-- +goose Up
CREATE TYPE api.flop_bid_event AS (
    bid_id NUMERIC,
    lot NUMERIC,
    bid_amount NUMERIC,
    act api.bid_act,
    block_height BIGINT,
    tx_idx INTEGER,
    contract_address TEXT
    );

COMMENT ON COLUMN api.flop_bid_event.block_height
    IS E'@omit';
COMMENT ON COLUMN api.flop_bid_event.tx_idx
    IS E'@omit';
COMMENT ON COLUMN api.flop_bid_event.contract_address
    IS E'@omit';

CREATE FUNCTION api.all_flop_bid_events(max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.flop_bid_event AS
$$
WITH address_id AS (
    SELECT address_id
    FROM maker.flop_kick
    LIMIT 1
),
     flop_address AS (
         SELECT address
         FROM maker.flop_kick
         JOIN addresses on addresses.id = flop_kick.address_id
         LIMIT 1
     ),
     deals AS (
         SELECT deal.bid_id,
                flop_bid_lot.lot,
                flop_bid_bid.bid                                           AS bid_amount,
                'deal'::api.bid_act                                        AS act,
                headers.block_number                                       AS block_height,
                tx_idx,
                (SELECT * FROM flop_address) AS contract_address
         FROM maker.deal
                  LEFT JOIN headers ON deal.header_id = headers.id
                  LEFT JOIN maker.flop_bid_bid
                            ON deal.bid_id = flop_bid_bid.bid_id
                                AND flop_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flop_bid_lot
                            ON deal.bid_id = flop_bid_lot.bid_id
                                AND flop_bid_lot.block_number = headers.block_number
         WHERE deal.address_id = (SELECT * FROM address_id)
     ),
     yanks AS (
         SELECT yank.bid_id,
                flop_bid_lot.lot,
                flop_bid_bid.bid                                           AS bid_amount,
                'yank'::api.bid_act                                        AS act,
                headers.block_number                                       AS block_height,
                tx_idx,
                (SELECT * FROM flop_address) AS contract_address
         FROM maker.yank
                  LEFT JOIN headers ON yank.header_id = headers.id
                  LEFT JOIN maker.flop_bid_bid
                            ON yank.bid_id = flop_bid_bid.bid_id
                                AND flop_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flop_bid_lot
                            ON yank.bid_id = flop_bid_lot.bid_id
                                AND flop_bid_lot.block_number = headers.block_number
         WHERE yank.address_id = (SELECT * FROM address_id)
         ORDER BY block_height DESC
     ),
     ticks AS (
         SELECT tick.bid_id,
                flop_bid_lot.lot,
                flop_bid_bid.bid                                           AS bid_amount,
                'tick'::api.bid_act                                        AS act,
                headers.block_number                                       AS block_height,
                tx_idx,
                (SELECT * FROM flop_address) AS contract_address
         FROM maker.tick
                  LEFT JOIN headers on tick.header_id = headers.id
                  LEFT JOIN maker.flop_bid_bid
                            ON tick.bid_id = flop_bid_bid.bid_id
                                AND flop_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flop_bid_lot
                            ON tick.bid_id = flop_bid_lot.bid_id
                                AND flop_bid_lot.block_number = headers.block_number
         WHERE tick.address_id = (SELECT * FROM address_id)
     )

SELECT flop_kick.bid_id,
       lot,
       bid                                                             AS bid_amount,
       'kick'::api.bid_act                                             AS act,
       block_number                                                    AS block_height,
       tx_idx,
       (SELECT * FROM flop_address) AS contract_address
FROM maker.flop_kick
         LEFT JOIN headers ON flop_kick.header_id = headers.id
UNION
SELECT bid_id,
       lot,
       bid                                                        AS bid_amount,
       'dent'::api.bid_act                                        AS act,
       block_number                                               AS block_height,
       tx_idx,
       (SELECT * FROM flop_address) AS contract_address
FROM maker.dent
         LEFT JOIN headers ON dent.header_id = headers.id
WHERE dent.address_id = (SELECT * FROM address_id)
UNION
SELECT *
FROM deals
UNION
SELECT *
FROM yanks
UNION
SELECT *
FROM ticks
ORDER BY block_height DESC
LIMIT all_flop_bid_events.max_results OFFSET all_flop_bid_events.result_offset
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
DROP FUNCTION api.all_flop_bid_events(INTEGER, INTEGER);
DROP TYPE api.flop_bid_event CASCADE;