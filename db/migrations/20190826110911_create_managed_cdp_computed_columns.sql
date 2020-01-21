-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend managed_cdp with ilk_state
CREATE FUNCTION api.managed_cdp_ilk(cdp api.managed_cdp)
    RETURNS api.historical_ilk_state AS
$$
SELECT *
FROM api.historical_ilk_state i
WHERE i.ilk_identifier = cdp.ilk_identifier
ORDER BY block_number DESC
LIMIT 1
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
