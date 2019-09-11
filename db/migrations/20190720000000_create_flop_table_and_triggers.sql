-- +goose Up
CREATE TABLE maker.flop
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT  DEFAULT NULL,
    block_hash   TEXT    DEFAULT NULL,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id       NUMERIC DEFAULT NULL,
    guy          TEXT    DEFAULT NULL,
    tic          BIGINT  DEFAULT NULL,
    "end"        BIGINT  DEFAULT NULL,
    lot          NUMERIC DEFAULT NULL,
    bid          NUMERIC DEFAULT NULL,
    created      TIMESTAMP,
    updated      TIMESTAMP,
    UNIQUE (block_number, bid_id)
);


CREATE FUNCTION get_latest_flop_bid_bid(bid_id numeric) RETURNS NUMERIC AS
$$
SELECT bid
FROM maker.flop
WHERE bid IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;


CREATE FUNCTION get_latest_flop_bid_guy(bid_id numeric) RETURNS TEXT AS
$$
SELECT guy
FROM maker.flop
WHERE guy IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

CREATE FUNCTION get_latest_flop_bid_tic(bid_id numeric) RETURNS BIGINT AS
$$
SELECT tic
FROM maker.flop
WHERE tic IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

CREATE FUNCTION get_latest_flop_bid_end(bid_id numeric) RETURNS BIGINT AS
$$
SELECT "end"
FROM maker.flop
WHERE "end" IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

CREATE FUNCTION get_latest_flop_bid_lot(bid_id numeric) RETURNS NUMERIC AS
$$
SELECT lot
FROM maker.flop
WHERE lot IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flop_bid() RETURNS TRIGGER
AS
$$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flop
        WHERE flop.bid_id = NEW.bid_id
        ORDER BY flop.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flop(bid_id, address_id, block_number, block_hash, bid, guy, tic, "end", lot, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW.bid,
            (SELECT get_latest_flop_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flop_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flop_bid_end(NEW.bid_id)),
            (SELECT get_latest_flop_bid_lot(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET bid = NEW.bid;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flop_guy() RETURNS TRIGGER
AS
$$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flop
        WHERE flop.bid_id = NEW.bid_id
        ORDER BY flop.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flop(bid_id, address_id, block_number, block_hash, guy, bid, tic, "end", lot, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW.guy,
            (SELECT get_latest_flop_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flop_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flop_bid_end(NEW.bid_id)),
            (SELECT get_latest_flop_bid_lot(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET guy = NEW.guy;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flop_tic() RETURNS TRIGGER
AS
$$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flop
        WHERE flop.bid_id = NEW.bid_id
        ORDER BY flop.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flop(bid_id, address_id, block_number, block_hash, tic, bid, guy, "end", lot, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW.tic,
            (SELECT get_latest_flop_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flop_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flop_bid_end(NEW.bid_id)),
            (SELECT get_latest_flop_bid_lot(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET tic = NEW.tic;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flop_end() RETURNS TRIGGER
AS
$$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flop
        WHERE flop.bid_id = NEW.bid_id
        ORDER BY flop.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flop(bid_id, address_id, block_number, block_hash, "end", bid, guy, tic, lot, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW."end",
            (SELECT get_latest_flop_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flop_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flop_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flop_bid_lot(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET "end" = NEW."end";
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flop_lot() RETURNS TRIGGER
AS
$$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flop
        WHERE flop.bid_id = NEW.bid_id
        ORDER BY flop.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flop(bid_id, address_id, block_number, block_hash, lot, bid, guy, tic, "end", updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW.lot,
            (SELECT get_latest_flop_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flop_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flop_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flop_bid_end(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET lot = NEW.lot;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flop_created() RETURNS TRIGGER
AS
$$
BEGIN
    WITH block_info AS (
        SELECT block_number, hash, api.epoch_to_datetime(headers.block_timestamp) AS datetime
        FROM public.headers
        WHERE headers.id = NEW.header_id
        LIMIT 1
    )
    INSERT
    INTO maker.flop(bid_id, address_id, block_number, block_hash, created, updated, bid, guy, tic, "end", lot)
    VALUES (NEW.bid_id, NEW.address_id,
            (SELECT block_number FROM block_info),
            (SELECT hash FROM block_info),
            (SELECT datetime FROM block_info),
            (SELECT datetime FROM block_info),
            (SELECT get_latest_flop_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flop_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flop_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flop_bid_end(NEW.bid_id)),
            (SELECT get_latest_flop_bid_lot(NEW.bid_id)))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET created = (SELECT datetime FROM block_info),
                                                     updated = (SELECT datetime FROM block_info);
    return NEW;
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
DROP FUNCTION get_latest_flop_bid_guy(numeric);
DROP FUNCTION get_latest_flop_bid_bid(numeric);
DROP FUNCTION get_latest_flop_bid_tic(numeric);
DROP FUNCTION get_latest_flop_bid_end(numeric);
DROP FUNCTION get_latest_flop_bid_lot(numeric);
DROP TABLE maker.flop;