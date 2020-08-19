-- +goose Up
CREATE TABLE api.urn_snapshot
(
    urn_identifier TEXT,
    ilk_identifier TEXT,
    block_height   BIGINT,
    ink            NUMERIC   DEFAULT NULL,
    art            NUMERIC   DEFAULT NULL,
    created        TIMESTAMP DEFAULT NULL,
    updated        TIMESTAMP NOT NULL,
    PRIMARY KEY (urn_identifier, ilk_identifier, block_height)
);

CREATE FUNCTION api.max_block()
    RETURNS BIGINT AS
$$
SELECT max(block_number)
FROM public.headers
$$
    LANGUAGE SQL
    STABLE;


CREATE FUNCTION urn_ink_before_block(urn_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id)
SELECT ink
FROM maker.vat_urn_ink
         LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
WHERE vat_urn_ink.urn_id = urn_ink_before_block.urn_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION urn_ink_before_block
    IS E'@omit';


CREATE FUNCTION urn_art_before_block(urn_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id)
SELECT art
FROM maker.vat_urn_art
         LEFT JOIN public.headers ON vat_urn_art.header_id = headers.id
WHERE vat_urn_art.urn_id = urn_art_before_block.urn_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION urn_art_before_block
    IS E'@omit';


CREATE FUNCTION urn_time_created(urn_id INTEGER) RETURNS TIMESTAMP AS
$$
SELECT api.epoch_to_datetime(MIN(block_timestamp))
FROM maker.vat_urn_ink
         LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
WHERE vat_urn_ink.urn_id = urn_time_created.urn_id
$$
    LANGUAGE sql;

COMMENT ON FUNCTION urn_time_created
    IS E'@omit';


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.delete_obsolete_urn_snapshot(urn_id INTEGER, header_id INTEGER) RETURNS api.urn_snapshot
    -- deletes row if there are no longer any diffs in that block for the associated urn
AS
$$
DECLARE
    urn_snapshot_block_number BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = header_id);
BEGIN
    DELETE
    FROM api.urn_snapshot
        USING maker.urns LEFT JOIN maker.ilks ON urns.ilk_id = ilks.id
    WHERE urn_snapshot.urn_identifier = urns.identifier
      AND urn_snapshot.ilk_identifier = ilks.identifier
      AND urns.id = urn_id
      AND urn_snapshot.block_height = urn_snapshot_block_number
      AND NOT (EXISTS(
            SELECT *
            FROM maker.vat_urn_ink
            WHERE vat_urn_ink.urn_id = delete_obsolete_urn_snapshot.urn_id
              AND vat_urn_ink.header_id = delete_obsolete_urn_snapshot.header_id))
      AND NOT (EXISTS(
            SELECT *
            FROM maker.vat_urn_art
            WHERE vat_urn_art.urn_id = delete_obsolete_urn_snapshot.urn_id
              AND vat_urn_art.header_id = delete_obsolete_urn_snapshot.header_id));
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


COMMENT ON FUNCTION maker.delete_obsolete_urn_snapshot
    IS E'@omit';


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_urn_ink(new_diff maker.vat_urn_ink) RETURNS maker.vat_urn_ink
AS
$$
BEGIN
    WITH urn AS (
        SELECT urns.identifier AS urn_identifier, ilks.identifier AS ilk_identifier
        FROM maker.urns
                 LEFT JOIN maker.ilks ON urns.ilk_id = ilks.id
        WHERE urns.id = new_diff.urn_id),
         new_diff_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp, block_number
             FROM public.headers
             WHERE id = new_diff.header_id)
    INSERT
    INTO api.urn_snapshot (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)
    VALUES ((SELECT urn_identifier FROM urn),
            (SELECT ilk_identifier FROM urn),
            (SELECT block_number FROM new_diff_header),
            new_diff.ink,
            urn_art_before_block(new_diff.urn_id, new_diff.header_id),
            urn_time_created(new_diff.urn_id),
            (SELECT block_timestamp FROM new_diff_header))
    ON CONFLICT (urn_identifier, ilk_identifier, block_height)
        DO UPDATE SET ink = new_diff.ink;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_urn_ink
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_urn_inks_until_next_diff(start_at_diff maker.vat_urn_ink, new_ink NUMERIC) RETURNS maker.vat_urn_ink
AS
$$
DECLARE
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_rate_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_urn_ink
                 LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
        WHERE vat_urn_ink.urn_id = start_at_diff.urn_id
          AND block_number > diff_block_number);
BEGIN
    WITH urn AS (
        SELECT urns.identifier AS urn_identifier, ilks.identifier AS ilk_identifier
        FROM maker.urns
                 LEFT JOIN maker.ilks ON urns.ilk_id = ilks.id
        WHERE urns.id = start_at_diff.urn_id)
    UPDATE api.urn_snapshot
    SET ink = new_ink
    FROM urn
    WHERE urn_snapshot.urn_identifier = urn.urn_identifier
      AND urn_snapshot.ilk_identifier = urn.ilk_identifier
      AND urn_snapshot.block_height >= diff_block_number
      AND (next_rate_diff_block IS NULL
        OR urn_snapshot.block_height < next_rate_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_urn_inks_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_urn_created(urn_id INTEGER) RETURNS maker.vat_urn_ink
AS
$$
BEGIN
    UPDATE api.urn_snapshot
    SET created = urn_time_created(urn_id)
    FROM maker.urns
             LEFT JOIN maker.ilks ON urns.ilk_id = ilks.id
    WHERE urns.identifier = urn_snapshot.urn_identifier
      AND ilks.identifier = urn_snapshot.ilk_identifier
      AND urns.id = urn_id;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_urn_created
    IS E'@omit';


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_urn_inks() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_urn_ink(NEW);
        PERFORM maker.update_urn_inks_until_next_diff(NEW, NEW.ink);
        PERFORM maker.update_urn_created(NEW.urn_id);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_urn_inks_until_next_diff(OLD, urn_ink_before_block(OLD.urn_id, OLD.header_id));
        PERFORM maker.delete_obsolete_urn_snapshot(OLD.urn_id, OLD.header_id);
        PERFORM maker.update_urn_created(OLD.urn_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER urn_ink
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.vat_urn_ink
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_urn_inks();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_urn_art(new_diff maker.vat_urn_art) RETURNS maker.vat_urn_art
AS
$$
BEGIN
    WITH urn AS (
        SELECT urns.identifier AS urn_identifier, ilks.identifier AS ilk_identifier
        FROM maker.urns
                 LEFT JOIN maker.ilks ON urns.ilk_id = ilks.id
        WHERE urns.id = new_diff.urn_id),
         new_diff_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp, block_number
             FROM public.headers
             WHERE id = new_diff.header_id)
    INSERT
    INTO api.urn_snapshot (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)
    VALUES ((SELECT urn_identifier FROM urn),
            (SELECT ilk_identifier FROM urn),
            (SELECT block_number FROM new_diff_header),
            urn_ink_before_block(new_diff.urn_id, new_diff.header_id),
            new_diff.art,
            urn_time_created(new_diff.urn_id),
            (SELECT block_timestamp FROM new_diff_header))
    ON CONFLICT (urn_identifier, ilk_identifier, block_height)
        DO UPDATE SET art = new_diff.art;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_urn_art
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_urn_arts_until_next_diff(start_at_diff maker.vat_urn_art, new_art NUMERIC) RETURNS maker.vat_urn_art
AS
$$
DECLARE
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_rate_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_urn_art
                 LEFT JOIN public.headers ON vat_urn_art.header_id = headers.id
        WHERE vat_urn_art.urn_id = start_at_diff.urn_id
          AND block_number > diff_block_number);
BEGIN
    WITH urn AS (
        SELECT urns.identifier AS urn_identifier, ilks.identifier AS ilk_identifier
        FROM maker.urns
                 LEFT JOIN maker.ilks ON urns.ilk_id = ilks.id
        WHERE urns.id = start_at_diff.urn_id)
    UPDATE api.urn_snapshot
    SET art = new_art
    FROM urn
    WHERE urn_snapshot.urn_identifier = urn.urn_identifier
      AND urn_snapshot.ilk_identifier = urn.ilk_identifier
      AND urn_snapshot.block_height >= diff_block_number
      AND (next_rate_diff_block IS NULL
        OR urn_snapshot.block_height < next_rate_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_urn_arts_until_next_diff
    IS E'@omit';


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_urn_arts() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_urn_art(NEW);
        PERFORM maker.update_urn_arts_until_next_diff(NEW, NEW.art);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_urn_arts_until_next_diff(OLD, urn_art_before_block(OLD.urn_id, OLD.header_id));
        PERFORM maker.delete_obsolete_urn_snapshot(OLD.urn_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER urn_art
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.vat_urn_art
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_urn_arts();


-- +goose Down
DROP TRIGGER urn_art ON maker.vat_urn_art;
DROP TRIGGER urn_ink ON maker.vat_urn_ink;

DROP FUNCTION maker.update_urn_arts();
DROP FUNCTION maker.update_urn_inks();

DROP FUNCTION maker.update_urn_created(INTEGER);
DROP FUNCTION maker.update_urn_arts_until_next_diff(maker.vat_urn_art, NUMERIC);
DROP FUNCTION maker.update_urn_inks_until_next_diff(maker.vat_urn_ink, NUMERIC);

DROP FUNCTION maker.insert_urn_art(maker.vat_urn_art);
DROP FUNCTION maker.insert_urn_ink(maker.vat_urn_ink);

DROP FUNCTION maker.delete_obsolete_urn_snapshot(INTEGER, INTEGER);

DROP FUNCTION urn_time_created(INTEGER);
DROP FUNCTION urn_art_before_block(INTEGER, INTEGER);
DROP FUNCTION urn_ink_before_block(INTEGER, INTEGER);

DROP FUNCTION api.max_block();

DROP TABLE api.urn_snapshot;
