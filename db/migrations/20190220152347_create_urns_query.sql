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
  urns AS (
    SELECT urns.id AS urn_id, ilks.id AS ilk_id, ilks.ilk, urns.guy
    FROM maker.urns urns
    LEFT JOIN maker.ilks ilks
    ON urns.ilk = ilks.id
  ),

  inks AS ( -- Latest ink for each urn
    SELECT DISTINCT ON (urn) urn, ink, block_number
    FROM maker.vat_urn_ink
    WHERE block_number <= block_height
    ORDER BY urn, block_number DESC
  ),

  arts AS ( -- Latest art for each urn
    SELECT DISTINCT ON (urn) urn, art, block_number
    FROM maker.vat_urn_art
    WHERE block_number <= block_height
    ORDER BY urn, block_number DESC
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
    SELECT urns.ilk, urns.guy, inks.ink, spots.spot, arts.art, rates.rate
    FROM inks
      JOIN urns ON inks.urn = urns.urn_id
      JOIN arts ON arts.urn = inks.urn
      JOIN spots ON spots.ilk = urns.ilk_id
      JOIN rates ON rates.ilk = spots.ilk
  ),

  ratios AS (
    SELECT ilk, guy, ((1.0 * ink * spot) / NULLIF(art * rate, 0)) AS ratio FROM ratio_data
  ),

  safe AS (
    SELECT ilk, guy, (ratio >= 1) AS safe FROM ratios
  ),

  created AS (
    SELECT urn, block_timestamp AS created
    FROM
      (
        SELECT DISTINCT ON (urn) urn, block_hash FROM maker.vat_urn_ink
        ORDER BY urn, block_number ASC
      ) earliest_blocks
        LEFT JOIN public.headers ON hash = block_hash
  ),

  updated AS (
    SELECT DISTINCT ON (urn) urn, headers.block_timestamp AS updated
    FROM
      (
        (SELECT DISTINCT ON (urn) urn, block_hash FROM maker.vat_urn_ink
         WHERE block_number <= block_height
         ORDER BY urn, block_number DESC)
        UNION
        (SELECT DISTINCT ON (urn) urn, block_hash FROM maker.vat_urn_art
         WHERE block_number <= block_height
         ORDER BY urn, block_number DESC)
      ) last_blocks
        LEFT JOIN public.headers ON headers.hash = last_blocks.block_hash
    ORDER BY urn, headers.block_timestamp DESC
  )

SELECT urns.guy, urns.ilk, $1, inks.ink, arts.art, ratios.ratio,
       COALESCE(safe.safe, arts.art = 0), created.created, updated.updated
FROM inks
  LEFT JOIN arts     ON arts.urn = inks.urn
  LEFT JOIN urns     ON arts.urn = urns.urn_id
  LEFT JOIN ratios   ON ratios.guy = urns.guy
  LEFT JOIN safe     ON safe.guy = ratios.guy
  LEFT JOIN created  ON created.urn = urns.urn_id
  LEFT JOIN updated  ON updated.urn = urns.urn_id
  -- Add collections of frob and bite events?
$body$
LANGUAGE SQL;


-- +goose Down
DROP FUNCTION IF EXISTS maker.get_all_urn_states_at_block(numeric);
DROP TYPE IF EXISTS maker.urn_state CASCADE;
