-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.sin_act AS ENUM (
    'flog',
    'fess'
    );

CREATE TYPE api.sin_queue_event AS (
    era NUMERIC,
    act api.sin_act,
    block_height BIGINT,
    log_id BIGINT
    -- tx
    );

COMMENT ON COLUMN api.sin_queue_event.block_height
    IS E'@omit';
COMMENT ON COLUMN api.sin_queue_event.log_id
    IS E'@omit';

CREATE FUNCTION api.all_sin_queue_events(era NUMERIC, max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.sin_queue_event AS
$$
SELECT block_timestamp AS era, 'fess' :: api.sin_act AS act, block_number AS block_height, log_id
FROM maker.vow_fess
         LEFT JOIN headers ON vow_fess.header_id = headers.id
WHERE block_timestamp = all_sin_queue_events.era
UNION
SELECT era, 'flog' :: api.sin_act AS act, block_number AS block_height, log_id
FROM maker.vow_flog
         LEFT JOIN headers ON vow_flog.header_id = headers.id
WHERE vow_flog.era = all_sin_queue_events.era
ORDER BY block_height DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
OFFSET
all_sin_queue_events.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.all_sin_queue_events(NUMERIC, INTEGER, INTEGER);
DROP TYPE api.sin_queue_event;
DROP TYPE api.sin_act;
