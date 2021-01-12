-- +goose Up

CREATE OR REPLACE FUNCTION api.get_block_heights_for_transformed_diffs()
    RETURNS SETOF BIGINT
    LANGUAGE sql
    STABLE
AS
$$
SELECT MAX (block_height) FROM public.storage_diff WHERE status = 'transformed'
$$;

-- +goose Down

DROP FUNCTION api.get_block_heights_for_transformed_diffs();
