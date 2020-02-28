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
    PRIMARY KEY (address_id, bid_id, block_number)
);

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
CREATE OR REPLACE FUNCTION maker.delete_obsolete_flap(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS VOID AS
$$
DECLARE
    flap_block      BIGINT     := (
        SELECT block_number
        FROM public.headers
        WHERE id = header_id);
    flap_state      maker.flap := (
        SELECT (flap.address_id, block_number, flap.bid_id, guy, tic, "end", lot, bid, created, updated)
        FROM maker.flap
        WHERE flap.bid_id = delete_obsolete_flap.bid_id
          AND flap.address_id = delete_obsolete_flap.address_id
          AND flap.block_number = flap_block);
    prev_flap_state maker.flap := (
        SELECT (flap.address_id, block_number, flap.bid_id, guy, tic, "end", lot, bid, created, updated)
        FROM maker.flap
        WHERE flap.bid_id = delete_obsolete_flap.bid_id
          AND flap.address_id = delete_obsolete_flap.address_id
          AND flap.block_number < flap_block
        ORDER BY flap.block_number DESC
        LIMIT 1);
BEGIN
    DELETE
    FROM maker.flap
    WHERE flap.bid_id = delete_obsolete_flap.bid_id
      AND flap.address_id = delete_obsolete_flap.address_id
      AND flap.block_number = flap_block
      AND (flap_state.guy IS NULL OR flap_state.guy = prev_flap_state.guy)
      AND (flap_state.tic IS NULL OR flap_state.tic = prev_flap_state.tic)
      AND (flap_state."end" IS NULL OR flap_state."end" = prev_flap_state."end")
      AND (flap_state.lot IS NULL OR flap_state.lot = prev_flap_state.lot)
      AND (flap_state.bid IS NULL OR flap_state.bid = prev_flap_state.bid);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.delete_obsolete_flap
    IS E'@omit';


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
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET guy = new_diff.guy
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_flap_guy(new_diff maker.flap_bid_guy)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_guys_until_next_diff(start_at_diff maker.flap_bid_guy, new_guy TEXT) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_guy_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_guy
                 LEFT JOIN public.headers ON flap_bid_guy.header_id = headers.id
        WHERE flap_bid_guy.bid_id = start_at_diff.bid_id
          AND flap_bid_guy.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET guy = new_guy
    WHERE flap.bid_id = start_at_diff.bid_id
      AND flap.address_id = start_at_diff.address_id
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
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flap_guy(NEW);
        PERFORM maker.update_flap_guys_until_next_diff(NEW, NEW.guy);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flap_guys_until_next_diff(
                OLD,
                flap_bid_guy_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flap(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_guy
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.flap_bid_guy
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flap_guys();

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
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET tic = new_diff.tic
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_flap_tic(new_diff maker.flap_bid_tic)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_tics_until_next_diff(start_at_diff maker.flap_bid_tic, new_tic NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_tic_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_tic
                 LEFT JOIN public.headers ON flap_bid_tic.header_id = headers.id
        WHERE flap_bid_tic.bid_id = start_at_diff.bid_id
          AND flap_bid_tic.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET tic = new_tic
    WHERE flap.bid_id = start_at_diff.bid_id
      AND flap.address_id = start_at_diff.address_id
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
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flap_tic(NEW);
        PERFORM maker.update_flap_tics_until_next_diff(NEW, NEW.tic);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flap_tics_until_next_diff(
                OLD,
                flap_bid_tic_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flap(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_tic
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.flap_bid_tic
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flap_tics();

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
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET "end" = new_diff."end"
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_flap_end(new_diff maker.flap_bid_end)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_ends_until_next_diff(start_at_diff maker.flap_bid_end, new_end NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_end_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_end
                 LEFT JOIN public.headers ON flap_bid_end.header_id = headers.id
        WHERE flap_bid_end.bid_id = start_at_diff.bid_id
          AND flap_bid_end.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET "end" = new_end
    WHERE flap.bid_id = start_at_diff.bid_id
      AND flap.address_id = start_at_diff.address_id
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
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flap_end(NEW);
        PERFORM maker.update_flap_ends_until_next_diff(NEW, NEW."end");
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flap_ends_until_next_diff(
                OLD,
                flap_bid_end_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flap(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_end
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.flap_bid_end
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flap_ends();

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
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET lot = new_diff.lot
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_flap_lot(new_diff maker.flap_bid_lot)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_lots_until_next_diff(start_at_diff maker.flap_bid_lot, new_lot NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_lot_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_lot
                 LEFT JOIN public.headers ON flap_bid_lot.header_id = headers.id
        WHERE flap_bid_lot.bid_id = start_at_diff.bid_id
          AND flap_bid_lot.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET lot = new_lot
    WHERE flap.bid_id = start_at_diff.bid_id
      AND flap.address_id = start_at_diff.address_id
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
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flap_lot(NEW);
        PERFORM maker.update_flap_lots_until_next_diff(NEW, NEW.lot);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flap_lots_until_next_diff(
                OLD,
                flap_bid_lot_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flap(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_lot
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.flap_bid_lot
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flap_lots();

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
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET bid = new_diff.bid
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_flap_bid(new_diff maker.flap_bid_bid)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_bids_until_next_diff(start_at_diff maker.flap_bid_bid, new_bid NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_bid_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_bid
                 LEFT JOIN public.headers ON flap_bid_bid.header_id = headers.id
        WHERE flap_bid_bid.bid_id = start_at_diff.bid_id
          AND flap_bid_bid.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET bid = new_bid
    WHERE flap.bid_id = start_at_diff.bid_id
      AND flap.address_id = start_at_diff.address_id
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
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flap_bid(NEW);
        PERFORM maker.update_flap_bids_until_next_diff(NEW, NEW.bid);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flap_bids_until_next_diff(
                OLD,
                flap_bid_bid_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flap(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_bid
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.flap_bid_bid
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flap_bids();

CREATE OR REPLACE FUNCTION maker.insert_flap_created(new_event maker.flap_kick) RETURNS VOID
AS
$$
UPDATE maker.flap
SET created = api.epoch_to_datetime(headers.block_timestamp)
FROM public.headers
WHERE headers.id = new_event.header_id
  AND flap.address_id = new_event.address_id
  AND flap.bid_id = new_event.bid_id
  AND flap.created IS NULL
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_flap_created
    IS E'@omit';

CREATE OR REPLACE FUNCTION maker.clear_flap_created(old_event maker.flap_kick) RETURNS VOID
AS
$$
UPDATE maker.flap
SET created = flap_bid_time_created(old_event.address_id, old_event.bid_id)
WHERE flap.address_id = old_event.address_id
  AND flap.bid_id = old_event.bid_id
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.clear_flap_created
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flap_created() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        PERFORM maker.insert_flap_created(NEW);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.clear_flap_created(OLD);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flap_created
    IS E'@omit';

CREATE TRIGGER flap_created_trigger
    AFTER INSERT OR DELETE
    ON maker.flap_kick
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flap_created();

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
DROP FUNCTION maker.insert_flap_created(maker.flap_kick);
DROP FUNCTION maker.update_flap_guys_until_next_diff(maker.flap_bid_guy, TEXT);
DROP FUNCTION maker.update_flap_tics_until_next_diff(maker.flap_bid_tic, NUMERIC);
DROP FUNCTION maker.update_flap_ends_until_next_diff(maker.flap_bid_end, NUMERIC);
DROP FUNCTION maker.update_flap_lots_until_next_diff(maker.flap_bid_lot, NUMERIC);
DROP FUNCTION maker.update_flap_bids_until_next_diff(maker.flap_bid_bid, NUMERIC);
DROP FUNCTION maker.clear_flap_created(maker.flap_kick);
DROP FUNCTION maker.update_flap_guys();
DROP FUNCTION maker.update_flap_tics();
DROP FUNCTION maker.update_flap_ends();
DROP FUNCTION maker.update_flap_lots();
DROP FUNCTION maker.update_flap_bids();
DROP FUNCTION maker.update_flap_created();
DROP FUNCTION maker.delete_obsolete_flap(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flap_bid_guy_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flap_bid_tic_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flap_bid_end_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flap_bid_lot_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flap_bid_bid_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flap_bid_time_created(INTEGER, NUMERIC);

DROP INDEX maker.flap_address_index;
DROP TABLE maker.flap;
