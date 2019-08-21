-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE OR REPLACE FUNCTION api.all_managed_cdps() RETURNS SETOF api.managed_cdp AS
-- +goose StatementBegin
$BODY$
BEGIN
    RETURN QUERY (
        WITH cdpis AS (
            SELECT DISTINCT cdpi
            FROM maker.cdp_manager_cdpi
            ORDER BY cdpi)
        SELECT cdp.*
        FROM cdpis,
             LATERAL api.get_managed_cdp(cdpis.cdpi) cdp
    );
END
$BODY$
    LANGUAGE plpgsql
    STABLE;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.all_managed_cdps();
