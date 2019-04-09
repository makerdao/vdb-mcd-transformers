-- +goose Up
create type maker.relevant_block AS (
  block_number bigint,
  block_hash   text,
  ilk_id       integer
);

create or replace function maker.get_ilk_blocks_before(block_number bigint, ilk_id int)
  returns setof maker.relevant_block as $$
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.vat_ilk_rate
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.vat_ilk_art
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.vat_ilk_spot
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.vat_ilk_line
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.vat_ilk_dust
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.cat_ilk_chop
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.cat_ilk_lump
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.cat_ilk_flip
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.jug_ilk_rho
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.jug_ilk_tax
WHERE block_number <= $1
      AND ilk_id = $2
$$
LANGUAGE sql STABLE;

CREATE TYPE maker.ilk_state AS (
  ilk_id       integer,
  ilk          text,
  block_number bigint,
  rate         numeric,
  art          numeric,
  spot         numeric,
  line         numeric,
  dust         numeric,
  chop         numeric,
  lump         numeric,
  flip         text,
  rho          numeric,
  tax          numeric,
  created      numeric,
  updated      numeric
);

CREATE FUNCTION maker.get_ilk_at_block_number(block_number bigint, ilkId int)
  RETURNS maker.ilk_state
AS $$
WITH rates AS (
    SELECT
      rate,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_rate
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), arts AS (
    SELECT
      art,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_art
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), spots AS (
    SELECT
      spot,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_spot
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), lines AS (
    SELECT
      line,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_line
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), dusts AS (
    SELECT
      dust,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_dust
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), chops AS (
    SELECT
      chop,
      ilk_id,
      block_hash
    FROM maker.cat_ilk_chop
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), lumps AS (
    SELECT
      lump,
      ilk_id,
      block_hash
    FROM maker.cat_ilk_lump
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), flips AS (
    SELECT
      flip,
      ilk_id,
      block_hash
    FROM maker.cat_ilk_flip
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), rhos AS (
    SELECT
      rho,
      ilk_id,
      block_hash
    FROM maker.jug_ilk_rho
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), taxes AS (
    SELECT
      tax,
      ilk_id,
      block_hash
    FROM maker.jug_ilk_tax
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), relevant_blocks AS (
  SELECT * FROM maker.get_ilk_blocks_before($1, $2)
), created AS (
    SELECT DISTINCT ON (relevant_blocks.ilk_id, relevant_blocks.block_number)
      relevant_blocks.block_number,
      relevant_blocks.block_hash,
      relevant_blocks.ilk_id,
      headers.block_timestamp
    FROM relevant_blocks
      LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
    ORDER BY block_number ASC
    LIMIT 1
), updated AS (
    SELECT DISTINCT ON (relevant_blocks.ilk_id, relevant_blocks.block_number)
      relevant_blocks.block_number,
      relevant_blocks.block_hash,
      relevant_blocks.ilk_id,
      headers.block_timestamp
    FROM relevant_blocks
      LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
    ORDER BY block_number DESC
    LIMIT 1
)

SELECT
  ilks.id,
  ilks.ilk,
  $1 block_number,
  rates.rate,
  arts.art,
  spots.spot,
  lines.line,
  dusts.dust,
  chops.chop,
  lumps.lump,
  flips.flip,
  rhos.rho,
  taxes.tax,
  created.block_timestamp,
  updated.block_timestamp
FROM maker.ilks AS ilks
  LEFT JOIN rates ON rates.ilk_id = ilks.id
  LEFT JOIN arts ON arts.ilk_id = ilks.id
  LEFT JOIN spots ON spots.ilk_id = ilks.id
  LEFT JOIN lines ON lines.ilk_id = ilks.id
  LEFT JOIN dusts ON dusts.ilk_id = ilks.id
  LEFT JOIN chops ON chops.ilk_id = ilks.id
  LEFT JOIN lumps ON lumps.ilk_id = ilks.id
  LEFT JOIN flips ON flips.ilk_id = ilks.id
  LEFT JOIN rhos ON rhos.ilk_id = ilks.id
  LEFT JOIN taxes ON taxes.ilk_id = ilks.id
  LEFT JOIN created ON created.ilk_id = ilks.id
  LEFT JOIN updated ON updated.ilk_id = ilks.id
WHERE (
  rates.rate is not null OR
  arts.art is not null OR
  spots.spot is not null OR
  lines.line is not null OR
  dusts.dust is not null OR
  chops.chop is not null OR
  lumps.lump is not null OR
  flips.flip is not null OR
  rhos.rho is not null OR
  taxes.tax is not null
)
$$
LANGUAGE SQL
STABLE;

-- +goose Down
DROP FUNCTION IF EXISTS maker.get_ilk_blocks_before(block_number bigint, ilk_id int);
DROP TYPE maker.relevant_block CASCADE;
DROP FUNCTION IF EXISTS maker.get_ilk_at_block_number(block_number bigint, ilk_id int );
DROP TYPE maker.ilk_state CASCADE;