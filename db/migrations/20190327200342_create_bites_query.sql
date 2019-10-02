-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.bite_event AS (
    ilk_identifier TEXT,
    -- ilk object
    urn_identifier TEXT,
    -- urn object
    bid_id NUMERIC,
    -- bid object
    ink NUMERIC,
    art NUMERIC,
    tab NUMERIC,
    block_height BIGINT,
    log_id BIGINT
    -- tx
    );

COMMENT ON COLUMN api.bite_event.block_height
    IS E'@omit';
COMMENT ON COLUMN api.bite_event.log_id
    IS E'@omit';

CREATE FUNCTION api.all_bites(ilk_identifier TEXT, max_results INTEGER DEFAULT NULL, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier)

SELECT ilk_identifier, identifier AS urn_identifier, bid_id, ink, art, tab, block_number, log_id
FROM maker.bite
         LEFT JOIN maker.urns ON bite.urn_id = urns.id
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE urns.ilk_id = (SELECT id FROM ilk)
ORDER BY urn_identifier, block_number DESC
LIMIT all_bites.max_results OFFSET all_bites.result_offset
$$
    LANGUAGE sql
    STABLE;


CREATE FUNCTION api.urn_bites(ilk_identifier TEXT, urn_identifier TEXT, max_results INTEGER DEFAULT NULL,
                              result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier),
     urn AS (SELECT id
             FROM maker.urns
             WHERE ilk_id = (SELECT id FROM ilk)
               AND identifier = urn_bites.urn_identifier)

SELECT ilk_identifier, urn_bites.urn_identifier, bid_id, ink, art, tab, block_number, log_id
FROM maker.bite
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE bite.urn_id = (SELECT id FROM urn)
ORDER BY block_number DESC
LIMIT urn_bites.max_results OFFSET urn_bites.result_offset
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.urn_bites(TEXT, TEXT, INTEGER, INTEGER);
DROP FUNCTION api.all_bites(TEXT, INTEGER, INTEGER);
DROP TYPE api.bite_event CASCADE;
