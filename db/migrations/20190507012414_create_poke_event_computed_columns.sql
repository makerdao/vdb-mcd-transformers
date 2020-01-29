-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend type poke_event with ilk field
CREATE FUNCTION api.poke_event_ilk(event api.poke_event)
    RETURNS api.ilk_snapshot AS
$$
SELECT i.*
FROM api.ilk_snapshot i
         LEFT JOIN maker.ilks ON ilks.identifier = i.ilk_identifier
WHERE ilks.id = event.ilk_id
  AND i.block_number <= event.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$
    LANGUAGE sql
    STABLE;

-- extend type poke_event with tx field
CREATE FUNCTION api.poke_event_tx(priceUpdate api.poke_event)
    RETURNS api.tx AS
$$
SELECT *
FROM get_tx_data(priceUpdate.block_height, priceUpdate.log_id)
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP FUNCTION api.poke_event_tx(api.poke_event);
DROP FUNCTION api.poke_event_ilk(api.poke_event);
