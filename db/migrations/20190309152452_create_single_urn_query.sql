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
         SELECT DISTINCT ON (urn_id) urn_id, ink, block_number, block_timestamp
         FROM maker.vat_urn_ink
                  LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
         WHERE urn_id = (SELECT urn_id from urn where identifier = urn_identifier)
           AND block_number <= get_urn.block_height
         ORDER BY urn_id, block_number DESC),
     art AS ( -- Latest art
         SELECT DISTINCT ON (urn_id) urn_id, art, block_number, block_timestamp
         FROM maker.vat_urn_art
                  LEFT JOIN public.headers ON vat_urn_art.header_id = headers.id
         WHERE urn_id = (SELECT urn_id from urn where identifier = urn_identifier)
           AND block_number <= get_urn.block_height
         ORDER BY urn_id, block_number DESC),
     created AS (SELECT urn_id, api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM (SELECT DISTINCT ON (urn_id) urn_id,
                                                   block_timestamp
                                                   -- TODO: should we be using urn ink for created?
                                                   -- Can a CDP exist before collateral is locked?
                       FROM maker.vat_urn_ink
                                LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
                       WHERE urn_id = (SELECT urn_id from urn where identifier = urn_identifier)
                       ORDER BY urn_id, block_number ASC) earliest_blocks),
     updated AS (SELECT DISTINCT ON (urn_id) urn_id, api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM (SELECT urn_id, block_number, block_timestamp
                       FROM ink
                       UNION
                       SELECT urn_id, block_number, block_timestamp
                       FROM art) last_blocks
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
