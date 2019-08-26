-- +goose Up
CREATE TYPE api.flip_bid_event AS (
    bid_id NUMERIC,
    lot NUMERIC,
    bid_amount NUMERIC,
    act api.bid_act,
    block_height BIGINT,
    tx_idx INTEGER,
    contract_address TEXT
    );

COMMENT ON COLUMN api.flip_bid_event.block_height
    IS E'@omit';
COMMENT ON COLUMN api.flip_bid_event.tx_idx
    IS E'@omit';
COMMENT ON COLUMN api.flip_bid_event.contract_address
    IS E'@omit';

CREATE FUNCTION api.all_flip_bid_events()
    RETURNS SETOF api.flip_bid_event AS
$$
WITH addresses AS (
    SELECT distinct contract_address
    FROM maker.flip_kick
),
     deals AS (
         SELECT deal.bid_id,
                flip_bid_lot.lot,
                flip_bid_bid.bid     AS bid_amount,
                'deal'::api.bid_act  AS act,
                headers.block_number AS block_height,
                tx_idx,
                deal.contract_address
         FROM maker.deal
                  LEFT JOIN headers ON deal.header_id = headers.id
                  LEFT JOIN maker.flip_bid_bid
                            ON deal.bid_id = flip_bid_bid.bid_id
                                AND flip_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flip_bid_lot
                            ON deal.bid_id = flip_bid_lot.bid_id
                                AND flip_bid_lot.block_number = headers.block_number
         WHERE deal.contract_address IN (SELECT * FROM addresses)
         ORDER BY block_height DESC
     ),
     yanks AS (
         SELECT yank.bid_id,
                flip_bid_lot.lot,
                flip_bid_bid.bid     AS bid_amount,
                'yank'::api.bid_act  AS act,
                headers.block_number AS block_height,
                tx_idx,
                yank.contract_address
         FROM maker.yank
                  LEFT JOIN headers ON yank.header_id = headers.id
                  LEFT JOIN maker.flip_bid_bid
                            ON yank.bid_id = flip_bid_bid.bid_id
                                AND flip_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flip_bid_lot
                            ON yank.bid_id = flip_bid_lot.bid_id
                                AND flip_bid_lot.block_number = headers.block_number
         WHERE yank.contract_address IN (SELECT * FROM addresses)
         ORDER BY block_height DESC
     ),
     ticks AS (
         SELECT flip_tick.bid_id,
                flip_bid_lot.lot,
                flip_bid_bid.bid     AS bid_amount,
                'tick'::api.bid_act  AS act,
                headers.block_number AS block_height,
                tx_idx,
                flip_tick.contract_address
         FROM maker.flip_tick
                  LEFT JOIN headers on flip_tick.header_id = headers.id
                  LEFT JOIN maker.flip_bid_bid
                            ON flip_tick.bid_id = flip_bid_bid.bid_id
                                AND flip_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flip_bid_lot
                            ON flip_tick.bid_id = flip_bid_lot.bid_id
                                AND flip_bid_lot.block_number = headers.block_number
         ORDER BY block_height DESC
     )

SELECT flip_kick.bid_id,
       lot,
       bid                 AS bid_amount,
       'kick'::api.bid_act AS act,
       block_number        AS block_height,
       tx_idx,
       contract_address
FROM maker.flip_kick
         LEFT JOIN headers ON flip_kick.header_id = headers.id
UNION
SELECT bid_id,
       lot,
       bid                 AS bid_amount,
       'tend'::api.bid_act AS act,
       block_number        AS block_height,
       tx_idx,
       contract_address
FROM maker.tend
         LEFT JOIN headers on tend.header_id = headers.id
WHERE tend.contract_address IN (SELECT * FROM addresses)
UNION
SELECT bid_id,
       lot,
       bid                 AS bid_amount,
       'dent'::api.bid_act AS act,
       block_number        AS block_height,
       tx_idx,
       dent.contract_address
FROM maker.dent
         LEFT JOIN headers on dent.header_id = headers.id
WHERE dent.contract_address IN (SELECT * FROM addresses)
UNION
SELECT *
from deals
UNION
SELECT *
from yanks
UNION
SELECT *
FROM ticks
$$
    LANGUAGE sql
    STABLE;
-- +goose Down
DROP FUNCTION api.all_flip_bid_events();
DROP TYPE api.flip_bid_event CASCADE ;
