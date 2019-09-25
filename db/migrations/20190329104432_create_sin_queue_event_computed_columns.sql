-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend type sin_queue_event with txs field
CREATE FUNCTION api.sin_queue_event_tx(event api.sin_queue_event)
    RETURNS api.tx AS
$$
SELECT * FROM get_tx_data(event.block_height, event.log_id)
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.sin_queue_event_tx(api.sin_queue_event);
