-- +goose Up

-- Extend queued sin with sin queue events
CREATE FUNCTION api.queued_sin_sin_queue_events(state api.queued_sin)
    RETURNS SETOF api.sin_queue_event AS
$$
SELECT *
FROM api.all_sin_queue_events(state.era)
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
DROP FUNCTION api.all_queued_sin_sin_queue_events(api.queued_sin)
