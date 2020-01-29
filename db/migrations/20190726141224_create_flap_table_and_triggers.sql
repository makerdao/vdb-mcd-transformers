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

CREATE FUNCTION flap_bid_guy_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS TEXT AS
$$
SELECT guy
FROM maker.flap_bid_guy
         LEFT JOIN public.headers ON flap_bid_guy.header_id = headers.id
WHERE flap_bid_guy.bid_id = flap_bid_guy_before_block.bid_id
  AND flap_bid_guy.address_id = flap_bid_guy_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flap_bid_guy_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flap_bid_guy_before_block
    IS E'@omit';

CREATE FUNCTION flap_bid_tic_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS BIGINT AS
$$
SELECT tic
FROM maker.flap_bid_tic
         LEFT JOIN public.headers ON flap_bid_tic.header_id = headers.id
WHERE flap_bid_tic.bid_id = flap_bid_tic_before_block.bid_id
  AND flap_bid_tic.address_id = flap_bid_tic_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flap_bid_tic_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flap_bid_tic_before_block
    IS E'@omit';

CREATE FUNCTION flap_bid_end_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS BIGINT AS
$$
SELECT "end"
FROM maker.flap_bid_end
         LEFT JOIN public.headers ON flap_bid_end.header_id = headers.id
WHERE flap_bid_end.bid_id = flap_bid_end_before_block.bid_id
  AND flap_bid_end.address_id = flap_bid_end_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flap_bid_end_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flap_bid_end_before_block
    IS E'@omit';

CREATE FUNCTION flap_bid_lot_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT lot
FROM maker.flap_bid_lot
         LEFT JOIN public.headers ON flap_bid_lot.header_id = headers.id
WHERE flap_bid_lot.bid_id = flap_bid_lot_before_block.bid_id
  AND flap_bid_lot.address_id = flap_bid_lot_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flap_bid_lot_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flap_bid_lot_before_block
    IS E'@omit';

CREATE FUNCTION flap_bid_bid_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT bid
FROM maker.flap_bid_bid
         LEFT JOIN public.headers ON flap_bid_bid.header_id = headers.id
WHERE flap_bid_bid.bid_id = flap_bid_bid_before_block.bid_id
  AND flap_bid_bid.address_id = flap_bid_bid_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flap_bid_bid_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flap_bid_bid_before_block
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
CREATE OR REPLACE FUNCTION maker.insert_new_flap_guy(new_diff maker.flap_bid_guy) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flap (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        new_diff.guy,
        flap_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flap_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET guy = new_diff.guy;
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_guys_until_next_diff(new_diff maker.flap_bid_guy) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_guy_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_guy
                 LEFT JOIN public.headers ON flap_bid_guy.header_id = headers.id
        WHERE flap_bid_guy.bid_id = new_diff.bid_id
          AND flap_bid_guy.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET guy = new_diff.guy
    WHERE flap.bid_id = new_diff.bid_id
      AND flap.address_id = new_diff.address_id
      AND flap.block_number >= diff_block_number
      AND (next_guy_diff_block IS NULL
        OR flap.block_number < next_guy_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flap_guys_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_guys() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flap_guy(NEW);
    PERFORM maker.update_flap_guys_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_guy
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_guy
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flap_guys();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flap_tic(new_diff maker.flap_bid_tic) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flap (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flap_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.tic,
        flap_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flap_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET tic = new_diff.tic;
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_tics_until_next_diff(new_diff maker.flap_bid_tic) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_tic_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_tic
                 LEFT JOIN public.headers ON flap_bid_tic.header_id = headers.id
        WHERE flap_bid_tic.bid_id = new_diff.bid_id
          AND flap_bid_tic.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET tic = new_diff.tic
    WHERE flap.bid_id = new_diff.bid_id
      AND flap.address_id = new_diff.address_id
      AND flap.block_number >= diff_block_number
      AND (next_tic_diff_block IS NULL
        OR flap.block_number < next_tic_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flap_tics_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_tics() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flap_tic(NEW);
    PERFORM maker.update_flap_tics_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_tic
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_tic
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flap_tics();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flap_end(new_diff maker.flap_bid_end) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flap (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flap_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff."end",
        flap_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flap_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET "end" = new_diff."end";
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_ends_until_next_diff(new_diff maker.flap_bid_end) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_end_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_end
                 LEFT JOIN public.headers ON flap_bid_end.header_id = headers.id
        WHERE flap_bid_end.bid_id = new_diff.bid_id
          AND flap_bid_end.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET "end" = new_diff."end"
    WHERE flap.bid_id = new_diff.bid_id
      AND flap.address_id = new_diff.address_id
      AND flap.block_number >= diff_block_number
      AND (next_end_diff_block IS NULL
        OR flap.block_number < next_end_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flap_ends_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_ends() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flap_end(NEW);
    PERFORM maker.update_flap_ends_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_end
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_end
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flap_ends();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flap_lot(new_diff maker.flap_bid_lot) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flap (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flap_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.lot,
        flap_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flap_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET lot = new_diff.lot;
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_lots_until_next_diff(new_diff maker.flap_bid_lot) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_lot_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_lot
                 LEFT JOIN public.headers ON flap_bid_lot.header_id = headers.id
        WHERE flap_bid_lot.bid_id = new_diff.bid_id
          AND flap_bid_lot.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET lot = new_diff.lot
    WHERE flap.bid_id = new_diff.bid_id
      AND flap.address_id = new_diff.address_id
      AND flap.block_number >= diff_block_number
      AND (next_lot_diff_block IS NULL
        OR flap.block_number < next_lot_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flap_lots_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_lots() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flap_lot(NEW);
    PERFORM maker.update_flap_lots_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_lot
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_lot
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flap_lots();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flap_bid(new_diff maker.flap_bid_bid) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flap (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flap_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.bid,
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flap_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET bid = new_diff.bid;
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_bids_until_next_diff(new_diff maker.flap_bid_bid) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_bid_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_bid
                 LEFT JOIN public.headers ON flap_bid_bid.header_id = headers.id
        WHERE flap_bid_bid.bid_id = new_diff.bid_id
          AND flap_bid_bid.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET bid = new_diff.bid
    WHERE flap.bid_id = new_diff.bid_id
      AND flap.address_id = new_diff.address_id
      AND flap.block_number >= diff_block_number
      AND (next_bid_diff_block IS NULL
        OR flap.block_number < next_bid_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flap_bids_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_bids() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flap_bid(NEW);
    PERFORM maker.update_flap_bids_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_bid
    AFTER INSERT OR UPDATE
    ON maker.flap_bid_bid
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flap_bids();

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

CREATE TRIGGER flap_created_trigger
    AFTER INSERT
    ON maker.flap_kick
    FOR EACH ROW
EXECUTE PROCEDURE maker.flap_created();

-- +goose Down
DROP TRIGGER flap_guy ON maker.flap_bid_guy;
DROP TRIGGER flap_tic ON maker.flap_bid_tic;
DROP TRIGGER flap_end ON maker.flap_bid_end;
DROP TRIGGER flap_lot ON maker.flap_bid_lot;
DROP TRIGGER flap_bid ON maker.flap_bid_bid;
DROP TRIGGER flap_created_trigger ON maker.flap_kick;

DROP FUNCTION maker.insert_new_flap_guy(maker.flap_bid_guy);
DROP FUNCTION maker.insert_new_flap_tic(maker.flap_bid_tic);
DROP FUNCTION maker.insert_new_flap_end(maker.flap_bid_end);
DROP FUNCTION maker.insert_new_flap_lot(maker.flap_bid_lot);
DROP FUNCTION maker.insert_new_flap_bid(maker.flap_bid_bid);
DROP FUNCTION maker.update_flap_guys_until_next_diff(maker.flap_bid_guy);
DROP FUNCTION maker.update_flap_tics_until_next_diff(maker.flap_bid_tic);
DROP FUNCTION maker.update_flap_ends_until_next_diff(maker.flap_bid_end);
DROP FUNCTION maker.update_flap_lots_until_next_diff(maker.flap_bid_lot);
DROP FUNCTION maker.update_flap_bids_until_next_diff(maker.flap_bid_bid);
DROP FUNCTION maker.update_flap_guys();
DROP FUNCTION maker.update_flap_tics();
DROP FUNCTION maker.update_flap_ends();
DROP FUNCTION maker.update_flap_lots();
DROP FUNCTION maker.update_flap_bids();
DROP FUNCTION maker.flap_created();
DROP FUNCTION flap_bid_guy_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flap_bid_tic_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flap_bid_end_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flap_bid_lot_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flap_bid_bid_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flap_bid_time_created(INTEGER, NUMERIC);

DROP INDEX maker.flap_address_index;
DROP TABLE maker.flap;
