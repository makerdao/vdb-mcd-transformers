-- +goose Up

-- replace urn_ink_before_block function to take in a diff_id param instead of a header_id
DROP FUNCTION urn_ink_before_block(urn_id INTEGER, header_id INTEGER);
CREATE FUNCTION urn_ink_before_block(urn_id INTEGER, diff_id BIGINT) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_height
    FROM public.storage_diff
    WHERE id = diff_id)
SELECT ink
FROM maker.vat_urn_ink
         LEFT JOIN public.storage_diff ON vat_urn_ink.diff_id = storage_diff.id
WHERE vat_urn_ink.urn_id = urn_ink_before_block.urn_id
  AND storage_diff.block_height < (SELECT block_height FROM passed_block_number)
ORDER BY block_height DESC
LIMIT 1
$$
    LANGUAGE sql;

-- replace urn_art_before_block function to take in a diff_id param instead of a header_id
DROP FUNCTION urn_art_before_block(urn_id integer, header_id integer);
CREATE FUNCTION urn_art_before_block(urn_id INTEGER, diff_id BIGINT) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_height
    FROM public.storage_diff
    WHERE id = diff_id)
SELECT art
FROM maker.vat_urn_art
         LEFT JOIN public.storage_diff ON vat_urn_art.diff_id = storage_diff.id
WHERE vat_urn_art.urn_id = urn_art_before_block.urn_id
  AND storage_diff.block_height < (SELECT block_height FROM passed_block_number)
ORDER BY block_height DESC
LIMIT 1
$$
    LANGUAGE sql;

-- replace delete_obsolete_urn_snapshot function to take in a diff_id param
DROP FUNCTION  maker.delete_obsolete_urn_snapshot(urn_id INTEGER, header_id INTEGER);
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.delete_obsolete_urn_snapshot(urn_id INTEGER, header_id INTEGER, diff_id BIGINT) RETURNS api.urn_snapshot
AS
$$
DECLARE
    urn_snapshot_block_number BIGINT := (
        SELECT block_height
        FROM public.storage_diff
        WHERE id = diff_id);
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

-- update insert_urn_ink function to use new urn_art_before_block function
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
            urn_art_before_block(new_diff.urn_id, new_diff.diff_id),
            urn_time_created(new_diff.urn_id),
            (SELECT block_timestamp FROM new_diff_header))
    ON CONFLICT (urn_identifier, ilk_identifier, block_height)
        DO UPDATE SET ink = new_diff.ink;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- update update_urn_inks_until_next_diff to use the snapshot's diff instead of its header to get block_number
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_urn_inks_until_next_diff(start_at_diff maker.vat_urn_ink, new_ink NUMERIC) RETURNS maker.vat_urn_ink
AS
$$
DECLARE
    diff_block_number    BIGINT := (
        SELECT block_height
        FROM public.storage_diff
        WHERE id = start_at_diff.diff_id);
    next_ink_diff_block BIGINT := (
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
      AND (next_ink_diff_block IS NULL
        OR urn_snapshot.block_height < next_ink_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- use the new maker.delete_obsolete_urn_snapshot function in update_urn_inks
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
        PERFORM maker.update_urn_inks_until_next_diff(OLD, urn_ink_before_block(OLD.urn_id, OLD.diff_id));
        PERFORM maker.delete_obsolete_urn_snapshot(OLD.urn_id, OLD.header_id, OLD.diff_id);
        PERFORM maker.update_urn_created(OLD.urn_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- update insert_urn_art function to use new urn_ink_before_block function
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
            urn_ink_before_block(new_diff.urn_id, new_diff.diff_id),
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

-- update update_urn_arts to use the new delete_obsolete_urn_snapshot function
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_urn_arts() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_urn_art(NEW);
        PERFORM maker.update_urn_arts_until_next_diff(NEW, NEW.art);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_urn_arts_until_next_diff(OLD, urn_art_before_block(OLD.urn_id, OLD.diff_id));
        PERFORM maker.delete_obsolete_urn_snapshot(OLD.urn_id, OLD.header_id, OLD.diff_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- update update_urn_arts_until_next_diff to use the snapshot's diff instead of its header to get block_number
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_urn_arts_until_next_diff(start_at_diff maker.vat_urn_art, new_art NUMERIC) RETURNS maker.vat_urn_art
AS
$$
DECLARE
    diff_block_number    BIGINT := (
        SELECT block_height
        FROM public.storage_diff
        WHERE id = start_at_diff.diff_id);
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

-- +goose Down



-- drop new urn_ink_before_block and replace with old function
DROP FUNCTION urn_ink_before_block(urn_id INTEGER, diff_id BIGINT);
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


-- drop new urn_art_before_block and replace with old function
DROP FUNCTION urn_art_before_block(urn_id INTEGER, diff_id BIGINT);
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



-- drop new delete_obsolete_urn_snapshot and replace with old function
DROP FUNCTION maker.delete_obsolete_urn_snapshot(urn_id INTEGER, header_id INTEGER, diff_id BIGINT);
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


-- put insert_urn_ink back to old implementation
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

-- put update_urn_inks_until_next_diff back to old implementation
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

-- put update_urn_inks back to old implementation
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

-- put insert_urn_art back to old implementation
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


-- put update_urn_arts back to old implementation
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


-- put update_urn_arts_until_next_diff back to old implementation
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
