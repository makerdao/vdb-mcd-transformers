-- +goose Up
CREATE OR REPLACE FUNCTION api.get_max_transformed_event_block()
    RETURNS BIGINT
    LANGUAGE sql
    STABLE
AS
$$
SELECT h.block_number
FROM public.event_logs
         JOIN headers h on h.id = event_logs.header_id
WHERE transformed = true
ORDER BY h.block_number DESC
LIMIT 1
$$;

-- +goose Down
DROP FUNCTION api.get_max_transformed_event_block();
