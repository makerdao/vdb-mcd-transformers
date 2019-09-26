-- +goose Up
CREATE TABLE api.current_ilk_state
(
    ilk_identifier TEXT      PRIMARY KEY,
    rate           NUMERIC   DEFAULT NULL,
    art            NUMERIC   DEFAULT NULL,
    spot           NUMERIC   DEFAULT NULL,
    line           NUMERIC   DEFAULT NULL,
    dust           NUMERIC   DEFAULT NULL,
    chop           NUMERIC   DEFAULT NULL,
    lump           NUMERIC   DEFAULT NULL,
    flip           TEXT      DEFAULT NULL,
    rho            NUMERIC   DEFAULT NULL,
    duty           NUMERIC   DEFAULT NULL,
    pip            TEXT      DEFAULT NULL,
    mat            NUMERIC   DEFAULT NULL,
    created        TIMESTAMP DEFAULT NULL,
    updated        TIMESTAMP DEFAULT NULL
);

COMMENT ON TABLE api.current_ilk_state IS '@omit create,update,delete';
COMMENT ON COLUMN api.current_ilk_state.ilk_identifier IS '@name id';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_rate() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, rate, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.rate,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET rate    = (
            CASE
                WHEN current_ilk_state.rate IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.rate
                ELSE current_ilk_state.rate END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_art() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, art, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.art,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET art     = (
            CASE
                WHEN current_ilk_state.art IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.art
                ELSE current_ilk_state.art END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_spot() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, spot, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.spot,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET spot    = (
            CASE
                WHEN current_ilk_state.spot IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.spot
                ELSE current_ilk_state.spot END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_line() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, line, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.line,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET line    = (
            CASE
                WHEN current_ilk_state.line IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.line
                ELSE current_ilk_state.line END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_dust() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, dust, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.dust,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET dust    = (
            CASE
                WHEN current_ilk_state.dust IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.dust
                ELSE current_ilk_state.dust END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_chop() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, chop, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.chop,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET chop    = (
            CASE
                WHEN current_ilk_state.chop IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.chop
                ELSE current_ilk_state.chop END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_lump() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, lump, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.lump,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET lump    = (
            CASE
                WHEN current_ilk_state.lump IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.lump
                ELSE current_ilk_state.lump END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_flip() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, flip, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.flip,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET flip    = (
            CASE
                WHEN current_ilk_state.flip IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.flip
                ELSE current_ilk_state.flip END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_rho() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, rho, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.rho,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET rho     = (
            CASE
                WHEN current_ilk_state.rho IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.rho
                ELSE current_ilk_state.rho END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_duty() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, duty, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.duty,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET duty    = (
            CASE
                WHEN current_ilk_state.duty IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.duty
                ELSE current_ilk_state.duty END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_pip() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, pip, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.pip,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET pip     = (
            CASE
                WHEN current_ilk_state.pip IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.pip
                ELSE current_ilk_state.pip END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_ilk_mat() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, mat, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.mat,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET mat     = (
            CASE
                WHEN current_ilk_state.mat IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.mat
                ELSE current_ilk_state.mat END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_rate
    AFTER INSERT OR UPDATE
    ON maker.vat_ilk_rate
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_rate();

CREATE TRIGGER ilk_art
    AFTER INSERT OR UPDATE
    ON maker.vat_ilk_art
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_art();

CREATE TRIGGER ilk_spot
    AFTER INSERT OR UPDATE
    ON maker.vat_ilk_spot
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_spot();

CREATE TRIGGER ilk_line
    AFTER INSERT OR UPDATE
    ON maker.vat_ilk_line
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_line();

CREATE TRIGGER ilk_dust
    AFTER INSERT OR UPDATE
    ON maker.vat_ilk_dust
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_dust();

CREATE TRIGGER ilk_chop
    AFTER INSERT OR UPDATE
    ON maker.cat_ilk_chop
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_chop();

CREATE TRIGGER ilk_lump
    AFTER INSERT OR UPDATE
    ON maker.cat_ilk_lump
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_lump();

CREATE TRIGGER ilk_flip
    AFTER INSERT OR UPDATE
    ON maker.cat_ilk_flip
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_flip();

CREATE TRIGGER ilk_rho
    AFTER INSERT OR UPDATE
    ON maker.jug_ilk_rho
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_rho();

CREATE TRIGGER ilk_duty
    AFTER INSERT OR UPDATE
    ON maker.jug_ilk_duty
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_duty();

CREATE TRIGGER ilk_pip
    AFTER INSERT OR UPDATE
    ON maker.spot_ilk_pip
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_pip();

CREATE TRIGGER ilk_mat
    AFTER INSERT OR UPDATE
    ON maker.spot_ilk_mat
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_ilk_mat();


-- +goose Down
DROP TRIGGER ilk_mat ON maker.spot_ilk_mat;
DROP TRIGGER ilk_pip ON maker.spot_ilk_pip;
DROP TRIGGER ilk_duty ON maker.jug_ilk_duty;
DROP TRIGGER ilk_rho ON maker.jug_ilk_rho;
DROP TRIGGER ilk_flip ON maker.cat_ilk_flip;
DROP TRIGGER ilk_lump ON maker.cat_ilk_lump;
DROP TRIGGER ilk_chop ON maker.cat_ilk_chop;
DROP TRIGGER ilk_dust ON maker.vat_ilk_dust;
DROP TRIGGER ilk_line ON maker.vat_ilk_line;
DROP TRIGGER ilk_spot ON maker.vat_ilk_spot;
DROP TRIGGER ilk_art ON maker.vat_ilk_art;
DROP TRIGGER ilk_rate ON maker.vat_ilk_rate;

DROP FUNCTION maker.insert_ilk_mat();
DROP FUNCTION maker.insert_ilk_pip();
DROP FUNCTION maker.insert_ilk_duty();
DROP FUNCTION maker.insert_ilk_rho();
DROP FUNCTION maker.insert_ilk_flip();
DROP FUNCTION maker.insert_ilk_lump();
DROP FUNCTION maker.insert_ilk_chop();
DROP FUNCTION maker.insert_ilk_dust();
DROP FUNCTION maker.insert_ilk_line();
DROP FUNCTION maker.insert_ilk_spot();
DROP FUNCTION maker.insert_ilk_art();
DROP FUNCTION maker.insert_ilk_rate();

DROP TABLE api.current_ilk_state;