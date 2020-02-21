-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.poke_event AS
(
    ilk_id       INTEGER,
    -- ilk object
    val          NUMERIC,
    spot         NUMERIC,
    block_height BIGINT,
    log_id       BIGINT
    -- tx
);

COMMENT ON TYPE api.poke_event
    IS E'Poke event emitted by Spot contract, with nested data regarding associated Ilk and Tx.';

COMMENT ON COLUMN api.poke_event.ilk_id
    IS E'@omit';
COMMENT ON COLUMN api.poke_event.block_height
    IS E'@omit';
COMMENT ON COLUMN api.poke_event.log_id
    IS E'@omit';

CREATE FUNCTION api.max_timestamp()
    RETURNS NUMERIC AS
$$
SELECT max(block_timestamp)
FROM public.headers
$$
    LANGUAGE SQL
    STABLE;

COMMENT ON FUNCTION api.max_timestamp()
    IS E'@omit';

CREATE FUNCTION api.all_poke_events(beginTime NUMERIC DEFAULT 0, endTime NUMERIC DEFAULT api.max_timestamp(),
                                    max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.poke_event AS
$body$
SELECT ilk_id, "value" AS val, spot, block_number AS block_height, log_id
FROM maker.spot_poke
         LEFT JOIN public.headers ON spot_poke.header_id = headers.id
WHERE block_timestamp BETWEEN beginTime AND endTime
ORDER BY block_height DESC
LIMIT all_poke_events.max_results
OFFSET
all_poke_events.result_offset
$body$
    LANGUAGE sql
    STABLE;

COMMENT ON FUNCTION api.all_poke_events(beginTime NUMERIC, endTime NUMERIC, max_results INTEGER, result_offset INTEGER)
    IS E'Get all Poke events in a given range of time. All arguments are optional. beginTime (e.g. 1582164329) defaults to 0. endTime (e.g. 1582164329) defaults to the timestamp of the most recent block. maxResults defaults to null (no max). resultOffset defaults to 0.';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP FUNCTION api.all_poke_events(NUMERIC, NUMERIC, INTEGER, INTEGER);
DROP FUNCTION api.max_timestamp();
DROP TYPE api.poke_event CASCADE;