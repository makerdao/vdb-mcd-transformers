-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Function returning state for a single urn as of given block
CREATE FUNCTION api.get_urn(ilk_identifier TEXT, urn_identifier TEXT, block_height BIGINT DEFAULT api.max_block())
    RETURNS api.urn_state
AS

$body$
WITH urn AS (SELECT urns.id AS urn_id, ilks.id AS ilk_id, ilks.ilk, urns.identifier
             FROM maker.urns urns
                      LEFT JOIN maker.ilks ilks ON urns.ilk_id = ilks.id
             WHERE ilks.identifier = ilk_identifier
               AND urns.identifier = urn_identifier),
     ink AS ( -- Latest ink
         SELECT DISTINCT ON (urn_id) urn_id, ink, block_number
         FROM maker.vat_urn_ink
         WHERE urn_id = (SELECT urn_id from urn where identifier = urn_identifier)
           AND block_number <= get_urn.block_height
         ORDER BY urn_id, block_number DESC),
     art AS ( -- Latest art
         SELECT DISTINCT ON (urn_id) urn_id, art, block_number
         FROM maker.vat_urn_art
         WHERE urn_id = (SELECT urn_id from urn where identifier = urn_identifier)
           AND block_number <= get_urn.block_height
         ORDER BY urn_id, block_number DESC),
     rate AS ( -- Latest rate for ilk
         SELECT DISTINCT ON (ilk_id) ilk_id, rate, block_number
         FROM maker.vat_ilk_rate
         WHERE ilk_id = (SELECT ilk_id FROM urn)
           AND block_number <= get_urn.block_height
         ORDER BY ilk_id, block_number DESC),
     spot AS ( -- Get latest price update for ilk. Problematic from update frequency, slow query?
         SELECT DISTINCT ON (ilk_id) ilk_id, spot, block_number
         FROM maker.vat_ilk_spot
         WHERE ilk_id = (SELECT ilk_id FROM urn)
           AND block_number <= get_urn.block_height
         ORDER BY ilk_id, block_number DESC),
     created AS (SELECT urn_id, api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM (SELECT DISTINCT ON (urn_id) urn_id, block_hash
                       FROM maker.vat_urn_ink
                       WHERE urn_id = (SELECT urn_id from urn where identifier = urn_identifier)
                       ORDER BY urn_id, block_number ASC) earliest_blocks
                          LEFT JOIN public.headers ON hash = block_hash),
     updated AS (SELECT DISTINCT ON (urn_id) urn_id, api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM (SELECT urn_id, block_number
                       FROM ink
                       UNION
                       SELECT urn_id, block_number
                       FROM art) last_blocks
                          LEFT JOIN public.headers ON headers.block_number = last_blocks.block_number
                 ORDER BY urn_id, block_timestamp DESC)

SELECT get_urn.urn_identifier,
       ilk_identifier,
       $3,
       ink.ink,
       COALESCE(art.art, 0),
       created.datetime,
       updated.datetime
FROM ink
         LEFT JOIN art ON art.urn_id = ink.urn_id
         LEFT JOIN urn ON urn.urn_id = ink.urn_id
         LEFT JOIN created ON created.urn_id = art.urn_id
         LEFT JOIN updated ON updated.urn_id = art.urn_id
WHERE ink.urn_id IS NOT NULL
$body$
    LANGUAGE SQL
    STABLE
    STRICT;

-- +goose Down
DROP FUNCTION api.get_urn(TEXT, TEXT, BIGINT);
