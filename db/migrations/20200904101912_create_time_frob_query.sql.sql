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

CREATE FUNCTION api.time_frob_totals(ilk_identifier TEXT, range_start TIMESTAMP, range_end TIMESTAMP, bucket_interval INTERVAL DEFAULT '1 day'::INTERVAL)
    RETURNS SETOF api.time_frob_total AS
$$
WITH buckets AS (SELECT generate_series(range_start, range_end - bucket_interval, bucket_interval) AS bucket_start),
buckets_precompute AS (
    SELECT
        bucket_start,
        (extract(epoch FROM bucket_start))::INTEGER AS epoch_start,
        (extract(epoch FROM bucket_start) + 60)::INTEGER AS epoch_start_eps,
        (extract(epoch FROM bucket_start + bucket_interval) - 60)::INTEGER AS epoch_end_eps,
        (extract(epoch FROM bucket_start + bucket_interval))::INTEGER AS epoch_end
    FROM buckets
),
buckets_headers AS (
    SELECT
        bucket_start,
        (SELECT AVG(id) + (COUNT(id)/2.0) - 0.5 FROM public.headers WHERE block_timestamp >= epoch_start AND block_timestamp <= epoch_start_eps) AS header_id_start,
        (SELECT AVG(id) + (COUNT(id)/2.0) - 0.5 FROM public.headers WHERE block_timestamp >= epoch_end_eps AND block_timestamp <= epoch_end) AS header_id_end
    FROM buckets_precompute
)
SELECT buckets_headers.bucket_start AS bucket_start,
    buckets_headers.bucket_start + bucket_interval AS bucket_end,
    bucket_interval AS bucket_interval,
    COUNT(dink) AS count,
    COALESCE(SUM(dink), 0) AS dink,
    COALESCE(SUM(dart), 0) AS dart,
    COALESCE(SUM(GREATEST(dink, 0)), 0) AS lock,
    COALESCE(SUM(GREATEST(-dink, 0)), 0) AS free,
    COALESCE(SUM(GREATEST(dart, 0)), 0) AS draw,
    COALESCE(SUM(GREATEST(-dart, 0)), 0) AS wipe
FROM buckets_headers
    LEFT JOIN maker.vat_frob ON (header_id >= header_id_start AND header_id <= header_id_end)
    LEFT JOIN maker.urns ON (urns.id = urn_id)
    LEFT JOIN maker.ilks ON (ilks.id = ilk_id)
    WHERE ilks.identifier = time_frob_totals.ilk_identifier
GROUP BY bucket_start, bucket_end, bucket_interval
ORDER BY bucket_start
$$
    LANGUAGE sql
    STRICT
    STABLE;

-- +goose Down
DROP FUNCTION api.time_frob_totals(TEXT, TIMESTAMP, TIMESTAMP, INTERVAL);
DROP TYPE api.time_frob_total CASCADE;
