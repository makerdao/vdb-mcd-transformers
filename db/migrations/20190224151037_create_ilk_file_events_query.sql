-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.ilk_file_event AS
(
    ilk_identifier TEXT,
    -- ilk
    what           TEXT,
    data           TEXT,
    block_height   BIGINT,
    log_id         BIGINT
    -- tx
);

COMMENT ON TYPE api.ilk_file_event
    IS E'File event associated with an Ilk emitted by Cat, Jug, Spot, or Vat contract, with nested data regarding associated Ilk and Tx.';

COMMENT ON COLUMN api.ilk_file_event.ilk_identifier
    IS E'@omit';
COMMENT ON COLUMN api.ilk_file_event.block_height
    IS E'@omit';
COMMENT ON COLUMN api.ilk_file_event.log_id
    IS E'@omit';

CREATE FUNCTION api.all_ilk_file_events(ilk_identifier TEXT, max_results INTEGER DEFAULT -1,
                                        result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.ilk_file_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier)

SELECT ilk_identifier, what, data :: text, block_number, log_id
FROM maker.cat_file_chop_lump
         LEFT JOIN headers ON cat_file_chop_lump.header_id = headers.id
WHERE cat_file_chop_lump.ilk_id = (SELECT id FROM ilk)
UNION
SELECT ilk_identifier, what, flip AS data, block_number, log_id
FROM maker.cat_file_flip
         LEFT JOIN headers ON cat_file_flip.header_id = headers.id
WHERE cat_file_flip.ilk_id = (SELECT id FROM ilk)
UNION
SELECT ilk_identifier, what, data :: text, block_number, log_id
FROM maker.jug_file_ilk
         LEFT JOIN headers ON jug_file_ilk.header_id = headers.id
WHERE jug_file_ilk.ilk_id = (SELECT id FROM ilk)
UNION
SELECT ilk_identifier, what, data :: text, block_number, log_id
FROM maker.spot_file_mat
         LEFT JOIN headers ON spot_file_mat.header_id = headers.id
WHERE spot_file_mat.ilk_id = (SELECT id FROM ilk)
UNION
SELECT ilk_identifier, what, pip AS data, block_number, log_id
FROM maker.spot_file_pip
         LEFT JOIN headers ON spot_file_pip.header_id = headers.id
WHERE spot_file_pip.ilk_id = (SELECT id FROM ilk)
UNION
SELECT ilk_identifier, what, data :: text, block_number, log_id
FROM maker.vat_file_ilk
         LEFT JOIN headers ON vat_file_ilk.header_id = headers.id
WHERE vat_file_ilk.ilk_id = (SELECT id FROM ilk)
ORDER BY block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
OFFSET
all_ilk_file_events.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;

COMMENT ON FUNCTION api.all_ilk_file_events(ilk_identifier TEXT, max_results INTEGER, result_offset INTEGER)
    IS E'Get all File events associated with an Ilk. ilkIdentifier (e.g. "ETH-A") is required. maxResults and resultOffset are optional, defaulting to no max/offset.';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.all_ilk_file_events(TEXT, INTEGER, INTEGER);
DROP TYPE api.ilk_file_event CASCADE;