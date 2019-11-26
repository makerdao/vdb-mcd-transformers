-- +goose Up

-- Function returning state for all ilks as of the given block height
CREATE FUNCTION api.all_ilks(block_height BIGINT DEFAULT api.max_block(), max_results INTEGER DEFAULT NULL,
                             result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.ilk_state
AS
$$
WITH rates AS (SELECT DISTINCT ON (ilk_id) rate, ilk_id
               FROM maker.vat_ilk_rate
                        LEFT JOIN public.headers ON vat_ilk_rate.header_id = headers.id
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     arts AS (SELECT DISTINCT ON (ilk_id) art, ilk_id
              FROM maker.vat_ilk_art
                       LEFT JOIN public.headers ON vat_ilk_art.header_id = headers.id
              WHERE block_number <= all_ilks.block_height
              ORDER BY ilk_id, block_number DESC),
     spots AS (SELECT DISTINCT ON (ilk_id) spot, ilk_id
               FROM maker.vat_ilk_spot
                        LEFT JOIN public.headers ON vat_ilk_spot.header_id = headers.id
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     lines AS (SELECT DISTINCT ON (ilk_id) line, ilk_id
               FROM maker.vat_ilk_line
                        LEFT JOIN public.headers ON vat_ilk_line.header_id = headers.id
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     dusts AS (SELECT DISTINCT ON (ilk_id) dust, ilk_id
               FROM maker.vat_ilk_dust
                        LEFT JOIN public.headers ON vat_ilk_dust.header_id = headers.id
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     chops AS (SELECT DISTINCT ON (ilk_id) chop, ilk_id
               FROM maker.cat_ilk_chop
                        LEFT JOIN public.headers ON cat_ilk_chop.header_id = headers.id
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     lumps AS (SELECT DISTINCT ON (ilk_id) lump, ilk_id
               FROM maker.cat_ilk_lump
                        LEFT JOIN public.headers ON cat_ilk_lump.header_id = headers.id
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     flips AS (SELECT DISTINCT ON (ilk_id) flip, ilk_id
               FROM maker.cat_ilk_flip
                        LEFT JOIN public.headers ON cat_ilk_flip.header_id = headers.id
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     rhos AS (SELECT DISTINCT ON (ilk_id) rho, ilk_id
              FROM maker.jug_ilk_rho
                       LEFT JOIN public.headers ON jug_ilk_rho.header_id = headers.id
              WHERE block_number <= all_ilks.block_height
              ORDER BY ilk_id, block_number DESC),
     duties AS (SELECT DISTINCT ON (ilk_id) duty, ilk_id
                FROM maker.jug_ilk_duty
                         LEFT JOIN public.headers ON jug_ilk_duty.header_id = headers.id
                WHERE block_number <= all_ilks.block_height
                ORDER BY ilk_id, block_number DESC),
     pips AS (SELECT DISTINCT ON (ilk_id) pip, ilk_id
              FROM maker.spot_ilk_pip
                       LEFT JOIN public.headers ON spot_ilk_pip.header_id = headers.id
              WHERE block_number <= all_ilks.block_height
              ORDER BY ilk_id, block_number DESC),
     mats AS (SELECT DISTINCT ON (ilk_id) mat, ilk_id
              FROM maker.spot_ilk_mat
                       LEFT JOIN public.headers ON spot_ilk_mat.header_id = headers.id
              WHERE block_number <= all_ilks.block_height
              ORDER BY ilk_id, block_number DESC)
SELECT ilks.identifier,
       all_ilks.block_height,
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
       pips.pip,
       mats.mat,
       (SELECT api.epoch_to_datetime(b.block_timestamp) AS created
        FROM api.get_ilk_blocks_before(ilks.identifier, all_ilks.block_height) b
        ORDER BY b.block_height ASC
        LIMIT 1),
       (SELECT api.epoch_to_datetime(b.block_timestamp) AS updated
        FROM api.get_ilk_blocks_before(ilks.identifier, all_ilks.block_height) b
        ORDER BY b.block_height DESC
        LIMIT 1)
FROM maker.ilks AS ilks
         LEFT JOIN rates on rates.ilk_id = ilks.id
         LEFT JOIN arts on arts.ilk_id = ilks.id
         LEFT JOIN spots on spots.ilk_id = ilks.id
         LEFT JOIN lines on lines.ilk_id = ilks.id
         LEFT JOIN dusts on dusts.ilk_id = ilks.id
         LEFT JOIN chops on chops.ilk_id = ilks.id
         LEFT JOIN lumps on lumps.ilk_id = ilks.id
         LEFT JOIN flips on flips.ilk_id = ilks.id
         LEFT JOIN rhos on rhos.ilk_id = ilks.id
         LEFT JOIN duties on duties.ilk_id = ilks.id
         LEFT JOIN pips on pips.ilk_id = ilks.id
         LEFT JOIN mats on mats.ilk_id = ilks.id
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
              duties.duty is not null OR
              pips.pip is not null OR
              mats.mat is not null
          )
ORDER BY updated DESC
LIMIT all_ilks.max_results
OFFSET
all_ilks.result_offset
$$
    LANGUAGE SQL
    STABLE;

-- +goose Down
DROP FUNCTION api.all_ilks(BIGINT, INTEGER, INTEGER);
