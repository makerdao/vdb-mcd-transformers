-- +goose Up
CREATE TYPE api.time_ilk_snapshot AS
(
    bucket_start TIMESTAMP,
    bucket_end TIMESTAMP,
    bucket_interval INTERVAL,
    ilk_identifier TEXT,
    block_number BIGINT,
    rate NUMERIC,
    art NUMERIC,
    spot NUMERIC,
    line NUMERIC,
    dust NUMERIC,
    chop NUMERIC,
    lump NUMERIC,
    flip TEXT,
    rho NUMERIC,
    duty NUMERIC,
    pip TEXT,
    mat NUMERIC,
    created TIMESTAMP,
    updated TIMESTAMP
);

CREATE FUNCTION api.time_ilk_snapshots(ilk_identifier TEXT, bucket_start TIMESTAMP, bucket_end TIMESTAMP, bucket_interval INTERVAL)
    RETURNS SETOF api.time_ilk_snapshot AS
$$
WITH buckets AS (SELECT generate_series(bucket_start, bucket_end - bucket_interval, bucket_interval) AS bucket_start)
SELECT buckets.bucket_start AS bucket_start,
       buckets.bucket_start + bucket_interval AS bucket_end,
       bucket_interval AS bucket_interval,
       ilk_identifier,
       block_number,
       rate,
       art,
       spot,
       line,
       dust,
       chop,
       lump,
       flip,
       rho,
       duty,
       pip,
       mat,
       created,
       updated
FROM buckets
    LEFT JOIN api.ilk_snapshot ON
    (
        ilk_snapshot.ilk_identifier = ilk_identifier AND
        block_number = (
            SELECT MAX(block_number)
            FROM api.ilk_snapshot
            WHERE ilk_snapshot.ilk_identifier = ilk_identifier AND updated < buckets.bucket_start + bucket_interval
        )
    )
ORDER BY bucket_start
$$
    LANGUAGE sql
    STRICT
    STABLE;

-- +goose Down
DROP FUNCTION api.time_ilk_snapshots(TEXT, TIMESTAMP, TIMESTAMP, INTERVAL);
DROP TYPE api.time_ilk_snapshot CASCADE;
