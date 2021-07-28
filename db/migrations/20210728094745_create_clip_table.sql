-- +goose Up
CREATE TABLE maker.clip
(
    address_id   BIGINT    NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    block_number BIGINT    NOT NULL,
    sale_id      NUMERIC   NOT NULL,
    pos          NUMERIC   DEFAULT NULL,
    tab          NUMERIC   DEFAULT NULL,
    lot          NUMERIC   DEFAULT NULL,
    usr          TEXT      DEFAULT NULL,
    tic          NUMERIC   DEFAULT NULL,
    "top"        NUMERIC   DEFAULT NULL,
    created      TIMESTAMP DEFAULT NULL,
    updated      TIMESTAMP NOT NULL,
    PRIMARY KEY (address_id, sale_id, block_number)
);

CREATE INDEX clip_address_index
    ON maker.clip (address_id);

CREATE FUNCTION clip_sale_pos_before_block(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT pos
FROM maker.clip_sale_pos
         LEFT JOIN public.headers ON clip_sale_pos.header_id = headers.id
WHERE clip_sale_pos.sale_id = clip_sale_pos_before_block.sale_id
  AND clip_sale_pos.address_id = clip_sale_pos_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = clip_sale_pos_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_pos_before_block
    IS E'@omit';

CREATE FUNCTION clip_sale_tab_before_block(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT tab
FROM maker.clip_sale_tab
         LEFT JOIN public.headers ON clip_sale_tab.header_id = headers.id
WHERE clip_sale_tab.sale_id = clip_sale_tab_before_block.sale_id
  AND clip_sale_tab.address_id = clip_sale_tab_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = clip_sale_tab_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_tab_before_block
    IS E'@omit';

CREATE FUNCTION clip_sale_lot_before_block(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT lot
FROM maker.clip_sale_lot
         LEFT JOIN public.headers ON clip_sale_lot.header_id = headers.id
WHERE clip_sale_lot.sale_id = clip_sale_lot_before_block.sale_id
  AND clip_sale_lot.address_id = clip_sale_lot_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = clip_sale_lot_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_lot_before_block
    IS E'@omit';

CREATE FUNCTION clip_sale_usr_before_block(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS TEXT AS
$$
SELECT usr
FROM maker.clip_sale_usr
         LEFT JOIN public.headers ON clip_sale_usr.header_id = headers.id
WHERE clip_sale_usr.sale_id = clip_sale_usr_before_block.sale_id
  AND clip_sale_usr.address_id = clip_sale_usr_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = clip_sale_usr_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_usr_before_block
    IS E'@omit';

CREATE FUNCTION clip_sale_tic_before_block(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT tic
FROM maker.clip_sale_tic
         LEFT JOIN public.headers ON clip_sale_tic.header_id = headers.id
WHERE clip_sale_tic.sale_id = clip_sale_tic_before_block.sale_id
  AND clip_sale_tic.address_id = clip_sale_tic_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = clip_sale_tic_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_tic_before_block
    IS E'@omit';

CREATE FUNCTION clip_sale_top_before_block(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT top
FROM maker.clip_sale_top
         LEFT JOIN public.headers ON clip_sale_top.header_id = headers.id
WHERE clip_sale_top.sale_id = clip_sale_top_before_block.sale_id
  AND clip_sale_top.address_id = clip_sale_top_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = clip_sale_top_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_top_before_block
    IS E'@omit';

CREATE FUNCTION clip_sale_time_created(address_id BIGINT, sale_id NUMERIC) RETURNS TIMESTAMP AS
$$
SELECT api.epoch_to_datetime(MIN(block_timestamp))
FROM public.headers
         LEFT JOIN maker.clip_kick ON clip_kick.header_id = headers.id
WHERE clip_kick.address_id = clip_sale_time_created.address_id
  AND clip_kick.sale_id = clip_sale_time_created.sale_id
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_time_created
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.delete_obsolete_clip(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS VOID AS
$$
DECLARE
    clip_block      BIGINT     := (
        SELECT block_number
        FROM public.headers
        WHERE id = header_id);
    clip_state      maker.clip := (
        SELECT (clip.address_id, block_number, clip.sale_id, pos, tab, lot, usr, tic, "top", created, updated)
        FROM maker.clip
        WHERE clip.sale_id = delete_obsolete_clip.sale_id
          AND clip.address_id = delete_obsolete_clip.address_id
          AND clip.block_number = clip_block);
    prev_clip_state maker.clip := (
        SELECT (clip.address_id, block_number, clip.sale_id, pos, tab, lot, usr, tic, "top", created, updated)
        FROM maker.clip
        WHERE clip.sale_id = delete_obsolete_clip.sale_id
          AND clip.address_id = delete_obsolete_clip.address_id
          AND clip.block_number < clip_block
        ORDER BY clip.block_number DESC
        LIMIT 1);
BEGIN
    DELETE
    FROM maker.clip
    WHERE clip.sale_id = delete_obsolete_clip.sale_id
      AND clip.address_id = delete_obsolete_clip.address_id
      AND clip.block_number = clip_block
      AND (clip_state.pos IS NULL OR clip_state.pos = prev_clip_state.pos)
      AND (clip_state.tab IS NULL OR clip_state.tab = prev_clip_state.tab)
      AND (clip_state.lot IS NULL OR clip_state.lot = prev_clip_state.lot)
      AND (clip_state.usr IS NULL OR clip_state.usr = prev_clip_state.usr)
      AND (clip_state.tic IS NULL OR clip_state.tic = prev_clip_state.tic)
      AND (clip_state."top" IS NULL OR clip_state."top" = prev_clip_state."top");
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.delete_obsolete_clip
    IS E'@omit';

CREATE OR REPLACE FUNCTION maker.insert_new_clip_pos(new_diff maker.clip_sale_pos) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.clip (sale_id, address_id, block_number, pos, tab, lot, usr, tic, "top", updated, created)
VALUES (new_diff.sale_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        new_diff.pos,
        clip_sale_tab_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_lot_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_usr_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_tic_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_top_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        clip_sale_time_created(new_diff.address_id, new_diff.sale_id))
ON CONFLICT (block_number, sale_id, address_id) DO UPDATE SET pos = new_diff.pos
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_clip_pos(new_diff maker.clip_sale_pos)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_pos_until_next_diff(start_at_diff maker.clip_sale_pos, new_pos NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_pos_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.clip_sale_pos
                 LEFT JOIN public.headers ON clip_sale_pos.header_id = headers.id
        WHERE clip_sale_pos.sale_id = start_at_diff.sale_id
          AND clip_sale_pos.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.clip
    SET pos = new_pos
    WHERE clip.sale_id = start_at_diff.sale_id
      AND clip.address_id = start_at_diff.address_id
      AND clip.block_number >= diff_block_number
      AND (next_pos_diff_block IS NULL
        OR clip.block_number < next_pos_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_clip_pos_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_pos() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_clip_pos(NEW);
        PERFORM maker.update_clip_pos_until_next_diff(NEW, NEW.pos);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_clip_pos_until_next_diff(
                OLD,
                clip_sale_pos_before_block(OLD.sale_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_clip(OLD.sale_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER clip_pos
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.clip_sale_pos
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_clip_pos();

CREATE OR REPLACE FUNCTION maker.insert_new_clip_tab(new_diff maker.clip_sale_tab) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.clip (sale_id, address_id, block_number, pos, tab, lot, usr, tic, "top", updated, created)
VALUES (new_diff.sale_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        clip_sale_pos_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        new_diff.tab,
        clip_sale_lot_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_usr_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_tic_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_top_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        clip_sale_time_created(new_diff.address_id, new_diff.sale_id))
ON CONFLICT (block_number, sale_id, address_id) DO UPDATE SET tab = new_diff.tab
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_clip_tab(new_diff maker.clip_sale_tab)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_tab_until_next_diff(start_at_diff maker.clip_sale_tab, new_tab NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_tab_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.clip_sale_tab
                 LEFT JOIN public.headers ON clip_sale_tab.header_id = headers.id
        WHERE clip_sale_tab.sale_id = start_at_diff.sale_id
          AND clip_sale_tab.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.clip
    SET tab = new_tab
    WHERE clip.sale_id = start_at_diff.sale_id
      AND clip.address_id = start_at_diff.address_id
      AND clip.block_number >= diff_block_number
      AND (next_tab_diff_block IS NULL
        OR clip.block_number < next_tab_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_clip_tab_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_tab() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_clip_tab(NEW);
        PERFORM maker.update_clip_tab_until_next_diff(NEW, NEW.tab);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_clip_tab_until_next_diff(
                OLD,
                clip_sale_tab_before_block(OLD.sale_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_clip(OLD.sale_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER clip_tab
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.clip_sale_tab
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_clip_tab();

CREATE OR REPLACE FUNCTION maker.insert_new_clip_lot(new_diff maker.clip_sale_lot) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.clip (sale_id, address_id, block_number, pos, tab, lot, usr, tic, "top", updated, created)
VALUES (new_diff.sale_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        clip_sale_pos_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_tab_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        new_diff.lot,
        clip_sale_usr_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_tic_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_top_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        clip_sale_time_created(new_diff.address_id, new_diff.sale_id))
ON CONFLICT (block_number, sale_id, address_id) DO UPDATE SET lot = new_diff.lot
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_clip_lot(new_diff maker.clip_sale_lot)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_lot_until_next_diff(start_at_diff maker.clip_sale_lot, new_lot NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_lot_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.clip_sale_lot
                 LEFT JOIN public.headers ON clip_sale_lot.header_id = headers.id
        WHERE clip_sale_lot.sale_id = start_at_diff.sale_id
          AND clip_sale_lot.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.clip
    SET lot = new_lot
    WHERE clip.sale_id = start_at_diff.sale_id
      AND clip.address_id = start_at_diff.address_id
      AND clip.block_number >= diff_block_number
      AND (next_lot_diff_block IS NULL
        OR clip.block_number < next_lot_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_clip_lot_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_lot() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_clip_lot(NEW);
        PERFORM maker.update_clip_lot_until_next_diff(NEW, NEW.lot);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_clip_lot_until_next_diff(
                OLD,
                clip_sale_lot_before_block(OLD.sale_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_clip(OLD.sale_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER clip_lot
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.clip_sale_lot
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_clip_lot();

CREATE OR REPLACE FUNCTION maker.insert_new_clip_usr(new_diff maker.clip_sale_usr) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.clip (sale_id, address_id, block_number, pos, tab, lot, usr, tic, "top", updated, created)
VALUES (new_diff.sale_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        clip_sale_pos_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_tab_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_lot_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        new_diff.usr,
        clip_sale_tic_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_top_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        clip_sale_time_created(new_diff.address_id, new_diff.sale_id))
ON CONFLICT (block_number, sale_id, address_id) DO UPDATE SET usr = new_diff.usr
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_clip_usr(new_diff maker.clip_sale_usr)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_usr_until_next_diff(start_at_diff maker.clip_sale_usr, new_usr TEXT) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_usr_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.clip_sale_usr
                 LEFT JOIN public.headers ON clip_sale_usr.header_id = headers.id
        WHERE clip_sale_usr.sale_id = start_at_diff.sale_id
          AND clip_sale_usr.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.clip
    SET usr = new_usr
    WHERE clip.sale_id = start_at_diff.sale_id
      AND clip.address_id = start_at_diff.address_id
      AND clip.block_number >= diff_block_number
      AND (next_usr_diff_block IS NULL
        OR clip.block_number < next_usr_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_clip_usr_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_usr() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_clip_usr(NEW);
        PERFORM maker.update_clip_usr_until_next_diff(NEW, NEW.usr);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_clip_usr_until_next_diff(
                OLD,
                clip_sale_usr_before_block(OLD.sale_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_clip(OLD.sale_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER clip_usr
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.clip_sale_usr
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_clip_usr();

CREATE OR REPLACE FUNCTION maker.insert_new_clip_tic(new_diff maker.clip_sale_tic) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.clip (sale_id, address_id, block_number, pos, tab, lot, usr, tic, "top", updated, created)
VALUES (new_diff.sale_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        clip_sale_pos_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_tab_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_lot_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_usr_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        new_diff.tic,
        clip_sale_top_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        clip_sale_time_created(new_diff.address_id, new_diff.sale_id))
ON CONFLICT (block_number, sale_id, address_id) DO UPDATE SET tic = new_diff.tic
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_clip_tic(new_diff maker.clip_sale_tic)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_tic_until_next_diff(start_at_diff maker.clip_sale_tic, new_tic NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_tic_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.clip_sale_tic
                 LEFT JOIN public.headers ON clip_sale_tic.header_id = headers.id
        WHERE clip_sale_tic.sale_id = start_at_diff.sale_id
          AND clip_sale_tic.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.clip
    SET tic = new_tic
    WHERE clip.sale_id = start_at_diff.sale_id
      AND clip.address_id = start_at_diff.address_id
      AND clip.block_number >= diff_block_number
      AND (next_tic_diff_block IS NULL
        OR clip.block_number < next_tic_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_clip_tic_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_tic() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_clip_tic(NEW);
        PERFORM maker.update_clip_tic_until_next_diff(NEW, NEW.tic);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_clip_tic_until_next_diff(
                OLD,
                clip_sale_tic_before_block(OLD.sale_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_clip(OLD.sale_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER clip_tic
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.clip_sale_tic
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_clip_tic();

CREATE OR REPLACE FUNCTION maker.insert_new_clip_top(new_diff maker.clip_sale_top) RETURNS VOID AS
$$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.clip (sale_id, address_id, block_number, pos, tab, lot, usr, tic, "top", updated, created)
VALUES (new_diff.sale_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        clip_sale_pos_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_tab_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_lot_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_usr_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        clip_sale_tic_before_block(new_diff.sale_id, new_diff.address_id, new_diff.header_id),
        new_diff.top,
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        clip_sale_time_created(new_diff.address_id, new_diff.sale_id))
ON CONFLICT (block_number, sale_id, address_id) DO UPDATE SET top = new_diff.top
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_new_clip_top(new_diff maker.clip_sale_top)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_top_until_next_diff(start_at_diff maker.clip_sale_top, new_top NUMERIC) RETURNS VOID
AS
$$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_top_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.clip_sale_top
                 LEFT JOIN public.headers ON clip_sale_top.header_id = headers.id
        WHERE clip_sale_top.sale_id = start_at_diff.sale_id
          AND clip_sale_top.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.clip
    SET "top" = new_top
    WHERE clip.sale_id = start_at_diff.sale_id
      AND clip.address_id = start_at_diff.address_id
      AND clip.block_number >= diff_block_number
      AND (next_top_diff_block IS NULL
        OR clip.block_number < next_top_diff_block);
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_clip_top_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_top() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_clip_top(NEW);
        PERFORM maker.update_clip_top_until_next_diff(NEW, NEW.top);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_clip_top_until_next_diff(
                OLD,
                clip_sale_top_before_block(OLD.sale_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_clip(OLD.sale_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER clip_top
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.clip_sale_top
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_clip_top();



CREATE OR REPLACE FUNCTION maker.insert_clip_created(new_event maker.clip_kick) RETURNS VOID
AS
$$
UPDATE maker.clip
SET created = api.epoch_to_datetime(headers.block_timestamp)
FROM public.headers
WHERE headers.id = new_event.header_id
  AND clip.address_id = new_event.address_id
  AND clip.sale_id = new_event.sale_id
  AND clip.created IS NULL
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_clip_created
    IS E'@omit';

CREATE OR REPLACE FUNCTION maker.clear_clip_created(old_event maker.clip_kick) RETURNS VOID
AS
$$
UPDATE maker.clip
SET created = clip_sale_time_created(old_event.address_id, old_event.sale_id)
WHERE clip.address_id = old_event.address_id
  AND clip.sale_id = old_event.sale_id
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.clear_clip_created
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_created() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        PERFORM maker.insert_clip_created(NEW);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.clear_clip_created(OLD);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_clip_created
    IS E'@omit';

CREATE TRIGGER clip_created_trigger
    AFTER INSERT OR DELETE
    ON maker.clip_kick
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_clip_created();

-- +goose Down
DROP TRIGGER clip_pos ON maker.clip_sale_pos;
DROP TRIGGER clip_tab ON maker.clip_sale_tab;
DROP TRIGGER clip_lot ON maker.clip_sale_lot;
DROP TRIGGER clip_usr ON maker.clip_sale_usr;
DROP TRIGGER clip_tic ON maker.clip_sale_tic;
DROP TRIGGER clip_top ON maker.clip_sale_top;
DROP TRIGGER clip_created_trigger ON maker.clip_kick;

DROP FUNCTION maker.insert_new_clip_pos(maker.clip_sale_pos);
DROP FUNCTION maker.insert_new_clip_tab(maker.clip_sale_tab);
DROP FUNCTION maker.insert_new_clip_lot(maker.clip_sale_lot);
DROP FUNCTION maker.insert_new_clip_usr(maker.clip_sale_usr);
DROP FUNCTION maker.insert_new_clip_tic(maker.clip_sale_tic);
DROP FUNCTION maker.insert_new_clip_top(maker.clip_sale_top);
DROP FUNCTION maker.insert_clip_created(maker.clip_kick);
DROP FUNCTION maker.update_clip_pos_until_next_diff(maker.clip_sale_pos, NUMERIC);
DROP FUNCTION maker.update_clip_tab_until_next_diff(maker.clip_sale_tab, NUMERIC);
DROP FUNCTION maker.update_clip_lot_until_next_diff(maker.clip_sale_lot, NUMERIC);
DROP FUNCTION maker.update_clip_usr_until_next_diff(maker.clip_sale_usr, TEXT);
DROP FUNCTION maker.update_clip_tic_until_next_diff(maker.clip_sale_tic, NUMERIC);
DROP FUNCTION maker.update_clip_top_until_next_diff(maker.clip_sale_top, NUMERIC);
DROP FUNCTION maker.clear_clip_created(maker.clip_kick);
DROP FUNCTION maker.update_clip_pos();
DROP FUNCTION maker.update_clip_tab();
DROP FUNCTION maker.update_clip_lot();
DROP FUNCTION maker.update_clip_usr();
DROP FUNCTION maker.update_clip_tic();
DROP FUNCTION maker.update_clip_top();
DROP FUNCTION maker.update_clip_created();
DROP FUNCTION clip_sale_pos_before_block(NUMERIC, BIGINT, INTEGER);
DROP FUNCTION clip_sale_tab_before_block(NUMERIC, BIGINT, INTEGER);
DROP FUNCTION clip_sale_lot_before_block(NUMERIC, BIGINT, INTEGER);
DROP FUNCTION clip_sale_usr_before_block(NUMERIC, BIGINT, INTEGER);
DROP FUNCTION clip_sale_tic_before_block(NUMERIC, BIGINT, INTEGER);
DROP FUNCTION clip_sale_top_before_block(NUMERIC, BIGINT, INTEGER);
DROP FUNCTION clip_sale_time_created(BIGINT, NUMERIC);

DROP INDEX maker.clip_address_index;
DROP TABLE maker.clip;