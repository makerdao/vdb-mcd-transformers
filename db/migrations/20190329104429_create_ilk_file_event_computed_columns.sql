-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend file_event with ilk_state
CREATE FUNCTION api.ilk_file_event_ilk(event api.ilk_file_event)
    RETURNS SETOF api.ilk_state AS
$$
SELECT *
FROM api.get_ilk(event.ilk_identifier, event.block_height)
$$
    LANGUAGE sql
    STABLE;

-- Extend file_event with txs
CREATE FUNCTION api.ilk_file_event_tx(event api.ilk_file_event)
    RETURNS api.tx AS
$$
SELECT * FROM get_tx_data(event.block_height, event.log_id)
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.ilk_file_event_ilk(api.ilk_file_event);
DROP FUNCTION api.ilk_file_event_tx(api.ilk_file_event);
