-- +goose Up
CREATE INDEX ON api.ilk_snapshot (ilk_identifier, updated);

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
    dunk NUMERIC,
    flip TEXT,
    rho NUMERIC,
    duty NUMERIC,
    pip TEXT,
    mat NUMERIC,
    created TIMESTAMP,
    updated TIMESTAMP
);

-- +goose StatementBegin
CREATE FUNCTION api.time_ilk_snapshots(ilk_identifier TEXT, range_start TIMESTAMP, range_end TIMESTAMP, bucket_interval INTERVAL DEFAULT '1 day'::INTERVAL)
    RETURNS SETOF api.time_ilk_snapshot AS
$$
DECLARE
    r api.time_ilk_snapshot%rowtype;
BEGIN
    ASSERT EXTRACT(EPOCH FROM (range_end - range_start)) / EXTRACT(EPOCH FROM bucket_interval) <= 1000, 'Please limit requests to at most 1000 buckets.';
    
    FOR r IN 
        WITH buckets AS (SELECT generate_series(range_start, range_end - bucket_interval, bucket_interval) AS bucket_start)
        SELECT buckets.bucket_start,
            buckets.bucket_start + bucket_interval AS bucket_end,
            bucket_interval,
            time_ilk_snapshots.ilk_identifier,
            block_number,
            rate,
            art,
            spot,
            line,
            dust,
            chop,
            lump,
            dunk,
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
                ilk_snapshot.ilk_identifier = time_ilk_snapshots.ilk_identifier AND
                block_number = (
                    SELECT block_number
                    FROM api.ilk_snapshot
                    WHERE ilk_snapshot.ilk_identifier = time_ilk_snapshots.ilk_identifier AND updated < buckets.bucket_start + bucket_interval
                    ORDER BY updated DESC
                    LIMIT 1
                )
            )
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
DROP FUNCTION api.time_ilk_snapshots(TEXT, TIMESTAMP, TIMESTAMP, INTERVAL);
DROP TYPE api.time_ilk_snapshot CASCADE;
DROP INDEX api.ilk_snapshot_ilk_identifier_updated_idx;
