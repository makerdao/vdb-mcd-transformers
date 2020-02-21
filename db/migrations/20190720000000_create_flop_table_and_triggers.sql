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
    PRIMARY KEY (address_id, bid_id, block_number)
);

CREATE INDEX flop_address_index
    ON maker.flop (address_id);

CREATE FUNCTION flop_bid_guy_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS TEXT AS
$$
SELECT guy
FROM maker.flop_bid_guy
         LEFT JOIN public.headers ON flop_bid_guy.header_id = headers.id
WHERE flop_bid_guy.bid_id = flop_bid_guy_before_block.bid_id
  AND flop_bid_guy.address_id = flop_bid_guy_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flop_bid_guy_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flop_bid_guy_before_block
    IS E'@omit';

CREATE FUNCTION flop_bid_tic_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS BIGINT AS
$$
SELECT tic
FROM maker.flop_bid_tic
         LEFT JOIN public.headers ON flop_bid_tic.header_id = headers.id
WHERE flop_bid_tic.bid_id = flop_bid_tic_before_block.bid_id
  AND flop_bid_tic.address_id = flop_bid_tic_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flop_bid_tic_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flop_bid_tic_before_block
    IS E'@omit';

CREATE FUNCTION flop_bid_end_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS BIGINT AS
$$
SELECT "end"
FROM maker.flop_bid_end
         LEFT JOIN public.headers ON flop_bid_end.header_id = headers.id
WHERE flop_bid_end.bid_id = flop_bid_end_before_block.bid_id
  AND flop_bid_end.address_id = flop_bid_end_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flop_bid_end_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flop_bid_end_before_block
    IS E'@omit';

CREATE FUNCTION flop_bid_lot_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT lot
FROM maker.flop_bid_lot
         LEFT JOIN public.headers ON flop_bid_lot.header_id = headers.id
WHERE flop_bid_lot.bid_id = flop_bid_lot_before_block.bid_id
  AND flop_bid_lot.address_id = flop_bid_lot_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flop_bid_lot_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flop_bid_lot_before_block
    IS E'@omit';

CREATE FUNCTION flop_bid_bid_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT bid
FROM maker.flop_bid_bid
         LEFT JOIN public.headers ON flop_bid_bid.header_id = headers.id
WHERE flop_bid_bid.bid_id = flop_bid_bid_before_block.bid_id
  AND flop_bid_bid.address_id = flop_bid_bid_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flop_bid_bid_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flop_bid_bid_before_block
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
CREATE OR REPLACE FUNCTION maker.delete_obsolete_flop(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS VOID AS
$$
DECLARE
    flop_block      BIGINT     := (
        SELECT block_number
        FROM public.headers
        WHERE id = header_id);
    flop_state      maker.flop := (
        SELECT (flop.address_id, block_number, flop.bid_id, guy, tic, "end", lot, bid, created, updated)
        FROM maker.flop
        WHERE flop.bid_id = delete_obsolete_flop.bid_id
          AND flop.address_id = delete_obsolete_flop.address_id
          AND flop.block_number = flop_block);
    prev_flop_state maker.flop := (
        SELECT (flop.address_id, block_number, flop.bid_id, guy, tic, "end", lot, bid, created, updated)
        FROM maker.flop
        WHERE flop.bid_id = delete_obsolete_flop.bid_id
          AND flop.address_id = delete_obsolete_flop.address_id
          AND flop.block_number < flop_block
        ORDER BY flop.block_number DESC
        LIMIT 1);
BEGIN
    DELETE
    FROM maker.flop
    WHERE flop.bid_id = delete_obsolete_flop.bid_id
      AND flop.address_id = delete_obsolete_flop.address_id
      AND flop.block_number = flop_block
      AND (flop_state.guy IS NULL OR flop_state.guy = prev_flop_state.guy)
      AND (flop_state.tic IS NULL OR flop_state.tic = prev_flop_state.tic)
      AND (flop_state."end" IS NULL OR flop_state."end" = prev_flop_state."end")
      AND (flop_state.lot IS NULL OR flop_state.lot = prev_flop_state.lot)
      AND (flop_state.bid IS NULL OR flop_state.bid = prev_flop_state.bid);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.delete_obsolete_flop
    IS E'@omit';


CREATE OR REPLACE FUNCTION maker.insert_new_flop_guy(new_diff maker.flop_bid_guy) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flop (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        new_diff.guy,
        flop_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flop_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET guy = new_diff.guy
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_flop_guy(new_diff maker.flop_bid_guy)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flop_guys_until_next_diff(start_at_diff maker.flop_bid_guy, new_guy TEXT) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_guy_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flop_bid_guy
                 LEFT JOIN public.headers ON flop_bid_guy.header_id = headers.id
        WHERE flop_bid_guy.bid_id = start_at_diff.bid_id
          AND flop_bid_guy.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flop
    SET guy = new_guy
    WHERE flop.bid_id = start_at_diff.bid_id
      AND flop.address_id = start_at_diff.address_id
      AND flop.block_number >= diff_block_number
      AND (next_guy_diff_block IS NULL
        OR flop.block_number < next_guy_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


COMMENT ON FUNCTION maker.update_flop_guys_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flop_guys() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flop_guy(NEW);
        PERFORM maker.update_flop_guys_until_next_diff(NEW, NEW.guy);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flop_guys_until_next_diff(
                OLD,
                flop_bid_guy_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flop(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flop_guy
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.flop_bid_guy
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flop_guys();

CREATE OR REPLACE FUNCTION maker.insert_new_flop_tic(new_diff maker.flop_bid_tic) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flop (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flop_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.tic,
        flop_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flop_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET tic = new_diff.tic
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_flop_tic(new_diff maker.flop_bid_tic)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flop_tics_until_next_diff(start_at_diff maker.flop_bid_tic, new_tic NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_tic_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flop_bid_tic
                 LEFT JOIN public.headers ON flop_bid_tic.header_id = headers.id
        WHERE flop_bid_tic.bid_id = start_at_diff.bid_id
          AND flop_bid_tic.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flop
    SET tic = new_tic
    WHERE flop.bid_id = start_at_diff.bid_id
      AND flop.address_id = start_at_diff.address_id
      AND flop.block_number >= diff_block_number
      AND (next_tic_diff_block IS NULL
        OR flop.block_number < next_tic_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flop_tics_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flop_tics() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flop_tic(NEW);
        PERFORM maker.update_flop_tics_until_next_diff(NEW, NEW.tic);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flop_tics_until_next_diff(
                OLD,
                flop_bid_tic_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flop(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flop_tic
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.flop_bid_tic
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flop_tics();

CREATE OR REPLACE FUNCTION maker.insert_new_flop_end(new_diff maker.flop_bid_end) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flop (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flop_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff."end",
        flop_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flop_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET "end" = new_diff."end"
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_flop_end(new_diff maker.flop_bid_end)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flop_ends_until_next_diff(start_at_diff maker.flop_bid_end, new_end NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_end_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flop_bid_end
                 LEFT JOIN public.headers ON flop_bid_end.header_id = headers.id
        WHERE flop_bid_end.bid_id = start_at_diff.bid_id
          AND flop_bid_end.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flop
    SET "end" = new_end
    WHERE flop.bid_id = start_at_diff.bid_id
      AND flop.address_id = start_at_diff.address_id
      AND flop.block_number >= diff_block_number
      AND (next_end_diff_block IS NULL
        OR flop.block_number < next_end_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flop_ends_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flop_ends() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flop_end(NEW);
        PERFORM maker.update_flop_ends_until_next_diff(NEW, NEW."end");
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flop_ends_until_next_diff(
                OLD,
                flop_bid_end_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flop(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flop_end
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.flop_bid_end
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flop_ends();

CREATE OR REPLACE FUNCTION maker.insert_new_flop_lot(new_diff maker.flop_bid_lot) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flop (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flop_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.lot,
        flop_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flop_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET lot = new_diff.lot
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_flop_lot(new_diff maker.flop_bid_lot)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flop_lots_until_next_diff(start_at_diff maker.flop_bid_lot, new_lot NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_lot_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flop_bid_lot
                 LEFT JOIN public.headers ON flop_bid_lot.header_id = headers.id
        WHERE flop_bid_lot.bid_id = start_at_diff.bid_id
          AND flop_bid_lot.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flop
    SET lot = new_lot
    WHERE flop.bid_id = start_at_diff.bid_id
      AND flop.address_id = start_at_diff.address_id
      AND flop.block_number >= diff_block_number
      AND (next_lot_diff_block IS NULL
        OR flop.block_number < next_lot_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flop_lots_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flop_lots() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flop_lot(NEW);
        PERFORM maker.update_flop_lots_until_next_diff(NEW, NEW.lot);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flop_lots_until_next_diff(
                OLD,
                flop_bid_lot_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flop(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flop_lot
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.flop_bid_lot
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flop_lots();

CREATE OR REPLACE FUNCTION maker.insert_new_flop_bid(new_diff maker.flop_bid_bid) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flop (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flop_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.bid,
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flop_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET bid = new_diff.bid
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_flop_bid(new_diff maker.flop_bid_bid)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flop_bids_until_next_diff(start_at_diff maker.flop_bid_bid, new_bid NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_bid_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flop_bid_bid
                 LEFT JOIN public.headers ON flop_bid_bid.header_id = headers.id
        WHERE flop_bid_bid.bid_id = start_at_diff.bid_id
          AND flop_bid_bid.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flop
    SET bid = new_bid
    WHERE flop.bid_id = start_at_diff.bid_id
      AND flop.address_id = start_at_diff.address_id
      AND flop.block_number >= diff_block_number
      AND (next_bid_diff_block IS NULL
        OR flop.block_number < next_bid_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flop_bids_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flop_bids() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flop_bid(NEW);
        PERFORM maker.update_flop_bids_until_next_diff(NEW, NEW.bid);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flop_bids_until_next_diff(
                OLD,
                flop_bid_bid_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flop(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flop_bid
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.flop_bid_bid
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flop_bids();

CREATE OR REPLACE FUNCTION maker.insert_flop_created(new_event maker.flop_kick) RETURNS VOID
AS
$$
UPDATE maker.flop
SET created = api.epoch_to_datetime(headers.block_timestamp)
FROM public.headers
WHERE headers.id = new_event.header_id
  AND flop.address_id = new_event.address_id
  AND flop.bid_id = new_event.bid_id
  AND flop.created IS NULL
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_flop_created
    IS E'@omit';

CREATE OR REPLACE FUNCTION maker.clear_flop_created(old_event maker.flop_kick) RETURNS VOID
AS
$$
UPDATE maker.flop
SET created = flop_bid_time_created(old_event.address_id, old_event.bid_id)
WHERE flop.address_id = old_event.address_id
  AND flop.bid_id = old_event.bid_id
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.clear_flop_created
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flop_created() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        PERFORM maker.insert_flop_created(NEW);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.clear_flop_created(OLD);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flop_created
    IS E'@omit';

CREATE TRIGGER flop_created_trigger
    AFTER INSERT OR DELETE
    ON maker.flop_kick
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flop_created();

-- +goose Down
DROP TRIGGER flop_guy ON maker.flop_bid_guy;
DROP TRIGGER flop_tic ON maker.flop_bid_tic;
DROP TRIGGER flop_end ON maker.flop_bid_end;
DROP TRIGGER flop_lot ON maker.flop_bid_lot;
DROP TRIGGER flop_bid ON maker.flop_bid_bid;
DROP TRIGGER flop_created_trigger ON maker.flop_kick;

DROP FUNCTION maker.insert_new_flop_guy(maker.flop_bid_guy);
DROP FUNCTION maker.insert_new_flop_tic(maker.flop_bid_tic);
DROP FUNCTION maker.insert_new_flop_end(maker.flop_bid_end);
DROP FUNCTION maker.insert_new_flop_lot(maker.flop_bid_lot);
DROP FUNCTION maker.insert_new_flop_bid(maker.flop_bid_bid);
DROP FUNCTION maker.insert_flop_created(maker.flop_kick);
DROP FUNCTION maker.update_flop_guys_until_next_diff(maker.flop_bid_guy, TEXT);
DROP FUNCTION maker.update_flop_tics_until_next_diff(maker.flop_bid_tic, NUMERIC);
DROP FUNCTION maker.update_flop_ends_until_next_diff(maker.flop_bid_end, NUMERIC);
DROP FUNCTION maker.update_flop_lots_until_next_diff(maker.flop_bid_lot, NUMERIC);
DROP FUNCTION maker.update_flop_bids_until_next_diff(maker.flop_bid_bid, NUMERIC);
DROP FUNCTION maker.clear_flop_created(maker.flop_kick);
DROP FUNCTION maker.update_flop_guys();
DROP FUNCTION maker.update_flop_tics();
DROP FUNCTION maker.update_flop_ends();
DROP FUNCTION maker.update_flop_lots();
DROP FUNCTION maker.update_flop_bids();
DROP FUNCTION maker.update_flop_created();
DROP FUNCTION maker.delete_obsolete_flop(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flop_bid_guy_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flop_bid_tic_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flop_bid_end_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flop_bid_lot_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flop_bid_bid_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flop_bid_time_created(INTEGER, NUMERIC);

DROP INDEX maker.flop_address_index;
DROP TABLE maker.flop;
