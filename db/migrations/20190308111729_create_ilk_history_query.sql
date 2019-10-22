-- +goose Up
-- +goose StatementBegin

-- Function returning the history of a given ilk as of the given block height
CREATE FUNCTION api.all_ilk_states(ilk_identifier TEXT, block_height BIGINT DEFAULT api.max_block(),
                                   max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.ilk_state AS
$$
BEGIN
    RETURN QUERY (
        WITH relevant_blocks AS (
            SELECT get_ilk_blocks_before.block_height
            FROM api.get_ilk_blocks_before(ilk_identifier, all_ilk_states.block_height)
        )
        SELECT r.*
        FROM relevant_blocks,
             LATERAL api.get_ilk(ilk_identifier, relevant_blocks.block_height) r
        LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
        OFFSET
        all_ilk_states.result_offset
    );
END;
$$
    LANGUAGE plpgsql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.all_ilk_states(TEXT, BIGINT, INTEGER, INTEGER);
