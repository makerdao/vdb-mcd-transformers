-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Function returning state for a single urn as of given block
CREATE OR REPLACE FUNCTION maker.get_urn_state_at_block(ilk TEXT, urn TEXT, block_height NUMERIC)
  RETURNS maker.urn_state
AS

$body$
WITH
  urn AS (
    SELECT urns.id AS urn_id, ilks.id AS ilk_id, ilks.ilk, urns.guy
    FROM maker.urns urns
    LEFT JOIN maker.ilks ilks
    ON urns.ilk = ilks.id
    WHERE ilks.ilk = $1 AND urns.guy = $2
  ),

  ink AS ( -- Latest ink
    SELECT DISTINCT ON (urn) urn, ink, block_number
    FROM maker.vat_urn_ink
    WHERE urn = (SELECT urn_id from urn where guy = $2) AND block_number <= block_height
    ORDER BY urn, block_number DESC
  ),

  art AS ( -- Latest art
    SELECT DISTINCT ON (urn) urn, art, block_number
    FROM maker.vat_urn_art
    WHERE urn = (SELECT urn_id from urn where guy = $2) AND block_number <= block_height
    ORDER BY urn, block_number DESC
  ),

  rate AS ( -- Latest rate for ilk
    SELECT DISTINCT ON (ilk) ilk, rate, block_number
    FROM maker.vat_ilk_rate
    WHERE ilk = (SELECT ilk_id from urn where ilk = $1) AND block_number <= block_height
    ORDER BY ilk, block_number DESC
  ),

  spot AS ( -- Get latest price update for ilk. Problematic from update frequency, slow query?
    SELECT DISTINCT ON (ilk) ilk, spot, block_number
    FROM maker.pit_ilk_spot
    WHERE ilk = (SELECT ilk_id from urn where ilk = $1) AND block_number <= block_height
    ORDER BY ilk, block_number DESC
  ),

  ratio_data AS (
    SELECT urn.ilk, urn.guy, ink, spot, art, rate
    FROM ink
      JOIN urn ON ink.urn = urn.urn_id
      JOIN art ON art.urn = ink.urn
      JOIN spot ON spot.ilk = urn.ilk_id
      JOIN rate ON rate.ilk = spot.ilk
  ),

  ratio AS (
    SELECT ilk, guy, ((1.0 * ink * spot) / NULLIF(art * rate, 0)) AS ratio
    FROM ratio_data
  ),

  safe AS (
    SELECT ilk, guy, (ratio >= 1) AS safe FROM ratio
  ),

  created AS (
    SELECT urn, block_timestamp AS created
    FROM
      (
        SELECT DISTINCT ON (urn) urn, block_hash FROM maker.vat_urn_ink
        WHERE urn = (SELECT urn_id from urn where guy = $2)
        ORDER BY urn, block_number ASC
      ) earliest_blocks
        LEFT JOIN public.headers ON hash = block_hash
  ),

  updated AS (
    SELECT DISTINCT ON (urn) urn, headers.block_timestamp AS updated
    FROM
      (
        SELECT urn, block_number FROM ink
        UNION
        SELECT urn, block_number FROM art
      ) last_blocks
        LEFT JOIN public.headers ON headers.block_number = last_blocks.block_number
    ORDER BY urn, headers.block_timestamp DESC
  )

SELECT $2 AS urnId, $1 AS ilkId, $3 AS block_height, ink.ink, art.art, ratio.ratio,
       COALESCE(safe.safe, art.art = 0), created.created, updated.updated
FROM ink
  LEFT JOIN art     ON art.urn = ink.urn
  LEFT JOIN urn     ON urn.urn_id = ink.urn
  LEFT JOIN ratio   ON ratio.ilk = urn.ilk   AND ratio.guy = urn.guy
  LEFT JOIN safe    ON safe.ilk = ratio.ilk  AND safe.guy = ratio.guy
  LEFT JOIN created ON created.urn = art.urn
  LEFT JOIN updated ON updated.urn = art.urn
  -- Add collections of frob and bite events?
WHERE ink.urn IS NOT NULL
$body$
  LANGUAGE SQL;


-- +goose Down
DROP FUNCTION IF EXISTS maker.get_urn_state_at_block(TEXT, TEXT, NUMERIC);
