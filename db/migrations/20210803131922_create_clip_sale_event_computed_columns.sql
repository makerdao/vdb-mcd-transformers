-- +goose Up

CREATE FUNCTION api.clip_sale_event_tx(event api.clip_sale_event) RETURNS SETOF api.tx
    LANGUAGE sql
    STABLE
AS
$$
SELECT *
FROM get_tx_data(event.block_height, event.log_id)
$$;

CREATE TYPE api.clip_sale_snapshot AS
(
    block_height BIGINT,
    sale_id      NUMERIC,
    ilk_id       INTEGER,
    urn_id       INTEGER,
    pos          NUMERIC,
    tab          NUMERIC,
    lot          NUMERIC,
    usr          TEXT,
    tic          NUMERIC,
    "top"        NUMERIC,

    created      TIMESTAMP WITHOUT TIME ZONE,
    updated      TIMESTAMP WITHOUT TIME ZONE,
    clip_address TEXT
);

CREATE FUNCTION api.get_clip_with_address(sale_id NUMERIC, clip_address TEXT,
                                          block_height BIGINT DEFAULT api.max_block()) RETURNS api.clip_sale_snapshot
    LANGUAGE sql
    STABLE STRICT
AS
$$
WITH address_id AS (SELECT id FROM public.addresses WHERE address = get_clip_with_address.clip_address),
     ilk_id AS (SELECT ilk_id FROM maker.dog_bark WHERE clip = (SELECT id FROM address_id)),
     urn_id AS (SELECT urn_id FROM maker.dog_bark WHERE clip = (SELECT id FROM address_id)),
     storage_values AS (
         SELECT pos,
                tab,
                lot,
                usr,
                tic,
                top,
                created,
                updated,
                block_number
         FROM maker.clip
         WHERE clip.sale_id = get_clip_with_address.sale_id
           AND clip.address_id = (SELECT id FROM address_id)
           AND block_number <= get_clip_with_address.block_height
         ORDER BY block_number DESC
         LIMIT 1
     )
SELECT storage_values.block_number,
       get_clip_with_address.sale_id,
       (SELECT ilk_id FROM ilk_id),
       (SELECT urn_id FROM urn_id),
       storage_values.pos,
       storage_values.tab,
       storage_values.lot,
       storage_values.usr,
       storage_values.tic,
       storage_values.top,
       storage_values.created,
       storage_values.updated,
       get_clip_with_address.clip_address
FROM storage_values
$$;

CREATE FUNCTION api.clip_sale_event_sale(event api.clip_sale_event) RETURNS api.clip_sale_snapshot
    LANGUAGE sql
    STABLE
AS
$$
SELECT *
FROM api.get_clip_with_address(event.sale_id, event.contract_address, event.block_height)
$$;


CREATE FUNCTION api.clip_sale_snapshot_sale_events(clip api.clip_sale_snapshot,
                                                   max_results integer DEFAULT NULL::integer,
                                                   result_offset integer DEFAULT 0) RETURNS SETOF api.clip_sale_event
    LANGUAGE sql
    STABLE
AS
$$

SELECT sale_id, act, block_height, events.log_id, events.contract_address
FROM api.all_clip_sale_events() AS events
WHERE sale_id = clip.sale_id
  AND contract_address = clip.clip_address
ORDER BY block_height DESC
LIMIT clip_sale_snapshot_sale_events.max_results OFFSET clip_sale_snapshot_sale_events.result_offset
$$;


-- +goose Down
DROP FUNCTION api.clip_sale_event_tx(event api.clip_sale_event);
DROP FUNCTION api.clip_sale_event_sale(event api.clip_sale_event);
DROP FUNCTION api.get_clip_with_address(sale_id numeric, clip_address text, block_height bigint);
DROP FUNCTION api.clip_sale_snapshot_sale_events(clip api.clip_sale_snapshot, max_results integer, result_offset integer);
DROP TYPE api.clip_sale_snapshot CASCADE;