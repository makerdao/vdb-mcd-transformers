-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.ilk_file_event AS (
  ilk_identifier  TEXT,
  -- ilk
  what            TEXT,
  data            TEXT,
  block_height    BIGINT,
  tx_idx          INTEGER
  -- tx
);

COMMENT ON COLUMN api.ilk_file_event.ilk_identifier IS E'@omit';
COMMENT ON COLUMN api.ilk_file_event.block_height IS E'@omit';
COMMENT ON COLUMN api.ilk_file_event.tx_idx IS E'@omit';

CREATE FUNCTION api.all_ilk_file_events(ilk_identifier TEXT)
  RETURNS SETOF api.ilk_file_event AS
$$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = ilk_identifier)

  SELECT ilk_identifier, what, data::text, block_number, tx_idx
  FROM maker.cat_file_chop_lump
  LEFT JOIN headers ON cat_file_chop_lump.header_id = headers.id
  WHERE cat_file_chop_lump.ilk_id = (SELECT id FROM ilk)
  UNION
  SELECT ilk_identifier, what, flip AS data, block_number, tx_idx
  FROM maker.cat_file_flip
  LEFT JOIN headers ON cat_file_flip.header_id = headers.id
  WHERE cat_file_flip.ilk_id = (SELECT id FROM ilk)
  UNION
  SELECT ilk_identifier, what, data::text, block_number, tx_idx
  FROM maker.jug_file_ilk
  LEFT JOIN headers ON jug_file_ilk.header_id = headers.id
  WHERE jug_file_ilk.ilk_id = (SELECT id FROM ilk)
  UNION
  SELECT ilk_identifier, what, data::text, block_number, tx_idx
  FROM maker.vat_file_ilk
  LEFT JOIN headers ON vat_file_ilk.header_id = headers.id
  WHERE vat_file_ilk.ilk_id = (SELECT id FROM ilk)
  ORDER BY block_number DESC
$$ LANGUAGE sql STABLE STRICT;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.all_ilk_file_events(TEXT);
DROP TYPE api.ilk_file_event CASCADE;