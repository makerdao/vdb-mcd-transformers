-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.file_event AS (
  id            TEXT,
  ilk_name      TEXT,
  -- ilk
  what          TEXT,
  data          TEXT,
  block_height  BIGINT,
  tx_idx        INTEGER
  -- tx
);

COMMENT ON COLUMN api.file_event.block_height IS E'@omit';
COMMENT ON COLUMN api.file_event.tx_idx IS E'@omit';

CREATE FUNCTION api.ilk_files(ilk_name TEXT)
  RETURNS SETOF api.file_event AS
$body$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = $1)

  SELECT cat_file_chop_lump.raw_log::json->>'address' AS id, $1 AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.cat_file_chop_lump
  LEFT JOIN headers ON cat_file_chop_lump.header_id = headers.id
  WHERE cat_file_chop_lump.ilk_id = (SELECT id FROM ilk)
  UNION
  SELECT cat_file_flip.raw_log::json->>'address' AS id, $1 AS ilk_name, what, flip AS data, block_number AS block_height, tx_idx
  FROM maker.cat_file_flip
  LEFT JOIN headers ON cat_file_flip.header_id = headers.id
  WHERE cat_file_flip.ilk_id = (SELECT id FROM ilk)
  UNION
  SELECT jug_file_ilk.raw_log::json->>'address' AS id, $1 AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.jug_file_ilk
  LEFT JOIN headers ON jug_file_ilk.header_id = headers.id
  WHERE jug_file_ilk.ilk_id = (SELECT id FROM ilk)
  UNION
  SELECT vat_file_ilk.raw_log::json->>'address' AS id, $1 AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.vat_file_ilk
  LEFT JOIN headers ON vat_file_ilk.header_id = headers.id
  WHERE vat_file_ilk.ilk_id = (SELECT id FROM ilk)
  ORDER BY block_height DESC
$body$
LANGUAGE sql STABLE STRICT;

-- expensive query, hitting lots of tables
-- probably can narrow it down for specific contracts
CREATE FUNCTION api.address_files(address TEXT)
  RETURNS SETOF api.file_event AS
$body$
  WITH
    lowerAddress AS (SELECT lower($1))

-- ilk files
  SELECT cat_file_chop_lump.raw_log::json->>'address' AS id, ilks.name AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.cat_file_chop_lump
  LEFT JOIN maker.ilks ON cat_file_chop_lump.ilk_id = ilks.id
  LEFT JOIN headers    ON cat_file_chop_lump.header_id = headers.id
  WHERE lower(cat_file_chop_lump.raw_log::json->>'address') = (SELECT * FROM lowerAddress)
  UNION
  SELECT cat_file_flip.raw_log::json->>'address' AS id, ilks.name AS ilk_name, what, flip AS data, block_number AS block_height, tx_idx
  FROM maker.cat_file_flip
  LEFT JOIN maker.ilks ON cat_file_flip.ilk_id = ilks.id
  LEFT JOIN headers ON cat_file_flip.header_id = headers.id
  WHERE lower(cat_file_flip.raw_log::json->>'address') = (SELECT * FROM lowerAddress)
  UNION
  SELECT jug_file_ilk.raw_log::json->>'address' AS id, ilks.name AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.jug_file_ilk
  LEFT JOIN maker.ilks ON jug_file_ilk.ilk_id = ilks.id
  LEFT JOIN headers ON jug_file_ilk.header_id = headers.id
  WHERE lower(jug_file_ilk.raw_log::json->>'address') = (SELECT * FROM lowerAddress)
  UNION
  SELECT vat_file_ilk.raw_log::json->>'address' AS id, ilks.name AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.vat_file_ilk
  LEFT JOIN maker.ilks ON vat_file_ilk.ilk_id = ilks.id
  LEFT JOIN headers ON vat_file_ilk.header_id = headers.id
  WHERE lower(vat_file_ilk.raw_log::json->>'address') = (SELECT * FROM lowerAddress)

-- contract files
  UNION
  SELECT cat_file_vow.raw_log::json->>'address' AS id, NULL AS ilk_name, what, data, block_number AS block_height, tx_idx
  FROM maker.cat_file_vow
  LEFT JOIN headers ON cat_file_vow.header_id = headers.id
  WHERE lower(cat_file_vow.raw_log::json->>'address') = (SELECT * FROM lowerAddress)
  UNION
  SELECT jug_file_base.raw_log::json->>'address' AS id, NULL AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.jug_file_base
  LEFT JOIN headers ON jug_file_base.header_id = headers.id
  WHERE lower(jug_file_base.raw_log::json->>'address') = (SELECT * FROM lowerAddress)
  UNION
  SELECT jug_file_vow.raw_log::json->>'address' AS id, NULL AS ilk_name, what, data, block_number AS block_height, tx_idx
  FROM maker.jug_file_vow
  LEFT JOIN headers ON jug_file_vow.header_id = headers.id
  WHERE lower(jug_file_vow.raw_log::json->>'address') = (SELECT * FROM lowerAddress)
  UNION
  SELECT vat_file_debt_ceiling.raw_log::json->>'address' AS id, NULL AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.vat_file_debt_ceiling
  LEFT JOIN headers on vat_file_debt_ceiling.header_id = headers.id
  WHERE lower(vat_file_debt_ceiling.raw_log::json->>'address') = (SELECT * FROM lowerAddress)

  ORDER BY block_height DESC
$body$
LANGUAGE sql STABLE STRICT;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.ilk_files(TEXT);
DROP FUNCTION api.address_files(TEXT);
DROP TYPE api.file_event CASCADE;