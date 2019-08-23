-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TYPE api.managed_cdp AS (
    usr TEXT,
    id NUMERIC,
    urn_identifier TEXT,
    ilk_identifier TEXT,
    created TIMESTAMP
    );

COMMENT ON COLUMN api.managed_cdp.urn_identifier
    IS E'@name urn_id';
COMMENT ON COLUMN api.managed_cdp.ilk_identifier
    IS E'@name ilk_id';

-- Function returning properties of a single CDP by ID
CREATE FUNCTION api.get_managed_cdp(id NUMERIC)
    RETURNS api.managed_cdp
AS
$$
WITH owners AS (
    SELECT cdp_manager_owns.owner, cdpi
    FROM maker.cdp_manager_owns
    WHERE cdpi = get_managed_cdp.id
    ORDER BY cdp_manager_owns.block_number DESC
    LIMIT 1),
     ilk AS (
         SELECT ilks.identifier, cdp_manager_ilks.cdpi
         FROM maker.cdp_manager_ilks
                  LEFT JOIN maker.ilks ON ilks.id = cdp_manager_ilks.ilk_id
         WHERE cdp_manager_ilks.cdpi = get_managed_cdp.id
         ORDER BY cdp_manager_ilks.block_number DESC
         LIMIT 1),
     urn AS (
         SELECT cdp_manager_urns.urn AS identifier, cdp_manager_urns.cdpi
         FROM maker.cdp_manager_urns
         WHERE cdp_manager_urns.cdpi = get_managed_cdp.id
         ORDER BY cdp_manager_urns.block_number DESC
         LIMIT 1),
     created AS (
         SELECT api.epoch_to_datetime(headers.block_timestamp) AS datetime, cdp_manager_cdpi.cdpi
         FROM headers
                  LEFT JOIN maker.cdp_manager_cdpi ON cdp_manager_cdpi.block_number = headers.block_number
         WHERE cdp_manager_cdpi.cdpi = get_managed_cdp.id
         LIMIT 1)
SELECT owners.owner     AS usr,
       get_managed_cdp.id,
       urn.identifier   AS urn_identifier,
       ilk.identifier   AS ilk_identifier,
       created.datetime AS created
FROM owners
         LEFT JOIN ilk ON ilk.cdpi = owners.cdpi
         LEFT JOIN urn ON urn.cdpi = owners.cdpi
         LEFT JOIN created ON created.cdpi = owners.cdpi
$$
    LANGUAGE SQL
    STABLE
    STRICT;

-- +goose Down
DROP FUNCTION api.get_managed_cdp(NUMERIC);
DROP TYPE api.managed_cdp CASCADE;
