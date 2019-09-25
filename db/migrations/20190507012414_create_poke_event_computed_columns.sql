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
$$
SELECT * FROM get_tx_data(priceUpdate.block_height, priceUpdate.log_id)
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP FUNCTION api.poke_event_tx(api.poke_event);
DROP FUNCTION api.poke_event_ilk(api.poke_event);
