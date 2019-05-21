-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.frob_event AS (
  ilk_name     TEXT,
  -- ilk object
  urn_guy      TEXT,
  -- urn object
  dink         NUMERIC,
  dart         NUMERIC,
  block_height BIGINT,
  tx_idx       INTEGER
  -- tx
);

COMMENT ON COLUMN api.frob_event.block_height IS E'@omit';
COMMENT ON COLUMN api.frob_event.tx_idx IS E'@omit';

CREATE FUNCTION api.urn_frobs(_ilk_name TEXT, _urn TEXT)
  RETURNS SETOF api.frob_event AS
$$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = _ilk_name),
    urn AS (
      SELECT id FROM maker.urns
      WHERE ilk_id = (SELECT id FROM ilk)
        AND guy = _urn
    )

  SELECT _ilk_name AS ilk_name, _urn AS urn_id, dink, dart, block_number AS block_height, tx_idx
  FROM maker.vat_frob LEFT JOIN headers ON vat_frob.header_id = headers.id
  WHERE vat_frob.urn_id = (SELECT id FROM urn)
  ORDER BY block_number DESC
$$ LANGUAGE sql STABLE STRICT;


CREATE FUNCTION api.all_frobs(_ilk_name TEXT)
  RETURNS SETOF api.frob_event AS
$$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = _ilk_name)

  SELECT _ilk_name AS ilk_name, guy AS urn_id, dink, dart, block_number AS block_height, tx_idx
  FROM maker.vat_frob
  LEFT JOIN maker.urns ON vat_frob.urn_id = urns.id
  LEFT JOIN headers    ON vat_frob.header_id = headers.id
  WHERE urns.ilk_id = (SELECT id FROM ilk)
  ORDER BY guy, block_number DESC
$$ LANGUAGE sql STABLE STRICT;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.urn_frobs(TEXT, TEXT);
DROP FUNCTION api.all_frobs(TEXT);
DROP TYPE api.frob_event CASCADE;