-- +goose Up
-- +goose StatementBegin
CREATE TYPE api.bid_act AS ENUM (
    'kick',
    'tick',
    'tend',
    'dent',
    'deal',
    'yank'
    );

CREATE TYPE api.flap_bid_event AS (
    bid_id NUMERIC,
    lot NUMERIC,
    bid_amount NUMERIC,
    act api.bid_act,
    block_height BIGINT,
    tx_idx INTEGER,
    contract_address TEXT
    );

COMMENT ON COLUMN api.flap_bid_event.block_height
    IS E'@omit';
COMMENT ON COLUMN api.flap_bid_event.tx_idx
    IS E'@omit';
COMMENT ON COLUMN api.flap_bid_event.contract_address
    IS E'@omit';

CREATE FUNCTION api.all_flap_bid_events()
    RETURNS SETOF api.flap_bid_event AS
$$
WITH address AS (
    SELECT contract_address
    FROM maker.flap_kick
    LIMIT 1
),
     deals AS (
         SELECT deal.bid_id,
                flap_bid_lot.lot,
                flap_bid_bid.bid     AS bid_amount,
                'deal'::api.bid_act  AS act,
                headers.block_number AS block_height,
                tx_idx,
                deal.contract_address
         FROM maker.deal
                  LEFT JOIN headers ON deal.header_id = headers.id
                  LEFT JOIN maker.flap_bid_bid
                            ON deal.bid_id = flap_bid_bid.bid_id
                                AND flap_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flap_bid_lot
                            ON deal.bid_id = flap_bid_lot.bid_id
                                AND flap_bid_lot.block_number = headers.block_number
         WHERE deal.contract_address = (SELECT * FROM address)
         ORDER BY block_height DESC
     ),
     yanks AS (
         SELECT yank.bid_id,
                flap_bid_lot.lot,
                flap_bid_bid.bid     AS bid_amount,
                'yank'::api.bid_act  AS act,
                headers.block_number AS block_height,
                tx_idx,
                yank.contract_address
         FROM maker.yank
                  LEFT JOIN headers ON yank.header_id = headers.id
                  LEFT JOIN maker.flap_bid_bid
                            ON yank.bid_id = flap_bid_bid.bid_id
                                AND flap_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flap_bid_lot
                            ON yank.bid_id = flap_bid_lot.bid_id
                                AND flap_bid_lot.block_number = headers.block_number
         WHERE yank.contract_address = (SELECT * FROM address)
         ORDER BY block_height DESC
     )

SELECT flap_kick.bid_id,
       lot,
       bid                 AS bid_amount,
       'kick'::api.bid_act AS act,
       block_number        AS block_height,
       tx_idx,
       contract_address
FROM maker.flap_kick
         LEFT JOIN headers ON flap_kick.header_id = headers.id
UNION
SELECT bid_id,
       lot,
       bid                 AS bid_amount,
       'tend'::api.bid_act AS act,
       block_number        AS block_height,
       tx_idx,
       contract_address
FROM maker.tend
         LEFT JOIN headers ON tend.header_id = headers.id
WHERE tend.contract_address = (SELECT * FROM address)
UNION
SELECT *
FROM deals
UNION
SELECT *
FROM yanks

$$
    LANGUAGE sql
    STABLE;

-- +goose StatementEnd
-- +goose Down
DROP FUNCTION api.all_flap_bid_events();
DROP TYPE api.flap_bid_event CASCADE;
DROP TYPE api.bid_act CASCADE;