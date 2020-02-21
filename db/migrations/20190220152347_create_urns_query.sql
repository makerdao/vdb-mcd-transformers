-- +goose Up

-- Spec: https://github.com/makerdao/vulcan.spec/blob/master/mcd.graphql

CREATE TYPE api.urn_state AS
(
    urn_identifier TEXT,
    ilk_identifier TEXT,
    block_height   BIGINT,
    -- ilk object
    ink            NUMERIC,
    art            NUMERIC,
    -- frobs
    -- bites
    created        TIMESTAMP,
    updated        TIMESTAMP
);

CREATE FUNCTION api.epoch_to_datetime(epoch NUMERIC)
    RETURNS TIMESTAMP AS
$$
SELECT TIMESTAMP 'epoch' + epoch * INTERVAL '1 second' AS datetime
$$
    LANGUAGE SQL
    IMMUTABLE;

CREATE FUNCTION api.max_block()
    RETURNS BIGINT AS
$$
SELECT max(block_number)
FROM public.headers
$$
    LANGUAGE SQL
    STABLE;

-- Function returning state for all urns as of given block
CREATE FUNCTION api.all_urns(block_height BIGINT DEFAULT api.max_block(), max_results INTEGER DEFAULT NULL,
                             result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.urn_state
AS

$body$
WITH urns AS (SELECT urns.id AS urn_id, ilks.id AS ilk_id, ilks.ilk, urns.identifier
              FROM maker.urns urns
                       LEFT JOIN maker.ilks ilks ON urns.ilk_id = ilks.id),
     inks AS ( -- Latest ink for each urn
         SELECT DISTINCT ON (urn_id) urn_id, ink, block_number
         FROM maker.vat_urn_ink
                  LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
         WHERE block_number <= all_urns.block_height
         ORDER BY urn_id, block_number DESC),
     arts AS ( -- Latest art for each urn
         SELECT DISTINCT ON (urn_id) urn_id, art, block_number
         FROM maker.vat_urn_art
                  LEFT JOIN public.headers ON vat_urn_art.header_id = headers.id
         WHERE block_number <= all_urns.block_height
         ORDER BY urn_id, block_number DESC),
     created AS (SELECT urn_id, api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM (SELECT DISTINCT ON (urn_id) urn_id, block_timestamp
                       FROM maker.vat_urn_ink
                                LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
                       ORDER BY urn_id, block_number ASC) earliest_blocks),
     updated AS (SELECT DISTINCT ON (urn_id) urn_id, api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM ((SELECT DISTINCT ON (urn_id) urn_id, block_timestamp
                        FROM maker.vat_urn_ink
                                 LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
                        WHERE block_number <= block_height
                        ORDER BY urn_id, block_number DESC)
                       UNION
                       (SELECT DISTINCT ON (urn_id) urn_id, block_timestamp
                        FROM maker.vat_urn_art
                                 LEFT JOIN public.headers ON vat_urn_art.header_id = headers.id
                        WHERE block_number <= block_height
                        ORDER BY urn_id, block_number DESC)) last_blocks
                 ORDER BY urn_id, block_timestamp DESC)

SELECT urns.identifier,
       ilks.identifier,
       all_urns.block_height,
       inks.ink,
       COALESCE(arts.art, 0),
       created.datetime,
       updated.datetime
FROM inks
         LEFT JOIN arts ON arts.urn_id = inks.urn_id
         LEFT JOIN urns ON inks.urn_id = urns.urn_id
         LEFT JOIN created ON created.urn_id = urns.urn_id
         LEFT JOIN updated ON updated.urn_id = urns.urn_id
         LEFT JOIN maker.ilks ON ilks.id = urns.ilk_id
ORDER BY updated DESC
LIMIT all_urns.max_results
OFFSET
all_urns.result_offset
$body$
    LANGUAGE SQL
    STABLE;

-- +goose Down
DROP FUNCTION api.all_urns(BIGINT, INTEGER, INTEGER);
DROP FUNCTION api.max_block();
DROP FUNCTION api.epoch_to_datetime(NUMERIC);
DROP TYPE api.urn_state CASCADE;
