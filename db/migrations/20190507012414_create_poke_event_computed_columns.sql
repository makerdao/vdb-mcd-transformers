-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend type poke_event with ilk field
CREATE FUNCTION api.poke_event_ilk(priceUpdate api.poke_event)
    RETURNS api.ilk_state AS
$$
WITH raw_ilk AS (SELECT * FROM maker.ilks WHERE ilks.id = priceUpdate.ilk_id)

SELECT *
FROM api.get_ilk((SELECT identifier FROM raw_ilk), priceUpdate.block_height)
$$
    LANGUAGE sql
    STABLE;

-- extend type poke_event with tx field
CREATE FUNCTION api.poke_event_tx(priceUpdate api.poke_event)
    RETURNS api.tx AS
$body$
SELECT txs.hash, txs.tx_index, headers.block_number, headers.hash, txs.tx_from, txs.tx_to
FROM public.header_sync_transactions txs
         LEFT JOIN headers ON txs.header_id = headers.id
WHERE headers.block_number = priceUpdate.block_height
  AND txs.tx_index = priceUpdate.tx_idx
ORDER BY headers.block_number DESC
$body$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP FUNCTION api.poke_event_tx(api.poke_event);
DROP FUNCTION api.poke_event_ilk(api.poke_event);
