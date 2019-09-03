-- +goose Up
CREATE TABLE api.managed_cdp
(
    id             SERIAL PRIMARY KEY,
    cdpi           NUMERIC   DEFAULT NULL UNIQUE,
    usr            TEXT      DEFAULT NULL,
    urn_identifier TEXT      DEFAULT NULL,
    ilk_identifier TEXT      DEFAULT NULL,
    created        TIMESTAMP DEFAULT NULL
);

COMMENT ON TABLE api.managed_cdp IS '@omit create,update,delete';
COMMENT ON COLUMN api.managed_cdp.id IS '@omit';
COMMENT ON COLUMN api.managed_cdp.cdpi IS '@name id';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_cdp_created() RETURNS TRIGGER
AS
$$
BEGIN
    WITH block_info AS (
        SELECT api.epoch_to_datetime(headers.block_timestamp) AS datetime
        FROM public.headers
        WHERE headers.block_number = NEW.block_number
        LIMIT 1)
    INSERT
    INTO api.managed_cdp (cdpi, created)
    VALUES (NEW.cdpi, (SELECT datetime FROM block_info))
    ON CONFLICT (cdpi)
        DO UPDATE SET created = (SELECT datetime FROM block_info)
    WHERE managed_cdp.created IS NULL;
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_cdp_usr() RETURNS TRIGGER
AS
$$
BEGIN
    INSERT
    INTO api.managed_cdp (cdpi, usr)
    VALUES (NEW.cdpi, NEW.owner)
    -- only update usr if the new owner is coming from the latest owns block we know about for the given cdpi
    ON CONFLICT (cdpi)
        DO UPDATE SET usr = NEW.owner
    WHERE NEW.block_number >= (
        SELECT MAX(block_number)
        FROM maker.cdp_manager_owns
        WHERE cdp_manager_owns.cdpi = NEW.cdpi);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_cdp_urn_identifier() RETURNS TRIGGER
AS
$$
BEGIN
    INSERT
    INTO api.managed_cdp (cdpi, urn_identifier)
    VALUES (NEW.cdpi, NEW.urn)
    ON CONFLICT (cdpi) DO UPDATE SET urn_identifier = NEW.urn;
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_cdp_ilk_identifier() RETURNS TRIGGER
AS
$$
BEGIN
    WITH ilk AS (
        SELECT ilks.identifier
        FROM maker.cdp_manager_ilks
                 LEFT JOIN maker.ilks ON ilks.id = cdp_manager_ilks.ilk_id
        WHERE cdp_manager_ilks.cdpi = NEW.cdpi
    )
    INSERT
    INTO api.managed_cdp (cdpi, ilk_identifier)
    VALUES (NEW.cdpi, (SELECT identifier FROM ilk))
    ON CONFLICT (cdpi) DO UPDATE SET ilk_identifier = (SELECT identifier FROM ilk);
    RETURN NEW;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER managed_cdp_cdpi
    AFTER INSERT OR UPDATE
    ON maker.cdp_manager_cdpi
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_cdp_created();

CREATE TRIGGER managed_cdp_usr
    AFTER INSERT OR UPDATE
    ON maker.cdp_manager_owns
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_cdp_usr();

CREATE TRIGGER managed_cdp_urn
    AFTER INSERT OR UPDATE
    ON maker.cdp_manager_urns
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_cdp_urn_identifier();

CREATE TRIGGER managed_cdp_ilk
    AFTER INSERT OR UPDATE
    ON maker.cdp_manager_ilks
    FOR EACH ROW
EXECUTE PROCEDURE maker.insert_cdp_ilk_identifier();

-- +goose Down
DROP TRIGGER managed_cdp_ilk ON maker.cdp_manager_ilks;
DROP TRIGGER managed_cdp_urn ON maker.cdp_manager_urns;
DROP TRIGGER managed_cdp_usr ON maker.cdp_manager_owns;
DROP TRIGGER managed_cdp_cdpi ON maker.cdp_manager_cdpi;

DROP FUNCTION maker.insert_cdp_ilk_identifier();
DROP FUNCTION maker.insert_cdp_urn_identifier();
DROP FUNCTION maker.insert_cdp_usr();
DROP FUNCTION maker.insert_cdp_created();

DROP TABLE api.managed_cdp;