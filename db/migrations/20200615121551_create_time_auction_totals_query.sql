-- +goose Up
CREATE TYPE api.time_bid_total AS
(
    bucket_start TIMESTAMP,
    bucket_end TIMESTAMP,
    bucket_interval INTERVAL,
    lot NUMERIC,
    bid_amount NUMERIC
);

CREATE FUNCTION api.time_flip_bid_totals(ilk_identifier TEXT, range_start TIMESTAMP, range_end TIMESTAMP, bucket_interval INTERVAL DEFAULT '1 day'::INTERVAL)
    RETURNS SETOF api.time_bid_total AS
$$
WITH buckets AS (SELECT generate_series(range_start, range_end - bucket_interval, bucket_interval) AS bucket_start),
bid_results AS (
    SELECT bid_id,
        MAX(block_timestamp) AS block_timestamp,
        MIN(lot) AS lot,
        MAX(bid_amount) AS bid_amount
    FROM maker.bid_event
    LEFT JOIN public.headers ON (headers.block_number = block_height)
    WHERE ilk_identifier = time_flip_bid_totals.ilk_identifier
    GROUP BY bid_id
)
SELECT buckets.bucket_start AS bucket_start,
       buckets.bucket_start + bucket_interval AS bucket_end,
       bucket_interval AS bucket_interval,
       COALESCE(SUM(lot), 0) AS lot,
       COALESCE(SUM(bid_amount), 0) AS bid_amount
FROM buckets
    LEFT JOIN bid_results ON (
        block_timestamp >= extract(epoch FROM buckets.bucket_start) AND
        block_timestamp < extract(epoch FROM buckets.bucket_start + bucket_interval)
    )
GROUP BY bucket_start, bucket_end, bucket_interval
ORDER BY bucket_start
$$
    LANGUAGE sql
    STRICT
    STABLE;

CREATE FUNCTION api.time_flap_bid_totals(range_start TIMESTAMP, range_end TIMESTAMP, bucket_interval INTERVAL DEFAULT '1 day'::INTERVAL)
    RETURNS SETOF api.time_bid_total AS
$$
WITH buckets AS (SELECT generate_series(range_start, range_end - bucket_interval, bucket_interval) AS bucket_start),
flap_address AS (
    SELECT address
    FROM maker.flap_kick
        JOIN addresses on addresses.id = flap_kick.address_id
    LIMIT 1
),
bid_results AS (
    SELECT bid_id,
        MAX(block_timestamp) AS block_timestamp,
        MIN(lot) AS lot,
        MAX(bid_amount) AS bid_amount
    FROM maker.bid_event
    LEFT JOIN public.headers ON (headers.block_number = block_height)
    WHERE contract_address = (SELECT * FROM flap_address)
    GROUP BY bid_id
)
SELECT buckets.bucket_start AS bucket_start,
       buckets.bucket_start + bucket_interval AS bucket_end,
       bucket_interval AS bucket_interval,
       COALESCE(SUM(lot), 0) AS lot,
       COALESCE(SUM(bid_amount), 0) AS bid_amount
FROM buckets
    LEFT JOIN bid_results ON (
        block_timestamp >= extract(epoch FROM buckets.bucket_start) AND
        block_timestamp < extract(epoch FROM buckets.bucket_start + bucket_interval)
    )
GROUP BY bucket_start, bucket_end, bucket_interval
ORDER BY bucket_start
$$
    LANGUAGE sql
    STRICT
    STABLE;

CREATE FUNCTION api.time_flop_bid_totals(range_start TIMESTAMP, range_end TIMESTAMP, bucket_interval INTERVAL DEFAULT '1 day'::INTERVAL)
    RETURNS SETOF api.time_bid_total AS
$$
WITH buckets AS (SELECT generate_series(range_start, range_end - bucket_interval, bucket_interval) AS bucket_start),
flop_address AS (
    SELECT address
    FROM maker.flop_kick
        JOIN addresses on addresses.id = flop_kick.address_id
    LIMIT 1
),
bid_results AS (
    SELECT bid_id,
        MAX(block_timestamp) AS block_timestamp,
        MIN(lot) AS lot,
        MAX(bid_amount) AS bid_amount
    FROM maker.bid_event
    LEFT JOIN public.headers ON (headers.block_number = block_height)
    WHERE contract_address = (SELECT * FROM flop_address)
    GROUP BY bid_id
)
SELECT buckets.bucket_start AS bucket_start,
       buckets.bucket_start + bucket_interval AS bucket_end,
       bucket_interval AS bucket_interval,
       COALESCE(SUM(lot), 0) AS lot,
       COALESCE(SUM(bid_amount), 0) AS bid_amount
FROM buckets
    LEFT JOIN bid_results ON (
        block_timestamp >= extract(epoch FROM buckets.bucket_start) AND
        block_timestamp < extract(epoch FROM buckets.bucket_start + bucket_interval)
    )
GROUP BY bucket_start, bucket_end, bucket_interval
ORDER BY bucket_start
$$
    LANGUAGE sql
    STRICT
    STABLE;

-- +goose Down
DROP FUNCTION api.time_flip_bid_totals(TEXT, TIMESTAMP, TIMESTAMP, INTERVAL);
DROP FUNCTION api.time_flap_bid_totals(TIMESTAMP, TIMESTAMP, INTERVAL);
DROP FUNCTION api.time_flop_bid_totals(TIMESTAMP, TIMESTAMP, INTERVAL);
DROP TYPE api.time_bid_total CASCADE;