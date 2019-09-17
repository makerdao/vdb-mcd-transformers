-- +goose Up
CREATE TABLE maker.flip
(
    id               SERIAL PRIMARY KEY,
    block_number     BIGINT  DEFAULT NULL,
    block_hash       TEXT    DEFAULT NULL,
    address_id       INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id           NUMERIC DEFAULT NULL,
    guy              TEXT    DEFAULT NULL,
    tic              BIGINT  DEFAULT NULL,
    "end"            BIGINT  DEFAULT NULL,
    lot              NUMERIC DEFAULT NULL,
    bid              NUMERIC DEFAULT NULL,
    gal              TEXT    DEFAULT NULL,
    tab              NUMERIC DEFAULT NULL,
    created          TIMESTAMP,
    updated          TIMESTAMP,
    UNIQUE (block_number, bid_id)
);

CREATE FUNCTION get_latest_flip_bid_guy(bid_id numeric) RETURNS TEXT AS
$$
SELECT guy
FROM maker.flip
WHERE guy IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

CREATE FUNCTION get_latest_flip_bid_tic(bid_id numeric) RETURNS BIGINT AS
$$
SELECT tic
FROM maker.flip
WHERE tic IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

CREATE FUNCTION get_latest_flip_bid_end(bid_id numeric) RETURNS BIGINT AS
$$
SELECT "end"
FROM maker.flip
WHERE "end" IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

CREATE FUNCTION get_latest_flip_bid_lot(bid_id numeric) RETURNS NUMERIC AS
$$
SELECT lot
FROM maker.flip
WHERE lot IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

CREATE FUNCTION get_latest_flip_bid_bid(bid_id numeric) RETURNS NUMERIC AS
$$
SELECT bid
FROM maker.flip
WHERE bid IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

CREATE FUNCTION get_latest_flip_bid_gal(bid_id numeric) RETURNS TEXT AS
$$
SELECT gal
FROM maker.flip
WHERE gal IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

CREATE FUNCTION get_latest_flip_bid_tab(bid_id numeric) RETURNS NUMERIC AS
$$
SELECT tab
FROM maker.flip
WHERE tab IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

CREATE FUNCTION get_block_timestamp(block_hash varchar) RETURNS TIMESTAMP AS
$$
SELECT api.epoch_to_datetime(headers.block_timestamp) AS datetime
FROM public.headers
WHERE headers.hash = block_hash
ORDER BY headers.block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_guy() RETURNS TRIGGER AS
$$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, guy, tic, "end", lot, bid, gal, tab, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW.guy,
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET guy = NEW.guy;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_tic() RETURNS TRIGGER AS
$$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, tic, guy, "end", lot, bid, gal, tab, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW.tic,
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET tic = NEW.tic;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_end() RETURNS TRIGGER AS
$$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, "end", guy, tic, lot, bid, gal, tab, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW."end",
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET "end" = NEW."end";
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_lot() RETURNS TRIGGER AS
$$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, lot, guy, tic, "end", bid, gal, tab, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW.lot,
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET lot = NEW.lot;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_bid() RETURNS TRIGGER AS
$$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, bid, guy, tic, "end", lot, gal, tab, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW.bid,
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET bid = NEW.bid;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_gal() RETURNS TRIGGER AS
$$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, gal, guy, tic, "end", lot, bid, tab, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW.gal,
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET gal = NEW.gal;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_tab() RETURNS TRIGGER AS
$$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, tab, guy, tic, "end", lot, bid, gal, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW.tab,
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET tab = NEW.tab;
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flip_created() RETURNS TRIGGER
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
    INTO maker.flip(bid_id, address_id, block_number, block_hash, created, updated, guy, tic, "end", lot, bid,
                    gal, tab)
    VALUES (NEW.bid_id, NEW.address_id,
            (SELECT block_number FROM block_info),
            (SELECT hash FROM block_info),
            (SELECT datetime FROM block_info),
            (SELECT datetime FROM block_info),
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET created = (SELECT datetime FROM block_info),
                                                     updated = (SELECT datetime FROM block_info);
    return NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_bid_guy
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_guy
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flip_guy();

CREATE TRIGGER flip_bid_tic
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_tic
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flip_tic();

CREATE TRIGGER flip_bid_end
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_end
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flip_end();

CREATE TRIGGER flip_bid_lot
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_lot
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flip_lot();

CREATE TRIGGER flip_bid_bid
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_bid
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flip_bid();

CREATE TRIGGER flip_bid_gal
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_gal
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flip_gal();

CREATE TRIGGER flip_bid_tab
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_tab
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_updated_flip_tab();

CREATE TRIGGER flip_created_trigger
    AFTER INSERT
    ON maker.flip_kick
    FOR EACH ROW
EXECUTE PROCEDURE maker.flip_created();


-- +goose Down
DROP TRIGGER flip_bid_guy ON maker.flip_bid_guy;
DROP TRIGGER flip_bid_tic ON maker.flip_bid_tic;
DROP TRIGGER flip_bid_end ON maker.flip_bid_end;
DROP TRIGGER flip_bid_lot ON maker.flip_bid_lot;
DROP TRIGGER flip_bid_bid ON maker.flip_bid_bid;
DROP TRIGGER flip_bid_gal ON maker.flip_bid_gal;
DROP TRIGGER flip_bid_tab ON maker.flip_bid_tab;
DROP TRIGGER flip_created_trigger ON maker.flip_kick;

DROP FUNCTION maker.insert_updated_flip_guy();
DROP FUNCTION maker.insert_updated_flip_tic();
DROP FUNCTION maker.insert_updated_flip_end();
DROP FUNCTION maker.insert_updated_flip_lot();
DROP FUNCTION maker.insert_updated_flip_bid();
DROP FUNCTION maker.insert_updated_flip_gal();
DROP FUNCTION maker.insert_updated_flip_tab();
DROP FUNCTION maker.flip_created();
DROP FUNCTION get_latest_flip_bid_guy(numeric);
DROP FUNCTION get_latest_flip_bid_bid(numeric);
DROP FUNCTION get_latest_flip_bid_tic(numeric);
DROP FUNCTION get_latest_flip_bid_end(numeric);
DROP FUNCTION get_latest_flip_bid_lot(numeric);
DROP FUNCTION get_latest_flip_bid_gal(numeric);
DROP FUNCTION get_latest_flip_bid_tab(numeric);
DROP FUNCTION get_block_timestamp(varchar);
DROP TABLE maker.flip;
