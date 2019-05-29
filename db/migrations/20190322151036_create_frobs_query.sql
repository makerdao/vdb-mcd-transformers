-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.frob_event AS (
  ilk_identifier TEXT,
  -- ilk object
  urn_guy        TEXT,
  -- urn object
  dink           NUMERIC,
  dart           NUMERIC,
  block_height   BIGINT,
  tx_idx         INTEGER
  -- tx
);

COMMENT ON COLUMN api.frob_event.block_height
IS E'@omit';
COMMENT ON COLUMN api.frob_event.tx_idx
IS E'@omit';

CREATE FUNCTION api.urn_frobs(ilk_identifier TEXT, urn_guy TEXT)
  RETURNS SETOF api.frob_event AS
$$
WITH
    ilk AS (SELECT id
            FROM maker.ilks
            WHERE ilks.identifier = ilk_identifier),
    urn AS (
      SELECT id
      FROM maker.urns
      WHERE ilk_id = (SELECT id
                      FROM ilk)
            AND guy = urn_guy
  )

SELECT
  ilk_identifier,
  urn_guy,
  dink,
  dart,
  block_number,
  tx_idx
FROM maker.vat_frob
  LEFT JOIN headers ON vat_frob.header_id = headers.id
WHERE vat_frob.urn_id = (SELECT id
                         FROM urn)
ORDER BY block_number DESC
$$
LANGUAGE sql
STABLE
STRICT;


CREATE FUNCTION api.all_frobs(ilk_identifier TEXT)
  RETURNS SETOF api.frob_event AS
$$
WITH
    ilk AS (SELECT id
            FROM maker.ilks
            WHERE ilks.identifier = ilk_identifier)

SELECT
  ilk_identifier,
  guy AS urn_id,
  dink,
  dart,
  block_number,
  tx_idx
FROM maker.vat_frob
  LEFT JOIN maker.urns ON vat_frob.urn_id = urns.id
  LEFT JOIN headers ON vat_frob.header_id = headers.id
WHERE urns.ilk_id = (SELECT id
                     FROM ilk)
ORDER BY guy, block_number DESC
$$
LANGUAGE sql
STABLE
STRICT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.urn_frobs( TEXT, TEXT );
DROP FUNCTION api.all_frobs( TEXT );
DROP TYPE api.frob_event CASCADE;