-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.frob_event AS (
    ilk_identifier TEXT,
    -- ilk object
    urn_identifier TEXT,
    -- urn object
    dink NUMERIC,
    dart NUMERIC,
    ilk_rate NUMERIC,
    block_height BIGINT,
    log_id BIGINT
    -- tx
    );

COMMENT ON COLUMN api.frob_event.block_height
    IS E'@omit';
COMMENT ON COLUMN api.frob_event.log_id
    IS E'@omit';

CREATE FUNCTION api.urn_frobs(ilk_identifier TEXT, urn_identifier TEXT, max_results INTEGER DEFAULT -1,
                              result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.frob_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier),
     urn AS (SELECT id
             FROM maker.urns
             WHERE ilk_id = (SELECT id FROM ilk)
               AND identifier = urn_identifier),
     rates AS (SELECT block_number, rate
               FROM maker.vat_ilk_rate
               WHERE ilk_id = (SELECT id FROM ilk)
               ORDER BY block_number DESC
     )

SELECT ilk_identifier,
       urn_identifier,
       dink,
       dart,
       (SELECT rate from rates WHERE block_number <= headers.block_number LIMIT 1) AS ilk_rate,
       headers.block_number,
       log_id
FROM maker.vat_frob
         LEFT JOIN headers ON vat_frob.header_id = headers.id
WHERE vat_frob.urn_id = (SELECT id FROM urn)
ORDER BY block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
OFFSET
urn_frobs.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;


CREATE FUNCTION api.all_frobs(ilk_identifier TEXT, max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.frob_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier),
     rates AS (SELECT block_number, rate
               FROM maker.vat_ilk_rate
               WHERE ilk_id = (SELECT id FROM ilk)
               ORDER BY block_number DESC
     )

SELECT ilk_identifier,
       urns.identifier                                                             AS urn_identifier,
       dink,
       dart,
       (SELECT rate from rates WHERE block_number <= headers.block_number LIMIT 1) AS ilk_rate,
       block_number,
       log_id
FROM maker.vat_frob
         LEFT JOIN maker.urns ON vat_frob.urn_id = urns.id
         LEFT JOIN headers ON vat_frob.header_id = headers.id
WHERE urns.ilk_id = (SELECT id FROM ilk)
ORDER BY block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
OFFSET
all_frobs.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.urn_frobs(TEXT, TEXT, INTEGER, INTEGER);
DROP FUNCTION api.all_frobs(TEXT, INTEGER, INTEGER);
DROP TYPE api.frob_event CASCADE;