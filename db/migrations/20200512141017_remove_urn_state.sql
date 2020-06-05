-- +goose Up
-- SQL in this section is executed when the migration is applied.
DROP FUNCTION IF EXISTS api.all_urn_states;

DROP FUNCTION IF EXISTS api.get_urn;
CREATE FUNCTION api.get_urn(ilk_identifier text, urn_identifier text, block_height bigint DEFAULT api.max_block()) RETURNS api.urn_snapshot
    LANGUAGE sql STABLE STRICT
    AS $_$

SELECT urn_identifier, ilk_identifier, get_urn.block_height, ink, art, created, updated
    FROM api.urn_snapshot
    WHERE ilk_identifier = get_urn.ilk_identifier
    AND urn_identifier = get_urn.urn_identifier
    AND block_height <= get_urn.block_height
    ORDER BY updated DESC
    LIMIT 1
$_$;

DROP FUNCTION IF EXISTS api.bite_event_urn;
CREATE FUNCTION api.bite_event_urn(event api.bite_event) RETURNS api.urn_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn(event.ilk_identifier, event.urn_identifier, event.block_height)
$$;

DROP FUNCTION IF EXISTS api.flip_bid_snapshot_urn;
CREATE FUNCTION api.flip_bid_snapshot_urn(flip api.flip_bid_snapshot) RETURNS api.urn_snapshot
       LANGUAGE sql STABLE
       AS $$
SELECT *
FROM api.get_urn(
     (SELECT identifier FROM maker.ilks WHERE ilks.id = flip.ilk_id),
     (SELECT identifier FROM maker.urns WHERE urns.id = flip.urn_id),
     flip.block_height)
$$;

DROP FUNCTION IF EXISTS api.frob_event_urn;
CREATE FUNCTION api.frob_event_urn(event api.frob_event) RETURNS api.urn_snapshot
       LANGUAGE sql STABLE
       AS $$
SELECT *
FROM api.get_urn(event.ilk_identifier, event.urn_identifier, event.block_height)
$$;

DROP FUNCTION IF EXISTS api.managed_cdp_urn;
CREATE FUNCTION api.managed_cdp_urn(cdp api.managed_cdp) RETURNS api.urn_snapshot
       LANGUAGE sql STABLE
       AS $$
SELECT *
FROM api.get_urn(cdp.ilk_identifier, cdp.urn_identifier, api.max_block())
$$;

DROP FUNCTION IF EXISTS api.urn_state_bites;
DROP FUNCTION IF EXISTS api.urn_state_frobs;
DROP FUNCTION IF EXISTS api.urn_state_ilk;
DROP TYPE IF EXISTS api.urn_state;

-- +goose Down
CREATE TYPE api.urn_state AS (
	urn_identifier text,
	ilk_identifier text,
	block_height bigint,
	ink numeric,
	art numeric,
	created timestamp without time zone,
	updated timestamp without time zone
);

DROP FUNCTION IF EXISTS api.get_urn;
CREATE FUNCTION api.get_urn(ilk_identifier text, urn_identifier text, block_height bigint DEFAULT api.max_block()) RETURNS api.urn_state
    LANGUAGE sql STABLE STRICT
    AS $_$
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
$_$;

-- +goose StatementBegin
CREATE FUNCTION api.all_urn_states(ilk_identifier TEXT, urn_identifier TEXT,
                                   block_height BIGINT DEFAULT api.max_block(),
                                   max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.urn_state AS
$$
BEGIN
    RETURN QUERY (
        WITH urn_id AS (
            SELECT id
            FROM maker.urns
            WHERE urns.identifier = all_urn_states.urn_identifier
              AND urns.ilk_id = (SELECT id
                                 FROM maker.ilks
                                 WHERE ilks.identifier = all_urn_states.ilk_identifier)
                                 ),
             relevant_blocks AS (
                 SELECT block_number
                 FROM maker.vat_urn_ink
                          LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
                 WHERE vat_urn_ink.urn_id = (SELECT * FROM urn_id)
                   AND block_number <= all_urn_states.block_height
                 UNION
                 SELECT block_number
                 FROM maker.vat_urn_art
                          LEFT JOIN public.headers ON vat_urn_art.header_id = headers.id
                 WHERE vat_urn_art.urn_id = (SELECT * FROM urn_id)
                   AND block_number <= all_urn_states.block_height)
        SELECT r.*
        FROM relevant_blocks,
             LATERAL api.get_urn(ilk_identifier, urn_identifier, relevant_blocks.block_number) r
        ORDER BY relevant_blocks.block_number DESC
        LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
        OFFSET
        all_urn_states.result_offset
    );
END;
$$
    LANGUAGE plpgsql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;
-- +goose StatementEnd

DROP FUNCTION IF EXISTS api.bite_event_urn;
CREATE FUNCTION api.bite_event_urn(event api.bite_event) RETURNS api.urn_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn(event.ilk_identifier, event.urn_identifier, event.block_height)
$$;

DROP FUNCTION IF EXISTS api.flip_bid_snapshot_urn;
CREATE FUNCTION api.flip_bid_snapshot_urn(flip api.flip_bid_snapshot) RETURNS SETOF api.urn_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn((SELECT identifier FROM maker.ilks WHERE ilks.id = flip.ilk_id),
                 (SELECT identifier FROM maker.urns WHERE urns.id = flip.urn_id), flip.block_height)
$$;

DROP FUNCTION IF EXISTS api.frob_event_urn;
CREATE FUNCTION api.frob_event_urn(event api.frob_event) RETURNS SETOF api.urn_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn(event.ilk_identifier, event.urn_identifier, event.block_height)
$$;

DROP FUNCTION IF EXISTS api.managed_cdp_urn;
CREATE FUNCTION api.managed_cdp_urn(cdp api.managed_cdp) RETURNS SETOF api.urn_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn(cdp.ilk_identifier, cdp.urn_identifier)
$$;

DROP FUNCTION IF EXISTS api.urn_state_bites;
CREATE FUNCTION api.urn_state_bites(state api.urn_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.bite_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.urn_bites(state.ilk_identifier, state.urn_identifier)
WHERE block_height <= state.block_height
LIMIT urn_state_bites.max_results
OFFSET
urn_state_bites.result_offset
$$;

DROP FUNCTION IF EXISTS api.urn_state_frobs;
CREATE FUNCTION api.urn_state_frobs(state api.urn_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.frob_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.urn_frobs(state.ilk_identifier, state.urn_identifier)
WHERE block_height <= state.block_height
LIMIT urn_state_frobs.max_results
OFFSET
urn_state_frobs.result_offset
$$;

DROP FUNCTION IF EXISTS api.urn_state_ilk;
CREATE FUNCTION api.urn_state_ilk(state api.urn_state) RETURNS api.ilk_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.ilk_snapshot i
WHERE i.ilk_identifier = state.ilk_identifier
  AND i.block_number <= state.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$;
