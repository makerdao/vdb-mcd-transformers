-- +goose Up

CREATE FUNCTION api.epoch_to_datetime(epoch NUMERIC)
    RETURNS TIMESTAMP AS
$$
SELECT TIMESTAMP 'epoch' + epoch * INTERVAL '1 second' AS datetime
$$
    LANGUAGE SQL
    IMMUTABLE;

CREATE FUNCTION api.max_block()
    RETURNS BIGINT AS
$$
SELECT max(block_number)
FROM public.headers
$$
    LANGUAGE SQL
    STABLE;

-- +goose Down
DROP FUNCTION api.max_block();
DROP FUNCTION api.epoch_to_datetime(NUMERIC);
