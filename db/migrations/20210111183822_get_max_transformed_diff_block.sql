-- +goose Up

CREATE OR REPLACE FUNCTION api.get_max_transformed_diff_block()
    RETURNS BIGINT
    LANGUAGE sql
    STABLE
AS
$$
SELECT MAX (block_height) FROM public.storage_diff WHERE status = 'transformed'
$$;

-- +goose Down

DROP FUNCTION api.get_max_transformed_diff_block();
