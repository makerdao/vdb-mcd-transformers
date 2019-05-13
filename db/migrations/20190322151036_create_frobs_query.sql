-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE maker.frob_event AS (
  ilk_name     TEXT,
  -- ilk object
  urn_id       TEXT,
  dink         NUMERIC,
  dart         NUMERIC,
  block_height BIGINT,
  tx_idx       INTEGER
  -- tx
);

COMMENT ON COLUMN maker.frob_event.block_height IS E'@omit';
COMMENT ON COLUMN maker.frob_event.tx_idx IS E'@omit';

CREATE OR REPLACE FUNCTION maker.urn_frobs(ilk_name TEXT, urn TEXT)
  RETURNS SETOF maker.frob_event AS
$body$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = $1),
    urn AS (
      SELECT id FROM maker.urns
      WHERE ilk_id = (SELECT id FROM ilk)
        AND guy = $2
    )

  SELECT $1 AS ilk_name, $2 AS urn_id, dink, dart, block_number AS block_height, tx_idx
  FROM maker.vat_frob LEFT JOIN headers ON vat_frob.header_id = headers.id
  WHERE vat_frob.urn_id = (SELECT id FROM urn)
  ORDER BY block_number DESC
$body$
LANGUAGE sql STABLE;


CREATE OR REPLACE FUNCTION maker.all_frobs(ilk_name TEXT)
  RETURNS SETOF maker.frob_event AS
$$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = $1)

  SELECT $1 AS ilk_name, guy AS urn_id, dink, dart, block_number AS block_height, tx_idx
  FROM maker.vat_frob
  LEFT JOIN maker.urns ON vat_frob.urn_id = urns.id
  LEFT JOIN headers    ON vat_frob.header_id = headers.id
  WHERE urns.ilk_id = (SELECT id FROM ilk)
  ORDER BY guy, block_number DESC
$$ LANGUAGE sql STABLE SECURITY DEFINER;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION IF EXISTS maker.urn_frobs(TEXT, TEXT);
DROP FUNCTION maker.all_frobs(TEXT);
DROP TYPE maker.frob_event CASCADE;