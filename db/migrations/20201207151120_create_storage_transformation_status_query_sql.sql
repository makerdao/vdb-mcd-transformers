-- +goose Up

CREATE OR REPLACE FUNCTION api.storage_transformation_status()
    RETURNS SETOF BIGINT
    LANGUAGE sql
    STABLE
AS
$$
SELECT block_height FROM public.storage_diff WHERE status = 'new' ORDER BY block_height ASC
$$;


-- +goose Down

DROP FUNCTION api.storage_transformation_status();
