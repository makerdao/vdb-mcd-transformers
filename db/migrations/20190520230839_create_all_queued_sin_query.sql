-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION api.all_queued_sin(max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.queued_sin AS
$$
BEGIN
    RETURN QUERY (
        WITH eras AS (
            SELECT DISTINCT era
            FROM maker.vow_sin_mapping
            ORDER BY era DESC
            LIMIT all_queued_sin.max_results
            OFFSET
            all_queued_sin.result_offset
        )
        SELECT sin.*
        FROM eras,
             LATERAL api.get_queued_sin(eras.era) sin
    );
END
$$
    LANGUAGE plpgsql
    STABLE;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.all_queued_sin(INTEGER, INTEGER);