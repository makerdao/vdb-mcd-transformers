-- +goose Up

-- Function returning state for all ilks as of the given block height
CREATE FUNCTION api.all_ilks(block_height BIGINT DEFAULT api.max_block())
    RETURNS SETOF api.ilk_snapshot
AS
$$
SELECT DISTINCT ON (ilk_identifier) *
FROM api.ilk_snapshot
WHERE block_number <= block_height
ORDER BY ilk_identifier, block_number DESC
$$
    LANGUAGE SQL
    STABLE;

-- +goose Down
DROP FUNCTION api.all_ilks(BIGINT);
