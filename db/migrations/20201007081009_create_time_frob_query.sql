-- +goose Up
CREATE TYPE api.time_frob_total AS
(
    bucket_start TIMESTAMP,
    bucket_end TIMESTAMP,
    bucket_interval INTERVAL,
    count BIGINT,
    dink NUMERIC,
    dart NUMERIC,
    lock NUMERIC,
    free NUMERIC,
    draw NUMERIC,
    wipe NUMERIC
);

-- +goose StatementBegin
CREATE FUNCTION api.time_frob_totals(ilk_identifier TEXT, range_start TIMESTAMP, range_end TIMESTAMP, bucket_interval INTERVAL DEFAULT '1 day'::INTERVAL)
    RETURNS SETOF api.time_frob_total AS
$$
DECLARE
    r api.time_frob_total%rowtype;
BEGIN
    ASSERT EXTRACT(EPOCH FROM (range_end - range_start)) / EXTRACT(EPOCH FROM bucket_interval) <= 1000, 'Please limit requests to at most 1000 buckets.';
    
    FOR r IN 
        SELECT TO_TIMESTAMP(EXTRACT(EPOCH FROM range_start) + ROUND((headers.block_timestamp - EXTRACT(EPOCH FROM range_start)) / EXTRACT(EPOCH FROM bucket_interval)) * EXTRACT(EPOCH FROM bucket_interval)) AS bucket_start,
            TO_TIMESTAMP(EXTRACT(EPOCH FROM range_start) + ROUND((headers.block_timestamp - EXTRACT(EPOCH FROM range_start)) / EXTRACT(EPOCH FROM bucket_interval)) * EXTRACT(EPOCH FROM bucket_interval)) + bucket_interval AS bucket_end,
            bucket_interval,
            COUNT(dink) AS count,
            COALESCE(SUM(dink), 0) AS dink,
            COALESCE(SUM(dart), 0) AS dart,
            COALESCE(SUM(GREATEST(dink, 0)), 0) AS lock,
            COALESCE(SUM(GREATEST(-dink, 0)), 0) AS free,
            COALESCE(SUM(GREATEST(dart, 0)), 0) AS draw,
            COALESCE(SUM(GREATEST(-dart, 0)), 0) AS wipe
        FROM maker.vat_frob
            LEFT JOIN public.headers ON (vat_frob.header_id = headers.id)
            LEFT JOIN maker.urns ON (urns.id = urn_id)
            LEFT JOIN maker.ilks ON (ilks.id = ilk_id)
            WHERE ilks.identifier = time_frob_totals.ilk_identifier AND
                headers.block_timestamp >= EXTRACT(EPOCH FROM range_start) and
                headers.block_timestamp < EXTRACT(EPOCH FROM range_end)
        GROUP BY bucket_start
        ORDER BY bucket_start
    LOOP
        return next r;
    END LOOP;
    return;
END;
$$
    LANGUAGE plpgsql
    STRICT
    STABLE;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.time_frob_totals(TEXT, TIMESTAMP, TIMESTAMP, INTERVAL);
DROP TYPE api.time_frob_total CASCADE;
