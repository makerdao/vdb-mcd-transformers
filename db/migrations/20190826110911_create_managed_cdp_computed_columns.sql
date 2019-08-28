-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend managed_cdp with ilk_state
CREATE FUNCTION api.managed_cdp_ilk(cdp api.managed_cdp)
    RETURNS api.ilk_state AS
$$
SELECT *
FROM api.get_ilk(cdp.ilk_identifier)
$$
    LANGUAGE sql
    STABLE;

-- Extend managed_cdp with urn_state
CREATE FUNCTION api.managed_cdp_urn(cdp api.managed_cdp)
    RETURNS SETOF api.urn_state AS
$$
SELECT *
FROM api.get_urn(cdp.ilk_identifier, cdp.urn_identifier)
$$
    LANGUAGE sql
    STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.managed_cdp_urn(api.managed_cdp);
DROP FUNCTION api.managed_cdp_ilk(api.managed_cdp);
