-- +goose Up
CREATE TABLE maker.flap
(
    address_id   INTEGER   NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    block_number BIGINT    NOT NULL,
    bid_id       NUMERIC   NOT NULL,
    guy          TEXT      DEFAULT NULL,
    tic          BIGINT    DEFAULT NULL,
    "end"        BIGINT    DEFAULT NULL,
    lot          NUMERIC   DEFAULT NULL,
    bid          NUMERIC   DEFAULT NULL,
    created      TIMESTAMP DEFAULT NULL,
    updated      TIMESTAMP NOT NULL,
    PRIMARY KEY (block_number, bid_id, address_id)
);

COMMENT ON TABLE maker.flap
    IS E'@name historicalFlapState';

CREATE INDEX flap_address_index
    ON maker.flap (address_id);

CREATE FUNCTION get_latest_flap_bid_guy(bid_id numeric) RETURNS TEXT AS
$$
SELECT guy
FROM maker.flap
WHERE guy IS NOT NULL
  AND flap.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION get_latest_flap_bid_guy
    IS E'@omit';

CREATE FUNCTION get_latest_flap_bid_bid(bid_id numeric) RETURNS NUMERIC AS
$$
SELECT bid
FROM maker.flap
WHERE bid IS NOT NULL
  AND flap.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION get_latest_flap_bid_bid
    IS E'@omit';

CREATE FUNCTION get_latest_flap_bid_tic(bid_id numeric) RETURNS BIGINT AS
$$
SELECT tic
FROM maker.flap
WHERE tic IS NOT NULL
  AND flap.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION get_latest_flap_bid_tic
    IS E'@omit';

CREATE FUNCTION get_latest_flap_bid_end(bid_id numeric) RETURNS BIGINT AS
$$
SELECT "end"
FROM maker.flap
WHERE "end" IS NOT NULL
  AND flap.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION get_latest_flap_bid_end
    IS E'@omit';

CREATE FUNCTION get_latest_flap_bid_lot(bid_id numeric) RETURNS NUMERIC AS
$$
SELECT lot
FROM maker.flap
WHERE lot IS NOT NULL
  AND flap.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION get_latest_flap_bid_lot
    IS E'@omit';

CREATE FUNCTION flap_bid_time_created(address_id INTEGER, bid_id NUMERIC) RETURNS TIMESTAMP AS
$$
SELECT api.epoch_to_datetime(MIN(block_timestamp))
FROM public.headers
         LEFT JOIN maker.flap_kick ON flap_kick.header_id = headers.id
WHERE flap_kick.address_id = flap_bid_time_created.address_id
  AND flap_kick.bid_id = flap_bid_time_created.bid_id
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flap_bid_time_created
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flap_bid() RETURNS TRIGGER
AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flap (bid_id, address_id, block_number, bid, guy, tic, "end", lot, updated, created)
    VALUES (NEW.bid_id, NEW.address_id, (SELECT block_number FROM diff_block), NEW.bid,
            get_latest_flap_bid_guy(NEW.bid_id),
            get_latest_flap_bid_tic(NEW.bid_id),
            get_latest_flap_bid_end(NEW.bid_id),
            get_latest_flap_bid_lot(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flap_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (address_id, bid_id, block_number) DO UPDATE SET bid = NEW.bid;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flap_guy() RETURNS TRIGGER
AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flap (bid_id, address_id, block_number, guy, bid, tic, "end", lot, updated, created)
    VALUES (NEW.bid_id, NEW.address_id, (SELECT block_number FROM diff_block), NEW.guy,
            get_latest_flap_bid_bid(NEW.bid_id),
            get_latest_flap_bid_tic(NEW.bid_id),
            get_latest_flap_bid_end(NEW.bid_id),
            get_latest_flap_bid_lot(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flap_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (address_id, bid_id, block_number) DO UPDATE SET guy = NEW.guy;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flap_tic() RETURNS TRIGGER
AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flap (bid_id, address_id, block_number, tic, bid, guy, "end", lot, updated, created)
    VALUES (NEW.bid_id, NEW.address_id, (SELECT block_number FROM diff_block), NEW.tic,
            get_latest_flap_bid_bid(NEW.bid_id),
            get_latest_flap_bid_guy(NEW.bid_id),
            get_latest_flap_bid_end(NEW.bid_id),
            get_latest_flap_bid_lot(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flap_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (address_id, bid_id, block_number) DO UPDATE SET tic = NEW.tic;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flap_end() RETURNS TRIGGER
AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flap (bid_id, address_id, block_number, "end", bid, guy, tic, lot, updated, created)
    VALUES (NEW.bid_id, NEW.address_id, (SELECT block_number FROM diff_block), NEW."end",
            get_latest_flap_bid_bid(NEW.bid_id),
            get_latest_flap_bid_guy(NEW.bid_id),
            get_latest_flap_bid_tic(NEW.bid_id),
            get_latest_flap_bid_lot(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flap_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (address_id, bid_id, block_number) DO UPDATE SET "end" = NEW."end";
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flap_lot() RETURNS TRIGGER
AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flap (bid_id, address_id, block_number, lot, bid, guy, tic, "end", updated, created)
    VALUES (NEW.bid_id, NEW.address_id, (SELECT block_number FROM diff_block), NEW.lot,
            get_latest_flap_bid_bid(NEW.bid_id),
            get_latest_flap_bid_guy(NEW.bid_id),
            get_latest_flap_bid_tic(NEW.bid_id),
            get_latest_flap_bid_end(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flap_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (address_id, bid_id, block_number) DO UPDATE SET lot = NEW.lot;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flap_created() RETURNS TRIGGER
AS
$$
DECLARE
    diff_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE headers.id = NEW.header_id);
BEGIN
    UPDATE maker.flap
    SET created = diff_timestamp
    WHERE flap.address_id = NEW.address_id
      AND flap.bid_id = NEW.bid_id
      AND flap.created IS NULL;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_bid_bid
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_bid
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flap_bid();

CREATE TRIGGER flap_bid_guy
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_guy
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flap_guy();

CREATE TRIGGER flap_bid_tic
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_tic
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flap_tic();

CREATE TRIGGER flap_bid_end
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_end
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flap_end();

CREATE TRIGGER flap_bid_lot
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_lot
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flap_lot();

CREATE TRIGGER flap_created_trigger
    AFTER INSERT
    ON maker.flap_kick
    FOR EACH ROW
EXECUTE PROCEDURE maker.flap_created();

-- +goose Down
DROP TRIGGER flap_bid_bid ON maker.flap_bid_bid;
DROP TRIGGER flap_bid_guy ON maker.flap_bid_guy;
DROP TRIGGER flap_bid_tic ON maker.flap_bid_tic;
DROP TRIGGER flap_bid_end ON maker.flap_bid_end;
DROP TRIGGER flap_bid_lot ON maker.flap_bid_lot;
DROP TRIGGER flap_created_trigger ON maker.flap_kick;

DROP FUNCTION maker.insert_updated_flap_bid();
DROP FUNCTION maker.insert_updated_flap_guy();
DROP FUNCTION maker.insert_updated_flap_tic();
DROP FUNCTION maker.insert_updated_flap_end();
DROP FUNCTION maker.insert_updated_flap_lot();
DROP FUNCTION maker.flap_created();
DROP FUNCTION get_latest_flap_bid_guy(numeric);
DROP FUNCTION get_latest_flap_bid_bid(numeric);
DROP FUNCTION get_latest_flap_bid_tic(numeric);
DROP FUNCTION get_latest_flap_bid_end(numeric);
DROP FUNCTION get_latest_flap_bid_lot(numeric);
DROP FUNCTION flap_bid_time_created(INTEGER, NUMERIC);

DROP INDEX maker.flap_address_index;
DROP TABLE maker.flap;