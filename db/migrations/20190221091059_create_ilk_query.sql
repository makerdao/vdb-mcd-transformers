-- +goose Up
CREATE TYPE api.relevant_block AS (
    block_height BIGINT,
    block_hash TEXT,
    ilk_id INTEGER
    );

CREATE FUNCTION api.get_ilk_blocks_before(ilk_identifier TEXT, block_height BIGINT)
    RETURNS SETOF api.relevant_block AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE identifier = ilk_identifier)
SELECT block_number AS block_height, block_hash, ilk_id
FROM maker.vat_ilk_rate
WHERE block_number <= get_ilk_blocks_before.block_height
  AND ilk_id = (SELECT id FROM ilk)
UNION
SELECT block_number AS block_height, block_hash, ilk_id
FROM maker.vat_ilk_art
WHERE block_number <= get_ilk_blocks_before.block_height
  AND ilk_id = (SELECT id FROM ilk)
UNION
SELECT block_number AS block_height, block_hash, ilk_id
FROM maker.vat_ilk_spot
WHERE block_number <= get_ilk_blocks_before.block_height
  AND ilk_id = (SELECT id FROM ilk)
UNION
SELECT block_number AS block_height, block_hash, ilk_id
FROM maker.vat_ilk_line
WHERE block_number <= get_ilk_blocks_before.block_height
  AND ilk_id = (SELECT id FROM ilk)
UNION
SELECT block_number AS block_height, block_hash, ilk_id
FROM maker.vat_ilk_dust
WHERE block_number <= get_ilk_blocks_before.block_height
  AND ilk_id = (SELECT id FROM ilk)
UNION
SELECT block_number AS block_height, block_hash, ilk_id
FROM maker.cat_ilk_chop
WHERE block_number <= get_ilk_blocks_before.block_height
  AND ilk_id = (SELECT id FROM ilk)
UNION
SELECT block_number AS block_height, block_hash, ilk_id
FROM maker.cat_ilk_lump
WHERE block_number <= get_ilk_blocks_before.block_height
  AND ilk_id = (SELECT id FROM ilk)
UNION
SELECT block_number AS block_height, block_hash, ilk_id
FROM maker.cat_ilk_flip
WHERE block_number <= get_ilk_blocks_before.block_height
  AND ilk_id = (SELECT id FROM ilk)
UNION
SELECT block_number AS block_height, block_hash, ilk_id
FROM maker.jug_ilk_rho
WHERE block_number <= get_ilk_blocks_before.block_height
  AND ilk_id = (SELECT id FROM ilk)
UNION
SELECT block_number AS block_height, block_hash, ilk_id
FROM maker.jug_ilk_duty
WHERE block_number <= get_ilk_blocks_before.block_height
  AND ilk_id = (SELECT id FROM ilk)
ORDER BY block_height DESC
$$
    LANGUAGE sql
    STABLE;

COMMENT ON FUNCTION api.get_ilk_blocks_before(TEXT, BIGINT)
    IS E'@omit';


CREATE TYPE api.ilk_state AS (
    ilk_identifier TEXT,
    block_height BIGINT,
    rate NUMERIC,
    art NUMERIC,
    spot NUMERIC,
    line NUMERIC,
    dust NUMERIC,
    chop NUMERIC,
    lump NUMERIC,
    flip TEXT,
    rho NUMERIC,
    duty NUMERIC,
    created TIMESTAMP,
    updated TIMESTAMP
    );

-- Function returning the state for a single ilk as of the given block height
CREATE FUNCTION api.get_ilk(ilk_identifier TEXT, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.ilk_state
AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE identifier = ilk_identifier),
     rates AS (SELECT rate, ilk_id, block_hash
               FROM maker.vat_ilk_rate
               WHERE ilk_id = (SELECT id FROM ilk)
                 AND block_number <= get_ilk.block_height
               ORDER BY ilk_id, block_number DESC
               LIMIT 1),
     arts AS (SELECT art, ilk_id, block_hash
              FROM maker.vat_ilk_art
              WHERE ilk_id = (SELECT id FROM ilk)
                AND block_number <= get_ilk.block_height
              ORDER BY ilk_id, block_number DESC
              LIMIT 1),
     spots AS (SELECT spot, ilk_id, block_hash
               FROM maker.vat_ilk_spot
               WHERE ilk_id = (SELECT id FROM ilk)
                 AND block_number <= get_ilk.block_height
               ORDER BY ilk_id, block_number DESC
               LIMIT 1),
     lines AS (SELECT line, ilk_id, block_hash
               FROM maker.vat_ilk_line
               WHERE ilk_id = (SELECT id FROM ilk)
                 AND block_number <= get_ilk.block_height
               ORDER BY ilk_id, block_number DESC
               LIMIT 1),
     dusts AS (SELECT dust, ilk_id, block_hash
               FROM maker.vat_ilk_dust
               WHERE ilk_id = (SELECT id FROM ilk)
                 AND block_number <= get_ilk.block_height
               ORDER BY ilk_id, block_number DESC
               LIMIT 1),
     chops AS (SELECT chop, ilk_id, block_hash
               FROM maker.cat_ilk_chop
               WHERE ilk_id = (SELECT id FROM ilk)
                 AND block_number <= get_ilk.block_height
               ORDER BY ilk_id, block_number DESC
               LIMIT 1),
     lumps AS (SELECT lump, ilk_id, block_hash
               FROM maker.cat_ilk_lump
               WHERE ilk_id = (SELECT id FROM ilk)
                 AND block_number <= get_ilk.block_height
               ORDER BY ilk_id, block_number DESC
               LIMIT 1),
     flips AS (SELECT flip, ilk_id, block_hash
               FROM maker.cat_ilk_flip
               WHERE ilk_id = (SELECT id FROM ilk)
                 AND block_number <= get_ilk.block_height
               ORDER BY ilk_id, block_number DESC
               LIMIT 1),
     rhos AS (SELECT rho, ilk_id, block_hash
              FROM maker.jug_ilk_rho
              WHERE ilk_id = (SELECT id FROM ilk)
                AND block_number <= get_ilk.block_height
              ORDER BY ilk_id, block_number DESC
              LIMIT 1),
     duties AS (SELECT duty, ilk_id, block_hash
                FROM maker.jug_ilk_duty
                WHERE ilk_id = (SELECT id FROM ilk)
                  AND block_number <= get_ilk.block_height
                ORDER BY ilk_id, block_number DESC
                LIMIT 1),
     relevant_blocks AS (SELECT * FROM api.get_ilk_blocks_before(ilk_identifier, get_ilk.block_height)),
     created AS (SELECT DISTINCT ON (relevant_blocks.ilk_id,
         relevant_blocks.block_height) relevant_blocks.block_height,
                                       relevant_blocks.block_hash,
                                       relevant_blocks.ilk_id,
                                       api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM relevant_blocks
                          LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
                 ORDER BY relevant_blocks.block_height ASC
                 LIMIT 1),
     updated AS (SELECT DISTINCT ON (relevant_blocks.ilk_id,
         relevant_blocks.block_height) relevant_blocks.block_height,
                                       relevant_blocks.block_hash,
                                       relevant_blocks.ilk_id,
                                       api.epoch_to_datetime(headers.block_timestamp) AS datetime
                 FROM relevant_blocks
                          LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
                 ORDER BY relevant_blocks.block_height DESC
                 LIMIT 1)

SELECT ilks.identifier,
       get_ilk.block_height,
       rates.rate,
       arts.art,
       spots.spot,
       lines.line,
       dusts.dust,
       chops.chop,
       lumps.lump,
       flips.flip,
       rhos.rho,
       duties.duty,
       created.datetime,
       updated.datetime
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
         LEFT JOIN duties ON duties.ilk_id = ilks.id
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
              duties.duty is not null
          )
$$
    LANGUAGE SQL
    STABLE
    STRICT;

-- +goose Down
DROP FUNCTION api.get_ilk_blocks_before(TEXT, BIGINT);
DROP TYPE api.relevant_block CASCADE;
DROP FUNCTION api.get_ilk(TEXT, BIGINT);
DROP TYPE api.ilk_state CASCADE;