-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.bite_event AS
(
    ilk_identifier TEXT,
    -- ilk object
    urn_identifier TEXT,
    -- urn object
    bid_id         NUMERIC,
    -- bid object
    ink            NUMERIC,
    art            NUMERIC,
    tab            NUMERIC,
    block_height   BIGINT,
    log_id         BIGINT
    -- tx
);

COMMENT ON TYPE api.bite_event
    IS E'Bite event emitted by Cat contract, with nested data regarding associated Ilk, Urn, Bid, and Tx.';

COMMENT ON COLUMN api.bite_event.block_height
    IS E'@omit';
COMMENT ON COLUMN api.bite_event.log_id
    IS E'@omit';

CREATE FUNCTION api.all_bites(ilk_identifier TEXT, max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier)

SELECT ilk_identifier,
       identifier AS urn_identifier,
       bid_id,
       ink,
       art,
       tab,
       block_number,
       log_id
FROM maker.bite
         LEFT JOIN maker.urns ON bite.urn_id = urns.id
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE urns.ilk_id = (SELECT id FROM ilk)
ORDER BY urn_identifier, block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
OFFSET
all_bites.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;

COMMENT ON FUNCTION api.all_bites(ilk_identifier TEXT, max_results INTEGER, result_offset INTEGER)
    IS E'Get Bite events associated with a given Ilk. ilkIdentifier (e.g. "ETH-A") is required. maxResults and resultOffset are optional, defaulting to no max/offset.';


CREATE FUNCTION api.urn_bites(ilk_identifier TEXT, urn_identifier TEXT, max_results INTEGER DEFAULT -1,
                              result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier),
     urn AS (SELECT id
             FROM maker.urns
             WHERE ilk_id = (SELECT id FROM ilk)
               AND identifier = urn_bites.urn_identifier)

SELECT ilk_identifier,
       urn_bites.urn_identifier,
       bid_id,
       ink,
       art,
       tab,
       block_number,
       log_id
FROM maker.bite
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE bite.urn_id = (SELECT id FROM urn)
ORDER BY block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
OFFSET
urn_bites.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;

COMMENT ON FUNCTION api.urn_bites(ilk_identifier TEXT, urn_identifier TEXT, max_results INTEGER, result_offset INTEGER)
    IS E'Get Bite events associated with a given Urn. ilkIdentifier (e.g. "ETH-A") and urnIdentifier (e.g. "0xC93C178EC17B06bddBa0CC798546161aF9D25e8A") are required. maxResults and resultOffset are optional, defaulting to no max/offset.';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.urn_bites(TEXT, TEXT, INTEGER, INTEGER);
DROP FUNCTION api.all_bites(TEXT, INTEGER, INTEGER);
DROP TYPE api.bite_event CASCADE;
