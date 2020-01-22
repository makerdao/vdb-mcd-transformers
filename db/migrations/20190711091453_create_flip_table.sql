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
    usr          TEXT      DEFAULT NULL,
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

CREATE FUNCTION flip_bid_guy_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS TEXT AS
$$
SELECT guy
FROM maker.flip_bid_guy
         LEFT JOIN public.headers ON flip_bid_guy.header_id = headers.id
WHERE flip_bid_guy.bid_id = flip_bid_guy_before_block.bid_id
  AND flip_bid_guy.address_id = flip_bid_guy_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_guy_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flip_bid_guy_before_block
    IS E'@omit';

CREATE FUNCTION flip_bid_tic_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS BIGINT AS
$$
SELECT tic
FROM maker.flip_bid_tic
         LEFT JOIN public.headers ON flip_bid_tic.header_id = headers.id
WHERE flip_bid_tic.bid_id = flip_bid_tic_before_block.bid_id
  AND flip_bid_tic.address_id = flip_bid_tic_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_tic_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flip_bid_tic_before_block
    IS E'@omit';

CREATE FUNCTION flip_bid_end_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS BIGINT AS
$$
SELECT "end"
FROM maker.flip_bid_end
         LEFT JOIN public.headers ON flip_bid_end.header_id = headers.id
WHERE flip_bid_end.bid_id = flip_bid_end_before_block.bid_id
  AND flip_bid_end.address_id = flip_bid_end_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_end_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flip_bid_end_before_block
    IS E'@omit';

CREATE FUNCTION flip_bid_lot_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT lot
FROM maker.flip_bid_lot
         LEFT JOIN public.headers ON flip_bid_lot.header_id = headers.id
WHERE flip_bid_lot.bid_id = flip_bid_lot_before_block.bid_id
  AND flip_bid_lot.address_id = flip_bid_lot_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_lot_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flip_bid_lot_before_block
    IS E'@omit';

CREATE FUNCTION flip_bid_bid_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT bid
FROM maker.flip_bid_bid
         LEFT JOIN public.headers ON flip_bid_bid.header_id = headers.id
WHERE flip_bid_bid.bid_id = flip_bid_bid_before_block.bid_id
  AND flip_bid_bid.address_id = flip_bid_bid_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_bid_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flip_bid_bid_before_block
    IS E'@omit';

CREATE FUNCTION flip_bid_usr_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS TEXT AS
$$
SELECT usr
FROM maker.flip_bid_usr
         LEFT JOIN public.headers ON flip_bid_usr.header_id = headers.id
WHERE flip_bid_usr.bid_id = flip_bid_usr_before_block.bid_id
  AND flip_bid_usr.address_id = flip_bid_usr_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_usr_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flip_bid_usr_before_block
    IS E'@omit';

CREATE FUNCTION flip_bid_gal_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS TEXT AS
$$
SELECT gal
FROM maker.flip_bid_gal
         LEFT JOIN public.headers ON flip_bid_gal.header_id = headers.id
WHERE flip_bid_gal.bid_id = flip_bid_gal_before_block.bid_id
  AND flip_bid_gal.address_id = flip_bid_gal_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_gal_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flip_bid_gal_before_block
    IS E'@omit';

CREATE FUNCTION flip_bid_tab_before_block(bid_id NUMERIC, address_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT tab
FROM maker.flip_bid_tab
         LEFT JOIN public.headers ON flip_bid_tab.header_id = headers.id
WHERE flip_bid_tab.bid_id = flip_bid_tab_before_block.bid_id
  AND flip_bid_tab.address_id = flip_bid_tab_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_tab_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION flip_bid_tab_before_block
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
CREATE OR REPLACE FUNCTION maker.insert_new_flip_guy(new_diff maker.flip_bid_guy) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        new_diff.guy,
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET guy = new_diff.guy;
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_guys_until_next_diff(new_diff maker.flip_bid_guy) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_guy_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_guy
                 LEFT JOIN public.headers ON flip_bid_guy.header_id = headers.id
        WHERE flip_bid_guy.bid_id = new_diff.bid_id
          AND flip_bid_guy.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET guy = new_diff.guy
    WHERE flip.bid_id = new_diff.bid_id
      AND flip.address_id = new_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_guy_diff_block IS NULL
        OR flip.block_number < next_guy_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flip_guys_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_guys() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flip_guy(NEW);
    PERFORM maker.update_flip_guys_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_guy
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_guy
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flip_guys();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flip_tic(new_diff maker.flip_bid_tic) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.tic,
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET tic = new_diff.tic;
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_tics_until_next_diff(new_diff maker.flip_bid_tic) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_tic_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_tic
                 LEFT JOIN public.headers ON flip_bid_tic.header_id = headers.id
        WHERE flip_bid_tic.bid_id = new_diff.bid_id
          AND flip_bid_tic.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET tic = new_diff.tic
    WHERE flip.bid_id = new_diff.bid_id
      AND flip.address_id = new_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_tic_diff_block IS NULL
        OR flip.block_number < next_tic_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flip_tics_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_tics() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flip_tic(NEW);
    PERFORM maker.update_flip_tics_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_tic
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_tic
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flip_tics();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flip_end(new_diff maker.flip_bid_end) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff."end",
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET "end" = new_diff."end";
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_ends_until_next_diff(new_diff maker.flip_bid_end) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_end_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_end
                 LEFT JOIN public.headers ON flip_bid_end.header_id = headers.id
        WHERE flip_bid_end.bid_id = new_diff.bid_id
          AND flip_bid_end.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET "end" = new_diff."end"
    WHERE flip.bid_id = new_diff.bid_id
      AND flip.address_id = new_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_end_diff_block IS NULL
        OR flip.block_number < next_end_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flip_ends_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_ends() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flip_end(NEW);
    PERFORM maker.update_flip_ends_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_end
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_end
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flip_ends();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flip_lot(new_diff maker.flip_bid_lot) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.lot,
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET lot = new_diff.lot;
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_lots_until_next_diff(new_diff maker.flip_bid_lot) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_lot_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_lot
                 LEFT JOIN public.headers ON flip_bid_lot.header_id = headers.id
        WHERE flip_bid_lot.bid_id = new_diff.bid_id
          AND flip_bid_lot.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET lot = new_diff.lot
    WHERE flip.bid_id = new_diff.bid_id
      AND flip.address_id = new_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_lot_diff_block IS NULL
        OR flip.block_number < next_lot_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flip_lots_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_lots() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flip_lot(NEW);
    PERFORM maker.update_flip_lots_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_lot
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_lot
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flip_lots();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flip_bid(new_diff maker.flip_bid_bid) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.bid,
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET bid = new_diff.bid;
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_bids_until_next_diff(new_diff maker.flip_bid_bid) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_bid_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_bid
                 LEFT JOIN public.headers ON flip_bid_bid.header_id = headers.id
        WHERE flip_bid_bid.bid_id = new_diff.bid_id
          AND flip_bid_bid.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET bid = new_diff.bid
    WHERE flip.bid_id = new_diff.bid_id
      AND flip.address_id = new_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_bid_diff_block IS NULL
        OR flip.block_number < next_bid_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flip_bids_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_bids() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flip_bid(NEW);
    PERFORM maker.update_flip_bids_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_bid
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_bid
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flip_bids();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flip_usr(new_diff maker.flip_bid_usr) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.usr,
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET usr = new_diff.usr;
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_usrs_until_next_diff(new_diff maker.flip_bid_usr) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_usr_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_usr
                 LEFT JOIN public.headers ON flip_bid_usr.header_id = headers.id
        WHERE flip_bid_usr.bid_id = new_diff.bid_id
          AND flip_bid_usr.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET usr = new_diff.usr
    WHERE flip.bid_id = new_diff.bid_id
      AND flip.address_id = new_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_usr_diff_block IS NULL
        OR flip.block_number < next_usr_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flip_usrs_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_usrs() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flip_usr(NEW);
    PERFORM maker.update_flip_usrs_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_usr
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_usr
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flip_usrs();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flip_gal(new_diff maker.flip_bid_gal) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.gal,
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET gal = new_diff.gal;
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_gals_until_next_diff(new_diff maker.flip_bid_gal) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_gal_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_gal
                 LEFT JOIN public.headers ON flip_bid_gal.header_id = headers.id
        WHERE flip_bid_gal.bid_id = new_diff.bid_id
          AND flip_bid_gal.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET gal = new_diff.gal
    WHERE flip.bid_id = new_diff.bid_id
      AND flip.address_id = new_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_gal_diff_block IS NULL
        OR flip.block_number < next_gal_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flip_gals_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_gals() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flip_gal(NEW);
    PERFORM maker.update_flip_gals_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_gal
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_gal
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flip_gals();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flip_tab(new_diff maker.flip_bid_tab) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.tab,
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET tab = new_diff.tab;
$$
    LANGUAGE sql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_tabs_until_next_diff(new_diff maker.flip_bid_tab) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
    next_tab_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_tab
                 LEFT JOIN public.headers ON flip_bid_tab.header_id = headers.id
        WHERE flip_bid_tab.bid_id = new_diff.bid_id
          AND flip_bid_tab.address_id = new_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET tab = new_diff.tab
    WHERE flip.bid_id = new_diff.bid_id
      AND flip.address_id = new_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_tab_diff_block IS NULL
        OR flip.block_number < next_tab_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flip_tabs_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flip_tabs() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_new_flip_tab(NEW);
    PERFORM maker.update_flip_tabs_until_next_diff(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_tab
    AFTER INSERT OR UPDATE
    ON maker.flip_bid_tab
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_flip_tabs();

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

CREATE TRIGGER flip_created_trigger
    AFTER INSERT
    ON maker.flip_kick
    FOR EACH ROW
EXECUTE PROCEDURE maker.flip_created();


-- +goose Down
DROP TRIGGER flip_guy ON maker.flip_bid_guy;
DROP TRIGGER flip_tic ON maker.flip_bid_tic;
DROP TRIGGER flip_end ON maker.flip_bid_end;
DROP TRIGGER flip_lot ON maker.flip_bid_lot;
DROP TRIGGER flip_bid ON maker.flip_bid_bid;
DROP TRIGGER flip_usr ON maker.flip_bid_usr;
DROP TRIGGER flip_gal ON maker.flip_bid_gal;
DROP TRIGGER flip_tab ON maker.flip_bid_tab;
DROP TRIGGER flip_created_trigger ON maker.flip_kick;

DROP FUNCTION maker.insert_new_flip_guy(maker.flip_bid_guy);
DROP FUNCTION maker.insert_new_flip_tic(maker.flip_bid_tic);
DROP FUNCTION maker.insert_new_flip_end(maker.flip_bid_end);
DROP FUNCTION maker.insert_new_flip_lot(maker.flip_bid_lot);
DROP FUNCTION maker.insert_new_flip_bid(maker.flip_bid_bid);
DROP FUNCTION maker.insert_new_flip_usr(maker.flip_bid_usr);
DROP FUNCTION maker.insert_new_flip_gal(maker.flip_bid_gal);
DROP FUNCTION maker.insert_new_flip_tab(maker.flip_bid_tab);
DROP FUNCTION maker.update_flip_guys_until_next_diff(maker.flip_bid_guy);
DROP FUNCTION maker.update_flip_tics_until_next_diff(maker.flip_bid_tic);
DROP FUNCTION maker.update_flip_ends_until_next_diff(maker.flip_bid_end);
DROP FUNCTION maker.update_flip_lots_until_next_diff(maker.flip_bid_lot);
DROP FUNCTION maker.update_flip_bids_until_next_diff(maker.flip_bid_bid);
DROP FUNCTION maker.update_flip_usrs_until_next_diff(maker.flip_bid_usr);
DROP FUNCTION maker.update_flip_gals_until_next_diff(maker.flip_bid_gal);
DROP FUNCTION maker.update_flip_tabs_until_next_diff(maker.flip_bid_tab);
DROP FUNCTION maker.update_flip_guys();
DROP FUNCTION maker.update_flip_tics();
DROP FUNCTION maker.update_flip_ends();
DROP FUNCTION maker.update_flip_lots();
DROP FUNCTION maker.update_flip_bids();
DROP FUNCTION maker.update_flip_usrs();
DROP FUNCTION maker.update_flip_gals();
DROP FUNCTION maker.update_flip_tabs();
DROP FUNCTION maker.flip_created();
DROP FUNCTION flip_bid_guy_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flip_bid_tic_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flip_bid_end_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flip_bid_lot_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flip_bid_bid_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flip_bid_usr_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flip_bid_gal_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flip_bid_tab_before_block(NUMERIC, INTEGER, INTEGER);
DROP FUNCTION flip_bid_time_created(INTEGER, NUMERIC);

DROP INDEX maker.flip_address_index;
DROP TABLE maker.flip;
