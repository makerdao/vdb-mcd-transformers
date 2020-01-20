-- +goose Up
CREATE TABLE maker.flop
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

CREATE INDEX flop_address_index
    ON maker.flop (address_id);

COMMENT ON TABLE maker.flop
    IS E'@name historicalFlopState';

CREATE FUNCTION get_latest_flop_bid_bid(bid_id NUMERIC) RETURNS NUMERIC AS
$$
SELECT bid
FROM maker.flop
WHERE bid IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION get_latest_flop_bid_bid
    IS E'@omit';


CREATE FUNCTION get_latest_flop_bid_guy(bid_id NUMERIC) RETURNS TEXT AS
$$
SELECT guy
FROM maker.flop
WHERE guy IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION get_latest_flop_bid_guy
    IS E'@omit';

CREATE FUNCTION get_latest_flop_bid_tic(bid_id NUMERIC) RETURNS BIGINT AS
$$
SELECT tic
FROM maker.flop
WHERE tic IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION get_latest_flop_bid_tic
    IS E'@omit';

CREATE FUNCTION get_latest_flop_bid_end(bid_id NUMERIC) RETURNS BIGINT AS
$$
SELECT "end"
FROM maker.flop
WHERE "end" IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION get_latest_flop_bid_end
    IS E'@omit';

CREATE FUNCTION get_latest_flop_bid_lot(bid_id NUMERIC) RETURNS NUMERIC AS
$$
SELECT lot
FROM maker.flop
WHERE lot IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION get_latest_flop_bid_lot
    IS E'@omit';

CREATE FUNCTION flop_bid_time_created(address_id INTEGER, bid_id NUMERIC) RETURNS TIMESTAMP AS
$$
SELECT api.epoch_to_datetime(MIN(block_timestamp))
FROM public.headers
         LEFT JOIN maker.flop_kick ON flop_kick.header_id = headers.id
WHERE flop_kick.address_id = flop_bid_time_created.address_id
  AND flop_kick.bid_id = flop_bid_time_created.bid_id
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flop_bid_time_created
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flop_bid() RETURNS TRIGGER
AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, hash, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flop (address_id, bid_id, block_number, bid, guy, tic, "end", lot, updated, created)
    VALUES (NEW.address_id, NEW.bid_id,
            (SELECT block_number FROM diff_block), NEW.bid,
            get_latest_flop_bid_guy(NEW.bid_id),
            get_latest_flop_bid_tic(NEW.bid_id),
            get_latest_flop_bid_end(NEW.bid_id),
            get_latest_flop_bid_lot(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flop_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (address_id, bid_id, block_number) DO UPDATE SET bid = NEW.bid;
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flop_guy() RETURNS TRIGGER
AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, hash, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flop (address_id, bid_id, block_number, guy, bid, tic, "end", lot, updated, created)
    VALUES (NEW.address_id, NEW.bid_id,
            (SELECT block_number FROM diff_block), NEW.guy,
            get_latest_flop_bid_bid(NEW.bid_id),
            get_latest_flop_bid_tic(NEW.bid_id),
            get_latest_flop_bid_end(NEW.bid_id),
            get_latest_flop_bid_lot(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flop_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (address_id, bid_id, block_number) DO UPDATE SET guy = NEW.guy;
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flop_tic() RETURNS TRIGGER
AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, hash, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flop (address_id, bid_id, block_number, tic, bid, guy, "end", lot, updated, created)
    VALUES (NEW.address_id, NEW.bid_id,
            (SELECT block_number FROM diff_block), NEW.tic,
            get_latest_flop_bid_bid(NEW.bid_id),
            get_latest_flop_bid_guy(NEW.bid_id),
            get_latest_flop_bid_end(NEW.bid_id),
            get_latest_flop_bid_lot(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flop_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (address_id, bid_id, block_number) DO UPDATE SET tic = NEW.tic;
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flop_end() RETURNS TRIGGER
AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, hash, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flop (address_id, bid_id, block_number, "end", bid, guy, tic, lot, updated, created)
    VALUES (NEW.address_id, NEW.bid_id,
            (SELECT block_number FROM diff_block), NEW."end",
            get_latest_flop_bid_bid(NEW.bid_id),
            get_latest_flop_bid_guy(NEW.bid_id),
            get_latest_flop_bid_tic(NEW.bid_id),
            get_latest_flop_bid_lot(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flop_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (address_id, bid_id, block_number) DO UPDATE SET "end" = NEW."end";
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flop_lot() RETURNS TRIGGER
AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, hash, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flop (address_id, bid_id, block_number, lot, bid, guy, tic, "end", updated, created)
    VALUES (NEW.address_id, NEW.bid_id,
            (SELECT block_number FROM diff_block), NEW.lot,
            get_latest_flop_bid_bid(NEW.bid_id),
            get_latest_flop_bid_guy(NEW.bid_id),
            get_latest_flop_bid_tic(NEW.bid_id),
            get_latest_flop_bid_end(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flop_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (address_id, bid_id, block_number) DO UPDATE SET lot = NEW.lot;
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flop_created() RETURNS TRIGGER
AS
$$
DECLARE
    diff_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE headers.id = NEW.header_id);
BEGIN
    UPDATE maker.flop
    SET created = diff_timestamp
    WHERE flop.address_id = NEW.address_id
      AND flop.bid_id = NEW.bid_id
      AND flop.created IS NULL;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


CREATE TRIGGER flop_bid_bid
    AFTER INSERT OR UPDATE
    ON maker.flop_bid_bid
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flop_bid();

CREATE TRIGGER flop_bid_guy
    AFTER INSERT OR UPDATE
    ON maker.flop_bid_guy
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flop_guy();

CREATE TRIGGER flop_bid_tic
    AFTER INSERT OR UPDATE
    ON maker.flop_bid_tic
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flop_tic();

CREATE TRIGGER flop_bid_end
    AFTER INSERT OR UPDATE
    ON maker.flop_bid_end
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flop_end();

CREATE TRIGGER flop_bid_lot
    AFTER INSERT OR UPDATE
    ON maker.flop_bid_lot
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flop_lot();

CREATE TRIGGER flop_created_trigger
    AFTER INSERT
    ON maker.flop_kick
    FOR EACH ROW
EXECUTE PROCEDURE maker.flop_created();

-- +goose Down
DROP TRIGGER flop_bid_bid ON maker.flop_bid_bid;
DROP TRIGGER flop_bid_guy ON maker.flop_bid_guy;
DROP TRIGGER flop_bid_tic ON maker.flop_bid_tic;
DROP TRIGGER flop_bid_end ON maker.flop_bid_end;
DROP TRIGGER flop_bid_lot ON maker.flop_bid_lot;
DROP TRIGGER flop_created_trigger ON maker.flop_kick;

DROP FUNCTION maker.insert_updated_flop_bid();
DROP FUNCTION maker.insert_updated_flop_guy();
DROP FUNCTION maker.insert_updated_flop_tic();
DROP FUNCTION maker.insert_updated_flop_end();
DROP FUNCTION maker.insert_updated_flop_lot();
DROP FUNCTION maker.flop_created();
DROP FUNCTION get_latest_flop_bid_guy(NUMERIC);
DROP FUNCTION get_latest_flop_bid_bid(NUMERIC);
DROP FUNCTION get_latest_flop_bid_tic(NUMERIC);
DROP FUNCTION get_latest_flop_bid_end(NUMERIC);
DROP FUNCTION get_latest_flop_bid_lot(NUMERIC);
DROP FUNCTION flop_bid_time_created(INTEGER, NUMERIC);

DROP INDEX maker.flop_address_index;
DROP TABLE maker.flop;