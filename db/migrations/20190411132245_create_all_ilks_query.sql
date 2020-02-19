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

COMMENT ON FUNCTION api.all_ilks(block_height bigint)
    IS E'Get the state of every Ilk as of a given block height. blockHeight argument is optional and defaults to the most recent block.';

-- +goose Down
DROP FUNCTION api.all_ilks(BIGINT);
