-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend managed_cdp with ilk_snapshot
CREATE FUNCTION api.managed_cdp_ilk(cdp api.managed_cdp)
    RETURNS api.ilk_snapshot AS
$$
SELECT *
FROM api.ilk_snapshot i
WHERE i.ilk_identifier = cdp.ilk_identifier
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql
    STABLE;

--- Extend managed_cdp with urn_snapshot
CREATE FUNCTION api.managed_cdp_urn(cdp api.managed_cdp) RETURNS api.urn_snapshot
       LANGUAGE sql STABLE
       AS $$
SELECT *
FROM api.get_urn(cdp.ilk_identifier, cdp.urn_identifier, api.max_block())
$$;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.managed_cdp_urn(api.managed_cdp);
DROP FUNCTION api.managed_cdp_ilk(api.managed_cdp);
