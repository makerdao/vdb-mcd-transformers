-- +goose Up

DROP FUNCTION maker.delete_obsolete_urn_state;

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

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP FUNCTION maker.delete_obsolete_urn_snapshot;

-- +goose StatementBegin
CREATE FUNCTION maker.delete_obsolete_urn_state(urn_id INTEGER, header_id INTEGER) RETURNS api.urn_snapshot
    -- deletes row if there are no longer any diffs in that block for the associated urn
AS
$$
DECLARE
    urn_state_block_number BIGINT := (
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
      AND urn_snapshot.block_height = urn_state_block_number
      AND NOT (EXISTS(
            SELECT *
            FROM maker.vat_urn_ink
            WHERE vat_urn_ink.urn_id = delete_obsolete_urn_state.urn_id
              AND vat_urn_ink.header_id = delete_obsolete_urn_state.header_id))
      AND NOT (EXISTS(
            SELECT *
            FROM maker.vat_urn_art
            WHERE vat_urn_art.urn_id = delete_obsolete_urn_state.urn_id
              AND vat_urn_art.header_id = delete_obsolete_urn_state.header_id));
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.delete_obsolete_urn_state
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
        PERFORM maker.delete_obsolete_urn_state(OLD.urn_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

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
        PERFORM maker.delete_obsolete_urn_state(OLD.urn_id, OLD.header_id);
        PERFORM maker.update_urn_created(OLD.urn_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

