-- +goose Up

CREATE OR REPLACE FUNCTION api.new_storage_block_heights()
    RETURNS SETOF BIGINT
    LANGUAGE sql
    STABLE
AS
$$
SELECT block_height FROM public.storage_diff WHERE status = 'new' ORDER BY block_height ASC
$$;


-- +goose Down

DROP FUNCTION api.new_storage_block_heights();
