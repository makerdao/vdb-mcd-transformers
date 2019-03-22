-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Function returning state for a single urn as of given block
CREATE OR REPLACE FUNCTION maker.get_urn_state_at_block(ilk TEXT, urn TEXT, block_height NUMERIC)
  RETURNS maker.urn_state
AS

$body$
WITH
  ilk AS ( SELECT id, ilk FROM maker.ilks WHERE ilk = $1 ),

  ink AS ( -- Latest ink
    SELECT DISTINCT ON (ilk, urn) ilk, urn, ink, block_number
    FROM maker.vat_urn_ink
    WHERE ilk = (SELECT id from ilk where ilk = $1) AND urn = $2 AND block_number <= block_height
    ORDER BY ilk, urn, block_number DESC
  ),

  art AS ( -- Latest art
    SELECT DISTINCT ON (ilk, urn) ilk, urn, art, block_number
    FROM maker.vat_urn_art
    WHERE ilk = (SELECT id from ilk where ilk = $1) AND urn = $2 AND block_number <= block_height
    ORDER BY ilk, urn, block_number DESC
  ),

  rate AS ( -- Latest rate for ilk
    SELECT DISTINCT ON (ilk) ilk, rate, block_number
    FROM maker.vat_ilk_rate
    WHERE ilk = (SELECT id from ilk where ilk = $1) AND block_number <= block_height
    ORDER BY ilk, block_number DESC
  ),

  spot AS ( -- Get latest price update for ilk. Problematic from update frequency, slow query?
    SELECT DISTINCT ON (ilk) ilk, spot, block_number
    FROM maker.pit_ilk_spot
    WHERE ilk = (SELECT id from ilk where ilk = $1) AND block_number <= block_height
    ORDER BY ilk, block_number DESC
  ),

  ratio_data AS (
    SELECT ink.ilk, ink.urn, ink, spot, art, rate
    FROM ink
           JOIN art  ON art.ilk  = ink.ilk AND art.urn = ink.urn
           JOIN spot ON spot.ilk = art.ilk
           JOIN rate ON rate.ilk = art.ilk
  ),

  ratio AS (
    SELECT ilk, urn, ((1.0 * ink * spot) / NULLIF(art * rate, 0)) AS ratio
    FROM ratio_data
  ),

  safe AS (
    SELECT ilk, urn, (ratio >= 1) AS safe FROM ratio
  ),

  created AS (
    SELECT ilk, urn, block_timestamp AS created
    FROM
      (
        SELECT DISTINCT ON (ilk, urn) ilk, urn, block_hash FROM maker.vat_urn_ink
        WHERE ilk = (SELECT id from ilk where ilk = $1) AND urn = $2
        ORDER BY ilk, urn, block_number ASC
      ) earliest_blocks
        LEFT JOIN public.headers ON hash = block_hash
  ),

  updated AS (
    SELECT DISTINCT ON (ilk, urn) ilk, urn, headers.block_timestamp AS updated
    FROM
      (
        SELECT ilk, urn, block_number FROM ink
        UNION
        SELECT ilk, urn, block_number FROM art
      ) last_blocks
        LEFT JOIN public.headers ON headers.block_number = last_blocks.block_number
    ORDER BY ilk, urn, headers.block_timestamp DESC
  )

SELECT $2 AS urnId, $1 AS ilkId, $3 AS block_height, ink.ink, art.art, ratio.ratio,
       COALESCE(safe.safe, art.art = 0), created.created, updated.updated
FROM ink
  LEFT JOIN art     ON art.ilk = ink.ilk     AND art.urn = ink.urn
  LEFT JOIN ilk     ON ilk.id = art.ilk
  LEFT JOIN ratio   ON ratio.ilk = art.ilk   AND ratio.urn = art.urn
  LEFT JOIN safe    ON safe.ilk = art.ilk    AND safe.urn = art.urn
  LEFT JOIN created ON created.ilk = art.ilk AND created.urn = art.urn
  LEFT JOIN updated ON updated.ilk = art.ilk AND updated.urn = art.urn
  -- Add collections of frob and bite events?
WHERE ink.urn IS NOT NULL
$body$
  LANGUAGE SQL;


-- +goose Down
DROP FUNCTION IF EXISTS maker.get_urn_state_at_block(TEXT, TEXT, NUMERIC);
