-- +goose Up
CREATE TYPE api.time_bid_total AS
(
    bucket_start TIMESTAMP,
    bucket_end TIMESTAMP,
    bucket_interval INTERVAL,
    count BIGINT,
    lot_start NUMERIC,
    lot_end NUMERIC,
    bid_amount_start NUMERIC,
    bid_amount_end NUMERIC
);

CREATE FUNCTION api.time_flip_bid_totals(ilk_identifier TEXT, range_start TIMESTAMP, range_end TIMESTAMP, bucket_interval INTERVAL DEFAULT '1 day'::INTERVAL)
    RETURNS SETOF api.time_bid_total AS
$$
DECLARE
    r api.time_bid_total%rowtype;
BEGIN
    ASSERT EXTRACT(EPOCH FROM (range_end - range_start)) / EXTRACT(EPOCH FROM bucket_interval) <= 100, 'Please limit requests to at most 100 buckets.';
    
    FOR r IN 
        WITH buckets AS (
            SELECT generate_series(range_start, range_end - bucket_interval, bucket_interval) AS bucket_start,
            extract(epoch FROM generate_series(range_start, range_end - bucket_interval, bucket_interval)) AS bucket_start_epoch,
            extract(epoch FROM generate_series(range_start + bucket_interval, range_end, bucket_interval)) AS bucket_end_epoch
        ),
        bid_results AS (
            SELECT bid_id,
                contract_address,
                MIN(block_timestamp) AS block_timestamp,
                MAX(lot) AS lot_start,
                MIN(lot) AS lot_end,
                MIN(bid_amount) AS bid_amount_start,
                MAX(bid_amount) AS bid_amount_end
            FROM maker.bid_event
            LEFT JOIN public.headers ON (headers.block_number = block_height)
            WHERE bid_event.ilk_identifier = time_flip_bid_totals.ilk_identifier
            GROUP BY contract_address, bid_id
        )
        SELECT buckets.bucket_start AS bucket_start,
            buckets.bucket_start + bucket_interval AS bucket_end,
            bucket_interval AS bucket_interval,
            COUNT(lot_start) AS count,
            COALESCE(SUM(lot_start), 0) AS lot_start,
            COALESCE(SUM(lot_end), 0) AS lot_end,
            COALESCE(SUM(bid_amount_start), 0) AS bid_amount_start,
            COALESCE(SUM(bid_amount_end), 0) AS bid_amount_end
        FROM buckets
            LEFT JOIN bid_results ON (
                block_timestamp >= bucket_start_epoch AND
                block_timestamp < bucket_end_epoch
            )
        GROUP BY bucket_start, bucket_end, bucket_interval
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

CREATE FUNCTION api.time_flap_bid_totals(range_start TIMESTAMP, range_end TIMESTAMP, bucket_interval INTERVAL DEFAULT '1 day'::INTERVAL)
    RETURNS SETOF api.time_bid_total AS
$$
DECLARE
    r api.time_bid_total%rowtype;
BEGIN
    ASSERT EXTRACT(EPOCH FROM (range_end - range_start)) / EXTRACT(EPOCH FROM bucket_interval) <= 100, 'Please limit requests to at most 100 buckets.';
    
    FOR r IN 
        WITH buckets AS (
            SELECT generate_series(range_start, range_end - bucket_interval, bucket_interval) AS bucket_start,
            extract(epoch FROM generate_series(range_start, range_end - bucket_interval, bucket_interval)) AS bucket_start_epoch,
            extract(epoch FROM generate_series(range_start + bucket_interval, range_end, bucket_interval)) AS bucket_end_epoch
        ),
        flap_address AS (
            SELECT DISTINCT(address)
            FROM maker.flap_kick
                JOIN addresses on addresses.id = flap_kick.address_id
        ),
        bid_results AS (
            SELECT bid_id,
                contract_address,
                MIN(block_timestamp) AS block_timestamp,
                MAX(lot) AS lot_start,
                MIN(lot) AS lot_end,
                MIN(bid_amount) AS bid_amount_start,
                MAX(bid_amount) AS bid_amount_end
            FROM maker.bid_event
            LEFT JOIN public.headers ON (headers.block_number = block_height)
            WHERE contract_address IN (SELECT * FROM flap_address)
            GROUP BY contract_address, bid_id
        )
        SELECT buckets.bucket_start AS bucket_start,
            buckets.bucket_start + bucket_interval AS bucket_end,
            bucket_interval AS bucket_interval,
            COUNT(lot_start) AS count,
            COALESCE(SUM(lot_start), 0) AS lot_start,
            COALESCE(SUM(lot_end), 0) AS lot_end,
            COALESCE(SUM(bid_amount_start), 0) AS bid_amount_start,
            COALESCE(SUM(bid_amount_end), 0) AS bid_amount_end
        FROM buckets
            LEFT JOIN bid_results ON (
                block_timestamp >= bucket_start_epoch AND
                block_timestamp < bucket_end_epoch
            )
        GROUP BY bucket_start, bucket_end, bucket_interval
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

CREATE FUNCTION api.time_flop_bid_totals(range_start TIMESTAMP, range_end TIMESTAMP, bucket_interval INTERVAL DEFAULT '1 day'::INTERVAL)
    RETURNS SETOF api.time_bid_total AS
$$
DECLARE
    r api.time_bid_total%rowtype;
BEGIN
    ASSERT EXTRACT(EPOCH FROM (range_end - range_start)) / EXTRACT(EPOCH FROM bucket_interval) <= 100, 'Please limit requests to at most 100 buckets.';
    
    FOR r IN 
        WITH buckets AS (
            SELECT generate_series(range_start, range_end - bucket_interval, bucket_interval) AS bucket_start,
            extract(epoch FROM generate_series(range_start, range_end - bucket_interval, bucket_interval)) AS bucket_start_epoch,
            extract(epoch FROM generate_series(range_start + bucket_interval, range_end, bucket_interval)) AS bucket_end_epoch
        ),
        flop_address AS (
            SELECT DISTINCT(address)
            FROM maker.flop_kick
                JOIN addresses on addresses.id = flop_kick.address_id
        ),
        bid_results AS (
            SELECT bid_id,
                contract_address,
                MIN(block_timestamp) AS block_timestamp,
                MAX(lot) AS lot_start,
                MIN(lot) AS lot_end,
                MIN(bid_amount) AS bid_amount_start,
                MAX(bid_amount) AS bid_amount_end
            FROM maker.bid_event
            LEFT JOIN public.headers ON (headers.block_number = block_height)
            WHERE contract_address IN (SELECT * FROM flop_address)
            GROUP BY contract_address, bid_id
        )
        SELECT buckets.bucket_start AS bucket_start,
            buckets.bucket_start + bucket_interval AS bucket_end,
            bucket_interval AS bucket_interval,
            COUNT(lot_start) AS count,
            COALESCE(SUM(lot_start), 0) AS lot_start,
            COALESCE(SUM(lot_end), 0) AS lot_end,
            COALESCE(SUM(bid_amount_start), 0) AS bid_amount_start,
            COALESCE(SUM(bid_amount_end), 0) AS bid_amount_end
        FROM buckets
            LEFT JOIN bid_results ON (
                block_timestamp >= bucket_start_epoch AND
                block_timestamp < bucket_end_epoch
            )
        GROUP BY bucket_start, bucket_end, bucket_interval
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

CREATE TYPE api.time_bite_total AS
(
    bucket_start TIMESTAMP,
    bucket_end TIMESTAMP,
    bucket_interval INTERVAL,
    count BIGINT,
    ink NUMERIC,
    art NUMERIC,
    tab NUMERIC
);

CREATE FUNCTION api.time_bite_totals(ilk_identifier TEXT, range_start TIMESTAMP, range_end TIMESTAMP, bucket_interval INTERVAL DEFAULT '1 day'::INTERVAL)
    RETURNS SETOF api.time_bite_total AS
$$
DECLARE
    r api.time_bite_total%rowtype;
BEGIN
    ASSERT EXTRACT(EPOCH FROM (range_end - range_start)) / EXTRACT(EPOCH FROM bucket_interval) <= 100, 'Please limit requests to at most 100 buckets.';
    
    FOR r IN 
        WITH buckets AS (
            SELECT generate_series(range_start, range_end - bucket_interval, bucket_interval) AS bucket_start,
            extract(epoch FROM generate_series(range_start, range_end - bucket_interval, bucket_interval)) AS bucket_start_epoch,
            extract(epoch FROM generate_series(range_start + bucket_interval, range_end, bucket_interval)) AS bucket_end_epoch
        ),
        bite_results AS (
            SELECT *
            FROM maker.bite
            LEFT JOIN public.headers ON (headers.id = bite.header_id)
            LEFT JOIN maker.urns ON (urns.id = bite.urn_id)
            LEFT JOIN maker.ilks ON (ilks.id = urns.ilk_id)
            WHERE ilks.identifier = time_bite_totals.ilk_identifier
        )
        SELECT buckets.bucket_start AS bucket_start,
            buckets.bucket_start + bucket_interval AS bucket_end,
            bucket_interval AS bucket_interval,
            COUNT(ink) AS count,
            COALESCE(SUM(ink), 0) AS ink,
            COALESCE(SUM(art), 0) AS art,
            COALESCE(SUM(tab), 0) AS tab
        FROM buckets
            LEFT JOIN bite_results ON (
                block_timestamp >= bucket_start_epoch AND
                block_timestamp < bucket_end_epoch
            )
        GROUP BY bucket_start, bucket_end, bucket_interval
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

-- +goose Down
DROP FUNCTION api.time_flip_bid_totals(TEXT, TIMESTAMP, TIMESTAMP, INTERVAL);
DROP FUNCTION api.time_flap_bid_totals(TIMESTAMP, TIMESTAMP, INTERVAL);
DROP FUNCTION api.time_flop_bid_totals(TIMESTAMP, TIMESTAMP, INTERVAL);
DROP TYPE api.time_bid_total CASCADE;
DROP FUNCTION api.time_bite_totals(TEXT, TIMESTAMP, TIMESTAMP, INTERVAL);
DROP TYPE api.time_bite_total CASCADE;