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

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.managed_cdp_ilk(api.managed_cdp);
