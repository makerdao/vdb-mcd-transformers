-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE maker.frob_event AS (
  ilkId        TEXT,
  -- ilk object
  urnId        TEXT,
  dink         NUMERIC,
  dart         NUMERIC,
  block_number BIGINT
  -- tx
);


CREATE OR REPLACE FUNCTION maker.urn_frobs(ilk TEXT, urn TEXT)
  RETURNS SETOF maker.frob_event AS
$body$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.ilk = $1),
    urn AS (
      SELECT id FROM maker.urns
      WHERE ilk_id = (SELECT id FROM ilk)
        AND guy = $2
    )

  SELECT $1 AS ilkId, $2 AS urnId, dink, dart, block_number
  FROM maker.vat_frob LEFT JOIN headers ON vat_frob.header_id = headers.id
  WHERE vat_frob.urn_id = (SELECT id FROM urn)
  ORDER BY block_number DESC
$body$
LANGUAGE sql STABLE;


CREATE OR REPLACE FUNCTION maker.all_frobs(ilk TEXT)
  RETURNS SETOF maker.frob_event AS
$$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.ilk = $1)

  SELECT $1 AS ilkId, guy AS urnId, dink, dart, block_number
  FROM maker.vat_frob
  LEFT JOIN maker.urns ON vat_frob.urn_id = urns.id
  LEFT JOIN headers    ON vat_frob.header_id = headers.id
  WHERE urns.ilk_id = (SELECT id FROM ilk)
  ORDER BY guy, block_number DESC
$$ LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION IF EXISTS maker.urn_frobs(TEXT, TEXT);
DROP FUNCTION maker.all_frobs(TEXT);
DROP TYPE maker.frob_event CASCADE;