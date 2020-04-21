-- +goose Up
DROP TRIGGER urn_ink ON maker.vat_urn_ink;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_urn_inks() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_urn_ink(NEW);
        PERFORM maker.update_urn_inks_until_next_diff(NEW, NEW.ink);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_urn_inks_until_next_diff(OLD, urn_ink_before_block(OLD.urn_id, OLD.header_id));
        PERFORM maker.delete_obsolete_urn_state(OLD.urn_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

DROP FUNCTION maker.update_urn_created(INTEGER);

CREATE TRIGGER urn_ink
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.vat_urn_ink
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_urn_inks();

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER urn_ink ON maker.vat_urn_ink;

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

CREATE TRIGGER urn_ink
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.vat_urn_ink
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_urn_inks();