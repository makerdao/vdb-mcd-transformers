-- +goose Up

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_urn_created(urn_id INTEGER) RETURNS maker.vat_urn_ink
AS
$$
BEGIN
    WITH utc AS (select urn_time_created(urn_id) as utc)
    UPDATE api.urn_snapshot
    SET created = (SELECT utc FROM utc)
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

-- +goose Down

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
