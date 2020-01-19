-- +goose Up
CREATE TABLE maker.flip
(
    address_id   INTEGER   NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    block_number BIGINT    NOT NULL,
    bid_id       NUMERIC   NOT NULL,
    guy          TEXT      DEFAULT NULL,
    tic          BIGINT    DEFAULT NULL,
    "end"        BIGINT    DEFAULT NULL,
    lot          NUMERIC   DEFAULT NULL,
    bid          NUMERIC   DEFAULT NULL,
    gal          TEXT      DEFAULT NULL,
    tab          NUMERIC   DEFAULT NULL,
    created      TIMESTAMP DEFAULT NULL,
    updated      TIMESTAMP NOT NULL,
    PRIMARY KEY (block_number, bid_id, address_id)
);

CREATE INDEX flip_address_index
    ON maker.flip (address_id);

COMMENT ON TABLE maker.flip
    IS E'@name historicalFlipState';

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

COMMENT ON FUNCTION get_latest_flip_bid_guy
    IS E'@omit';

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

COMMENT ON FUNCTION get_latest_flip_bid_tic
    IS E'@omit';

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

COMMENT ON FUNCTION get_latest_flip_bid_end
    IS E'@omit';

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

COMMENT ON FUNCTION get_latest_flip_bid_lot
    IS E'@omit';

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

COMMENT ON FUNCTION get_latest_flip_bid_bid
    IS E'@omit';

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

COMMENT ON FUNCTION get_latest_flip_bid_gal
    IS E'@omit';

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

COMMENT ON FUNCTION get_latest_flip_bid_tab
    IS E'@omit';

CREATE FUNCTION flip_bid_time_created(address_id INTEGER, bid_id NUMERIC) RETURNS TIMESTAMP AS
$$
SELECT api.epoch_to_datetime(MIN(block_timestamp))
FROM public.headers
         LEFT JOIN maker.flip_kick ON flip_kick.header_id = headers.id
WHERE flip_kick.address_id = flip_bid_time_created.address_id
  AND flip_kick.bid_id = flip_bid_time_created.bid_id
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flip_bid_time_created
    IS E'@omit';


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_guy() RETURNS TRIGGER AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, hash, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, gal, tab, updated, created)
    VALUES (NEW.bid_id, NEW.address_id, (SELECT block_number FROM diff_block), NEW.guy,
            get_latest_flip_bid_tic(NEW.bid_id),
            get_latest_flip_bid_end(NEW.bid_id),
            get_latest_flip_bid_lot(NEW.bid_id),
            get_latest_flip_bid_bid(NEW.bid_id),
            get_latest_flip_bid_gal(NEW.bid_id),
            get_latest_flip_bid_tab(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flip_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET guy = NEW.guy;
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_tic() RETURNS TRIGGER AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, hash, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flip (bid_id, address_id, block_number, tic, guy, "end", lot, bid, gal, tab, updated, created)
    VALUES (NEW.bid_id, NEW.address_id, (SELECT block_number FROM diff_block), NEW.tic,
            get_latest_flip_bid_guy(NEW.bid_id),
            get_latest_flip_bid_end(NEW.bid_id),
            get_latest_flip_bid_lot(NEW.bid_id),
            get_latest_flip_bid_bid(NEW.bid_id),
            get_latest_flip_bid_gal(NEW.bid_id),
            get_latest_flip_bid_tab(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flip_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET tic = NEW.tic;
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_end() RETURNS TRIGGER AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, hash, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flip (bid_id, address_id, block_number, "end", guy, tic, lot, bid, gal, tab, updated, created)
    VALUES (NEW.bid_id, NEW.address_id, (SELECT block_number FROM diff_block), NEW."end",
            get_latest_flip_bid_guy(NEW.bid_id),
            get_latest_flip_bid_tic(NEW.bid_id),
            get_latest_flip_bid_lot(NEW.bid_id),
            get_latest_flip_bid_bid(NEW.bid_id),
            get_latest_flip_bid_gal(NEW.bid_id),
            get_latest_flip_bid_tab(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flip_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET "end" = NEW."end";
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_lot() RETURNS TRIGGER AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, hash, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flip (bid_id, address_id, block_number, lot, guy, tic, "end", bid, gal, tab, updated, created)
    VALUES (NEW.bid_id, NEW.address_id, (SELECT block_number FROM diff_block), NEW.lot,
            get_latest_flip_bid_guy(NEW.bid_id),
            get_latest_flip_bid_tic(NEW.bid_id),
            get_latest_flip_bid_end(NEW.bid_id),
            get_latest_flip_bid_bid(NEW.bid_id),
            get_latest_flip_bid_gal(NEW.bid_id),
            get_latest_flip_bid_tab(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flip_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET lot = NEW.lot;
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_bid() RETURNS TRIGGER AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, hash, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flip (bid_id, address_id, block_number, bid, guy, tic, "end", lot, gal, tab, updated, created)
    VALUES (NEW.bid_id, NEW.address_id, (SELECT block_number FROM diff_block), NEW.bid,
            get_latest_flip_bid_guy(NEW.bid_id),
            get_latest_flip_bid_tic(NEW.bid_id),
            get_latest_flip_bid_end(NEW.bid_id),
            get_latest_flip_bid_lot(NEW.bid_id),
            get_latest_flip_bid_gal(NEW.bid_id),
            get_latest_flip_bid_tab(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flip_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET bid = NEW.bid;
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_gal() RETURNS TRIGGER AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, hash, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flip (bid_id, address_id, block_number, gal, guy, tic, "end", lot, bid, tab, updated, created)
    VALUES (NEW.bid_id, NEW.address_id, (SELECT block_number FROM diff_block), NEW.gal,
            get_latest_flip_bid_guy(NEW.bid_id),
            get_latest_flip_bid_tic(NEW.bid_id),
            get_latest_flip_bid_end(NEW.bid_id),
            get_latest_flip_bid_lot(NEW.bid_id),
            get_latest_flip_bid_bid(NEW.bid_id),
            get_latest_flip_bid_tab(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flip_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET gal = NEW.gal;
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_updated_flip_tab() RETURNS TRIGGER AS
$$
BEGIN
    WITH diff_block AS (
        SELECT block_number, hash, block_timestamp
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO maker.flip (bid_id, address_id, block_number, tab, guy, tic, "end", lot, bid, gal, updated, created)
    VALUES (NEW.bid_id, NEW.address_id, (SELECT block_number FROM diff_block), NEW.tab,
            get_latest_flip_bid_guy(NEW.bid_id),
            get_latest_flip_bid_tic(NEW.bid_id),
            get_latest_flip_bid_end(NEW.bid_id),
            get_latest_flip_bid_lot(NEW.bid_id),
            get_latest_flip_bid_bid(NEW.bid_id),
            get_latest_flip_bid_gal(NEW.bid_id),
            (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
            flip_bid_time_created(NEW.address_id, NEW.bid_id))
    ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET tab = NEW.tab;
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.flip_created() RETURNS TRIGGER
AS
$$
DECLARE
    diff_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE headers.id = NEW.header_id);
BEGIN
    UPDATE maker.flip
    SET created = diff_timestamp
    WHERE flip.address_id = NEW.address_id
      AND flip.bid_id = NEW.bid_id
      AND flip.created IS NULL;
    RETURN NULL;
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
DROP FUNCTION get_latest_flip_bid_guy(NUMERIC);
DROP FUNCTION get_latest_flip_bid_bid(NUMERIC);
DROP FUNCTION get_latest_flip_bid_tic(NUMERIC);
DROP FUNCTION get_latest_flip_bid_end(NUMERIC);
DROP FUNCTION get_latest_flip_bid_lot(NUMERIC);
DROP FUNCTION get_latest_flip_bid_gal(NUMERIC);
DROP FUNCTION get_latest_flip_bid_tab(NUMERIC);
DROP FUNCTION flip_bid_time_created(INTEGER, NUMERIC);

DROP INDEX maker.flip_address_index;
DROP TABLE maker.flip;
