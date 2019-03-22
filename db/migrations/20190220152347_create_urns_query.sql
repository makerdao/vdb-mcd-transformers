-- +goose Up

-- Spec: https://github.com/makerdao/vulcan.spec/blob/master/mcd.graphql

CREATE TYPE maker.urn_state AS (
  urnId       TEXT,
  ilkId       TEXT,
  blockHeight NUMERIC,
  -- ilk object
  ink         NUMERIC,
  art         NUMERIC,
  ratio       NUMERIC,
  safe        BOOLEAN,
  -- frobs
  -- bites
  created     NUMERIC,
  updated     NUMERIC
  );

-- Function returning state for all urns as of given block
CREATE OR REPLACE FUNCTION maker.get_all_urn_states_at_block(block_height numeric)
  RETURNS SETOF maker.urn_state
AS

$body$
WITH
  ilks AS ( SELECT id, ilk FROM maker.ilks ),

  inks AS ( -- Latest ink for each urn
    SELECT DISTINCT ON (ilk, urn) ilk, urn, ink, block_number
    FROM maker.vat_urn_ink
    WHERE block_number <= block_height
    ORDER BY ilk, urn, block_number DESC
  ),

  arts AS ( -- Latest art for each urn
    SELECT DISTINCT ON (ilk, urn) ilk, urn, art, block_number
    FROM maker.vat_urn_art
    WHERE block_number <= block_height
    ORDER BY ilk, urn, block_number DESC
  ),

  rates AS ( -- Latest rate for each ilk
    SELECT DISTINCT ON (ilk) ilk, rate, block_number
    FROM maker.vat_ilk_rate
    WHERE block_number <= block_height
    ORDER BY ilk, block_number DESC
  ),

  spots AS ( -- Get latest price update for ilk. Problematic from update frequency, slow query?
    SELECT DISTINCT ON (ilk) ilk, spot, block_number
    FROM maker.pit_ilk_spot
    WHERE block_number <= block_height
    ORDER BY ilk, block_number DESC
  ),

  ratio_data AS (
    SELECT inks.ilk, inks.urn, ink, spot, art, rate
    FROM inks
      JOIN arts ON arts.ilk = inks.ilk AND arts.urn = inks.urn
      JOIN spots ON spots.ilk = arts.ilk
      JOIN rates ON rates.ilk = arts.ilk
  ),

  ratios AS (
    SELECT ilk, urn, ((1.0 * ink * spot) / NULLIF(art * rate, 0)) AS ratio FROM ratio_data
  ),

  safe AS (
    SELECT ilk, urn, (ratio >= 1) AS safe FROM ratios
  ),

  created AS (
    SELECT ilk, urn, block_timestamp AS created
    FROM
      (
        SELECT DISTINCT ON (ilk, urn) ilk, urn, block_hash FROM maker.vat_urn_ink
        ORDER BY ilk, urn, block_number ASC
      ) earliest_blocks
        LEFT JOIN public.headers ON hash = block_hash
  ),

  updated AS (
    SELECT DISTINCT ON (ilk, urn) ilk, urn, headers.block_timestamp AS updated
    FROM
      (
        (SELECT DISTINCT ON (ilk, urn) ilk, urn, block_hash FROM maker.vat_urn_ink
         WHERE block_number <= block_height
         ORDER BY ilk, urn, block_number DESC)
        UNION
        (SELECT DISTINCT ON (ilk, urn) ilk, urn, block_hash FROM maker.vat_urn_art
         WHERE block_number <= block_height
         ORDER BY ilk, urn, block_number DESC)
      ) last_blocks
        LEFT JOIN public.headers ON headers.hash = last_blocks.block_hash
    ORDER BY ilk, urn, headers.block_timestamp DESC
  )

SELECT inks.urn, ilks.ilk, $1, inks.ink, arts.art, ratios.ratio,
       COALESCE(safe.safe, arts.art = 0), created.created, updated.updated
FROM inks
  LEFT JOIN arts     ON arts.ilk = inks.ilk    AND arts.urn = inks.urn
  LEFT JOIN ilks     ON ilks.id = arts.ilk
  LEFT JOIN ratios   ON ratios.ilk = arts.ilk  AND ratios.urn = arts.urn
  LEFT JOIN safe     ON safe.ilk = arts.ilk    AND safe.urn = arts.urn
  LEFT JOIN created  ON created.ilk = arts.ilk AND created.urn = arts.urn
  LEFT JOIN updated  ON updated.ilk = arts.ilk AND updated.urn = arts.urn
  -- Add collections of frob and bite events?
$body$
LANGUAGE SQL;


-- +goose Down
DROP FUNCTION IF EXISTS maker.get_all_urn_states_at_block(numeric);
DROP TYPE IF EXISTS maker.urn_state CASCADE;
