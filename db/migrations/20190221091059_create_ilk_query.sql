-- +goose Up
create type maker.relevant_block AS (
  block_number bigint,
  block_hash   text,
  ilk          integer
);

create or replace function maker.get_ilk_blocks_before(block_number numeric, ilk_id int)
  returns setof maker.relevant_block as $$
SELECT
  block_number,
  block_hash,
  ilk
FROM maker.vat_ilk_take
WHERE block_number <= $1
      AND ilk = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk
FROM maker.vat_ilk_rate
WHERE block_number <= $1
      AND ilk = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk
FROM maker.vat_ilk_ink
WHERE block_number <= $1
      AND ilk = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk
FROM maker.vat_ilk_art
WHERE block_number <= $1
      AND ilk = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk
FROM maker.pit_ilk_spot
WHERE block_number <= $1
      AND ilk = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk
FROM maker.pit_ilk_line
WHERE block_number <= $1
      AND ilk = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk
FROM maker.cat_ilk_chop
WHERE block_number <= $1
      AND ilk = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk
FROM maker.cat_ilk_lump
WHERE block_number <= $1
      AND ilk = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk
FROM maker.cat_ilk_flip
WHERE block_number <= $1
      AND ilk = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk
FROM maker.jug_ilk_rho
WHERE block_number <= $1
      AND ilk = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk
FROM maker.jug_ilk_tax
WHERE block_number <= $1
      AND ilk = $2
$$
LANGUAGE sql;

CREATE TYPE maker.ilk_state AS (
  ilk_id       integer,
  ilk          text,
  block_number numeric,
  take         numeric,
  rate         numeric,
  ink          numeric,
  art          numeric,
  spot         numeric,
  line         numeric,
  chop         numeric,
  lump         numeric,
  flip         text,
  rho          numeric,
  tax          numeric,
  created      numeric,
  updated      numeric
);

CREATE FUNCTION maker.get_ilk_at_block_number(block_number numeric, ilk_id int)
  RETURNS maker.ilk_state
AS $$
WITH takes AS (
    SELECT
      take,
      ilk,
      block_hash
    FROM maker.vat_ilk_take
    WHERE ilk = ilk_id
          AND block_number <= $1
    ORDER BY ilk, block_number DESC
    LIMIT 1
), rates AS (
    SELECT
      rate,
      ilk,
      block_hash
    FROM maker.vat_ilk_rate
    WHERE ilk = ilk_id
          AND block_number <= $1
    ORDER BY ilk, block_number DESC
    LIMIT 1
), inks AS (
    SELECT
      ink,
      ilk,
      block_hash
    FROM maker.vat_ilk_ink
    WHERE ilk = ilk_id
          AND block_number <= $1
    ORDER BY ilk, block_number DESC
    LIMIT 1
), arts AS (
    SELECT
      art,
      ilk,
      block_hash
    FROM maker.vat_ilk_art
    WHERE ilk = ilk_id
          AND block_number <= $1
    ORDER BY ilk, block_number DESC
    LIMIT 1
), spots AS (
    SELECT
      spot,
      ilk,
      block_hash
    FROM maker.pit_ilk_spot
    WHERE ilk = ilk_id
          AND block_number <= $1
    ORDER BY ilk, block_number DESC
    LIMIT 1
), lines AS (
    SELECT
      line,
      ilk,
      block_hash
    FROM maker.pit_ilk_line
    WHERE ilk = ilk_id
          AND block_number <= $1
    ORDER BY ilk, block_number DESC
    LIMIT 1
), chops AS (
    SELECT
      chop,
      ilk,
      block_hash
    FROM maker.cat_ilk_chop
    WHERE ilk = ilk_id
          AND block_number <= $1
    ORDER BY ilk, block_number DESC
    LIMIT 1
), lumps AS (
    SELECT
      lump,
      ilk,
      block_hash
    FROM maker.cat_ilk_lump
    WHERE ilk = ilk_id
          AND block_number <= $1
    ORDER BY ilk, block_number DESC
    LIMIT 1
), flips AS (
    SELECT
      flip,
      ilk,
      block_hash
    FROM maker.cat_ilk_flip
    WHERE ilk = ilk_id
          AND block_number <= $1
    ORDER BY ilk, block_number DESC
    LIMIT 1
), rhos AS (
    SELECT
      rho,
      ilk,
      block_hash
    FROM maker.jug_ilk_rho
    WHERE ilk = ilk_id
          AND block_number <= $1
    ORDER BY ilk, block_number DESC
    LIMIT 1
), taxes AS (
    SELECT
      tax,
      ilk,
      block_hash
    FROM maker.jug_ilk_tax
    WHERE ilk = ilk_id
          AND block_number <= $1
    ORDER BY ilk, block_number DESC
    LIMIT 1
), relevant_blocks AS (
  SELECT * FROM maker.get_ilk_blocks_before($1, $2)
), created AS (
    SELECT DISTINCT ON (relevant_blocks.ilk, relevant_blocks.block_number)
      relevant_blocks.block_number,
      relevant_blocks.block_hash,
      relevant_blocks.ilk,
      headers.block_timestamp
    FROM relevant_blocks
      LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
    ORDER BY block_number ASC
    LIMIT 1
), updated AS (
    SELECT DISTINCT ON (relevant_blocks.ilk, relevant_blocks.block_number)
      relevant_blocks.block_number,
      relevant_blocks.block_hash,
      relevant_blocks.ilk,
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
  takes.take,
  rates.rate,
  inks.ink,
  arts.art,
  spots.spot,
  lines.line,
  chops.chop,
  lumps.lump,
  flips.flip,
  rhos.rho,
  taxes.tax,
  created.block_timestamp,
  updated.block_timestamp
FROM maker.ilks AS ilks
  LEFT JOIN takes ON takes.ilk = ilks.id
  LEFT JOIN rates ON rates.ilk = ilks.id
  LEFT JOIN inks ON inks.ilk = ilks.id
  LEFT JOIN arts ON arts.ilk = ilks.id
  LEFT JOIN spots ON spots.ilk = ilks.id
  LEFT JOIN lines ON lines.ilk = ilks.id
  LEFT JOIN chops ON chops.ilk = ilks.id
  LEFT JOIN lumps ON lumps.ilk = ilks.id
  LEFT JOIN flips ON flips.ilk = ilks.id
  LEFT JOIN rhos ON rhos.ilk = ilks.id
  LEFT JOIN taxes ON taxes.ilk = ilks.id
  LEFT JOIN created ON created.ilk = ilks.id
  LEFT JOIN updated ON updated.ilk = ilks.id
WHERE (
  takes.take is not null OR
  rates.rate is not null OR
  inks.ink is not null OR
  arts.art is not null OR
  spots.spot is not null OR
  lines.line is not null OR
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
DROP FUNCTION IF EXISTS maker.get_relevent_ilk_blocks(block_number numeric, ilk_id int);
DROP TYPE maker.relevant_block CASCADE;
DROP FUNCTION IF EXISTS maker.get_ilk_at_block_number(block_number numeric, ilk_id int );
DROP TYPE maker.ilk_state CASCADE;