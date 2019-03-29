-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE maker.frob_event AS (
  ilkId TEXT,
  -- ilk object
  urnId TEXT,
  dink  NUMERIC,
  dart  NUMERIC
  -- tx
);

CREATE OR REPLACE FUNCTION maker.frobs(ilk TEXT, urn TEXT)
  RETURNS SETOF maker.frob_event AS
$body$
WITH
  ilk AS (SELECT id FROM maker.ilks WHERE ilks.ilk = $1),
  urn AS (
    SELECT id FROM maker.urns
    WHERE ilk_id = (SELECT id FROM ilk)
      AND guy = $2
  )

SELECT $1 AS ilkId, $2 AS urnId, dink, dart
  FROM maker.vat_frob LEFT JOIN headers ON vat_frob.header_id = headers.id
WHERE vat_frob.urn_id = (SELECT id FROM urn)
ORDER BY block_number DESC

$body$
LANGUAGE SQL;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION IF EXISTS maker.frobs(TEXT, TEXT);
DROP TYPE maker.frob_event CASCADE;