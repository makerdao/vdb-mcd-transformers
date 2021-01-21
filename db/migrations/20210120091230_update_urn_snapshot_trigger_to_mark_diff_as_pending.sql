-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.mark_noncanonical_transformed_diff_as_pending(diff_id BIGINT) RETURNS VOID
AS
    $$
    BEGIN
        UPDATE public.storage_diff SET status = 'pending' WHERE id = diff_id AND status = 'transformed';
    END
    $$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- use mark_diff_pending_if_noncanonical if vat_urn_ink is being deleted
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
        PERFORM maker.mark_noncanonical_transformed_diff_as_pending(OLD.diff_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

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
        PERFORM maker.mark_noncanonical_transformed_diff_as_pending(OLD.diff_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down

-- put update_urn_inks back to previous implementation
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

-- put update_urn_arts back to previous implementation
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

DROP FUNCTION maker.mark_noncanonical_transformed_diff_as_pending(diff_id BIGINT);

