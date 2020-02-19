-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION api.all_urn_states(ilk_identifier TEXT, urn_identifier TEXT,
                                   block_height BIGINT DEFAULT api.max_block(),
                                   max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.urn_state AS
$$
BEGIN
    RETURN QUERY (
        WITH urn_id AS (
            SELECT id
            FROM maker.urns
            WHERE urns.identifier = all_urn_states.urn_identifier
              AND urns.ilk_id = (SELECT id
                                 FROM maker.ilks
                                 WHERE ilks.identifier = all_urn_states.ilk_identifier)
        ),
             relevant_blocks AS (
                 SELECT block_number
                 FROM maker.vat_urn_ink
                          LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
                 WHERE vat_urn_ink.urn_id = (SELECT * FROM urn_id)
                   AND block_number <= all_urn_states.block_height
                 UNION
                 SELECT block_number
                 FROM maker.vat_urn_art
                          LEFT JOIN public.headers ON vat_urn_art.header_id = headers.id
                 WHERE vat_urn_art.urn_id = (SELECT * FROM urn_id)
                   AND block_number <= all_urn_states.block_height)
        SELECT r.*
        FROM relevant_blocks,
             LATERAL api.get_urn(ilk_identifier, urn_identifier, relevant_blocks.block_number) r
        ORDER BY relevant_blocks.block_number DESC
        LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
        OFFSET
        all_urn_states.result_offset
    );
END;
$$
    LANGUAGE plpgsql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;
-- +goose StatementEnd

COMMENT ON FUNCTION api.all_urn_states(ilk_identifier TEXT, urn_identifier TEXT, block_height BIGINT, max_results INTEGER, result_offset INTEGER)
    IS E'Get all historical states for an Urn prior to a given block. ilkIdentifier (e.g. "ETH-A") and urnIdentifier (e.g. "0xC93C178EC17B06bddBa0CC798546161aF9D25e8A") are required. blockHeight, maxResults, and resultOffset are optional. blockHeight defaults to most recent block. maxResults defaults to null (no max). resultOffset defaults to 0.';

-- +goose Down
DROP FUNCTION api.all_urn_states(TEXT, TEXT, BIGINT, INTEGER, INTEGER);
