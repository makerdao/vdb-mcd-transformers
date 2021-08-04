-- +goose Up
CREATE TYPE api.sale_act AS ENUM (
    'kick',
    'take',
    'redo',
    'yank'
    );

CREATE TYPE api.clip_sale_event AS
(
    sale_id          numeric,
--     top              numeric,
--     tab              numeric,
--     lot              numeric,
--     usr              bigint,
--     kpr              bigint,
--     coin             numeric,
--     max              numeric,
--     price            numeric,
--     owe              numeric,
    act              api.sale_act,
    block_height     bigint,
    log_id           bigint,
    contract_address text
);

CREATE FUNCTION api.all_clip_sale_events(max_results integer DEFAULT NULL::integer,
                                         result_offset integer DEFAULT 0) RETURNS SETOF api.clip_sale_event
    LANGUAGE sql
    STABLE
    AS $$
WITH address_ids AS (
    SELECT distinct address_id
    FROM maker.clip_kick
)

SELECT sale_id,
--        top,
--        tab,
--        lot,
--        usr,
--        kpr,
--        coin,
       'kick'::api.sale_act AS                                          act,
       block_number        AS                                          block_height,
       log_id,
       (SELECT address FROM addresses WHERE id = clip_kick.address_id)
FROM maker.clip_kick
         LEFT JOIN headers ON clip_kick.header_id = headers.id
UNION
SELECT sale_id,
--        max,
--        price,
--        owe,
--        tab,
--        lot,
--        usr,
       'take'::api.sale_act AS act,
       block_number        AS block_height,
       log_id,
       (SELECT address FROM addresses WHERE id = clip_take.address_id)
FROM maker.clip_take
         LEFT JOIN headers on clip_take.header_id = headers.id
WHERE clip_take.address_id IN (SELECT * FROM address_ids)
UNION
SELECT sale_id,
--        top,
--        tab,
--        lot,
--        usr,
--        kpr,
--        coin,
       'redo'::api.sale_act AS                                          act,
       block_number        AS                                          block_height,
       log_id,
       (SELECT address FROM addresses WHERE id = clip_redo.address_id)
FROM maker.clip_redo
         LEFT JOIN headers ON clip_redo.header_id = headers.id
WHERE clip_redo.address_id IN (SELECT * FROM address_ids)
UNION
SELECT sale_id,
       'yank'::api.sale_act AS                                             act,
       block_number        AS                                          block_height,
       log_id,
       (SELECT address FROM addresses WHERE id = clip_yank.address_id)
FROM maker.clip_yank
         LEFT JOIN headers ON clip_yank.header_id = headers.id
WHERE clip_yank.address_id IN (SELECT * FROM address_ids)
ORDER BY block_height DESC
LIMIT all_clip_sale_events.max_results
    OFFSET
    all_clip_sale_events.result_offset
$$;



-- +goose Down
DROP FUNCTION api.all_clip_sale_events(INTEGER, INTEGER);
DROP TYPE api.sale_act CASCADE;
DROP TYPE api.clip_sale_event CASCADE;