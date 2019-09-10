-- +goose Up

-- Extend queued sin with sin queue events
CREATE FUNCTION api.queued_sin_sin_queue_events(state api.queued_sin, max_results INTEGER DEFAULT NULL)
    RETURNS SETOF api.sin_queue_event AS
$$
SELECT *
FROM api.all_sin_queue_events(state.era)
LIMIT queued_sin_sin_queue_events.max_results
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
DROP FUNCTION api.queued_sin_sin_queue_events(api.queued_sin, INTEGER);
