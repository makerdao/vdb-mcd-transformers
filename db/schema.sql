--
-- PostgreSQL database dump
--

-- Dumped from database version 11.5
-- Dumped by pg_dump version 11.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: api; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA api;


--
-- Name: maker; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA maker;


--
-- Name: citext; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;


--
-- Name: EXTENSION citext; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';


--
-- Name: bid_act; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.bid_act AS ENUM (
    'kick',
    'tick',
    'tend',
    'dent',
    'deal',
    'yank'
);


--
-- Name: bite_event; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.bite_event AS (
	ilk_identifier text,
	urn_identifier text,
	bid_id numeric,
	ink numeric,
	art numeric,
	tab numeric,
	block_height bigint,
	log_id bigint
);


--
-- Name: COLUMN bite_event.block_height; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.bite_event.block_height IS '@omit';


--
-- Name: COLUMN bite_event.log_id; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.bite_event.log_id IS '@omit';


--
-- Name: era; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.era AS (
	epoch bigint,
	iso timestamp without time zone
);


--
-- Name: flap_bid_event; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.flap_bid_event AS (
	bid_id numeric,
	lot numeric,
	bid_amount numeric,
	act api.bid_act,
	block_height bigint,
	log_id bigint,
	contract_address text
);


--
-- Name: COLUMN flap_bid_event.block_height; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.flap_bid_event.block_height IS '@omit';


--
-- Name: COLUMN flap_bid_event.log_id; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.flap_bid_event.log_id IS '@omit';


--
-- Name: COLUMN flap_bid_event.contract_address; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.flap_bid_event.contract_address IS '@omit';


--
-- Name: flap_state; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.flap_state AS (
	bid_id numeric,
	guy text,
	tic bigint,
	"end" bigint,
	lot numeric,
	bid numeric,
	dealt boolean,
	created timestamp without time zone,
	updated timestamp without time zone
);


--
-- Name: flip_bid_event; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.flip_bid_event AS (
	bid_id numeric,
	lot numeric,
	bid_amount numeric,
	act api.bid_act,
	block_height bigint,
	log_id bigint,
	contract_address text
);


--
-- Name: COLUMN flip_bid_event.block_height; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.flip_bid_event.block_height IS '@omit';


--
-- Name: COLUMN flip_bid_event.log_id; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.flip_bid_event.log_id IS '@omit';


--
-- Name: COLUMN flip_bid_event.contract_address; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.flip_bid_event.contract_address IS '@omit';


--
-- Name: flip_state; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.flip_state AS (
	block_height bigint,
	bid_id numeric,
	ilk_id integer,
	urn_id integer,
	guy text,
	tic bigint,
	"end" bigint,
	lot numeric,
	bid numeric,
	gal text,
	dealt boolean,
	tab numeric,
	created timestamp without time zone,
	updated timestamp without time zone
);


--
-- Name: COLUMN flip_state.block_height; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.flip_state.block_height IS '@omit';


--
-- Name: COLUMN flip_state.ilk_id; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.flip_state.ilk_id IS '@omit';


--
-- Name: COLUMN flip_state.urn_id; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.flip_state.urn_id IS '@omit';


--
-- Name: flop_bid_event; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.flop_bid_event AS (
	bid_id numeric,
	lot numeric,
	bid_amount numeric,
	act api.bid_act,
	block_height bigint,
	log_id bigint,
	contract_address text
);


--
-- Name: COLUMN flop_bid_event.block_height; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.flop_bid_event.block_height IS '@omit';


--
-- Name: COLUMN flop_bid_event.log_id; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.flop_bid_event.log_id IS '@omit';


--
-- Name: COLUMN flop_bid_event.contract_address; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.flop_bid_event.contract_address IS '@omit';


--
-- Name: flop_state; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.flop_state AS (
	bid_id numeric,
	guy text,
	tic bigint,
	"end" bigint,
	lot numeric,
	bid numeric,
	dealt boolean,
	created timestamp without time zone,
	updated timestamp without time zone
);


--
-- Name: frob_event; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.frob_event AS (
	ilk_identifier text,
	urn_identifier text,
	dink numeric,
	dart numeric,
	ilk_rate numeric,
	block_height bigint,
	log_id bigint
);


--
-- Name: COLUMN frob_event.block_height; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.frob_event.block_height IS '@omit';


--
-- Name: COLUMN frob_event.log_id; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.frob_event.log_id IS '@omit';


--
-- Name: ilk_file_event; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.ilk_file_event AS (
	ilk_identifier text,
	what text,
	data text,
	block_height bigint,
	log_id bigint
);


--
-- Name: COLUMN ilk_file_event.ilk_identifier; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.ilk_file_event.ilk_identifier IS '@omit';


--
-- Name: COLUMN ilk_file_event.block_height; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.ilk_file_event.block_height IS '@omit';


--
-- Name: COLUMN ilk_file_event.log_id; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.ilk_file_event.log_id IS '@omit';


--
-- Name: ilk_state; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.ilk_state AS (
	ilk_identifier text,
	block_height bigint,
	rate numeric,
	art numeric,
	spot numeric,
	line numeric,
	dust numeric,
	chop numeric,
	lump numeric,
	flip text,
	rho numeric,
	duty numeric,
	pip text,
	mat numeric,
	created timestamp without time zone,
	updated timestamp without time zone
);


--
-- Name: poke_event; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.poke_event AS (
	ilk_id integer,
	val numeric,
	spot numeric,
	block_height bigint,
	log_id bigint
);


--
-- Name: COLUMN poke_event.ilk_id; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.poke_event.ilk_id IS '@omit';


--
-- Name: COLUMN poke_event.block_height; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.poke_event.block_height IS '@omit';


--
-- Name: COLUMN poke_event.log_id; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.poke_event.log_id IS '@omit';


--
-- Name: queued_sin; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.queued_sin AS (
	era numeric,
	tab numeric,
	flogged boolean,
	created timestamp without time zone,
	updated timestamp without time zone
);


--
-- Name: relevant_block; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.relevant_block AS (
	block_height bigint,
	block_hash text,
	ilk_id integer
);


--
-- Name: sin_act; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.sin_act AS ENUM (
    'flog',
    'fess'
);


--
-- Name: sin_queue_event; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.sin_queue_event AS (
	era numeric,
	act api.sin_act,
	block_height bigint,
	log_id bigint
);


--
-- Name: COLUMN sin_queue_event.block_height; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.sin_queue_event.block_height IS '@omit';


--
-- Name: COLUMN sin_queue_event.log_id; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.sin_queue_event.log_id IS '@omit';


--
-- Name: tx; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.tx AS (
	transaction_hash text,
	transaction_index integer,
	block_height bigint,
	block_hash text,
	tx_from text,
	tx_to text
);


--
-- Name: urn_state; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.urn_state AS (
	urn_identifier text,
	ilk_identifier text,
	block_height bigint,
	ink numeric,
	art numeric,
	ratio numeric,
	safe boolean,
	created timestamp without time zone,
	updated timestamp without time zone
);


--
-- Name: all_bites(text, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_bites(ilk_identifier text, max_results integer DEFAULT '-1'::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.bite_event
    LANGUAGE sql STABLE STRICT
    AS $$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier)

SELECT ilk_identifier,
       identifier AS urn_identifier,
       bid_id,
       ink,
       art,
       tab,
       block_number,
       log_id
FROM maker.bite
         LEFT JOIN maker.urns ON bite.urn_id = urns.id
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE urns.ilk_id = (SELECT id FROM ilk)
ORDER BY urn_identifier, block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
OFFSET
all_bites.result_offset
$$;


--
-- Name: all_flap_bid_events(integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_flap_bid_events(max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flap_bid_event
    LANGUAGE sql STABLE
    AS $$
WITH address_id AS (
    SELECT address_id
    FROM maker.flap_kick
    LIMIT 1
),
     flap_address AS (
         SELECT address
         FROM maker.flap_kick
                  JOIN addresses on addresses.id = flap_kick.address_id
         LIMIT 1
     ),
     deals AS (
         SELECT deal.bid_id,
                flap_bid_lot.lot,
                flap_bid_bid.bid             AS bid_amount,
                'deal'::api.bid_act          AS act,
                headers.block_number         AS block_height,
                deal.log_id,
                (SELECT * FROM flap_address) AS contract_address
         FROM maker.deal
                  LEFT JOIN headers ON deal.header_id = headers.id
                  LEFT JOIN maker.flap_bid_bid
                            ON deal.bid_id = flap_bid_bid.bid_id
                                AND flap_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flap_bid_lot
                            ON deal.bid_id = flap_bid_lot.bid_id
                                AND flap_bid_lot.block_number = headers.block_number
         WHERE deal.address_id = (SELECT * FROM address_id)
     ),
     yanks AS (
         SELECT yank.bid_id,
                flap_bid_lot.lot,
                flap_bid_bid.bid             AS bid_amount,
                'yank'::api.bid_act          AS act,
                headers.block_number         AS block_height,
                yank.log_id,
                (SELECT * FROM flap_address) AS contract_address
         FROM maker.yank
                  LEFT JOIN headers ON yank.header_id = headers.id
                  LEFT JOIN maker.flap_bid_bid
                            ON yank.bid_id = flap_bid_bid.bid_id
                                AND flap_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flap_bid_lot
                            ON yank.bid_id = flap_bid_lot.bid_id
                                AND flap_bid_lot.block_number = headers.block_number
         WHERE yank.address_id = (SELECT * FROM address_id)
     ),
     ticks AS (
         SELECT tick.bid_id,
                flap_bid_lot.lot,
                flap_bid_bid.bid             AS bid_amount,
                'tick'::api.bid_act          AS act,
                headers.block_number         AS block_height,
                log_id,
                (SELECT * FROM flap_address) AS contract_address
         FROM maker.tick
                  LEFT JOIN headers on tick.header_id = headers.id
                  LEFT JOIN maker.flap_bid_bid
                            ON tick.bid_id = flap_bid_bid.bid_id
                                AND flap_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flap_bid_lot
                            ON tick.bid_id = flap_bid_lot.bid_id
                                AND flap_bid_lot.block_number = headers.block_number
         WHERE tick.address_id = (SELECT * FROM address_id)
     )

SELECT flap_kick.bid_id,
       lot,
       bid                          AS bid_amount,
       'kick'::api.bid_act          AS act,
       block_number                 AS block_height,
       log_id,
       (SELECT * FROM flap_address) AS contract_address
FROM maker.flap_kick
         LEFT JOIN headers ON flap_kick.header_id = headers.id
UNION
SELECT bid_id,
       lot,
       bid                          AS bid_amount,
       'tend'::api.bid_act          AS act,
       block_number                 AS block_height,
       log_id,
       (SELECT * FROM flap_address) AS contract_address
FROM maker.tend
         LEFT JOIN headers ON tend.header_id = headers.id
WHERE tend.address_id = (SELECT * FROM address_id)
UNION
SELECT *
FROM ticks
UNION
SELECT *
FROM deals
UNION
SELECT *
FROM yanks
ORDER BY block_height DESC
LIMIT all_flap_bid_events.max_results
OFFSET
all_flap_bid_events.result_offset
$$;


--
-- Name: all_flaps(integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_flaps(max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flap_state
    LANGUAGE plpgsql STABLE
    AS $$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT bid_id
            FROM maker.flap
            ORDER BY bid_id DESC
            LIMIT all_flaps.max_results OFFSET all_flaps.result_offset
        )
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flap(bid_ids.bid_id) f
    );
END
$$;


--
-- Name: all_flip_bid_events(integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_flip_bid_events(max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flip_bid_event
    LANGUAGE sql STABLE
    AS $$
WITH address_ids AS (
    SELECT distinct address_id
    FROM maker.flip_kick
),
     deals AS (
         SELECT deal.bid_id,
                flip_bid_lot.lot,
                flip_bid_bid.bid                                           AS bid_amount,
                'deal'::api.bid_act                                        AS act,
                headers.block_number                                       AS block_height,
                log_id,
                (SELECT address FROM addresses WHERE id = deal.address_id) AS contract_address
         FROM maker.deal
                  LEFT JOIN headers ON deal.header_id = headers.id
                  LEFT JOIN maker.flip_bid_bid
                            ON deal.bid_id = flip_bid_bid.bid_id
                                AND flip_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flip_bid_lot
                            ON deal.bid_id = flip_bid_lot.bid_id
                                AND flip_bid_lot.block_number = headers.block_number
         WHERE deal.address_id IN (SELECT * FROM address_ids)
     ),
     yanks AS (
         SELECT yank.bid_id,
                flip_bid_lot.lot,
                flip_bid_bid.bid     AS bid_amount,
                'yank'::api.bid_act  AS act,
                headers.block_number AS block_height,
                log_id,
                (SELECT address FROM addresses WHERE id = yank.address_id)
         FROM maker.yank
                  LEFT JOIN headers ON yank.header_id = headers.id
                  LEFT JOIN maker.flip_bid_bid
                            ON yank.bid_id = flip_bid_bid.bid_id
                                AND flip_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flip_bid_lot
                            ON yank.bid_id = flip_bid_lot.bid_id
                                AND flip_bid_lot.block_number = headers.block_number
         WHERE yank.address_id IN (SELECT * FROM address_ids)
     ),
     ticks AS (
         SELECT tick.bid_id,
                flip_bid_lot.lot,
                flip_bid_bid.bid     AS bid_amount,
                'tick'::api.bid_act  AS act,
                headers.block_number AS block_height,
                log_id,
                (SELECT address FROM addresses WHERE id = tick.address_id)
         FROM maker.tick
                  LEFT JOIN headers on tick.header_id = headers.id
                  LEFT JOIN maker.flip_bid_bid
                            ON tick.bid_id = flip_bid_bid.bid_id
                                AND flip_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flip_bid_lot
                            ON tick.bid_id = flip_bid_lot.bid_id
                                AND flip_bid_lot.block_number = headers.block_number
         WHERE tick.address_id IN (SELECT * FROM address_ids)
     )

SELECT flip_kick.bid_id,
       lot,
       bid                 AS                                          bid_amount,
       'kick'::api.bid_act AS                                          act,
       block_number        AS                                          block_height,
       log_id,
       (SELECT address FROM addresses WHERE id = flip_kick.address_id) s
FROM maker.flip_kick
         LEFT JOIN headers ON flip_kick.header_id = headers.id
UNION
SELECT bid_id,
       lot,
       bid                 AS bid_amount,
       'tend'::api.bid_act AS act,
       block_number        AS block_height,
       log_id,
       (SELECT address FROM addresses WHERE id = tend.address_id)
FROM maker.tend
         LEFT JOIN headers on tend.header_id = headers.id
WHERE tend.address_id IN (SELECT * FROM address_ids)
UNION
SELECT bid_id,
       lot,
       bid                 AS bid_amount,
       'dent'::api.bid_act AS act,
       block_number        AS block_height,
       log_id,
       (SELECT address FROM addresses WHERE id = dent.address_id)
FROM maker.dent
         LEFT JOIN headers on dent.header_id = headers.id
WHERE dent.address_id IN (SELECT * FROM address_ids)
UNION
SELECT *
from deals
UNION
SELECT *
from yanks
UNION
SELECT *
FROM ticks
ORDER BY block_height DESC
LIMIT all_flip_bid_events.max_results OFFSET all_flip_bid_events.result_offset
$$;


--
-- Name: all_flips(text, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_flips(ilk text, max_results integer DEFAULT '-1'::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flip_state
    LANGUAGE plpgsql STABLE STRICT
    AS $$
BEGIN
    RETURN QUERY (
        WITH ilk_ids AS (SELECT id
                         FROM maker.ilks
                         WHERE identifier = all_flips.ilk),
             address AS (
                 SELECT DISTINCT address_id
                 FROM maker.flip_ilk
                 WHERE flip_ilk.ilk_id = (SELECT id FROM ilk_ids)
                 LIMIT 1),
             bid_ids AS (
                 SELECT DISTINCT bid_id
                 FROM maker.flip
                 WHERE address_id = (SELECT * FROM address)
                 ORDER BY bid_id DESC
                 LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
                 OFFSET
                 all_flips.result_offset)
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flip(bid_ids.bid_id, all_flips.ilk) f
    );
END
$$;


--
-- Name: all_flop_bid_events(integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_flop_bid_events(max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flop_bid_event
    LANGUAGE sql STABLE
    AS $$
WITH address_id AS (
    SELECT address_id
    FROM maker.flop_kick
    LIMIT 1
),
     flop_address AS (
         SELECT address
         FROM maker.flop_kick
                  JOIN addresses on addresses.id = flop_kick.address_id
         LIMIT 1
     ),
     deals AS (
         SELECT deal.bid_id,
                flop_bid_lot.lot,
                flop_bid_bid.bid             AS bid_amount,
                'deal'::api.bid_act          AS act,
                headers.block_number         AS block_height,
                deal.log_id,
                (SELECT * FROM flop_address) AS contract_address
         FROM maker.deal
                  LEFT JOIN headers ON deal.header_id = headers.id
                  LEFT JOIN maker.flop_bid_bid
                            ON deal.bid_id = flop_bid_bid.bid_id
                                AND flop_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flop_bid_lot
                            ON deal.bid_id = flop_bid_lot.bid_id
                                AND flop_bid_lot.block_number = headers.block_number
         WHERE deal.address_id = (SELECT * FROM address_id)
     ),
     yanks AS (
         SELECT yank.bid_id,
                flop_bid_lot.lot,
                flop_bid_bid.bid             AS bid_amount,
                'yank'::api.bid_act          AS act,
                headers.block_number         AS block_height,
                yank.log_id,
                (SELECT * FROM flop_address) AS contract_address
         FROM maker.yank
                  LEFT JOIN headers ON yank.header_id = headers.id
                  LEFT JOIN maker.flop_bid_bid
                            ON yank.bid_id = flop_bid_bid.bid_id
                                AND flop_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flop_bid_lot
                            ON yank.bid_id = flop_bid_lot.bid_id
                                AND flop_bid_lot.block_number = headers.block_number
         WHERE yank.address_id = (SELECT * FROM address_id)
         ORDER BY block_height DESC
     ),
     ticks AS (
         SELECT tick.bid_id,
                flop_bid_lot.lot,
                flop_bid_bid.bid             AS bid_amount,
                'tick'::api.bid_act          AS act,
                headers.block_number         AS block_height,
                log_id,
                (SELECT * FROM flop_address) AS contract_address
         FROM maker.tick
                  LEFT JOIN headers on tick.header_id = headers.id
                  LEFT JOIN maker.flop_bid_bid
                            ON tick.bid_id = flop_bid_bid.bid_id
                                AND flop_bid_bid.block_number = headers.block_number
                  LEFT JOIN maker.flop_bid_lot
                            ON tick.bid_id = flop_bid_lot.bid_id
                                AND flop_bid_lot.block_number = headers.block_number
         WHERE tick.address_id = (SELECT * FROM address_id)
     )

SELECT flop_kick.bid_id,
       lot,
       bid                          AS bid_amount,
       'kick'::api.bid_act          AS act,
       block_number                 AS block_height,
       log_id,
       (SELECT * FROM flop_address) AS contract_address
FROM maker.flop_kick
         LEFT JOIN headers ON flop_kick.header_id = headers.id
UNION
SELECT bid_id,
       lot,
       bid                          AS bid_amount,
       'dent'::api.bid_act          AS act,
       block_number                 AS block_height,
       log_id,
       (SELECT * FROM flop_address) AS contract_address
FROM maker.dent
         LEFT JOIN headers ON dent.header_id = headers.id
WHERE dent.address_id = (SELECT * FROM address_id)
UNION
SELECT *
FROM deals
UNION
SELECT *
FROM yanks
UNION
SELECT *
FROM ticks
ORDER BY block_height DESC
LIMIT all_flop_bid_events.max_results OFFSET all_flop_bid_events.result_offset
$$;


--
-- Name: all_flops(integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_flops(max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flop_state
    LANGUAGE plpgsql STABLE
    AS $$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT bid_id
            FROM maker.flop
            ORDER BY bid_id DESC
            LIMIT all_flops.max_results OFFSET all_flops.result_offset
        )
        SELECT f.*
        FROM bid_ids,
             LATERAL api.get_flop(bid_ids.bid_id) f
    );
END
$$;


--
-- Name: all_frobs(text, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_frobs(ilk_identifier text, max_results integer DEFAULT '-1'::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.frob_event
    LANGUAGE sql STABLE STRICT
    AS $$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier),
     rates AS (SELECT block_number, rate
               FROM maker.vat_ilk_rate
               WHERE ilk_id = (SELECT id FROM ilk)
               ORDER BY block_number DESC
     )

SELECT ilk_identifier,
       urns.identifier                                                             AS urn_identifier,
       dink,
       dart,
       (SELECT rate from rates WHERE block_number <= headers.block_number LIMIT 1) AS ilk_rate,
       block_number,
       log_id
FROM maker.vat_frob
         LEFT JOIN maker.urns ON vat_frob.urn_id = urns.id
         LEFT JOIN headers ON vat_frob.header_id = headers.id
WHERE urns.ilk_id = (SELECT id FROM ilk)
ORDER BY block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
OFFSET
all_frobs.result_offset
$$;


--
-- Name: all_ilk_file_events(text, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_ilk_file_events(ilk_identifier text, max_results integer DEFAULT '-1'::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.ilk_file_event
    LANGUAGE sql STABLE STRICT
    AS $$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier)

SELECT ilk_identifier, what, data :: text, block_number, log_id
FROM maker.cat_file_chop_lump
         LEFT JOIN headers ON cat_file_chop_lump.header_id = headers.id
WHERE cat_file_chop_lump.ilk_id = (SELECT id FROM ilk)
UNION
SELECT ilk_identifier, what, flip AS data, block_number, log_id
FROM maker.cat_file_flip
         LEFT JOIN headers ON cat_file_flip.header_id = headers.id
WHERE cat_file_flip.ilk_id = (SELECT id FROM ilk)
UNION
SELECT ilk_identifier, what, data :: text, block_number, log_id
FROM maker.jug_file_ilk
         LEFT JOIN headers ON jug_file_ilk.header_id = headers.id
WHERE jug_file_ilk.ilk_id = (SELECT id FROM ilk)
UNION
SELECT ilk_identifier, what, data :: text, block_number, log_id
FROM maker.spot_file_mat
         LEFT JOIN headers ON spot_file_mat.header_id = headers.id
WHERE spot_file_mat.ilk_id = (SELECT id FROM ilk)
UNION
SELECT ilk_identifier, what, pip AS data, block_number, log_id
FROM maker.spot_file_pip
         LEFT JOIN headers ON spot_file_pip.header_id = headers.id
WHERE spot_file_pip.ilk_id = (SELECT id FROM ilk)
UNION
SELECT ilk_identifier, what, data :: text, block_number, log_id
FROM maker.vat_file_ilk
         LEFT JOIN headers ON vat_file_ilk.header_id = headers.id
WHERE vat_file_ilk.ilk_id = (SELECT id FROM ilk)
ORDER BY block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
OFFSET
all_ilk_file_events.result_offset
$$;


--
-- Name: max_block(); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.max_block() RETURNS bigint
    LANGUAGE sql STABLE
    AS $$
SELECT max(block_number)
FROM public.headers
$$;


--
-- Name: FUNCTION max_block(); Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON FUNCTION api.max_block() IS '@omit';


--
-- Name: all_ilk_states(text, bigint, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_ilk_states(ilk_identifier text, block_height bigint DEFAULT api.max_block(), max_results integer DEFAULT '-1'::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.ilk_state
    LANGUAGE plpgsql STABLE STRICT
    AS $$
BEGIN
    RETURN QUERY (
        WITH relevant_blocks AS (
            SELECT get_ilk_blocks_before.block_height
            FROM api.get_ilk_blocks_before(ilk_identifier, all_ilk_states.block_height)
        )
        SELECT r.*
        FROM relevant_blocks,
             LATERAL api.get_ilk(ilk_identifier, relevant_blocks.block_height) r
        LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
        OFFSET
        all_ilk_states.result_offset
    );
END;
$$;


--
-- Name: all_ilks(bigint, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_ilks(block_height bigint DEFAULT api.max_block(), max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.ilk_state
    LANGUAGE sql STABLE
    AS $$
WITH rates AS (SELECT DISTINCT ON (ilk_id) rate, ilk_id, block_hash
               FROM maker.vat_ilk_rate
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     arts AS (SELECT DISTINCT ON (ilk_id) art, ilk_id, block_hash
              FROM maker.vat_ilk_art
              WHERE block_number <= all_ilks.block_height
              ORDER BY ilk_id, block_number DESC),
     spots AS (SELECT DISTINCT ON (ilk_id) spot, ilk_id, block_hash
               FROM maker.vat_ilk_spot
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     lines AS (SELECT DISTINCT ON (ilk_id) line, ilk_id, block_hash
               FROM maker.vat_ilk_line
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     dusts AS (SELECT DISTINCT ON (ilk_id) dust, ilk_id, block_hash
               FROM maker.vat_ilk_dust
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     chops AS (SELECT DISTINCT ON (ilk_id) chop, ilk_id, block_hash
               FROM maker.cat_ilk_chop
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     lumps AS (SELECT DISTINCT ON (ilk_id) lump, ilk_id, block_hash
               FROM maker.cat_ilk_lump
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     flips AS (SELECT DISTINCT ON (ilk_id) flip, ilk_id, block_hash
               FROM maker.cat_ilk_flip
               WHERE block_number <= all_ilks.block_height
               ORDER BY ilk_id, block_number DESC),
     rhos AS (SELECT DISTINCT ON (ilk_id) rho, ilk_id, block_hash
              FROM maker.jug_ilk_rho
              WHERE block_number <= all_ilks.block_height
              ORDER BY ilk_id, block_number DESC),
     duties AS (SELECT DISTINCT ON (ilk_id) duty, ilk_id, block_hash
                FROM maker.jug_ilk_duty
                WHERE block_number <= all_ilks.block_height
                ORDER BY ilk_id, block_number DESC),
     pips AS (SELECT DISTINCT ON (ilk_id) pip, ilk_id, block_hash
              FROM maker.spot_ilk_pip
              WHERE block_number <= all_ilks.block_height
              ORDER BY ilk_id, block_number DESC),
     mats AS (SELECT DISTINCT ON (ilk_id) mat, ilk_id, block_hash
              FROM maker.spot_ilk_mat
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
       (SELECT api.epoch_to_datetime(h.block_timestamp) AS created
        FROM api.get_ilk_blocks_before(ilks.identifier, all_ilks.block_height) b
                 JOIN headers h on h.block_number = b.block_height
        ORDER BY h.block_number ASC
        LIMIT 1),
       (SELECT api.epoch_to_datetime(h.block_timestamp) AS updated
        FROM api.get_ilk_blocks_before(ilks.identifier, all_ilks.block_height) b
                 JOIN headers h on h.block_number = b.block_height
        ORDER BY h.block_number DESC
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
LIMIT all_ilks.max_results OFFSET all_ilks.result_offset
$$;


--
-- Name: max_timestamp(); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.max_timestamp() RETURNS numeric
    LANGUAGE sql STABLE
    AS $$
SELECT max(block_timestamp)
FROM public.headers
$$;


--
-- Name: all_poke_events(numeric, numeric, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_poke_events(begintime numeric DEFAULT 0, endtime numeric DEFAULT api.max_timestamp(), max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.poke_event
    LANGUAGE sql STABLE
    AS $$
SELECT ilk_id, "value" AS val, spot, block_number AS block_height, log_id
FROM maker.spot_poke
         LEFT JOIN public.headers ON spot_poke.header_id = headers.id
WHERE block_timestamp BETWEEN beginTime AND endTime
ORDER BY block_height DESC
LIMIT all_poke_events.max_results OFFSET all_poke_events.result_offset
$$;


--
-- Name: all_queued_sin(integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_queued_sin(max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.queued_sin
    LANGUAGE plpgsql STABLE
    AS $$
BEGIN
    RETURN QUERY (
        WITH eras AS (
            SELECT DISTINCT era
            FROM maker.vow_sin_mapping
            ORDER BY era DESC
            LIMIT all_queued_sin.max_results OFFSET all_queued_sin.result_offset
        )
        SELECT sin.*
        FROM eras,
             LATERAL api.get_queued_sin(eras.era) sin
    );
END
$$;


--
-- Name: all_sin_queue_events(numeric, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_sin_queue_events(era numeric, max_results integer DEFAULT '-1'::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.sin_queue_event
    LANGUAGE sql STABLE STRICT
    AS $$
SELECT block_timestamp AS era, 'fess' :: api.sin_act AS act, block_number AS block_height, log_id
FROM maker.vow_fess
         LEFT JOIN headers ON vow_fess.header_id = headers.id
WHERE block_timestamp = all_sin_queue_events.era
UNION
SELECT era, 'flog' :: api.sin_act AS act, block_number AS block_height, log_id
FROM maker.vow_flog
         LEFT JOIN headers ON vow_flog.header_id = headers.id
WHERE vow_flog.era = all_sin_queue_events.era
ORDER BY block_height DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
OFFSET
all_sin_queue_events.result_offset
$$;


--
-- Name: all_urn_states(text, text, bigint, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_urn_states(ilk_identifier text, urn_identifier text, block_height bigint DEFAULT api.max_block(), max_results integer DEFAULT '-1'::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.urn_state
    LANGUAGE plpgsql STABLE STRICT
    AS $$
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
                 WHERE vat_urn_ink.urn_id = (SELECT * FROM urn_id)
                   AND block_number <= all_urn_states.block_height
                 UNION
                 SELECT block_number
                 FROM maker.vat_urn_art
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
$$;


--
-- Name: all_urns(bigint, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_urns(block_height bigint DEFAULT api.max_block(), max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.urn_state
    LANGUAGE sql STABLE
    AS $$
WITH urns AS (SELECT urns.id AS urn_id, ilks.id AS ilk_id, ilks.ilk, urns.identifier
              FROM maker.urns urns
                       LEFT JOIN maker.ilks ilks ON urns.ilk_id = ilks.id),
     inks AS ( -- Latest ink for each urn
         SELECT DISTINCT ON (urn_id) urn_id, ink, block_number
         FROM maker.vat_urn_ink
         WHERE block_number <= all_urns.block_height
         ORDER BY urn_id, block_number DESC),
     arts AS ( -- Latest art for each urn
         SELECT DISTINCT ON (urn_id) urn_id, art, block_number
         FROM maker.vat_urn_art
         WHERE block_number <= all_urns.block_height
         ORDER BY urn_id, block_number DESC),
     rates AS ( -- Latest rate for each ilk
         SELECT DISTINCT ON (ilk_id) ilk_id, rate, block_number
         FROM maker.vat_ilk_rate
         WHERE block_number <= all_urns.block_height
         ORDER BY ilk_id, block_number DESC),
     spots AS ( -- Get latest price update for ilk. Problematic from update frequency, slow query?
         SELECT DISTINCT ON (ilk_id) ilk_id, spot, block_number
         FROM maker.vat_ilk_spot
         WHERE block_number <= all_urns.block_height
         ORDER BY ilk_id, block_number DESC),
     ratio_data AS (SELECT urns.ilk, urns.identifier, inks.ink, spots.spot, arts.art, rates.rate
                    FROM inks
                             JOIN urns ON inks.urn_id = urns.urn_id
                             JOIN arts ON arts.urn_id = inks.urn_id
                             JOIN spots ON spots.ilk_id = urns.ilk_id
                             JOIN rates ON rates.ilk_id = spots.ilk_id),
     ratios AS (SELECT ilk, identifier as urn_identifier, ((1.0 * ink * spot) / NULLIF(art * rate, 0)) AS ratio
                FROM ratio_data),
     safe AS (SELECT ilk, urn_identifier, (ratio >= 1) AS safe FROM ratios),
     created AS (SELECT urn_id, api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM (SELECT DISTINCT ON (urn_id) urn_id, block_hash
                       FROM maker.vat_urn_ink
                       ORDER BY urn_id, block_number ASC) earliest_blocks
                          LEFT JOIN public.headers ON hash = block_hash),
     updated AS (SELECT DISTINCT ON (urn_id) urn_id, api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM ((SELECT DISTINCT ON (urn_id) urn_id, block_hash
                        FROM maker.vat_urn_ink
                        WHERE block_number <= block_height
                        ORDER BY urn_id, block_number DESC)
                       UNION
                       (SELECT DISTINCT ON (urn_id) urn_id, block_hash
                        FROM maker.vat_urn_art
                        WHERE block_number <= block_height
                        ORDER BY urn_id, block_number DESC)) last_blocks
                          LEFT JOIN public.headers ON headers.hash = last_blocks.block_hash
                 ORDER BY urn_id, headers.block_timestamp DESC)

SELECT urns.identifier,
       ilks.identifier,
       all_urns.block_height,
       inks.ink,
       COALESCE(arts.art, 0),
       ratios.ratio,
       COALESCE(safe.safe, COALESCE(arts.art, 0) = 0),
       created.datetime,
       updated.datetime
FROM inks
         LEFT JOIN arts ON arts.urn_id = inks.urn_id
         LEFT JOIN urns ON inks.urn_id = urns.urn_id
         LEFT JOIN ratios ON ratios.urn_identifier = urns.identifier
         LEFT JOIN safe ON safe.urn_identifier = urns.identifier
         LEFT JOIN created ON created.urn_id = urns.urn_id
         LEFT JOIN updated ON updated.urn_id = urns.urn_id
         LEFT JOIN maker.ilks ON ilks.id = urns.ilk_id
ORDER BY updated DESC
LIMIT all_urns.max_results OFFSET all_urns.result_offset
$$;


--
-- Name: bite_event_bid(api.bite_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bite_event_bid(event api.bite_event) RETURNS api.flip_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_flip(event.bid_id, event.ilk_identifier, event.block_height)
$$;


--
-- Name: bite_event_ilk(api.bite_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bite_event_ilk(event api.bite_event) RETURNS api.ilk_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_ilk(event.ilk_identifier, event.block_height)
$$;


--
-- Name: bite_event_tx(api.bite_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bite_event_tx(event api.bite_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT txs.hash, txs.tx_index, headers.block_number, headers.hash, tx_from, tx_to
FROM public.header_sync_transactions txs
         LEFT JOIN headers ON txs.header_id = headers.id
         LEFT JOIN header_sync_logs ON txs.tx_index = header_sync_logs.tx_index
WHERE headers.block_number <= event.block_height
  AND header_sync_logs.id = event.log_id
ORDER BY block_number DESC
$$;


--
-- Name: bite_event_urn(api.bite_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bite_event_urn(event api.bite_event) RETURNS api.urn_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn(event.ilk_identifier, event.urn_identifier, event.block_height)
$$;


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: current_ilk_state; Type: TABLE; Schema: api; Owner: -
--

CREATE TABLE api.current_ilk_state (
    ilk_identifier text NOT NULL,
    rate numeric,
    art numeric,
    spot numeric,
    line numeric,
    dust numeric,
    chop numeric,
    lump numeric,
    flip text,
    rho numeric,
    duty numeric,
    pip text,
    mat numeric,
    created timestamp without time zone,
    updated timestamp without time zone
);


--
-- Name: TABLE current_ilk_state; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON TABLE api.current_ilk_state IS '@omit create,update,delete';


--
-- Name: COLUMN current_ilk_state.ilk_identifier; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.current_ilk_state.ilk_identifier IS '@name id';


--
-- Name: current_ilk_state_bites(api.current_ilk_state, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.current_ilk_state_bites(state api.current_ilk_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.bite_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_bites(state.ilk_identifier)
ORDER BY block_height DESC
LIMIT max_results
OFFSET
result_offset
$$;


--
-- Name: current_ilk_state_frobs(api.current_ilk_state, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.current_ilk_state_frobs(state api.current_ilk_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.frob_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_frobs(state.ilk_identifier)
ORDER BY block_height DESC
LIMIT max_results
OFFSET
result_offset
$$;


--
-- Name: current_ilk_state_ilk_file_events(api.current_ilk_state, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.current_ilk_state_ilk_file_events(state api.current_ilk_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.ilk_file_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_ilk_file_events(state.ilk_identifier)
LIMIT max_results
OFFSET
result_offset
$$;


--
-- Name: epoch_to_datetime(numeric); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.epoch_to_datetime(epoch numeric) RETURNS timestamp without time zone
    LANGUAGE sql IMMUTABLE
    AS $$
SELECT TIMESTAMP 'epoch' + epoch * INTERVAL '1 second' AS datetime
$$;


--
-- Name: FUNCTION epoch_to_datetime(epoch numeric); Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON FUNCTION api.epoch_to_datetime(epoch numeric) IS '@omit';


--
-- Name: flap_bid_event_bid(api.flap_bid_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flap_bid_event_bid(event api.flap_bid_event) RETURNS api.flap_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_flap(event.bid_id, event.block_height)
$$;


--
-- Name: flap_bid_event_tx(api.flap_bid_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flap_bid_event_tx(event api.flap_bid_event) RETURNS SETOF api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM get_tx_data(event.block_height, event.log_id)
$$;


--
-- Name: flap_state_bid_events(api.flap_state, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flap_state_bid_events(flap api.flap_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flap_bid_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_flap_bid_events() bids
WHERE bid_id = flap.bid_id
ORDER BY bids.block_height DESC
LIMIT flap_state_bid_events.max_results OFFSET flap_state_bid_events.result_offset
$$;


--
-- Name: flip_bid_event_bid(api.flip_bid_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flip_bid_event_bid(event api.flip_bid_event) RETURNS api.flip_state
    LANGUAGE sql STABLE
    AS $$
WITH ilks AS (
    SELECT ilks.identifier
    FROM maker.flip_ilk
             LEFT JOIN maker.ilks ON ilks.id = flip_ilk.ilk_id
    WHERE flip_ilk.address_id = (SELECT id FROM addresses WHERE address = event.contract_address)
    LIMIT 1
)
SELECT *
FROM api.get_flip(event.bid_id, (SELECT identifier FROM ilks))
$$;


--
-- Name: flip_bid_event_tx(api.flip_bid_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flip_bid_event_tx(event api.flip_bid_event) RETURNS SETOF api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM get_tx_data(event.block_height, event.log_id)
$$;


--
-- Name: flip_state_bid_events(api.flip_state, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flip_state_bid_events(flip api.flip_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flip_bid_event
    LANGUAGE sql STABLE
    AS $$
WITH address_ids AS ( -- get the contract address from flip_ilk table using the ilk_id from flip
    SELECT address_id
    FROM maker.flip_ilk
    WHERE ilk_id = flip.ilk_id
    LIMIT 1
),
     addresses AS (
         SELECT address
         FROM public.addresses
         WHERE id = (SELECT address_id FROM address_ids)
     )
SELECT bid_id, lot, bid_amount, act, block_height, events.log_id, events.contract_address
FROM api.all_flip_bid_events() AS events
WHERE bid_id = flip.bid_id
  AND contract_address = (SELECT address FROM addresses)
ORDER BY block_height DESC
LIMIT flip_state_bid_events.max_results OFFSET flip_state_bid_events.result_offset
$$;


--
-- Name: flip_state_ilk(api.flip_state); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flip_state_ilk(flip api.flip_state) RETURNS api.ilk_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_ilk((SELECT identifier FROM maker.ilks WHERE ilks.id = flip.ilk_id), flip.block_height)
$$;


--
-- Name: flip_state_urn(api.flip_state); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flip_state_urn(flip api.flip_state) RETURNS SETOF api.urn_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn((SELECT identifier FROM maker.ilks WHERE ilks.id = flip.ilk_id),
                 (SELECT identifier FROM maker.urns WHERE urns.id = flip.urn_id), flip.block_height)
$$;


--
-- Name: flop_bid_event_bid(api.flop_bid_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flop_bid_event_bid(event api.flop_bid_event) RETURNS api.flop_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_flop(event.bid_id, event.block_height)
$$;


--
-- Name: flop_bid_event_tx(api.flop_bid_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flop_bid_event_tx(event api.flop_bid_event) RETURNS SETOF api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM get_tx_data(event.block_height, event.log_id)
$$;


--
-- Name: flop_state_bid_events(api.flop_state, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flop_state_bid_events(flop api.flop_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flop_bid_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_flop_bid_events() bids
WHERE bid_id = flop.bid_id
ORDER BY bids.block_height DESC
LIMIT flop_state_bid_events.max_results OFFSET flop_state_bid_events.result_offset
$$;


--
-- Name: frob_event_ilk(api.frob_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.frob_event_ilk(event api.frob_event) RETURNS api.ilk_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_ilk(event.ilk_identifier, event.block_height)
$$;


--
-- Name: frob_event_tx(api.frob_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.frob_event_tx(event api.frob_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT * FROM get_tx_data(event.block_height, event.log_id)
$$;


--
-- Name: frob_event_urn(api.frob_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.frob_event_urn(event api.frob_event) RETURNS SETOF api.urn_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn(event.ilk_identifier, event.urn_identifier, event.block_height)
$$;


--
-- Name: get_flap(numeric, bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_flap(bid_id numeric, block_height bigint DEFAULT api.max_block()) RETURNS api.flap_state
    LANGUAGE sql STABLE STRICT
    AS $$
WITH address_id AS (
    SELECT address_id
    FROM maker.flap
    WHERE flap.bid_id = get_flap.bid_id
      AND block_number <= block_height
    LIMIT 1
),
     storage_values AS (
         SELECT bid_id,
                guy,
                tic,
                "end",
                lot,
                bid,
                created,
                updated
         FROM maker.flap
         WHERE bid_id = get_flap.bid_id
           AND block_number <= block_height
         ORDER BY block_number DESC
         LIMIT 1
     ),
     deal AS (
         SELECT deal, bid_id
         FROM maker.deal
                  LEFT JOIN public.headers ON deal.header_id = headers.id
         WHERE deal.bid_id = get_flap.bid_id
           AND deal.address_id = (SELECT * FROM address_id)
           AND headers.block_number <= block_height
         ORDER BY bid_id, block_number DESC
         LIMIT 1
     )

SELECT get_flap.bid_id,
       storage_values.guy,
       storage_values.tic,
       storage_values."end",
       storage_values.lot,
       storage_values.bid,
       CASE (SELECT COUNT(*) FROM deal)
           WHEN 0 THEN FALSE
           ELSE TRUE
           END AS dealt,
       storage_values.created,
       storage_values.updated
FROM storage_values
$$;


--
-- Name: get_flip(numeric, text, bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_flip(bid_id numeric, ilk text, block_height bigint DEFAULT api.max_block()) RETURNS api.flip_state
    LANGUAGE sql STABLE STRICT
    AS $$
WITH ilk_ids AS (SELECT id FROM maker.ilks WHERE ilks.identifier = get_flip.ilk),
     -- there should only ever be 1 address for a given ilk, which is why there's a LIMIT with no ORDER BY
     address_id AS (SELECT address_id
                    FROM maker.flip_ilk
                    WHERE flip_ilk.ilk_id = (SELECT id FROM ilk_ids)
                      AND block_number <= block_height
                    LIMIT 1),
     kicks AS (SELECT usr
               FROM maker.flip_kick
               WHERE flip_kick.bid_id = get_flip.bid_id
                 AND address_id = (SELECT * FROM address_id)
               LIMIT 1),
     urn_id AS (SELECT id
                FROM maker.urns
                WHERE urns.ilk_id = (SELECT id FROM ilk_ids)
                  AND urns.identifier = (SELECT usr FROM kicks)),

     storage_values AS (
         SELECT guy,
                tic,
                "end",
                lot,
                bid,
                gal,
                tab,
                created,
                updated
         FROM maker.flip
         WHERE bid_id = get_flip.bid_id
           AND block_number <= block_height
         ORDER BY block_number DESC
         LIMIT 1
     ),
     deals AS (SELECT deal.bid_id
               FROM maker.deal
                        LEFT JOIN public.headers ON deal.header_id = headers.id
               WHERE deal.bid_id = get_flip.bid_id
                 AND deal.address_id = (SELECT * FROM address_id)
                 AND headers.block_number <= block_height)

SELECT get_flip.block_height,
       get_flip.bid_id,
       (SELECT id FROM ilk_ids),
       (SELECT id FROM urn_id),
       storage_values.guy,
       storage_values.tic,
       storage_values."end",
       storage_values.lot,
       storage_values.bid,
       storage_values.gal,
       CASE (SELECT COUNT(*) FROM deals)
           WHEN 0 THEN FALSE
           ELSE TRUE END,
       storage_values.tab,
       storage_values.created,
       storage_values.updated
FROM storage_values
$$;


--
-- Name: get_flop(numeric, bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_flop(bid_id numeric, block_height bigint DEFAULT api.max_block()) RETURNS api.flop_state
    LANGUAGE sql STABLE STRICT
    AS $$
WITH address_id AS (
    SELECT address_id
    FROM maker.flop
    WHERE flop.bid_id = get_flop.bid_id
      AND block_number <= block_height
    LIMIT 1
),
     storage_values AS (
         SELECT bid_id,
                guy,
                tic,
                "end",
                lot,
                bid,
                created,
                updated
         FROM maker.flop
         WHERE bid_id = get_flop.bid_id
           AND block_number <= block_height
         ORDER BY block_number DESC
         LIMIT 1
     ),
     deal AS (
         SELECT deal.bid_id
         FROM maker.deal
                  LEFT JOIN public.headers ON deal.header_id = headers.id
         WHERE deal.bid_id = get_flop.bid_id
           AND deal.address_id = (SELECT address_id FROM address_id)
           AND headers.block_number <= block_height
         ORDER BY bid_id, block_number DESC
         LIMIT 1
     )

SELECT get_flop.bid_id,
       storage_values.guy,
       storage_values.tic,
       storage_values."end",
       storage_values.lot,
       storage_values.bid,
       CASE (SELECT COUNT(*) FROM deal)
           WHEN 0 THEN FALSE
           ELSE TRUE
           END AS dealt,
       storage_values.created,
       storage_values.updated
FROM storage_values
$$;


--
-- Name: get_ilk(text, bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_ilk(ilk_identifier text, block_height bigint DEFAULT api.max_block()) RETURNS api.ilk_state
    LANGUAGE sql STABLE STRICT
    AS $$
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
     pips AS (SELECT pip, ilk_id, block_hash
              FROM maker.spot_ilk_pip
              WHERE ilk_id = (SELECT id FROM ilk)
                AND block_number <= get_ilk.block_height
              ORDER BY ilk_id, block_number DESC
              LIMIT 1),
     mats AS (SELECT mat, ilk_id, block_hash
              FROM maker.spot_ilk_mat
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
       pips.pip,
       mats.mat,
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
         LEFT JOIN pips ON pips.ilk_id = ilks.id
         LEFT JOIN mats ON mats.ilk_id = ilks.id
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
              duties.duty is not null OR
              pips.pip is not null OR
              mats.mat is not null
          )
$$;


--
-- Name: get_ilk_blocks_before(text, bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_ilk_blocks_before(ilk_identifier text, block_height bigint) RETURNS SETOF api.relevant_block
    LANGUAGE sql STABLE
    AS $$
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
$$;


--
-- Name: FUNCTION get_ilk_blocks_before(ilk_identifier text, block_height bigint); Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON FUNCTION api.get_ilk_blocks_before(ilk_identifier text, block_height bigint) IS '@omit';


--
-- Name: get_queued_sin(numeric); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_queued_sin(era numeric) RETURNS api.queued_sin
    LANGUAGE sql STABLE STRICT
    AS $$
WITH created AS (SELECT era, vow_sin_mapping.block_number, api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM maker.vow_sin_mapping
                          LEFT JOIN public.headers h ON h.block_number = vow_sin_mapping.block_number
                 WHERE era = get_queued_sin.era
                 ORDER BY vow_sin_mapping.block_number ASC
                 LIMIT 1),
     updated AS (SELECT era, vow_sin_mapping.block_number, api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM maker.vow_sin_mapping
                          LEFT JOIN public.headers h ON h.block_number = vow_sin_mapping.block_number
                 WHERE era = get_queued_sin.era
                 ORDER BY vow_sin_mapping.block_number DESC
                 LIMIT 1)

SELECT get_queued_sin.era,
       tab,
       (SELECT EXISTS(SELECT id FROM maker.vow_flog WHERE vow_flog.era = get_queued_sin.era)) AS flogged,
       created.datetime,
       updated.datetime
FROM maker.vow_sin_mapping
         LEFT JOIN created ON created.era = vow_sin_mapping.era
         LEFT JOIN updated ON updated.era = vow_sin_mapping.era
WHERE vow_sin_mapping.era = get_queued_sin.era
ORDER BY vow_sin_mapping.block_number DESC
$$;


--
-- Name: get_urn(text, text, bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_urn(ilk_identifier text, urn_identifier text, block_height bigint DEFAULT api.max_block()) RETURNS api.urn_state
    LANGUAGE sql STABLE STRICT
    AS $_$
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
     ratio_data AS (SELECT urn.ilk, urn.identifier, ink, spot, art, rate
                    FROM ink
                             JOIN urn ON ink.urn_id = urn.urn_id
                             JOIN art ON art.urn_id = ink.urn_id
                             JOIN spot ON spot.ilk_id = urn.ilk_id
                             JOIN rate ON rate.ilk_id = spot.ilk_id),
     ratio AS (SELECT ilk, identifier as urn_identifier, ((1.0 * ink * spot) / NULLIF(art * rate, 0)) AS ratio FROM ratio_data),
     safe AS (SELECT ilk, urn_identifier, (ratio >= 1) AS safe FROM ratio),
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
       ratio.ratio,
       COALESCE(safe.safe, COALESCE(art.art, 0) = 0),
       created.datetime,
       updated.datetime
FROM ink
         LEFT JOIN art ON art.urn_id = ink.urn_id
         LEFT JOIN urn ON urn.urn_id = ink.urn_id
         LEFT JOIN ratio ON ratio.ilk = urn.ilk AND ratio.urn_identifier = urn.identifier
         LEFT JOIN safe ON safe.ilk = ratio.ilk AND safe.urn_identifier = ratio.urn_identifier
         LEFT JOIN created ON created.urn_id = art.urn_id
         LEFT JOIN updated ON updated.urn_id = art.urn_id
WHERE ink.urn_id IS NOT NULL
$_$;


--
-- Name: ilk_file_event_ilk(api.ilk_file_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_file_event_ilk(event api.ilk_file_event) RETURNS SETOF api.ilk_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_ilk(event.ilk_identifier, event.block_height)
$$;


--
-- Name: ilk_file_event_tx(api.ilk_file_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_file_event_tx(event api.ilk_file_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT * FROM get_tx_data(event.block_height, event.log_id)
$$;


--
-- Name: ilk_state_bites(api.ilk_state, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_state_bites(state api.ilk_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.bite_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_bites(state.ilk_identifier)
WHERE block_height <= state.block_height
ORDER BY block_height DESC
LIMIT ilk_state_bites.max_results OFFSET ilk_state_bites.result_offset
$$;


--
-- Name: ilk_state_frobs(api.ilk_state, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_state_frobs(state api.ilk_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.frob_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_frobs(state.ilk_identifier)
WHERE block_height <= state.block_height
ORDER BY block_height DESC
LIMIT ilk_state_frobs.max_results OFFSET ilk_state_frobs.result_offset
$$;


--
-- Name: ilk_state_ilk_file_events(api.ilk_state, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_state_ilk_file_events(state api.ilk_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.ilk_file_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_ilk_file_events(state.ilk_identifier)
WHERE block_height <= state.block_height
LIMIT ilk_state_ilk_file_events.max_results OFFSET ilk_state_ilk_file_events.result_offset
$$;


--
-- Name: managed_cdp; Type: TABLE; Schema: api; Owner: -
--

CREATE TABLE api.managed_cdp (
    id integer NOT NULL,
    cdpi numeric,
    usr text,
    urn_identifier text,
    ilk_identifier text,
    created timestamp without time zone
);


--
-- Name: TABLE managed_cdp; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON TABLE api.managed_cdp IS '@omit create,update,delete';


--
-- Name: COLUMN managed_cdp.id; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.managed_cdp.id IS '@omit';


--
-- Name: COLUMN managed_cdp.cdpi; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.managed_cdp.cdpi IS '@name id';


--
-- Name: managed_cdp_ilk(api.managed_cdp); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.managed_cdp_ilk(cdp api.managed_cdp) RETURNS api.ilk_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_ilk(cdp.ilk_identifier)
$$;


--
-- Name: managed_cdp_urn(api.managed_cdp); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.managed_cdp_urn(cdp api.managed_cdp) RETURNS SETOF api.urn_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn(cdp.ilk_identifier, cdp.urn_identifier)
$$;


--
-- Name: poke_event_ilk(api.poke_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.poke_event_ilk(priceupdate api.poke_event) RETURNS api.ilk_state
    LANGUAGE sql STABLE
    AS $$
WITH raw_ilk AS (SELECT * FROM maker.ilks WHERE ilks.id = priceUpdate.ilk_id)

SELECT *
FROM api.get_ilk((SELECT identifier FROM raw_ilk), priceUpdate.block_height)
$$;


--
-- Name: poke_event_tx(api.poke_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.poke_event_tx(priceupdate api.poke_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT * FROM get_tx_data(priceUpdate.block_height, priceUpdate.log_id)
$$;


--
-- Name: queued_sin_sin_queue_events(api.queued_sin, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.queued_sin_sin_queue_events(state api.queued_sin, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.sin_queue_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_sin_queue_events(state.era)
LIMIT queued_sin_sin_queue_events.max_results OFFSET queued_sin_sin_queue_events.result_offset
$$;


--
-- Name: sin_queue_event_tx(api.sin_queue_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.sin_queue_event_tx(event api.sin_queue_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT * FROM get_tx_data(event.block_height, event.log_id)
$$;


--
-- Name: total_ink(text, bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.total_ink(ilk_identifier text, block_height bigint DEFAULT api.max_block()) RETURNS numeric
    LANGUAGE sql STABLE STRICT
    AS $$
SELECT SUM(latest_ink_by_urn.ink)
FROM (SELECT DISTINCT ON (vat_urn_ink.urn_id) vat_urn_ink.ink
      FROM maker.ilks
               LEFT JOIN maker.urns ON urns.ilk_id = ilks.id
               LEFT JOIN maker.vat_urn_ink ON vat_urn_ink.urn_id = urns.id
      WHERE ilks.identifier = total_ink.ilk_identifier
        AND vat_urn_ink.block_number <= total_ink.block_height
      ORDER BY vat_urn_ink.urn_id, vat_urn_ink.block_number DESC) latest_ink_by_urn
$$;


--
-- Name: tx_era(api.tx); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.tx_era(tx api.tx) RETURNS api.era
    LANGUAGE sql STABLE
    AS $$
SELECT block_timestamp :: BIGINT AS "epoch", api.epoch_to_datetime(block_timestamp) AS iso
FROM headers
WHERE block_number = tx.block_height
$$;


--
-- Name: urn_bites(text, text, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.urn_bites(ilk_identifier text, urn_identifier text, max_results integer DEFAULT '-1'::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.bite_event
    LANGUAGE sql STABLE STRICT
    AS $$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier),
     urn AS (SELECT id
             FROM maker.urns
             WHERE ilk_id = (SELECT id FROM ilk)
               AND identifier = urn_bites.urn_identifier)

SELECT ilk_identifier,
       urn_bites.urn_identifier,
       bid_id,
       ink,
       art,
       tab,
       block_number,
       log_id
FROM maker.bite
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE bite.urn_id = (SELECT id FROM urn)
ORDER BY block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
OFFSET
urn_bites.result_offset
$$;


--
-- Name: urn_frobs(text, text, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.urn_frobs(ilk_identifier text, urn_identifier text, max_results integer DEFAULT '-1'::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.frob_event
    LANGUAGE sql STABLE STRICT
    AS $$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier),
     urn AS (SELECT id
             FROM maker.urns
             WHERE ilk_id = (SELECT id FROM ilk)
               AND identifier = urn_identifier),
     rates AS (SELECT block_number, rate
               FROM maker.vat_ilk_rate
               WHERE ilk_id = (SELECT id FROM ilk)
               ORDER BY block_number DESC
     )

SELECT ilk_identifier,
       urn_identifier,
       dink,
       dart,
       (SELECT rate from rates WHERE block_number <= headers.block_number LIMIT 1) AS ilk_rate,
       headers.block_number,
       log_id
FROM maker.vat_frob
         LEFT JOIN headers ON vat_frob.header_id = headers.id
WHERE vat_frob.urn_id = (SELECT id FROM urn)
ORDER BY block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
OFFSET
urn_frobs.result_offset
$$;


--
-- Name: urn_state_bites(api.urn_state, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.urn_state_bites(state api.urn_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.bite_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.urn_bites(state.ilk_identifier, state.urn_identifier)
WHERE block_height <= state.block_height
LIMIT urn_state_bites.max_results OFFSET urn_state_bites.result_offset
$$;


--
-- Name: urn_state_frobs(api.urn_state, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.urn_state_frobs(state api.urn_state, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.frob_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.urn_frobs(state.ilk_identifier, state.urn_identifier)
WHERE block_height <= state.block_height
LIMIT urn_state_frobs.max_results OFFSET urn_state_frobs.result_offset
$$;


--
-- Name: urn_state_ilk(api.urn_state); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.urn_state_ilk(state api.urn_state) RETURNS api.ilk_state
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_ilk(state.ilk_identifier, state.block_height)
$$;


--
-- Name: flap_created(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.flap_created() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH block_info AS (
        SELECT block_number, hash, api.epoch_to_datetime(headers.block_timestamp) AS datetime
        FROM public.headers
        WHERE headers.id = NEW.header_id
        LIMIT 1
    )
    INSERT
    INTO maker.flap(bid_id, address_id, block_number, block_hash, created, updated, bid, guy, tic, "end", lot)
    VALUES (NEW.bid_id, NEW.address_id,
            (SELECT block_number FROM block_info),
            (SELECT hash FROM block_info),
            (SELECT datetime FROM block_info),
            (SELECT datetime FROM block_info),
            (SELECT get_latest_flap_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flap_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flap_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flap_bid_end(NEW.bid_id)),
            (SELECT get_latest_flap_bid_lot(NEW.bid_id)))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET created = (SELECT datetime FROM block_info),
                                                     updated = (SELECT datetime FROM block_info);
    return NEW;
END
$$;


--
-- Name: flip_created(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.flip_created() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH block_info AS (
        SELECT block_number, hash, api.epoch_to_datetime(headers.block_timestamp) AS datetime
        FROM public.headers
        WHERE headers.id = NEW.header_id
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, created, updated, guy, tic, "end", lot, bid,
                    gal, tab)
    VALUES (NEW.bid_id, NEW.address_id,
            (SELECT block_number FROM block_info),
            (SELECT hash FROM block_info),
            (SELECT datetime FROM block_info),
            (SELECT datetime FROM block_info),
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET created = (SELECT datetime FROM block_info),
                                                     updated = (SELECT datetime FROM block_info);
    return NEW;
END
$$;


--
-- Name: flop_created(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.flop_created() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH block_info AS (
        SELECT block_number, hash, api.epoch_to_datetime(headers.block_timestamp) AS datetime
        FROM public.headers
        WHERE headers.id = NEW.header_id
        LIMIT 1
    )
    INSERT
    INTO maker.flop(bid_id, address_id, block_number, block_hash, created, updated, bid, guy, tic, "end", lot)
    VALUES (NEW.bid_id, NEW.address_id,
            (SELECT block_number FROM block_info),
            (SELECT hash FROM block_info),
            (SELECT datetime FROM block_info),
            (SELECT datetime FROM block_info),
            (SELECT get_latest_flop_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flop_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flop_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flop_bid_end(NEW.bid_id)),
            (SELECT get_latest_flop_bid_lot(NEW.bid_id)))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET created = (SELECT datetime FROM block_info),
                                                     updated = (SELECT datetime FROM block_info);
    return NEW;
END
$$;


--
-- Name: insert_cdp_created(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_cdp_created() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH block_info AS (
        SELECT api.epoch_to_datetime(headers.block_timestamp) AS datetime
        FROM public.headers
        WHERE headers.block_number = NEW.block_number
        LIMIT 1)
    INSERT
    INTO api.managed_cdp (cdpi, created)
    VALUES (NEW.cdpi, (SELECT datetime FROM block_info))
    ON CONFLICT (cdpi)
        DO UPDATE SET created = (SELECT datetime FROM block_info)
    WHERE managed_cdp.created IS NULL;
    RETURN NEW;
END
$$;


--
-- Name: insert_cdp_ilk_identifier(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_cdp_ilk_identifier() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (
        SELECT ilks.identifier
        FROM maker.cdp_manager_ilks
                 LEFT JOIN maker.ilks ON ilks.id = cdp_manager_ilks.ilk_id
        WHERE cdp_manager_ilks.cdpi = NEW.cdpi
    )
    INSERT
    INTO api.managed_cdp (cdpi, ilk_identifier)
    VALUES (NEW.cdpi, (SELECT identifier FROM ilk))
    ON CONFLICT (cdpi) DO UPDATE SET ilk_identifier = (SELECT identifier FROM ilk);
    RETURN NEW;
END
$$;


--
-- Name: insert_cdp_urn_identifier(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_cdp_urn_identifier() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    INSERT
    INTO api.managed_cdp (cdpi, urn_identifier)
    VALUES (NEW.cdpi, NEW.urn)
    ON CONFLICT (cdpi) DO UPDATE SET urn_identifier = NEW.urn;
    RETURN NEW;
END
$$;


--
-- Name: insert_cdp_usr(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_cdp_usr() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    INSERT
    INTO api.managed_cdp (cdpi, usr)
    VALUES (NEW.cdpi, NEW.owner)
           -- only update usr if the new owner is coming from the latest owns block we know about for the given cdpi
    ON CONFLICT (cdpi)
        DO UPDATE SET usr = NEW.owner
    WHERE NEW.block_number >= (
        SELECT MAX(block_number)
        FROM maker.cdp_manager_owns
        WHERE cdp_manager_owns.cdpi = NEW.cdpi);
    RETURN NEW;
END
$$;


--
-- Name: insert_ilk_art(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_ilk_art() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, art, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.art,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET art     = (
            CASE
                WHEN current_ilk_state.art IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.art
                ELSE current_ilk_state.art END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$;


--
-- Name: insert_ilk_chop(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_ilk_chop() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, chop, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.chop,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET chop    = (
            CASE
                WHEN current_ilk_state.chop IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.chop
                ELSE current_ilk_state.chop END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$;


--
-- Name: insert_ilk_dust(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_ilk_dust() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, dust, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.dust,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET dust    = (
            CASE
                WHEN current_ilk_state.dust IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.dust
                ELSE current_ilk_state.dust END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$;


--
-- Name: insert_ilk_duty(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_ilk_duty() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, duty, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.duty,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET duty    = (
            CASE
                WHEN current_ilk_state.duty IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.duty
                ELSE current_ilk_state.duty END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$;


--
-- Name: insert_ilk_flip(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_ilk_flip() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, flip, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.flip,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET flip    = (
            CASE
                WHEN current_ilk_state.flip IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.flip
                ELSE current_ilk_state.flip END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$;


--
-- Name: insert_ilk_line(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_ilk_line() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, line, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.line,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET line    = (
            CASE
                WHEN current_ilk_state.line IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.line
                ELSE current_ilk_state.line END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$;


--
-- Name: insert_ilk_lump(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_ilk_lump() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, lump, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.lump,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET lump    = (
            CASE
                WHEN current_ilk_state.lump IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.lump
                ELSE current_ilk_state.lump END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$;


--
-- Name: insert_ilk_mat(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_ilk_mat() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, mat, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.mat,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET mat     = (
            CASE
                WHEN current_ilk_state.mat IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.mat
                ELSE current_ilk_state.mat END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$;


--
-- Name: insert_ilk_pip(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_ilk_pip() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, pip, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.pip,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET pip     = (
            CASE
                WHEN current_ilk_state.pip IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.pip
                ELSE current_ilk_state.pip END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$;


--
-- Name: insert_ilk_rate(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_ilk_rate() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, rate, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.rate,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET rate    = (
            CASE
                WHEN current_ilk_state.rate IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.rate
                ELSE current_ilk_state.rate END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$;


--
-- Name: insert_ilk_rho(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_ilk_rho() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, rho, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.rho,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET rho     = (
            CASE
                WHEN current_ilk_state.rho IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.rho
                ELSE current_ilk_state.rho END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$;


--
-- Name: insert_ilk_spot(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_ilk_spot() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH ilk AS (SELECT identifier FROM maker.ilks WHERE id = NEW.ilk_id),
         block_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp
             FROM public.headers
             WHERE block_number = NEW.block_number)
    INSERT
    INTO api.current_ilk_state (ilk_identifier, spot, created, updated)
    VALUES ((SELECT identifier FROM ilk),
            NEW.spot,
            (SELECT block_timestamp FROM block_header),
            (SELECT block_timestamp FROM block_header))
    ON CONFLICT (ilk_identifier)
        DO UPDATE
        SET spot    = (
            CASE
                WHEN current_ilk_state.spot IS NULL OR
                     current_ilk_state.updated < (SELECT block_timestamp FROM block_header)
                    THEN NEW.spot
                ELSE current_ilk_state.spot END),
            created = LEAST((SELECT block_timestamp FROM block_header), current_ilk_state.created),
            updated = GREATEST((SELECT block_timestamp FROM block_header), current_ilk_state.updated);
    RETURN NEW;
END
$$;


--
-- Name: insert_updated_flap_bid(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flap_bid() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flap
        WHERE flap.bid_id = NEW.bid_id
        ORDER BY flap.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flap(bid_id, address_id, block_number, block_hash, bid, guy, tic, "end", lot, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW.bid,
            (SELECT get_latest_flap_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flap_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flap_bid_end(NEW.bid_id)),
            (SELECT get_latest_flap_bid_lot(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET bid = NEW.bid;
    return NEW;
END
$$;


--
-- Name: insert_updated_flap_end(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flap_end() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flap
        WHERE flap.bid_id = NEW.bid_id
        ORDER BY flap.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flap(bid_id, address_id, block_number, block_hash, "end", bid, guy, tic, lot, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW."end",
            (SELECT get_latest_flap_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flap_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flap_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flap_bid_lot(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET "end" = NEW."end";
    return NEW;
END
$$;


--
-- Name: insert_updated_flap_guy(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flap_guy() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flap
        WHERE flap.bid_id = NEW.bid_id
        ORDER BY flap.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flap(bid_id, address_id, block_number, block_hash, guy, bid, tic, "end", lot, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW.guy,
            (SELECT get_latest_flap_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flap_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flap_bid_end(NEW.bid_id)),
            (SELECT get_latest_flap_bid_lot(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET guy = NEW.guy;
    return NEW;
END
$$;


--
-- Name: insert_updated_flap_lot(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flap_lot() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flap
        WHERE flap.bid_id = NEW.bid_id
        ORDER BY flap.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flap(bid_id, address_id, block_number, block_hash, lot, bid, guy, tic, "end", updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW.lot,
            (SELECT get_latest_flap_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flap_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flap_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flap_bid_end(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET lot = NEW.lot;
    return NEW;
END
$$;


--
-- Name: insert_updated_flap_tic(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flap_tic() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flap
        WHERE flap.bid_id = NEW.bid_id
        ORDER BY flap.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flap(bid_id, address_id, block_number, block_hash, tic, bid, guy, "end", lot, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW.tic,
            (SELECT get_latest_flap_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flap_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flap_bid_end(NEW.bid_id)),
            (SELECT get_latest_flap_bid_lot(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET tic = NEW.tic;
    return NEW;
END
$$;


--
-- Name: insert_updated_flip_bid(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flip_bid() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, bid, guy, tic, "end", lot, gal, tab, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW.bid,
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET bid = NEW.bid;
    return NEW;
END
$$;


--
-- Name: insert_updated_flip_end(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flip_end() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, "end", guy, tic, lot, bid, gal, tab, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW."end",
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET "end" = NEW."end";
    return NEW;
END
$$;


--
-- Name: insert_updated_flip_gal(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flip_gal() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, gal, guy, tic, "end", lot, bid, tab, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW.gal,
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET gal = NEW.gal;
    return NEW;
END
$$;


--
-- Name: insert_updated_flip_guy(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flip_guy() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, guy, tic, "end", lot, bid, gal, tab, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW.guy,
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET guy = NEW.guy;
    return NEW;
END
$$;


--
-- Name: insert_updated_flip_lot(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flip_lot() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, lot, guy, tic, "end", bid, gal, tab, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW.lot,
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET lot = NEW.lot;
    return NEW;
END
$$;


--
-- Name: insert_updated_flip_tab(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flip_tab() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, tab, guy, tic, "end", lot, bid, gal, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW.tab,
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET tab = NEW.tab;
    return NEW;
END
$$;


--
-- Name: insert_updated_flip_tic(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flip_tic() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flip
        WHERE flip.bid_id = NEW.bid_id
        ORDER BY flip.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flip(bid_id, address_id, block_number, block_hash, tic, guy, "end", lot, bid, gal, tab, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, new.block_hash, NEW.tic,
            (SELECT get_latest_flip_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flip_bid_end(NEW.bid_id)),
            (SELECT get_latest_flip_bid_lot(NEW.bid_id)),
            (SELECT get_latest_flip_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flip_bid_gal(NEW.bid_id)),
            (SELECT get_latest_flip_bid_tab(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET tic = NEW.tic;
    return NEW;
END
$$;


--
-- Name: insert_updated_flop_bid(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flop_bid() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flop
        WHERE flop.bid_id = NEW.bid_id
        ORDER BY flop.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flop(bid_id, address_id, block_number, block_hash, bid, guy, tic, "end", lot, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW.bid,
            (SELECT get_latest_flop_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flop_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flop_bid_end(NEW.bid_id)),
            (SELECT get_latest_flop_bid_lot(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET bid = NEW.bid;
    return NEW;
END
$$;


--
-- Name: insert_updated_flop_end(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flop_end() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flop
        WHERE flop.bid_id = NEW.bid_id
        ORDER BY flop.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flop(bid_id, address_id, block_number, block_hash, "end", bid, guy, tic, lot, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW."end",
            (SELECT get_latest_flop_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flop_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flop_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flop_bid_lot(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET "end" = NEW."end";
    return NEW;
END
$$;


--
-- Name: insert_updated_flop_guy(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flop_guy() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flop
        WHERE flop.bid_id = NEW.bid_id
        ORDER BY flop.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flop(bid_id, address_id, block_number, block_hash, guy, bid, tic, "end", lot, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW.guy,
            (SELECT get_latest_flop_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flop_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flop_bid_end(NEW.bid_id)),
            (SELECT get_latest_flop_bid_lot(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET guy = NEW.guy;
    return NEW;
END
$$;


--
-- Name: insert_updated_flop_lot(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flop_lot() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flop
        WHERE flop.bid_id = NEW.bid_id
        ORDER BY flop.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flop(bid_id, address_id, block_number, block_hash, lot, bid, guy, tic, "end", updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW.lot,
            (SELECT get_latest_flop_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flop_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flop_bid_tic(NEW.bid_id)),
            (SELECT get_latest_flop_bid_end(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET lot = NEW.lot;
    return NEW;
END
$$;


--
-- Name: insert_updated_flop_tic(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_updated_flop_tic() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH created AS (
        SELECT created
        FROM maker.flop
        WHERE flop.bid_id = NEW.bid_id
        ORDER BY flop.block_number
        LIMIT 1
    )
    INSERT
    INTO maker.flop(bid_id, address_id, block_number, block_hash, tic, bid, guy, "end", lot, updated,
                    created)
    VALUES (NEW.bid_id, NEW.address_id, NEW.block_number, NEW.block_hash, NEW.tic,
            (SELECT get_latest_flop_bid_bid(NEW.bid_id)),
            (SELECT get_latest_flop_bid_guy(NEW.bid_id)),
            (SELECT get_latest_flop_bid_end(NEW.bid_id)),
            (SELECT get_latest_flop_bid_lot(NEW.bid_id)),
            (SELECT get_block_timestamp(NEW.block_hash)),
            (SELECT created FROM created))
    ON CONFLICT (bid_id, block_number) DO UPDATE SET tic = NEW.tic;
    return NEW;
END
$$;


--
-- Name: get_block_timestamp(character varying); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_block_timestamp(block_hash character varying) RETURNS timestamp without time zone
    LANGUAGE sql
    AS $$
SELECT api.epoch_to_datetime(headers.block_timestamp) AS datetime
FROM public.headers
WHERE headers.hash = block_hash
ORDER BY headers.block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flap_bid_bid(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flap_bid_bid(bid_id numeric) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT bid
FROM maker.flap
WHERE bid IS NOT NULL
  AND flap.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flap_bid_end(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flap_bid_end(bid_id numeric) RETURNS bigint
    LANGUAGE sql
    AS $$
SELECT "end"
FROM maker.flap
WHERE "end" IS NOT NULL
  AND flap.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flap_bid_guy(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flap_bid_guy(bid_id numeric) RETURNS text
    LANGUAGE sql
    AS $$
SELECT guy
FROM maker.flap
WHERE guy IS NOT NULL
  AND flap.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flap_bid_lot(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flap_bid_lot(bid_id numeric) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT lot
FROM maker.flap
WHERE lot IS NOT NULL
  AND flap.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flap_bid_tic(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flap_bid_tic(bid_id numeric) RETURNS bigint
    LANGUAGE sql
    AS $$
SELECT tic
FROM maker.flap
WHERE tic IS NOT NULL
  AND flap.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flip_bid_bid(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flip_bid_bid(bid_id numeric) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT bid
FROM maker.flip
WHERE bid IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flip_bid_end(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flip_bid_end(bid_id numeric) RETURNS bigint
    LANGUAGE sql
    AS $$
SELECT "end"
FROM maker.flip
WHERE "end" IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flip_bid_gal(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flip_bid_gal(bid_id numeric) RETURNS text
    LANGUAGE sql
    AS $$
SELECT gal
FROM maker.flip
WHERE gal IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flip_bid_guy(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flip_bid_guy(bid_id numeric) RETURNS text
    LANGUAGE sql
    AS $$
SELECT guy
FROM maker.flip
WHERE guy IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flip_bid_lot(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flip_bid_lot(bid_id numeric) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT lot
FROM maker.flip
WHERE lot IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flip_bid_tab(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flip_bid_tab(bid_id numeric) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT tab
FROM maker.flip
WHERE tab IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flip_bid_tic(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flip_bid_tic(bid_id numeric) RETURNS bigint
    LANGUAGE sql
    AS $$
SELECT tic
FROM maker.flip
WHERE tic IS NOT NULL
  AND flip.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flop_bid_bid(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flop_bid_bid(bid_id numeric) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT bid
FROM maker.flop
WHERE bid IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flop_bid_end(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flop_bid_end(bid_id numeric) RETURNS bigint
    LANGUAGE sql
    AS $$
SELECT "end"
FROM maker.flop
WHERE "end" IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flop_bid_guy(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flop_bid_guy(bid_id numeric) RETURNS text
    LANGUAGE sql
    AS $$
SELECT guy
FROM maker.flop
WHERE guy IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flop_bid_lot(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flop_bid_lot(bid_id numeric) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT lot
FROM maker.flop
WHERE lot IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_latest_flop_bid_tic(numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_latest_flop_bid_tic(bid_id numeric) RETURNS bigint
    LANGUAGE sql
    AS $$
SELECT tic
FROM maker.flop
WHERE tic IS NOT NULL
  AND flop.bid_id = bid_id
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: get_tx_data(bigint, bigint); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_tx_data(block_height bigint, log_id bigint) RETURNS SETOF api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT txs.hash, txs.tx_index, headers.block_number, headers.hash, tx_from, tx_to
FROM public.header_sync_transactions txs
         LEFT JOIN headers ON txs.header_id = headers.id
         LEFT JOIN header_sync_logs ON txs.tx_index = header_sync_logs.tx_index
WHERE headers.block_number <= block_height
  AND header_sync_logs.id = log_id
ORDER BY block_number DESC

$$;


--
-- Name: managed_cdp_id_seq; Type: SEQUENCE; Schema: api; Owner: -
--

CREATE SEQUENCE api.managed_cdp_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: managed_cdp_id_seq; Type: SEQUENCE OWNED BY; Schema: api; Owner: -
--

ALTER SEQUENCE api.managed_cdp_id_seq OWNED BY api.managed_cdp.id;


--
-- Name: bite; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.bite (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    urn_id integer NOT NULL,
    ink numeric,
    art numeric,
    tab numeric,
    flip text,
    bid_id numeric
);


--
-- Name: TABLE bite; Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON TABLE maker.bite IS '@name raw_bites';


--
-- Name: COLUMN bite.id; Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON COLUMN maker.bite.id IS '@omit';


--
-- Name: bite_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.bite_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bite_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.bite_id_seq OWNED BY maker.bite.id;


--
-- Name: cat_file_chop_lump; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_file_chop_lump (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL,
    what text,
    data numeric
);


--
-- Name: cat_file_chop_lump_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_file_chop_lump_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_file_chop_lump_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_file_chop_lump_id_seq OWNED BY maker.cat_file_chop_lump.id;


--
-- Name: cat_file_flip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_file_flip (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL,
    what text,
    flip text
);


--
-- Name: cat_file_flip_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_file_flip_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_file_flip_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_file_flip_id_seq OWNED BY maker.cat_file_flip.id;


--
-- Name: cat_file_vow; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_file_vow (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    what text,
    data text
);


--
-- Name: cat_file_vow_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_file_vow_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_file_vow_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_file_vow_id_seq OWNED BY maker.cat_file_vow.id;


--
-- Name: cat_ilk_chop; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_ilk_chop (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    chop numeric NOT NULL
);


--
-- Name: cat_ilk_chop_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_ilk_chop_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_ilk_chop_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_ilk_chop_id_seq OWNED BY maker.cat_ilk_chop.id;


--
-- Name: cat_ilk_flip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_ilk_flip (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    flip text
);


--
-- Name: cat_ilk_flip_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_ilk_flip_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_ilk_flip_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_ilk_flip_id_seq OWNED BY maker.cat_ilk_flip.id;


--
-- Name: cat_ilk_lump; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_ilk_lump (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    lump numeric NOT NULL
);


--
-- Name: cat_ilk_lump_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_ilk_lump_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_ilk_lump_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_ilk_lump_id_seq OWNED BY maker.cat_ilk_lump.id;


--
-- Name: cat_live; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_live (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    live numeric NOT NULL
);


--
-- Name: cat_live_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_live_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_live_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_live_id_seq OWNED BY maker.cat_live.id;


--
-- Name: cat_vat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_vat (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    vat text
);


--
-- Name: cat_vat_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_vat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_vat_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_vat_id_seq OWNED BY maker.cat_vat.id;


--
-- Name: cat_vow; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_vow (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    vow text
);


--
-- Name: cat_vow_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_vow_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_vow_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_vow_id_seq OWNED BY maker.cat_vow.id;


--
-- Name: cdp_manager_cdpi; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cdp_manager_cdpi (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    cdpi numeric NOT NULL
);


--
-- Name: cdp_manager_cdpi_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cdp_manager_cdpi_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cdp_manager_cdpi_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cdp_manager_cdpi_id_seq OWNED BY maker.cdp_manager_cdpi.id;


--
-- Name: cdp_manager_count; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cdp_manager_count (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    owner text,
    count numeric NOT NULL
);


--
-- Name: cdp_manager_count_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cdp_manager_count_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cdp_manager_count_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cdp_manager_count_id_seq OWNED BY maker.cdp_manager_count.id;


--
-- Name: cdp_manager_first; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cdp_manager_first (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    owner text,
    first numeric NOT NULL
);


--
-- Name: cdp_manager_first_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cdp_manager_first_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cdp_manager_first_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cdp_manager_first_id_seq OWNED BY maker.cdp_manager_first.id;


--
-- Name: cdp_manager_ilks; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cdp_manager_ilks (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    cdpi numeric NOT NULL,
    ilk_id integer NOT NULL
);


--
-- Name: cdp_manager_ilks_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cdp_manager_ilks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cdp_manager_ilks_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cdp_manager_ilks_id_seq OWNED BY maker.cdp_manager_ilks.id;


--
-- Name: cdp_manager_last; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cdp_manager_last (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    owner text,
    last numeric NOT NULL
);


--
-- Name: cdp_manager_last_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cdp_manager_last_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cdp_manager_last_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cdp_manager_last_id_seq OWNED BY maker.cdp_manager_last.id;


--
-- Name: cdp_manager_list_next; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cdp_manager_list_next (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    cdpi numeric NOT NULL,
    next numeric NOT NULL
);


--
-- Name: cdp_manager_list_next_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cdp_manager_list_next_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cdp_manager_list_next_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cdp_manager_list_next_id_seq OWNED BY maker.cdp_manager_list_next.id;


--
-- Name: cdp_manager_list_prev; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cdp_manager_list_prev (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    cdpi numeric NOT NULL,
    prev numeric NOT NULL
);


--
-- Name: cdp_manager_list_prev_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cdp_manager_list_prev_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cdp_manager_list_prev_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cdp_manager_list_prev_id_seq OWNED BY maker.cdp_manager_list_prev.id;


--
-- Name: cdp_manager_owns; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cdp_manager_owns (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    cdpi numeric NOT NULL,
    owner text
);


--
-- Name: cdp_manager_owns_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cdp_manager_owns_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cdp_manager_owns_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cdp_manager_owns_id_seq OWNED BY maker.cdp_manager_owns.id;


--
-- Name: cdp_manager_urns; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cdp_manager_urns (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    cdpi numeric NOT NULL,
    urn text
);


--
-- Name: cdp_manager_urns_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cdp_manager_urns_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cdp_manager_urns_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cdp_manager_urns_id_seq OWNED BY maker.cdp_manager_urns.id;


--
-- Name: cdp_manager_vat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cdp_manager_vat (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    vat text
);


--
-- Name: cdp_manager_vat_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cdp_manager_vat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cdp_manager_vat_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cdp_manager_vat_id_seq OWNED BY maker.cdp_manager_vat.id;


--
-- Name: deal; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.deal (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    bid_id numeric NOT NULL,
    address_id integer NOT NULL
);


--
-- Name: deal_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.deal_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: deal_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.deal_id_seq OWNED BY maker.deal.id;


--
-- Name: dent; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.dent (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric,
    bid numeric,
    address_id integer NOT NULL
);


--
-- Name: dent_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.dent_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: dent_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.dent_id_seq OWNED BY maker.dent.id;


--
-- Name: flap; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric,
    guy text,
    tic bigint,
    "end" bigint,
    lot numeric,
    bid numeric,
    created timestamp without time zone,
    updated timestamp without time zone
);


--
-- Name: flap_beg; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_beg (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    beg numeric NOT NULL
);


--
-- Name: flap_beg_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_beg_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_beg_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_beg_id_seq OWNED BY maker.flap_beg.id;


--
-- Name: flap_bid_bid; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_bid_bid (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    bid numeric NOT NULL
);


--
-- Name: flap_bid_bid_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_bid_bid_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_bid_bid_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_bid_bid_id_seq OWNED BY maker.flap_bid_bid.id;


--
-- Name: flap_bid_end; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_bid_end (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    "end" bigint NOT NULL
);


--
-- Name: flap_bid_end_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_bid_end_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_bid_end_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_bid_end_id_seq OWNED BY maker.flap_bid_end.id;


--
-- Name: flap_bid_guy; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_bid_guy (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    guy text NOT NULL
);


--
-- Name: flap_bid_guy_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_bid_guy_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_bid_guy_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_bid_guy_id_seq OWNED BY maker.flap_bid_guy.id;


--
-- Name: flap_bid_lot; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_bid_lot (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric NOT NULL
);


--
-- Name: flap_bid_lot_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_bid_lot_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_bid_lot_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_bid_lot_id_seq OWNED BY maker.flap_bid_lot.id;


--
-- Name: flap_bid_tic; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_bid_tic (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    tic bigint NOT NULL
);


--
-- Name: flap_bid_tic_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_bid_tic_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_bid_tic_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_bid_tic_id_seq OWNED BY maker.flap_bid_tic.id;


--
-- Name: flap_gem; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_gem (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    gem text NOT NULL
);


--
-- Name: flap_gem_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_gem_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_gem_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_gem_id_seq OWNED BY maker.flap_gem.id;


--
-- Name: flap_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_id_seq OWNED BY maker.flap.id;


--
-- Name: flap_kick; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_kick (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric NOT NULL,
    bid numeric NOT NULL,
    address_id integer NOT NULL
);


--
-- Name: TABLE flap_kick; Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON TABLE maker.flap_kick IS '@name flapKickEvent';


--
-- Name: flap_kick_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_kick_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_kick_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_kick_id_seq OWNED BY maker.flap_kick.id;


--
-- Name: flap_kicks; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_kicks (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    kicks numeric NOT NULL
);


--
-- Name: TABLE flap_kicks; Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON TABLE maker.flap_kicks IS '@name flapKicksStorage';


--
-- Name: flap_kicks_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_kicks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_kicks_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_kicks_id_seq OWNED BY maker.flap_kicks.id;


--
-- Name: flap_live; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_live (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    live numeric NOT NULL
);


--
-- Name: flap_live_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_live_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_live_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_live_id_seq OWNED BY maker.flap_live.id;


--
-- Name: flap_tau; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_tau (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    tau integer NOT NULL
);


--
-- Name: flap_tau_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_tau_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_tau_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_tau_id_seq OWNED BY maker.flap_tau.id;


--
-- Name: flap_ttl; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_ttl (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    ttl integer NOT NULL
);


--
-- Name: flap_ttl_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_ttl_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_ttl_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_ttl_id_seq OWNED BY maker.flap_ttl.id;


--
-- Name: flap_vat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_vat (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    vat text NOT NULL
);


--
-- Name: flap_vat_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flap_vat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flap_vat_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flap_vat_id_seq OWNED BY maker.flap_vat.id;


--
-- Name: flip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric,
    guy text,
    tic bigint,
    "end" bigint,
    lot numeric,
    bid numeric,
    gal text,
    tab numeric,
    created timestamp without time zone,
    updated timestamp without time zone
);


--
-- Name: flip_beg; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_beg (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    beg numeric NOT NULL
);


--
-- Name: flip_beg_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_beg_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_beg_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_beg_id_seq OWNED BY maker.flip_beg.id;


--
-- Name: flip_bid_bid; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_bid (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    bid numeric NOT NULL
);


--
-- Name: flip_bid_bid_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_bid_bid_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_bid_bid_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_bid_bid_id_seq OWNED BY maker.flip_bid_bid.id;


--
-- Name: flip_bid_end; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_end (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    "end" bigint NOT NULL
);


--
-- Name: flip_bid_end_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_bid_end_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_bid_end_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_bid_end_id_seq OWNED BY maker.flip_bid_end.id;


--
-- Name: flip_bid_gal; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_gal (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    gal text
);


--
-- Name: flip_bid_gal_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_bid_gal_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_bid_gal_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_bid_gal_id_seq OWNED BY maker.flip_bid_gal.id;


--
-- Name: flip_bid_guy; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_guy (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    guy text
);


--
-- Name: flip_bid_guy_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_bid_guy_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_bid_guy_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_bid_guy_id_seq OWNED BY maker.flip_bid_guy.id;


--
-- Name: flip_bid_lot; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_lot (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric NOT NULL
);


--
-- Name: flip_bid_lot_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_bid_lot_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_bid_lot_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_bid_lot_id_seq OWNED BY maker.flip_bid_lot.id;


--
-- Name: flip_bid_tab; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_tab (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    tab numeric NOT NULL
);


--
-- Name: flip_bid_tab_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_bid_tab_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_bid_tab_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_bid_tab_id_seq OWNED BY maker.flip_bid_tab.id;


--
-- Name: flip_bid_tic; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_tic (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    tic bigint NOT NULL
);


--
-- Name: flip_bid_tic_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_bid_tic_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_bid_tic_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_bid_tic_id_seq OWNED BY maker.flip_bid_tic.id;


--
-- Name: flip_bid_usr; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_usr (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    usr text
);


--
-- Name: flip_bid_usr_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_bid_usr_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_bid_usr_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_bid_usr_id_seq OWNED BY maker.flip_bid_usr.id;


--
-- Name: flip_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_id_seq OWNED BY maker.flip.id;


--
-- Name: flip_ilk; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_ilk (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    ilk_id integer NOT NULL
);


--
-- Name: flip_ilk_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_ilk_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_ilk_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_ilk_id_seq OWNED BY maker.flip_ilk.id;


--
-- Name: flip_kick; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_kick (
    id integer NOT NULL,
    header_id integer NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric,
    bid numeric,
    tab numeric,
    usr text,
    gal text,
    address_id integer NOT NULL,
    log_id bigint NOT NULL
);


--
-- Name: TABLE flip_kick; Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON TABLE maker.flip_kick IS '@name flipKickEvent';


--
-- Name: flip_kick_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_kick_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_kick_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_kick_id_seq OWNED BY maker.flip_kick.id;


--
-- Name: flip_kicks; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_kicks (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    kicks numeric NOT NULL
);


--
-- Name: TABLE flip_kicks; Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON TABLE maker.flip_kicks IS '@name flipKicksStorage';


--
-- Name: flip_kicks_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_kicks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_kicks_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_kicks_id_seq OWNED BY maker.flip_kicks.id;


--
-- Name: flip_tau; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_tau (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    tau numeric NOT NULL
);


--
-- Name: flip_tau_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_tau_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_tau_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_tau_id_seq OWNED BY maker.flip_tau.id;


--
-- Name: flip_ttl; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_ttl (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    ttl numeric NOT NULL
);


--
-- Name: flip_ttl_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_ttl_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_ttl_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_ttl_id_seq OWNED BY maker.flip_ttl.id;


--
-- Name: flip_vat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_vat (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    vat text
);


--
-- Name: flip_vat_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flip_vat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flip_vat_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flip_vat_id_seq OWNED BY maker.flip_vat.id;


--
-- Name: flop; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric,
    guy text,
    tic bigint,
    "end" bigint,
    lot numeric,
    bid numeric,
    created timestamp without time zone,
    updated timestamp without time zone
);


--
-- Name: flop_beg; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_beg (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    beg numeric NOT NULL
);


--
-- Name: flop_beg_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_beg_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_beg_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_beg_id_seq OWNED BY maker.flop_beg.id;


--
-- Name: flop_bid_bid; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_bid_bid (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    bid numeric NOT NULL
);


--
-- Name: flop_bid_bid_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_bid_bid_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_bid_bid_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_bid_bid_id_seq OWNED BY maker.flop_bid_bid.id;


--
-- Name: flop_bid_end; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_bid_end (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    "end" bigint NOT NULL
);


--
-- Name: flop_bid_end_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_bid_end_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_bid_end_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_bid_end_id_seq OWNED BY maker.flop_bid_end.id;


--
-- Name: flop_bid_guy; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_bid_guy (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    guy text
);


--
-- Name: flop_bid_guy_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_bid_guy_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_bid_guy_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_bid_guy_id_seq OWNED BY maker.flop_bid_guy.id;


--
-- Name: flop_bid_lot; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_bid_lot (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric NOT NULL
);


--
-- Name: flop_bid_lot_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_bid_lot_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_bid_lot_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_bid_lot_id_seq OWNED BY maker.flop_bid_lot.id;


--
-- Name: flop_bid_tic; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_bid_tic (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    tic bigint NOT NULL
);


--
-- Name: flop_bid_tic_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_bid_tic_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_bid_tic_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_bid_tic_id_seq OWNED BY maker.flop_bid_tic.id;


--
-- Name: flop_gem; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_gem (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    gem text
);


--
-- Name: flop_gem_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_gem_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_gem_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_gem_id_seq OWNED BY maker.flop_gem.id;


--
-- Name: flop_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_id_seq OWNED BY maker.flop.id;


--
-- Name: flop_kick; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_kick (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric NOT NULL,
    bid numeric NOT NULL,
    gal text,
    address_id integer NOT NULL
);


--
-- Name: TABLE flop_kick; Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON TABLE maker.flop_kick IS '@name flopKickEvent';


--
-- Name: flop_kick_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_kick_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_kick_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_kick_id_seq OWNED BY maker.flop_kick.id;


--
-- Name: flop_kicks; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_kicks (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    kicks numeric NOT NULL
);


--
-- Name: TABLE flop_kicks; Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON TABLE maker.flop_kicks IS '@name flopKicksStorage';


--
-- Name: flop_kicks_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_kicks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_kicks_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_kicks_id_seq OWNED BY maker.flop_kicks.id;


--
-- Name: flop_live; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_live (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    live numeric NOT NULL
);


--
-- Name: flop_live_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_live_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_live_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_live_id_seq OWNED BY maker.flop_live.id;


--
-- Name: flop_pad; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_pad (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    pad numeric NOT NULL
);


--
-- Name: flop_pad_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_pad_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_pad_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_pad_id_seq OWNED BY maker.flop_pad.id;


--
-- Name: flop_tau; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_tau (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    tau numeric NOT NULL
);


--
-- Name: flop_tau_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_tau_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_tau_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_tau_id_seq OWNED BY maker.flop_tau.id;


--
-- Name: flop_ttl; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_ttl (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    ttl numeric NOT NULL
);


--
-- Name: flop_ttl_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_ttl_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_ttl_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_ttl_id_seq OWNED BY maker.flop_ttl.id;


--
-- Name: flop_vat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_vat (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    address_id integer NOT NULL,
    vat text
);


--
-- Name: flop_vat_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_vat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_vat_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_vat_id_seq OWNED BY maker.flop_vat.id;


--
-- Name: ilks; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.ilks (
    id integer NOT NULL,
    ilk text NOT NULL,
    identifier text NOT NULL
);


--
-- Name: TABLE ilks; Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON TABLE maker.ilks IS '@name raw_ilks';


--
-- Name: ilks_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.ilks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ilks_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.ilks_id_seq OWNED BY maker.ilks.id;


--
-- Name: jug_base; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_base (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    base text
);


--
-- Name: jug_base_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_base_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_base_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_base_id_seq OWNED BY maker.jug_base.id;


--
-- Name: jug_drip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_drip (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL
);


--
-- Name: jug_drip_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_drip_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_drip_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_drip_id_seq OWNED BY maker.jug_drip.id;


--
-- Name: jug_file_base; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_file_base (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    what text,
    data numeric
);


--
-- Name: jug_file_base_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_file_base_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_file_base_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_file_base_id_seq OWNED BY maker.jug_file_base.id;


--
-- Name: jug_file_ilk; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_file_ilk (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL,
    what text,
    data numeric
);


--
-- Name: jug_file_ilk_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_file_ilk_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_file_ilk_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_file_ilk_id_seq OWNED BY maker.jug_file_ilk.id;


--
-- Name: jug_file_vow; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_file_vow (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    what text,
    data text
);


--
-- Name: jug_file_vow_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_file_vow_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_file_vow_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_file_vow_id_seq OWNED BY maker.jug_file_vow.id;


--
-- Name: jug_ilk_duty; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_ilk_duty (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    duty numeric NOT NULL
);


--
-- Name: jug_ilk_duty_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_ilk_duty_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_ilk_duty_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_ilk_duty_id_seq OWNED BY maker.jug_ilk_duty.id;


--
-- Name: jug_ilk_rho; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_ilk_rho (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    rho numeric NOT NULL
);


--
-- Name: jug_ilk_rho_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_ilk_rho_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_ilk_rho_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_ilk_rho_id_seq OWNED BY maker.jug_ilk_rho.id;


--
-- Name: jug_init; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_init (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL
);


--
-- Name: jug_init_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_init_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_init_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_init_id_seq OWNED BY maker.jug_init.id;


--
-- Name: jug_vat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_vat (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    vat text
);


--
-- Name: jug_vat_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_vat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_vat_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_vat_id_seq OWNED BY maker.jug_vat.id;


--
-- Name: jug_vow; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_vow (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    vow text
);


--
-- Name: jug_vow_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_vow_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_vow_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_vow_id_seq OWNED BY maker.jug_vow.id;


--
-- Name: new_cdp; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.new_cdp (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    usr text,
    own text,
    cdp numeric
);


--
-- Name: COLUMN new_cdp.id; Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON COLUMN maker.new_cdp.id IS '@omit';


--
-- Name: new_cdp_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.new_cdp_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: new_cdp_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.new_cdp_id_seq OWNED BY maker.new_cdp.id;


--
-- Name: spot_file_mat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_file_mat (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL,
    what text,
    data numeric
);


--
-- Name: spot_file_mat_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.spot_file_mat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: spot_file_mat_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.spot_file_mat_id_seq OWNED BY maker.spot_file_mat.id;


--
-- Name: spot_file_pip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_file_pip (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL,
    what text,
    pip text
);


--
-- Name: spot_file_pip_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.spot_file_pip_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: spot_file_pip_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.spot_file_pip_id_seq OWNED BY maker.spot_file_pip.id;


--
-- Name: spot_ilk_mat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_ilk_mat (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    mat numeric NOT NULL
);


--
-- Name: spot_ilk_mat_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.spot_ilk_mat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: spot_ilk_mat_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.spot_ilk_mat_id_seq OWNED BY maker.spot_ilk_mat.id;


--
-- Name: spot_ilk_pip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_ilk_pip (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    pip text
);


--
-- Name: spot_ilk_pip_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.spot_ilk_pip_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: spot_ilk_pip_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.spot_ilk_pip_id_seq OWNED BY maker.spot_ilk_pip.id;


--
-- Name: spot_par; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_par (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    par numeric NOT NULL
);


--
-- Name: spot_par_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.spot_par_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: spot_par_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.spot_par_id_seq OWNED BY maker.spot_par.id;


--
-- Name: spot_poke; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_poke (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL,
    value numeric,
    spot numeric
);


--
-- Name: spot_poke_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.spot_poke_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: spot_poke_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.spot_poke_id_seq OWNED BY maker.spot_poke.id;


--
-- Name: spot_vat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_vat (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    vat text
);


--
-- Name: spot_vat_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.spot_vat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: spot_vat_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.spot_vat_id_seq OWNED BY maker.spot_vat.id;


--
-- Name: tend; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.tend (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric,
    bid numeric,
    address_id integer NOT NULL
);


--
-- Name: tend_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.tend_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tend_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.tend_id_seq OWNED BY maker.tend.id;


--
-- Name: tick; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.tick (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    bid_id numeric NOT NULL,
    address_id integer NOT NULL
);


--
-- Name: tick_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.tick_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tick_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.tick_id_seq OWNED BY maker.tick.id;


--
-- Name: urns; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.urns (
    id integer NOT NULL,
    ilk_id integer NOT NULL,
    identifier public.citext
);


--
-- Name: TABLE urns; Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON TABLE maker.urns IS '@name raw_urns';


--
-- Name: urns_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.urns_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: urns_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.urns_id_seq OWNED BY maker.urns.id;


--
-- Name: vat_dai; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_dai (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    guy text,
    dai numeric NOT NULL
);


--
-- Name: vat_dai_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_dai_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_dai_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_dai_id_seq OWNED BY maker.vat_dai.id;


--
-- Name: vat_debt; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_debt (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    debt numeric NOT NULL
);


--
-- Name: vat_debt_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_debt_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_debt_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_debt_id_seq OWNED BY maker.vat_debt.id;


--
-- Name: vat_file_debt_ceiling; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_file_debt_ceiling (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    what text,
    data numeric
);


--
-- Name: vat_file_debt_ceiling_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_file_debt_ceiling_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_file_debt_ceiling_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_file_debt_ceiling_id_seq OWNED BY maker.vat_file_debt_ceiling.id;


--
-- Name: vat_file_ilk; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_file_ilk (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL,
    what text,
    data numeric
);


--
-- Name: vat_file_ilk_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_file_ilk_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_file_ilk_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_file_ilk_id_seq OWNED BY maker.vat_file_ilk.id;


--
-- Name: vat_flux; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_flux (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL,
    src text,
    dst text,
    wad numeric
);


--
-- Name: vat_flux_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_flux_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_flux_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_flux_id_seq OWNED BY maker.vat_flux.id;


--
-- Name: vat_fold; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_fold (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    urn_id integer NOT NULL,
    rate numeric
);


--
-- Name: vat_fold_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_fold_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_fold_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_fold_id_seq OWNED BY maker.vat_fold.id;


--
-- Name: vat_fork; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_fork (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL,
    src text,
    dst text,
    dink numeric,
    dart numeric
);


--
-- Name: vat_fork_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_fork_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_fork_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_fork_id_seq OWNED BY maker.vat_fork.id;


--
-- Name: vat_frob; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_frob (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    urn_id integer NOT NULL,
    v text,
    w text,
    dink numeric,
    dart numeric
);


--
-- Name: vat_frob_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_frob_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_frob_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_frob_id_seq OWNED BY maker.vat_frob.id;


--
-- Name: vat_gem; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_gem (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    guy text,
    gem numeric NOT NULL
);


--
-- Name: vat_gem_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_gem_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_gem_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_gem_id_seq OWNED BY maker.vat_gem.id;


--
-- Name: vat_grab; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_grab (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    urn_id integer NOT NULL,
    v text,
    w text,
    dink numeric,
    dart numeric
);


--
-- Name: vat_grab_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_grab_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_grab_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_grab_id_seq OWNED BY maker.vat_grab.id;


--
-- Name: vat_heal; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_heal (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    rad numeric
);


--
-- Name: vat_heal_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_heal_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_heal_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_heal_id_seq OWNED BY maker.vat_heal.id;


--
-- Name: vat_ilk_art; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_ilk_art (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    art numeric NOT NULL
);


--
-- Name: vat_ilk_art_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_ilk_art_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_ilk_art_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_ilk_art_id_seq OWNED BY maker.vat_ilk_art.id;


--
-- Name: vat_ilk_dust; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_ilk_dust (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    dust numeric NOT NULL
);


--
-- Name: vat_ilk_dust_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_ilk_dust_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_ilk_dust_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_ilk_dust_id_seq OWNED BY maker.vat_ilk_dust.id;


--
-- Name: vat_ilk_line; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_ilk_line (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    line numeric NOT NULL
);


--
-- Name: vat_ilk_line_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_ilk_line_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_ilk_line_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_ilk_line_id_seq OWNED BY maker.vat_ilk_line.id;


--
-- Name: vat_ilk_rate; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_ilk_rate (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    rate numeric NOT NULL
);


--
-- Name: vat_ilk_rate_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_ilk_rate_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_ilk_rate_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_ilk_rate_id_seq OWNED BY maker.vat_ilk_rate.id;


--
-- Name: vat_ilk_spot; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_ilk_spot (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    spot numeric NOT NULL
);


--
-- Name: vat_ilk_spot_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_ilk_spot_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_ilk_spot_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_ilk_spot_id_seq OWNED BY maker.vat_ilk_spot.id;


--
-- Name: vat_init; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_init (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL
);


--
-- Name: vat_init_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_init_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_init_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_init_id_seq OWNED BY maker.vat_init.id;


--
-- Name: vat_line; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_line (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    line numeric NOT NULL
);


--
-- Name: vat_line_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_line_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_line_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_line_id_seq OWNED BY maker.vat_line.id;


--
-- Name: vat_live; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_live (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    live numeric NOT NULL
);


--
-- Name: vat_live_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_live_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_live_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_live_id_seq OWNED BY maker.vat_live.id;


--
-- Name: vat_move; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_move (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    src text NOT NULL,
    dst text NOT NULL,
    rad numeric NOT NULL
);


--
-- Name: vat_move_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_move_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_move_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_move_id_seq OWNED BY maker.vat_move.id;


--
-- Name: vat_sin; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_sin (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    guy text,
    sin numeric NOT NULL
);


--
-- Name: vat_sin_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_sin_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_sin_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_sin_id_seq OWNED BY maker.vat_sin.id;


--
-- Name: vat_slip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_slip (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL,
    usr text,
    wad numeric
);


--
-- Name: vat_slip_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_slip_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_slip_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_slip_id_seq OWNED BY maker.vat_slip.id;


--
-- Name: vat_suck; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_suck (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    u text,
    v text,
    rad numeric
);


--
-- Name: vat_suck_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_suck_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_suck_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_suck_id_seq OWNED BY maker.vat_suck.id;


--
-- Name: vat_urn_art; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_urn_art (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    urn_id integer NOT NULL,
    art numeric NOT NULL
);


--
-- Name: vat_urn_art_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_urn_art_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_urn_art_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_urn_art_id_seq OWNED BY maker.vat_urn_art.id;


--
-- Name: vat_urn_ink; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_urn_ink (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    urn_id integer NOT NULL,
    ink numeric NOT NULL
);


--
-- Name: vat_urn_ink_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_urn_ink_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_urn_ink_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_urn_ink_id_seq OWNED BY maker.vat_urn_ink.id;


--
-- Name: vat_vice; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_vice (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    vice numeric NOT NULL
);


--
-- Name: vat_vice_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_vice_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_vice_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_vice_id_seq OWNED BY maker.vat_vice.id;


--
-- Name: vow_ash; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_ash (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ash numeric
);


--
-- Name: vow_ash_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_ash_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_ash_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_ash_id_seq OWNED BY maker.vow_ash.id;


--
-- Name: vow_bump; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_bump (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    bump numeric
);


--
-- Name: vow_bump_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_bump_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_bump_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_bump_id_seq OWNED BY maker.vow_bump.id;


--
-- Name: vow_dump; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_dump (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    dump numeric
);


--
-- Name: vow_dump_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_dump_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_dump_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_dump_id_seq OWNED BY maker.vow_dump.id;


--
-- Name: vow_fess; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_fess (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    tab numeric NOT NULL
);


--
-- Name: vow_fess_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_fess_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_fess_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_fess_id_seq OWNED BY maker.vow_fess.id;


--
-- Name: vow_file; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_file (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    what text,
    data numeric
);


--
-- Name: vow_file_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_file_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_file_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_file_id_seq OWNED BY maker.vow_file.id;


--
-- Name: vow_flapper; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_flapper (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    flapper text
);


--
-- Name: vow_flapper_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_flapper_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_flapper_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_flapper_id_seq OWNED BY maker.vow_flapper.id;


--
-- Name: vow_flog; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_flog (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    era integer NOT NULL
);


--
-- Name: vow_flog_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_flog_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_flog_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_flog_id_seq OWNED BY maker.vow_flog.id;


--
-- Name: vow_flopper; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_flopper (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    flopper text
);


--
-- Name: vow_flopper_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_flopper_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_flopper_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_flopper_id_seq OWNED BY maker.vow_flopper.id;


--
-- Name: vow_hump; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_hump (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    hump numeric
);


--
-- Name: vow_hump_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_hump_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_hump_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_hump_id_seq OWNED BY maker.vow_hump.id;


--
-- Name: vow_sin_integer; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_sin_integer (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    sin numeric
);


--
-- Name: vow_sin_integer_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_sin_integer_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_sin_integer_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_sin_integer_id_seq OWNED BY maker.vow_sin_integer.id;


--
-- Name: vow_sin_mapping; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_sin_mapping (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    era numeric,
    tab numeric
);


--
-- Name: vow_sin_mapping_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_sin_mapping_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_sin_mapping_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_sin_mapping_id_seq OWNED BY maker.vow_sin_mapping.id;


--
-- Name: vow_sump; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_sump (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    sump numeric
);


--
-- Name: vow_sump_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_sump_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_sump_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_sump_id_seq OWNED BY maker.vow_sump.id;


--
-- Name: vow_vat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_vat (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    vat text
);


--
-- Name: vow_vat_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_vat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_vat_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_vat_id_seq OWNED BY maker.vow_vat.id;


--
-- Name: vow_wait; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_wait (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    wait numeric
);


--
-- Name: vow_wait_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_wait_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_wait_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_wait_id_seq OWNED BY maker.vow_wait.id;


--
-- Name: yank; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.yank (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    bid_id numeric NOT NULL,
    address_id integer NOT NULL
);


--
-- Name: yank_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.yank_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: yank_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.yank_id_seq OWNED BY maker.yank.id;


--
-- Name: addresses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.addresses (
    id integer NOT NULL,
    address character varying(42)
);


--
-- Name: addresses_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.addresses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.addresses_id_seq OWNED BY public.addresses.id;


--
-- Name: full_sync_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.full_sync_logs (
    id integer NOT NULL,
    block_number bigint,
    address character varying(66),
    tx_hash character varying(66),
    index bigint,
    topic0 character varying(66),
    topic1 character varying(66),
    topic2 character varying(66),
    topic3 character varying(66),
    data text,
    receipt_id integer
);


--
-- Name: block_stats; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.block_stats AS
 SELECT max(full_sync_logs.block_number) AS max_block,
    min(full_sync_logs.block_number) AS min_block
   FROM public.full_sync_logs;


--
-- Name: blocks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.blocks (
    id integer NOT NULL,
    difficulty bigint,
    extra_data character varying,
    gas_limit bigint,
    gas_used bigint,
    hash character varying(66),
    miner character varying(42),
    nonce character varying(20),
    number bigint,
    parent_hash character varying(66),
    reward numeric,
    uncles_reward numeric,
    size character varying,
    "time" bigint,
    is_final boolean,
    uncle_hash character varying(66),
    eth_node_id integer NOT NULL,
    eth_node_fingerprint character varying(128) NOT NULL
);


--
-- Name: blocks_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.blocks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: blocks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.blocks_id_seq OWNED BY public.blocks.id;


--
-- Name: checked_headers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.checked_headers (
    id integer NOT NULL,
    header_id integer NOT NULL,
    yank integer DEFAULT 0 NOT NULL,
    new_cdp integer DEFAULT 0 NOT NULL
);


--
-- Name: checked_headers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.checked_headers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: checked_headers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.checked_headers_id_seq OWNED BY public.checked_headers.id;


--
-- Name: eth_nodes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.eth_nodes (
    id integer NOT NULL,
    client_name character varying,
    genesis_block character varying(66),
    network_id numeric,
    eth_node_id character varying(128)
);


--
-- Name: full_sync_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.full_sync_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: full_sync_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.full_sync_logs_id_seq OWNED BY public.full_sync_logs.id;


--
-- Name: full_sync_receipts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.full_sync_receipts (
    id integer NOT NULL,
    contract_address_id integer NOT NULL,
    cumulative_gas_used numeric,
    gas_used numeric,
    state_root character varying(66),
    status integer,
    tx_hash character varying(66),
    block_id integer NOT NULL
);


--
-- Name: full_sync_receipts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.full_sync_receipts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: full_sync_receipts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.full_sync_receipts_id_seq OWNED BY public.full_sync_receipts.id;


--
-- Name: full_sync_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.full_sync_transactions (
    id integer NOT NULL,
    block_id integer NOT NULL,
    gas_limit numeric,
    gas_price numeric,
    hash character varying(66),
    input_data bytea,
    nonce numeric,
    raw bytea,
    tx_from character varying(66),
    tx_index integer,
    tx_to character varying(66),
    value numeric
);


--
-- Name: full_sync_transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.full_sync_transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: full_sync_transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.full_sync_transactions_id_seq OWNED BY public.full_sync_transactions.id;


--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.goose_db_version_id_seq OWNED BY public.goose_db_version.id;


--
-- Name: header_sync_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.header_sync_logs (
    id integer NOT NULL,
    header_id integer NOT NULL,
    address integer NOT NULL,
    topics bytea[],
    data bytea,
    block_number bigint,
    block_hash character varying(66),
    tx_hash character varying(66),
    tx_index integer,
    log_index integer,
    raw jsonb,
    transformed boolean DEFAULT false NOT NULL
);


--
-- Name: header_sync_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.header_sync_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: header_sync_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.header_sync_logs_id_seq OWNED BY public.header_sync_logs.id;


--
-- Name: header_sync_receipts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.header_sync_receipts (
    id integer NOT NULL,
    transaction_id integer NOT NULL,
    header_id integer NOT NULL,
    contract_address_id integer NOT NULL,
    cumulative_gas_used numeric,
    gas_used numeric,
    state_root character varying(66),
    status integer,
    tx_hash character varying(66),
    rlp bytea
);


--
-- Name: header_sync_receipts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.header_sync_receipts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: header_sync_receipts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.header_sync_receipts_id_seq OWNED BY public.header_sync_receipts.id;


--
-- Name: header_sync_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.header_sync_transactions (
    id integer NOT NULL,
    header_id integer NOT NULL,
    hash character varying(66),
    gas_limit numeric,
    gas_price numeric,
    input_data bytea,
    nonce numeric,
    raw bytea,
    tx_from character varying(44),
    tx_index integer,
    tx_to character varying(44),
    value numeric
);


--
-- Name: header_sync_transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.header_sync_transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: header_sync_transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.header_sync_transactions_id_seq OWNED BY public.header_sync_transactions.id;


--
-- Name: headers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.headers (
    id integer NOT NULL,
    hash character varying(66),
    block_number bigint,
    raw jsonb,
    block_timestamp numeric,
    check_count integer DEFAULT 0 NOT NULL,
    eth_node_id integer NOT NULL,
    eth_node_fingerprint character varying(128)
);


--
-- Name: headers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.headers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: headers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.headers_id_seq OWNED BY public.headers.id;


--
-- Name: log_filters; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.log_filters (
    id integer NOT NULL,
    name character varying NOT NULL,
    from_block bigint,
    to_block bigint,
    address character varying(66),
    topic0 character varying(66),
    topic1 character varying(66),
    topic2 character varying(66),
    topic3 character varying(66),
    CONSTRAINT log_filters_from_block_check CHECK ((from_block >= 0)),
    CONSTRAINT log_filters_name_check CHECK (((name)::text <> ''::text)),
    CONSTRAINT log_filters_to_block_check CHECK ((to_block >= 0))
);


--
-- Name: log_filters_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.log_filters_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_filters_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.log_filters_id_seq OWNED BY public.log_filters.id;


--
-- Name: nodes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.nodes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: nodes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.nodes_id_seq OWNED BY public.eth_nodes.id;


--
-- Name: queued_storage; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.queued_storage (
    id integer NOT NULL,
    block_height bigint,
    block_hash bytea,
    contract bytea,
    storage_key bytea,
    storage_value bytea
);


--
-- Name: queued_storage_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.queued_storage_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: queued_storage_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.queued_storage_id_seq OWNED BY public.queued_storage.id;


--
-- Name: uncles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.uncles (
    id integer NOT NULL,
    hash character varying(66) NOT NULL,
    block_id integer NOT NULL,
    reward numeric NOT NULL,
    miner character varying(42) NOT NULL,
    raw jsonb,
    block_timestamp numeric,
    eth_node_id integer NOT NULL,
    eth_node_fingerprint character varying(128)
);


--
-- Name: uncles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.uncles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: uncles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.uncles_id_seq OWNED BY public.uncles.id;


--
-- Name: watched_contracts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.watched_contracts (
    contract_id integer NOT NULL,
    contract_abi json,
    contract_hash character varying(66)
);


--
-- Name: watched_contracts_contract_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.watched_contracts_contract_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: watched_contracts_contract_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.watched_contracts_contract_id_seq OWNED BY public.watched_contracts.contract_id;


--
-- Name: watched_event_logs; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.watched_event_logs AS
 SELECT log_filters.name,
    full_sync_logs.id,
    full_sync_logs.block_number,
    full_sync_logs.address,
    full_sync_logs.tx_hash,
    full_sync_logs.index,
    full_sync_logs.topic0,
    full_sync_logs.topic1,
    full_sync_logs.topic2,
    full_sync_logs.topic3,
    full_sync_logs.data,
    full_sync_logs.receipt_id
   FROM ((public.log_filters
     CROSS JOIN public.block_stats)
     JOIN public.full_sync_logs ON ((((full_sync_logs.address)::text = (log_filters.address)::text) AND (full_sync_logs.block_number >= COALESCE(log_filters.from_block, block_stats.min_block)) AND (full_sync_logs.block_number <= COALESCE(log_filters.to_block, block_stats.max_block)))))
  WHERE ((((log_filters.topic0)::text = (full_sync_logs.topic0)::text) OR (log_filters.topic0 IS NULL)) AND (((log_filters.topic1)::text = (full_sync_logs.topic1)::text) OR (log_filters.topic1 IS NULL)) AND (((log_filters.topic2)::text = (full_sync_logs.topic2)::text) OR (log_filters.topic2 IS NULL)) AND (((log_filters.topic3)::text = (full_sync_logs.topic3)::text) OR (log_filters.topic3 IS NULL)));


--
-- Name: watched_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.watched_logs (
    id integer NOT NULL,
    contract_address character varying(42),
    topic_zero character varying(66)
);


--
-- Name: watched_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.watched_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: watched_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.watched_logs_id_seq OWNED BY public.watched_logs.id;


--
-- Name: managed_cdp id; Type: DEFAULT; Schema: api; Owner: -
--

ALTER TABLE ONLY api.managed_cdp ALTER COLUMN id SET DEFAULT nextval('api.managed_cdp_id_seq'::regclass);


--
-- Name: bite id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite ALTER COLUMN id SET DEFAULT nextval('maker.bite_id_seq'::regclass);


--
-- Name: cat_file_chop_lump id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_chop_lump ALTER COLUMN id SET DEFAULT nextval('maker.cat_file_chop_lump_id_seq'::regclass);


--
-- Name: cat_file_flip id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_flip ALTER COLUMN id SET DEFAULT nextval('maker.cat_file_flip_id_seq'::regclass);


--
-- Name: cat_file_vow id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_vow ALTER COLUMN id SET DEFAULT nextval('maker.cat_file_vow_id_seq'::regclass);


--
-- Name: cat_ilk_chop id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_chop ALTER COLUMN id SET DEFAULT nextval('maker.cat_ilk_chop_id_seq'::regclass);


--
-- Name: cat_ilk_flip id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_flip ALTER COLUMN id SET DEFAULT nextval('maker.cat_ilk_flip_id_seq'::regclass);


--
-- Name: cat_ilk_lump id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_lump ALTER COLUMN id SET DEFAULT nextval('maker.cat_ilk_lump_id_seq'::regclass);


--
-- Name: cat_live id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_live ALTER COLUMN id SET DEFAULT nextval('maker.cat_live_id_seq'::regclass);


--
-- Name: cat_vat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vat ALTER COLUMN id SET DEFAULT nextval('maker.cat_vat_id_seq'::regclass);


--
-- Name: cat_vow id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vow ALTER COLUMN id SET DEFAULT nextval('maker.cat_vow_id_seq'::regclass);


--
-- Name: cdp_manager_cdpi id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_cdpi ALTER COLUMN id SET DEFAULT nextval('maker.cdp_manager_cdpi_id_seq'::regclass);


--
-- Name: cdp_manager_count id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_count ALTER COLUMN id SET DEFAULT nextval('maker.cdp_manager_count_id_seq'::regclass);


--
-- Name: cdp_manager_first id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_first ALTER COLUMN id SET DEFAULT nextval('maker.cdp_manager_first_id_seq'::regclass);


--
-- Name: cdp_manager_ilks id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_ilks ALTER COLUMN id SET DEFAULT nextval('maker.cdp_manager_ilks_id_seq'::regclass);


--
-- Name: cdp_manager_last id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_last ALTER COLUMN id SET DEFAULT nextval('maker.cdp_manager_last_id_seq'::regclass);


--
-- Name: cdp_manager_list_next id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_next ALTER COLUMN id SET DEFAULT nextval('maker.cdp_manager_list_next_id_seq'::regclass);


--
-- Name: cdp_manager_list_prev id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_prev ALTER COLUMN id SET DEFAULT nextval('maker.cdp_manager_list_prev_id_seq'::regclass);


--
-- Name: cdp_manager_owns id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_owns ALTER COLUMN id SET DEFAULT nextval('maker.cdp_manager_owns_id_seq'::regclass);


--
-- Name: cdp_manager_urns id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_urns ALTER COLUMN id SET DEFAULT nextval('maker.cdp_manager_urns_id_seq'::regclass);


--
-- Name: cdp_manager_vat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_vat ALTER COLUMN id SET DEFAULT nextval('maker.cdp_manager_vat_id_seq'::regclass);


--
-- Name: deal id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deal ALTER COLUMN id SET DEFAULT nextval('maker.deal_id_seq'::regclass);


--
-- Name: dent id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.dent ALTER COLUMN id SET DEFAULT nextval('maker.dent_id_seq'::regclass);


--
-- Name: flap id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap ALTER COLUMN id SET DEFAULT nextval('maker.flap_id_seq'::regclass);


--
-- Name: flap_beg id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_beg ALTER COLUMN id SET DEFAULT nextval('maker.flap_beg_id_seq'::regclass);


--
-- Name: flap_bid_bid id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_bid ALTER COLUMN id SET DEFAULT nextval('maker.flap_bid_bid_id_seq'::regclass);


--
-- Name: flap_bid_end id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_end ALTER COLUMN id SET DEFAULT nextval('maker.flap_bid_end_id_seq'::regclass);


--
-- Name: flap_bid_guy id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_guy ALTER COLUMN id SET DEFAULT nextval('maker.flap_bid_guy_id_seq'::regclass);


--
-- Name: flap_bid_lot id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_lot ALTER COLUMN id SET DEFAULT nextval('maker.flap_bid_lot_id_seq'::regclass);


--
-- Name: flap_bid_tic id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_tic ALTER COLUMN id SET DEFAULT nextval('maker.flap_bid_tic_id_seq'::regclass);


--
-- Name: flap_gem id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_gem ALTER COLUMN id SET DEFAULT nextval('maker.flap_gem_id_seq'::regclass);


--
-- Name: flap_kick id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kick ALTER COLUMN id SET DEFAULT nextval('maker.flap_kick_id_seq'::regclass);


--
-- Name: flap_kicks id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kicks ALTER COLUMN id SET DEFAULT nextval('maker.flap_kicks_id_seq'::regclass);


--
-- Name: flap_live id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_live ALTER COLUMN id SET DEFAULT nextval('maker.flap_live_id_seq'::regclass);


--
-- Name: flap_tau id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_tau ALTER COLUMN id SET DEFAULT nextval('maker.flap_tau_id_seq'::regclass);


--
-- Name: flap_ttl id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_ttl ALTER COLUMN id SET DEFAULT nextval('maker.flap_ttl_id_seq'::regclass);


--
-- Name: flap_vat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_vat ALTER COLUMN id SET DEFAULT nextval('maker.flap_vat_id_seq'::regclass);


--
-- Name: flip id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip ALTER COLUMN id SET DEFAULT nextval('maker.flip_id_seq'::regclass);


--
-- Name: flip_beg id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_beg ALTER COLUMN id SET DEFAULT nextval('maker.flip_beg_id_seq'::regclass);


--
-- Name: flip_bid_bid id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_bid ALTER COLUMN id SET DEFAULT nextval('maker.flip_bid_bid_id_seq'::regclass);


--
-- Name: flip_bid_end id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_end ALTER COLUMN id SET DEFAULT nextval('maker.flip_bid_end_id_seq'::regclass);


--
-- Name: flip_bid_gal id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_gal ALTER COLUMN id SET DEFAULT nextval('maker.flip_bid_gal_id_seq'::regclass);


--
-- Name: flip_bid_guy id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_guy ALTER COLUMN id SET DEFAULT nextval('maker.flip_bid_guy_id_seq'::regclass);


--
-- Name: flip_bid_lot id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_lot ALTER COLUMN id SET DEFAULT nextval('maker.flip_bid_lot_id_seq'::regclass);


--
-- Name: flip_bid_tab id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tab ALTER COLUMN id SET DEFAULT nextval('maker.flip_bid_tab_id_seq'::regclass);


--
-- Name: flip_bid_tic id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tic ALTER COLUMN id SET DEFAULT nextval('maker.flip_bid_tic_id_seq'::regclass);


--
-- Name: flip_bid_usr id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_usr ALTER COLUMN id SET DEFAULT nextval('maker.flip_bid_usr_id_seq'::regclass);


--
-- Name: flip_ilk id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ilk ALTER COLUMN id SET DEFAULT nextval('maker.flip_ilk_id_seq'::regclass);


--
-- Name: flip_kick id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kick ALTER COLUMN id SET DEFAULT nextval('maker.flip_kick_id_seq'::regclass);


--
-- Name: flip_kicks id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kicks ALTER COLUMN id SET DEFAULT nextval('maker.flip_kicks_id_seq'::regclass);


--
-- Name: flip_tau id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_tau ALTER COLUMN id SET DEFAULT nextval('maker.flip_tau_id_seq'::regclass);


--
-- Name: flip_ttl id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ttl ALTER COLUMN id SET DEFAULT nextval('maker.flip_ttl_id_seq'::regclass);


--
-- Name: flip_vat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_vat ALTER COLUMN id SET DEFAULT nextval('maker.flip_vat_id_seq'::regclass);


--
-- Name: flop id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop ALTER COLUMN id SET DEFAULT nextval('maker.flop_id_seq'::regclass);


--
-- Name: flop_beg id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_beg ALTER COLUMN id SET DEFAULT nextval('maker.flop_beg_id_seq'::regclass);


--
-- Name: flop_bid_bid id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_bid ALTER COLUMN id SET DEFAULT nextval('maker.flop_bid_bid_id_seq'::regclass);


--
-- Name: flop_bid_end id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_end ALTER COLUMN id SET DEFAULT nextval('maker.flop_bid_end_id_seq'::regclass);


--
-- Name: flop_bid_guy id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_guy ALTER COLUMN id SET DEFAULT nextval('maker.flop_bid_guy_id_seq'::regclass);


--
-- Name: flop_bid_lot id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_lot ALTER COLUMN id SET DEFAULT nextval('maker.flop_bid_lot_id_seq'::regclass);


--
-- Name: flop_bid_tic id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_tic ALTER COLUMN id SET DEFAULT nextval('maker.flop_bid_tic_id_seq'::regclass);


--
-- Name: flop_gem id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_gem ALTER COLUMN id SET DEFAULT nextval('maker.flop_gem_id_seq'::regclass);


--
-- Name: flop_kick id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kick ALTER COLUMN id SET DEFAULT nextval('maker.flop_kick_id_seq'::regclass);


--
-- Name: flop_kicks id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kicks ALTER COLUMN id SET DEFAULT nextval('maker.flop_kicks_id_seq'::regclass);


--
-- Name: flop_live id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_live ALTER COLUMN id SET DEFAULT nextval('maker.flop_live_id_seq'::regclass);


--
-- Name: flop_pad id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_pad ALTER COLUMN id SET DEFAULT nextval('maker.flop_pad_id_seq'::regclass);


--
-- Name: flop_tau id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_tau ALTER COLUMN id SET DEFAULT nextval('maker.flop_tau_id_seq'::regclass);


--
-- Name: flop_ttl id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_ttl ALTER COLUMN id SET DEFAULT nextval('maker.flop_ttl_id_seq'::regclass);


--
-- Name: flop_vat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vat ALTER COLUMN id SET DEFAULT nextval('maker.flop_vat_id_seq'::regclass);


--
-- Name: ilks id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.ilks ALTER COLUMN id SET DEFAULT nextval('maker.ilks_id_seq'::regclass);


--
-- Name: jug_base id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_base ALTER COLUMN id SET DEFAULT nextval('maker.jug_base_id_seq'::regclass);


--
-- Name: jug_drip id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_drip ALTER COLUMN id SET DEFAULT nextval('maker.jug_drip_id_seq'::regclass);


--
-- Name: jug_file_base id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_base ALTER COLUMN id SET DEFAULT nextval('maker.jug_file_base_id_seq'::regclass);


--
-- Name: jug_file_ilk id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_ilk ALTER COLUMN id SET DEFAULT nextval('maker.jug_file_ilk_id_seq'::regclass);


--
-- Name: jug_file_vow id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow ALTER COLUMN id SET DEFAULT nextval('maker.jug_file_vow_id_seq'::regclass);


--
-- Name: jug_ilk_duty id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_duty ALTER COLUMN id SET DEFAULT nextval('maker.jug_ilk_duty_id_seq'::regclass);


--
-- Name: jug_ilk_rho id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_rho ALTER COLUMN id SET DEFAULT nextval('maker.jug_ilk_rho_id_seq'::regclass);


--
-- Name: jug_init id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_init ALTER COLUMN id SET DEFAULT nextval('maker.jug_init_id_seq'::regclass);


--
-- Name: jug_vat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vat ALTER COLUMN id SET DEFAULT nextval('maker.jug_vat_id_seq'::regclass);


--
-- Name: jug_vow id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vow ALTER COLUMN id SET DEFAULT nextval('maker.jug_vow_id_seq'::regclass);


--
-- Name: new_cdp id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.new_cdp ALTER COLUMN id SET DEFAULT nextval('maker.new_cdp_id_seq'::regclass);


--
-- Name: spot_file_mat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_mat ALTER COLUMN id SET DEFAULT nextval('maker.spot_file_mat_id_seq'::regclass);


--
-- Name: spot_file_pip id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_pip ALTER COLUMN id SET DEFAULT nextval('maker.spot_file_pip_id_seq'::regclass);


--
-- Name: spot_ilk_mat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_mat ALTER COLUMN id SET DEFAULT nextval('maker.spot_ilk_mat_id_seq'::regclass);


--
-- Name: spot_ilk_pip id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_pip ALTER COLUMN id SET DEFAULT nextval('maker.spot_ilk_pip_id_seq'::regclass);


--
-- Name: spot_par id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_par ALTER COLUMN id SET DEFAULT nextval('maker.spot_par_id_seq'::regclass);


--
-- Name: spot_poke id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_poke ALTER COLUMN id SET DEFAULT nextval('maker.spot_poke_id_seq'::regclass);


--
-- Name: spot_vat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_vat ALTER COLUMN id SET DEFAULT nextval('maker.spot_vat_id_seq'::regclass);


--
-- Name: tend id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tend ALTER COLUMN id SET DEFAULT nextval('maker.tend_id_seq'::regclass);


--
-- Name: tick id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tick ALTER COLUMN id SET DEFAULT nextval('maker.tick_id_seq'::regclass);


--
-- Name: urns id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.urns ALTER COLUMN id SET DEFAULT nextval('maker.urns_id_seq'::regclass);


--
-- Name: vat_dai id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_dai ALTER COLUMN id SET DEFAULT nextval('maker.vat_dai_id_seq'::regclass);


--
-- Name: vat_debt id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_debt ALTER COLUMN id SET DEFAULT nextval('maker.vat_debt_id_seq'::regclass);


--
-- Name: vat_file_debt_ceiling id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_debt_ceiling ALTER COLUMN id SET DEFAULT nextval('maker.vat_file_debt_ceiling_id_seq'::regclass);


--
-- Name: vat_file_ilk id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_ilk ALTER COLUMN id SET DEFAULT nextval('maker.vat_file_ilk_id_seq'::regclass);


--
-- Name: vat_flux id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_flux ALTER COLUMN id SET DEFAULT nextval('maker.vat_flux_id_seq'::regclass);


--
-- Name: vat_fold id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold ALTER COLUMN id SET DEFAULT nextval('maker.vat_fold_id_seq'::regclass);


--
-- Name: vat_fork id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fork ALTER COLUMN id SET DEFAULT nextval('maker.vat_fork_id_seq'::regclass);


--
-- Name: vat_frob id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob ALTER COLUMN id SET DEFAULT nextval('maker.vat_frob_id_seq'::regclass);


--
-- Name: vat_gem id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_gem ALTER COLUMN id SET DEFAULT nextval('maker.vat_gem_id_seq'::regclass);


--
-- Name: vat_grab id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_grab ALTER COLUMN id SET DEFAULT nextval('maker.vat_grab_id_seq'::regclass);


--
-- Name: vat_heal id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_heal ALTER COLUMN id SET DEFAULT nextval('maker.vat_heal_id_seq'::regclass);


--
-- Name: vat_ilk_art id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_art ALTER COLUMN id SET DEFAULT nextval('maker.vat_ilk_art_id_seq'::regclass);


--
-- Name: vat_ilk_dust id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_dust ALTER COLUMN id SET DEFAULT nextval('maker.vat_ilk_dust_id_seq'::regclass);


--
-- Name: vat_ilk_line id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_line ALTER COLUMN id SET DEFAULT nextval('maker.vat_ilk_line_id_seq'::regclass);


--
-- Name: vat_ilk_rate id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_rate ALTER COLUMN id SET DEFAULT nextval('maker.vat_ilk_rate_id_seq'::regclass);


--
-- Name: vat_ilk_spot id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_spot ALTER COLUMN id SET DEFAULT nextval('maker.vat_ilk_spot_id_seq'::regclass);


--
-- Name: vat_init id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_init ALTER COLUMN id SET DEFAULT nextval('maker.vat_init_id_seq'::regclass);


--
-- Name: vat_line id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_line ALTER COLUMN id SET DEFAULT nextval('maker.vat_line_id_seq'::regclass);


--
-- Name: vat_live id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_live ALTER COLUMN id SET DEFAULT nextval('maker.vat_live_id_seq'::regclass);


--
-- Name: vat_move id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_move ALTER COLUMN id SET DEFAULT nextval('maker.vat_move_id_seq'::regclass);


--
-- Name: vat_sin id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_sin ALTER COLUMN id SET DEFAULT nextval('maker.vat_sin_id_seq'::regclass);


--
-- Name: vat_slip id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_slip ALTER COLUMN id SET DEFAULT nextval('maker.vat_slip_id_seq'::regclass);


--
-- Name: vat_suck id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_suck ALTER COLUMN id SET DEFAULT nextval('maker.vat_suck_id_seq'::regclass);


--
-- Name: vat_urn_art id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_art ALTER COLUMN id SET DEFAULT nextval('maker.vat_urn_art_id_seq'::regclass);


--
-- Name: vat_urn_ink id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_ink ALTER COLUMN id SET DEFAULT nextval('maker.vat_urn_ink_id_seq'::regclass);


--
-- Name: vat_vice id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_vice ALTER COLUMN id SET DEFAULT nextval('maker.vat_vice_id_seq'::regclass);


--
-- Name: vow_ash id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_ash ALTER COLUMN id SET DEFAULT nextval('maker.vow_ash_id_seq'::regclass);


--
-- Name: vow_bump id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_bump ALTER COLUMN id SET DEFAULT nextval('maker.vow_bump_id_seq'::regclass);


--
-- Name: vow_dump id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_dump ALTER COLUMN id SET DEFAULT nextval('maker.vow_dump_id_seq'::regclass);


--
-- Name: vow_fess id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_fess ALTER COLUMN id SET DEFAULT nextval('maker.vow_fess_id_seq'::regclass);


--
-- Name: vow_file id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file ALTER COLUMN id SET DEFAULT nextval('maker.vow_file_id_seq'::regclass);


--
-- Name: vow_flapper id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flapper ALTER COLUMN id SET DEFAULT nextval('maker.vow_flapper_id_seq'::regclass);


--
-- Name: vow_flog id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flog ALTER COLUMN id SET DEFAULT nextval('maker.vow_flog_id_seq'::regclass);


--
-- Name: vow_flopper id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flopper ALTER COLUMN id SET DEFAULT nextval('maker.vow_flopper_id_seq'::regclass);


--
-- Name: vow_hump id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_hump ALTER COLUMN id SET DEFAULT nextval('maker.vow_hump_id_seq'::regclass);


--
-- Name: vow_sin_integer id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_integer ALTER COLUMN id SET DEFAULT nextval('maker.vow_sin_integer_id_seq'::regclass);


--
-- Name: vow_sin_mapping id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_mapping ALTER COLUMN id SET DEFAULT nextval('maker.vow_sin_mapping_id_seq'::regclass);


--
-- Name: vow_sump id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sump ALTER COLUMN id SET DEFAULT nextval('maker.vow_sump_id_seq'::regclass);


--
-- Name: vow_vat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_vat ALTER COLUMN id SET DEFAULT nextval('maker.vow_vat_id_seq'::regclass);


--
-- Name: vow_wait id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_wait ALTER COLUMN id SET DEFAULT nextval('maker.vow_wait_id_seq'::regclass);


--
-- Name: yank id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.yank ALTER COLUMN id SET DEFAULT nextval('maker.yank_id_seq'::regclass);


--
-- Name: addresses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.addresses ALTER COLUMN id SET DEFAULT nextval('public.addresses_id_seq'::regclass);


--
-- Name: blocks id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks ALTER COLUMN id SET DEFAULT nextval('public.blocks_id_seq'::regclass);


--
-- Name: checked_headers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers ALTER COLUMN id SET DEFAULT nextval('public.checked_headers_id_seq'::regclass);


--
-- Name: eth_nodes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes ALTER COLUMN id SET DEFAULT nextval('public.nodes_id_seq'::regclass);


--
-- Name: full_sync_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_logs ALTER COLUMN id SET DEFAULT nextval('public.full_sync_logs_id_seq'::regclass);


--
-- Name: full_sync_receipts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_receipts ALTER COLUMN id SET DEFAULT nextval('public.full_sync_receipts_id_seq'::regclass);


--
-- Name: full_sync_transactions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_transactions ALTER COLUMN id SET DEFAULT nextval('public.full_sync_transactions_id_seq'::regclass);


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: header_sync_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_logs ALTER COLUMN id SET DEFAULT nextval('public.header_sync_logs_id_seq'::regclass);


--
-- Name: header_sync_receipts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_receipts ALTER COLUMN id SET DEFAULT nextval('public.header_sync_receipts_id_seq'::regclass);


--
-- Name: header_sync_transactions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_transactions ALTER COLUMN id SET DEFAULT nextval('public.header_sync_transactions_id_seq'::regclass);


--
-- Name: headers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers ALTER COLUMN id SET DEFAULT nextval('public.headers_id_seq'::regclass);


--
-- Name: log_filters id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.log_filters ALTER COLUMN id SET DEFAULT nextval('public.log_filters_id_seq'::regclass);


--
-- Name: queued_storage id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage ALTER COLUMN id SET DEFAULT nextval('public.queued_storage_id_seq'::regclass);


--
-- Name: uncles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles ALTER COLUMN id SET DEFAULT nextval('public.uncles_id_seq'::regclass);


--
-- Name: watched_contracts contract_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_contracts ALTER COLUMN contract_id SET DEFAULT nextval('public.watched_contracts_contract_id_seq'::regclass);


--
-- Name: watched_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_logs ALTER COLUMN id SET DEFAULT nextval('public.watched_logs_id_seq'::regclass);


--
-- Name: current_ilk_state current_ilk_state_pkey; Type: CONSTRAINT; Schema: api; Owner: -
--

ALTER TABLE ONLY api.current_ilk_state
    ADD CONSTRAINT current_ilk_state_pkey PRIMARY KEY (ilk_identifier);


--
-- Name: managed_cdp managed_cdp_cdpi_key; Type: CONSTRAINT; Schema: api; Owner: -
--

ALTER TABLE ONLY api.managed_cdp
    ADD CONSTRAINT managed_cdp_cdpi_key UNIQUE (cdpi);


--
-- Name: managed_cdp managed_cdp_pkey; Type: CONSTRAINT; Schema: api; Owner: -
--

ALTER TABLE ONLY api.managed_cdp
    ADD CONSTRAINT managed_cdp_pkey PRIMARY KEY (id);


--
-- Name: bite bite_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite
    ADD CONSTRAINT bite_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: bite bite_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite
    ADD CONSTRAINT bite_pkey PRIMARY KEY (id);


--
-- Name: cat_file_chop_lump cat_file_chop_lump_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_chop_lump
    ADD CONSTRAINT cat_file_chop_lump_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: cat_file_chop_lump cat_file_chop_lump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_chop_lump
    ADD CONSTRAINT cat_file_chop_lump_pkey PRIMARY KEY (id);


--
-- Name: cat_file_flip cat_file_flip_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_flip
    ADD CONSTRAINT cat_file_flip_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: cat_file_flip cat_file_flip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_flip
    ADD CONSTRAINT cat_file_flip_pkey PRIMARY KEY (id);


--
-- Name: cat_file_vow cat_file_vow_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_vow
    ADD CONSTRAINT cat_file_vow_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: cat_file_vow cat_file_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_vow
    ADD CONSTRAINT cat_file_vow_pkey PRIMARY KEY (id);


--
-- Name: cat_ilk_chop cat_ilk_chop_block_number_block_hash_ilk_id_chop_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_chop
    ADD CONSTRAINT cat_ilk_chop_block_number_block_hash_ilk_id_chop_key UNIQUE (block_number, block_hash, ilk_id, chop);


--
-- Name: cat_ilk_chop cat_ilk_chop_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_chop
    ADD CONSTRAINT cat_ilk_chop_pkey PRIMARY KEY (id);


--
-- Name: cat_ilk_flip cat_ilk_flip_block_number_block_hash_ilk_id_flip_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_flip
    ADD CONSTRAINT cat_ilk_flip_block_number_block_hash_ilk_id_flip_key UNIQUE (block_number, block_hash, ilk_id, flip);


--
-- Name: cat_ilk_flip cat_ilk_flip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_flip
    ADD CONSTRAINT cat_ilk_flip_pkey PRIMARY KEY (id);


--
-- Name: cat_ilk_lump cat_ilk_lump_block_number_block_hash_ilk_id_lump_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_lump
    ADD CONSTRAINT cat_ilk_lump_block_number_block_hash_ilk_id_lump_key UNIQUE (block_number, block_hash, ilk_id, lump);


--
-- Name: cat_ilk_lump cat_ilk_lump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_lump
    ADD CONSTRAINT cat_ilk_lump_pkey PRIMARY KEY (id);


--
-- Name: cat_live cat_live_block_number_block_hash_live_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_live
    ADD CONSTRAINT cat_live_block_number_block_hash_live_key UNIQUE (block_number, block_hash, live);


--
-- Name: cat_live cat_live_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_live
    ADD CONSTRAINT cat_live_pkey PRIMARY KEY (id);


--
-- Name: cat_vat cat_vat_block_number_block_hash_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vat
    ADD CONSTRAINT cat_vat_block_number_block_hash_vat_key UNIQUE (block_number, block_hash, vat);


--
-- Name: cat_vat cat_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vat
    ADD CONSTRAINT cat_vat_pkey PRIMARY KEY (id);


--
-- Name: cat_vow cat_vow_block_number_block_hash_vow_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vow
    ADD CONSTRAINT cat_vow_block_number_block_hash_vow_key UNIQUE (block_number, block_hash, vow);


--
-- Name: cat_vow cat_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vow
    ADD CONSTRAINT cat_vow_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_cdpi cdp_manager_cdpi_block_number_block_hash_cdpi_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_cdpi
    ADD CONSTRAINT cdp_manager_cdpi_block_number_block_hash_cdpi_key UNIQUE (block_number, block_hash, cdpi);


--
-- Name: cdp_manager_cdpi cdp_manager_cdpi_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_cdpi
    ADD CONSTRAINT cdp_manager_cdpi_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_count cdp_manager_count_block_number_block_hash_owner_count_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_count
    ADD CONSTRAINT cdp_manager_count_block_number_block_hash_owner_count_key UNIQUE (block_number, block_hash, owner, count);


--
-- Name: cdp_manager_count cdp_manager_count_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_count
    ADD CONSTRAINT cdp_manager_count_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_first cdp_manager_first_block_number_block_hash_owner_first_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_first
    ADD CONSTRAINT cdp_manager_first_block_number_block_hash_owner_first_key UNIQUE (block_number, block_hash, owner, first);


--
-- Name: cdp_manager_first cdp_manager_first_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_first
    ADD CONSTRAINT cdp_manager_first_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_ilks cdp_manager_ilks_block_number_block_hash_cdpi_ilk_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_ilks
    ADD CONSTRAINT cdp_manager_ilks_block_number_block_hash_cdpi_ilk_id_key UNIQUE (block_number, block_hash, cdpi, ilk_id);


--
-- Name: cdp_manager_ilks cdp_manager_ilks_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_ilks
    ADD CONSTRAINT cdp_manager_ilks_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_last cdp_manager_last_block_number_block_hash_owner_last_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_last
    ADD CONSTRAINT cdp_manager_last_block_number_block_hash_owner_last_key UNIQUE (block_number, block_hash, owner, last);


--
-- Name: cdp_manager_last cdp_manager_last_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_last
    ADD CONSTRAINT cdp_manager_last_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_list_next cdp_manager_list_next_block_number_block_hash_cdpi_next_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_next
    ADD CONSTRAINT cdp_manager_list_next_block_number_block_hash_cdpi_next_key UNIQUE (block_number, block_hash, cdpi, next);


--
-- Name: cdp_manager_list_next cdp_manager_list_next_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_next
    ADD CONSTRAINT cdp_manager_list_next_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_list_prev cdp_manager_list_prev_block_number_block_hash_cdpi_prev_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_prev
    ADD CONSTRAINT cdp_manager_list_prev_block_number_block_hash_cdpi_prev_key UNIQUE (block_number, block_hash, cdpi, prev);


--
-- Name: cdp_manager_list_prev cdp_manager_list_prev_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_prev
    ADD CONSTRAINT cdp_manager_list_prev_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_owns cdp_manager_owns_block_number_block_hash_cdpi_owner_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_owns
    ADD CONSTRAINT cdp_manager_owns_block_number_block_hash_cdpi_owner_key UNIQUE (block_number, block_hash, cdpi, owner);


--
-- Name: cdp_manager_owns cdp_manager_owns_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_owns
    ADD CONSTRAINT cdp_manager_owns_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_urns cdp_manager_urns_block_number_block_hash_cdpi_urn_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_urns
    ADD CONSTRAINT cdp_manager_urns_block_number_block_hash_cdpi_urn_key UNIQUE (block_number, block_hash, cdpi, urn);


--
-- Name: cdp_manager_urns cdp_manager_urns_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_urns
    ADD CONSTRAINT cdp_manager_urns_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_vat cdp_manager_vat_block_number_block_hash_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_vat
    ADD CONSTRAINT cdp_manager_vat_block_number_block_hash_vat_key UNIQUE (block_number, block_hash, vat);


--
-- Name: cdp_manager_vat cdp_manager_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_vat
    ADD CONSTRAINT cdp_manager_vat_pkey PRIMARY KEY (id);


--
-- Name: deal deal_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deal
    ADD CONSTRAINT deal_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: deal deal_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deal
    ADD CONSTRAINT deal_pkey PRIMARY KEY (id);


--
-- Name: dent dent_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.dent
    ADD CONSTRAINT dent_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: dent dent_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.dent
    ADD CONSTRAINT dent_pkey PRIMARY KEY (id);


--
-- Name: flap_beg flap_beg_block_number_block_hash_address_id_beg_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_beg
    ADD CONSTRAINT flap_beg_block_number_block_hash_address_id_beg_key UNIQUE (block_number, block_hash, address_id, beg);


--
-- Name: flap_beg flap_beg_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_beg
    ADD CONSTRAINT flap_beg_pkey PRIMARY KEY (id);


--
-- Name: flap_bid_bid flap_bid_bid_block_number_block_hash_address_id_bid_id_bid_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_bid
    ADD CONSTRAINT flap_bid_bid_block_number_block_hash_address_id_bid_id_bid_key UNIQUE (block_number, block_hash, address_id, bid_id, bid);


--
-- Name: flap_bid_bid flap_bid_bid_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_bid
    ADD CONSTRAINT flap_bid_bid_pkey PRIMARY KEY (id);


--
-- Name: flap_bid_end flap_bid_end_block_number_block_hash_address_id_bid_id_end_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_end
    ADD CONSTRAINT flap_bid_end_block_number_block_hash_address_id_bid_id_end_key UNIQUE (block_number, block_hash, address_id, bid_id, "end");


--
-- Name: flap_bid_end flap_bid_end_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_end
    ADD CONSTRAINT flap_bid_end_pkey PRIMARY KEY (id);


--
-- Name: flap_bid_guy flap_bid_guy_block_number_block_hash_address_id_bid_id_guy_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_guy
    ADD CONSTRAINT flap_bid_guy_block_number_block_hash_address_id_bid_id_guy_key UNIQUE (block_number, block_hash, address_id, bid_id, guy);


--
-- Name: flap_bid_guy flap_bid_guy_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_guy
    ADD CONSTRAINT flap_bid_guy_pkey PRIMARY KEY (id);


--
-- Name: flap_bid_lot flap_bid_lot_block_number_block_hash_address_id_bid_id_lot_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_lot
    ADD CONSTRAINT flap_bid_lot_block_number_block_hash_address_id_bid_id_lot_key UNIQUE (block_number, block_hash, address_id, bid_id, lot);


--
-- Name: flap_bid_lot flap_bid_lot_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_lot
    ADD CONSTRAINT flap_bid_lot_pkey PRIMARY KEY (id);


--
-- Name: flap_bid_tic flap_bid_tic_block_number_block_hash_address_id_bid_id_tic_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_tic
    ADD CONSTRAINT flap_bid_tic_block_number_block_hash_address_id_bid_id_tic_key UNIQUE (block_number, block_hash, address_id, bid_id, tic);


--
-- Name: flap_bid_tic flap_bid_tic_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_tic
    ADD CONSTRAINT flap_bid_tic_pkey PRIMARY KEY (id);


--
-- Name: flap flap_block_number_bid_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap
    ADD CONSTRAINT flap_block_number_bid_id_key UNIQUE (block_number, bid_id);


--
-- Name: flap_gem flap_gem_block_number_block_hash_address_id_gem_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_gem
    ADD CONSTRAINT flap_gem_block_number_block_hash_address_id_gem_key UNIQUE (block_number, block_hash, address_id, gem);


--
-- Name: flap_gem flap_gem_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_gem
    ADD CONSTRAINT flap_gem_pkey PRIMARY KEY (id);


--
-- Name: flap_kick flap_kick_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kick
    ADD CONSTRAINT flap_kick_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: flap_kick flap_kick_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kick
    ADD CONSTRAINT flap_kick_pkey PRIMARY KEY (id);


--
-- Name: flap_kicks flap_kicks_block_number_block_hash_address_id_kicks_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kicks
    ADD CONSTRAINT flap_kicks_block_number_block_hash_address_id_kicks_key UNIQUE (block_number, block_hash, address_id, kicks);


--
-- Name: flap_kicks flap_kicks_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kicks
    ADD CONSTRAINT flap_kicks_pkey PRIMARY KEY (id);


--
-- Name: flap_live flap_live_block_number_block_hash_address_id_live_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_live
    ADD CONSTRAINT flap_live_block_number_block_hash_address_id_live_key UNIQUE (block_number, block_hash, address_id, live);


--
-- Name: flap_live flap_live_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_live
    ADD CONSTRAINT flap_live_pkey PRIMARY KEY (id);


--
-- Name: flap flap_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap
    ADD CONSTRAINT flap_pkey PRIMARY KEY (id);


--
-- Name: flap_tau flap_tau_block_number_block_hash_address_id_tau_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_tau
    ADD CONSTRAINT flap_tau_block_number_block_hash_address_id_tau_key UNIQUE (block_number, block_hash, address_id, tau);


--
-- Name: flap_tau flap_tau_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_tau
    ADD CONSTRAINT flap_tau_pkey PRIMARY KEY (id);


--
-- Name: flap_ttl flap_ttl_block_number_block_hash_address_id_ttl_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_ttl
    ADD CONSTRAINT flap_ttl_block_number_block_hash_address_id_ttl_key UNIQUE (block_number, block_hash, address_id, ttl);


--
-- Name: flap_ttl flap_ttl_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_ttl
    ADD CONSTRAINT flap_ttl_pkey PRIMARY KEY (id);


--
-- Name: flap_vat flap_vat_block_number_block_hash_address_id_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_vat
    ADD CONSTRAINT flap_vat_block_number_block_hash_address_id_vat_key UNIQUE (block_number, block_hash, address_id, vat);


--
-- Name: flap_vat flap_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_vat
    ADD CONSTRAINT flap_vat_pkey PRIMARY KEY (id);


--
-- Name: flip_beg flip_beg_block_number_block_hash_address_id_beg_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_beg
    ADD CONSTRAINT flip_beg_block_number_block_hash_address_id_beg_key UNIQUE (block_number, block_hash, address_id, beg);


--
-- Name: flip_beg flip_beg_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_beg
    ADD CONSTRAINT flip_beg_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_bid flip_bid_bid_block_number_block_hash_bid_id_address_id_bid_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_bid
    ADD CONSTRAINT flip_bid_bid_block_number_block_hash_bid_id_address_id_bid_key UNIQUE (block_number, block_hash, bid_id, address_id, bid);


--
-- Name: flip_bid_bid flip_bid_bid_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_bid
    ADD CONSTRAINT flip_bid_bid_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_end flip_bid_end_block_number_block_hash_bid_id_address_id_end_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_end
    ADD CONSTRAINT flip_bid_end_block_number_block_hash_bid_id_address_id_end_key UNIQUE (block_number, block_hash, bid_id, address_id, "end");


--
-- Name: flip_bid_end flip_bid_end_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_end
    ADD CONSTRAINT flip_bid_end_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_gal flip_bid_gal_block_number_block_hash_bid_id_address_id_gal_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_gal
    ADD CONSTRAINT flip_bid_gal_block_number_block_hash_bid_id_address_id_gal_key UNIQUE (block_number, block_hash, bid_id, address_id, gal);


--
-- Name: flip_bid_gal flip_bid_gal_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_gal
    ADD CONSTRAINT flip_bid_gal_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_guy flip_bid_guy_block_number_block_hash_bid_id_address_id_guy_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_guy
    ADD CONSTRAINT flip_bid_guy_block_number_block_hash_bid_id_address_id_guy_key UNIQUE (block_number, block_hash, bid_id, address_id, guy);


--
-- Name: flip_bid_guy flip_bid_guy_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_guy
    ADD CONSTRAINT flip_bid_guy_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_lot flip_bid_lot_block_number_block_hash_bid_id_address_id_lot_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_lot
    ADD CONSTRAINT flip_bid_lot_block_number_block_hash_bid_id_address_id_lot_key UNIQUE (block_number, block_hash, bid_id, address_id, lot);


--
-- Name: flip_bid_lot flip_bid_lot_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_lot
    ADD CONSTRAINT flip_bid_lot_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_tab flip_bid_tab_block_number_block_hash_bid_id_address_id_tab_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tab
    ADD CONSTRAINT flip_bid_tab_block_number_block_hash_bid_id_address_id_tab_key UNIQUE (block_number, block_hash, bid_id, address_id, tab);


--
-- Name: flip_bid_tab flip_bid_tab_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tab
    ADD CONSTRAINT flip_bid_tab_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_tic flip_bid_tic_block_number_block_hash_bid_id_address_id_tic_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tic
    ADD CONSTRAINT flip_bid_tic_block_number_block_hash_bid_id_address_id_tic_key UNIQUE (block_number, block_hash, bid_id, address_id, tic);


--
-- Name: flip_bid_tic flip_bid_tic_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tic
    ADD CONSTRAINT flip_bid_tic_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_usr flip_bid_usr_block_number_block_hash_bid_id_address_id_usr_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_usr
    ADD CONSTRAINT flip_bid_usr_block_number_block_hash_bid_id_address_id_usr_key UNIQUE (block_number, block_hash, bid_id, address_id, usr);


--
-- Name: flip_bid_usr flip_bid_usr_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_usr
    ADD CONSTRAINT flip_bid_usr_pkey PRIMARY KEY (id);


--
-- Name: flip flip_block_number_bid_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip
    ADD CONSTRAINT flip_block_number_bid_id_key UNIQUE (block_number, bid_id);


--
-- Name: flip_ilk flip_ilk_block_number_block_hash_address_id_ilk_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ilk
    ADD CONSTRAINT flip_ilk_block_number_block_hash_address_id_ilk_id_key UNIQUE (block_number, block_hash, address_id, ilk_id);


--
-- Name: flip_ilk flip_ilk_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ilk
    ADD CONSTRAINT flip_ilk_pkey PRIMARY KEY (id);


--
-- Name: flip_kick flip_kick_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kick
    ADD CONSTRAINT flip_kick_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: flip_kick flip_kick_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kick
    ADD CONSTRAINT flip_kick_pkey PRIMARY KEY (id);


--
-- Name: flip_kicks flip_kicks_block_number_block_hash_address_id_kicks_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kicks
    ADD CONSTRAINT flip_kicks_block_number_block_hash_address_id_kicks_key UNIQUE (block_number, block_hash, address_id, kicks);


--
-- Name: flip_kicks flip_kicks_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kicks
    ADD CONSTRAINT flip_kicks_pkey PRIMARY KEY (id);


--
-- Name: flip flip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip
    ADD CONSTRAINT flip_pkey PRIMARY KEY (id);


--
-- Name: flip_tau flip_tau_block_number_block_hash_address_id_tau_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_tau
    ADD CONSTRAINT flip_tau_block_number_block_hash_address_id_tau_key UNIQUE (block_number, block_hash, address_id, tau);


--
-- Name: flip_tau flip_tau_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_tau
    ADD CONSTRAINT flip_tau_pkey PRIMARY KEY (id);


--
-- Name: flip_ttl flip_ttl_block_number_block_hash_address_id_ttl_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ttl
    ADD CONSTRAINT flip_ttl_block_number_block_hash_address_id_ttl_key UNIQUE (block_number, block_hash, address_id, ttl);


--
-- Name: flip_ttl flip_ttl_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ttl
    ADD CONSTRAINT flip_ttl_pkey PRIMARY KEY (id);


--
-- Name: flip_vat flip_vat_block_number_block_hash_address_id_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_vat
    ADD CONSTRAINT flip_vat_block_number_block_hash_address_id_vat_key UNIQUE (block_number, block_hash, address_id, vat);


--
-- Name: flip_vat flip_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_vat
    ADD CONSTRAINT flip_vat_pkey PRIMARY KEY (id);


--
-- Name: flop_beg flop_beg_block_number_block_hash_address_id_beg_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_beg
    ADD CONSTRAINT flop_beg_block_number_block_hash_address_id_beg_key UNIQUE (block_number, block_hash, address_id, beg);


--
-- Name: flop_beg flop_beg_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_beg
    ADD CONSTRAINT flop_beg_pkey PRIMARY KEY (id);


--
-- Name: flop_bid_bid flop_bid_bid_block_number_block_hash_bid_id_address_id_bid_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_bid
    ADD CONSTRAINT flop_bid_bid_block_number_block_hash_bid_id_address_id_bid_key UNIQUE (block_number, block_hash, bid_id, address_id, bid);


--
-- Name: flop_bid_bid flop_bid_bid_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_bid
    ADD CONSTRAINT flop_bid_bid_pkey PRIMARY KEY (id);


--
-- Name: flop_bid_end flop_bid_end_block_number_block_hash_bid_id_address_id_end_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_end
    ADD CONSTRAINT flop_bid_end_block_number_block_hash_bid_id_address_id_end_key UNIQUE (block_number, block_hash, bid_id, address_id, "end");


--
-- Name: flop_bid_end flop_bid_end_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_end
    ADD CONSTRAINT flop_bid_end_pkey PRIMARY KEY (id);


--
-- Name: flop_bid_guy flop_bid_guy_block_number_block_hash_bid_id_address_id_guy_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_guy
    ADD CONSTRAINT flop_bid_guy_block_number_block_hash_bid_id_address_id_guy_key UNIQUE (block_number, block_hash, bid_id, address_id, guy);


--
-- Name: flop_bid_guy flop_bid_guy_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_guy
    ADD CONSTRAINT flop_bid_guy_pkey PRIMARY KEY (id);


--
-- Name: flop_bid_lot flop_bid_lot_block_number_block_hash_bid_id_address_id_lot_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_lot
    ADD CONSTRAINT flop_bid_lot_block_number_block_hash_bid_id_address_id_lot_key UNIQUE (block_number, block_hash, bid_id, address_id, lot);


--
-- Name: flop_bid_lot flop_bid_lot_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_lot
    ADD CONSTRAINT flop_bid_lot_pkey PRIMARY KEY (id);


--
-- Name: flop_bid_tic flop_bid_tic_block_number_block_hash_bid_id_address_id_tic_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_tic
    ADD CONSTRAINT flop_bid_tic_block_number_block_hash_bid_id_address_id_tic_key UNIQUE (block_number, block_hash, bid_id, address_id, tic);


--
-- Name: flop_bid_tic flop_bid_tic_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_tic
    ADD CONSTRAINT flop_bid_tic_pkey PRIMARY KEY (id);


--
-- Name: flop flop_block_number_bid_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop
    ADD CONSTRAINT flop_block_number_bid_id_key UNIQUE (block_number, bid_id);


--
-- Name: flop_gem flop_gem_block_number_block_hash_address_id_gem_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_gem
    ADD CONSTRAINT flop_gem_block_number_block_hash_address_id_gem_key UNIQUE (block_number, block_hash, address_id, gem);


--
-- Name: flop_gem flop_gem_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_gem
    ADD CONSTRAINT flop_gem_pkey PRIMARY KEY (id);


--
-- Name: flop_kick flop_kick_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kick
    ADD CONSTRAINT flop_kick_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: flop_kick flop_kick_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kick
    ADD CONSTRAINT flop_kick_pkey PRIMARY KEY (id);


--
-- Name: flop_kicks flop_kicks_block_number_block_hash_address_id_kicks_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kicks
    ADD CONSTRAINT flop_kicks_block_number_block_hash_address_id_kicks_key UNIQUE (block_number, block_hash, address_id, kicks);


--
-- Name: flop_kicks flop_kicks_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kicks
    ADD CONSTRAINT flop_kicks_pkey PRIMARY KEY (id);


--
-- Name: flop_live flop_live_block_number_block_hash_address_id_live_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_live
    ADD CONSTRAINT flop_live_block_number_block_hash_address_id_live_key UNIQUE (block_number, block_hash, address_id, live);


--
-- Name: flop_live flop_live_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_live
    ADD CONSTRAINT flop_live_pkey PRIMARY KEY (id);


--
-- Name: flop_pad flop_pad_block_number_block_hash_address_id_pad_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_pad
    ADD CONSTRAINT flop_pad_block_number_block_hash_address_id_pad_key UNIQUE (block_number, block_hash, address_id, pad);


--
-- Name: flop_pad flop_pad_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_pad
    ADD CONSTRAINT flop_pad_pkey PRIMARY KEY (id);


--
-- Name: flop flop_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop
    ADD CONSTRAINT flop_pkey PRIMARY KEY (id);


--
-- Name: flop_tau flop_tau_block_number_block_hash_address_id_tau_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_tau
    ADD CONSTRAINT flop_tau_block_number_block_hash_address_id_tau_key UNIQUE (block_number, block_hash, address_id, tau);


--
-- Name: flop_tau flop_tau_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_tau
    ADD CONSTRAINT flop_tau_pkey PRIMARY KEY (id);


--
-- Name: flop_ttl flop_ttl_block_number_block_hash_address_id_ttl_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_ttl
    ADD CONSTRAINT flop_ttl_block_number_block_hash_address_id_ttl_key UNIQUE (block_number, block_hash, address_id, ttl);


--
-- Name: flop_ttl flop_ttl_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_ttl
    ADD CONSTRAINT flop_ttl_pkey PRIMARY KEY (id);


--
-- Name: flop_vat flop_vat_block_number_block_hash_address_id_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vat
    ADD CONSTRAINT flop_vat_block_number_block_hash_address_id_vat_key UNIQUE (block_number, block_hash, address_id, vat);


--
-- Name: flop_vat flop_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vat
    ADD CONSTRAINT flop_vat_pkey PRIMARY KEY (id);


--
-- Name: ilks ilks_identifier_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.ilks
    ADD CONSTRAINT ilks_identifier_key UNIQUE (identifier);


--
-- Name: ilks ilks_ilk_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.ilks
    ADD CONSTRAINT ilks_ilk_key UNIQUE (ilk);


--
-- Name: ilks ilks_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.ilks
    ADD CONSTRAINT ilks_pkey PRIMARY KEY (id);


--
-- Name: jug_base jug_base_block_number_block_hash_base_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_base
    ADD CONSTRAINT jug_base_block_number_block_hash_base_key UNIQUE (block_number, block_hash, base);


--
-- Name: jug_base jug_base_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_base
    ADD CONSTRAINT jug_base_pkey PRIMARY KEY (id);


--
-- Name: jug_drip jug_drip_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_drip
    ADD CONSTRAINT jug_drip_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: jug_drip jug_drip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_drip
    ADD CONSTRAINT jug_drip_pkey PRIMARY KEY (id);


--
-- Name: jug_file_base jug_file_base_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_base
    ADD CONSTRAINT jug_file_base_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: jug_file_base jug_file_base_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_base
    ADD CONSTRAINT jug_file_base_pkey PRIMARY KEY (id);


--
-- Name: jug_file_ilk jug_file_ilk_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_ilk
    ADD CONSTRAINT jug_file_ilk_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: jug_file_ilk jug_file_ilk_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_ilk
    ADD CONSTRAINT jug_file_ilk_pkey PRIMARY KEY (id);


--
-- Name: jug_file_vow jug_file_vow_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow
    ADD CONSTRAINT jug_file_vow_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: jug_file_vow jug_file_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow
    ADD CONSTRAINT jug_file_vow_pkey PRIMARY KEY (id);


--
-- Name: jug_ilk_duty jug_ilk_duty_block_number_block_hash_ilk_id_duty_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_duty
    ADD CONSTRAINT jug_ilk_duty_block_number_block_hash_ilk_id_duty_key UNIQUE (block_number, block_hash, ilk_id, duty);


--
-- Name: jug_ilk_duty jug_ilk_duty_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_duty
    ADD CONSTRAINT jug_ilk_duty_pkey PRIMARY KEY (id);


--
-- Name: jug_ilk_rho jug_ilk_rho_block_number_block_hash_ilk_id_rho_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_rho
    ADD CONSTRAINT jug_ilk_rho_block_number_block_hash_ilk_id_rho_key UNIQUE (block_number, block_hash, ilk_id, rho);


--
-- Name: jug_ilk_rho jug_ilk_rho_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_rho
    ADD CONSTRAINT jug_ilk_rho_pkey PRIMARY KEY (id);


--
-- Name: jug_init jug_init_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_init
    ADD CONSTRAINT jug_init_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: jug_init jug_init_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_init
    ADD CONSTRAINT jug_init_pkey PRIMARY KEY (id);


--
-- Name: jug_vat jug_vat_block_number_block_hash_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vat
    ADD CONSTRAINT jug_vat_block_number_block_hash_vat_key UNIQUE (block_number, block_hash, vat);


--
-- Name: jug_vat jug_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vat
    ADD CONSTRAINT jug_vat_pkey PRIMARY KEY (id);


--
-- Name: jug_vow jug_vow_block_number_block_hash_vow_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vow
    ADD CONSTRAINT jug_vow_block_number_block_hash_vow_key UNIQUE (block_number, block_hash, vow);


--
-- Name: jug_vow jug_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vow
    ADD CONSTRAINT jug_vow_pkey PRIMARY KEY (id);


--
-- Name: new_cdp new_cdp_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.new_cdp
    ADD CONSTRAINT new_cdp_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: new_cdp new_cdp_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.new_cdp
    ADD CONSTRAINT new_cdp_pkey PRIMARY KEY (id);


--
-- Name: spot_file_mat spot_file_mat_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_mat
    ADD CONSTRAINT spot_file_mat_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: spot_file_mat spot_file_mat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_mat
    ADD CONSTRAINT spot_file_mat_pkey PRIMARY KEY (id);


--
-- Name: spot_file_pip spot_file_pip_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_pip
    ADD CONSTRAINT spot_file_pip_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: spot_file_pip spot_file_pip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_pip
    ADD CONSTRAINT spot_file_pip_pkey PRIMARY KEY (id);


--
-- Name: spot_ilk_mat spot_ilk_mat_block_number_block_hash_ilk_id_mat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_mat
    ADD CONSTRAINT spot_ilk_mat_block_number_block_hash_ilk_id_mat_key UNIQUE (block_number, block_hash, ilk_id, mat);


--
-- Name: spot_ilk_mat spot_ilk_mat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_mat
    ADD CONSTRAINT spot_ilk_mat_pkey PRIMARY KEY (id);


--
-- Name: spot_ilk_pip spot_ilk_pip_block_number_block_hash_ilk_id_pip_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_pip
    ADD CONSTRAINT spot_ilk_pip_block_number_block_hash_ilk_id_pip_key UNIQUE (block_number, block_hash, ilk_id, pip);


--
-- Name: spot_ilk_pip spot_ilk_pip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_pip
    ADD CONSTRAINT spot_ilk_pip_pkey PRIMARY KEY (id);


--
-- Name: spot_par spot_par_block_number_block_hash_par_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_par
    ADD CONSTRAINT spot_par_block_number_block_hash_par_key UNIQUE (block_number, block_hash, par);


--
-- Name: spot_par spot_par_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_par
    ADD CONSTRAINT spot_par_pkey PRIMARY KEY (id);


--
-- Name: spot_poke spot_poke_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_poke
    ADD CONSTRAINT spot_poke_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: spot_poke spot_poke_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_poke
    ADD CONSTRAINT spot_poke_pkey PRIMARY KEY (id);


--
-- Name: spot_vat spot_vat_block_number_block_hash_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_vat
    ADD CONSTRAINT spot_vat_block_number_block_hash_vat_key UNIQUE (block_number, block_hash, vat);


--
-- Name: spot_vat spot_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_vat
    ADD CONSTRAINT spot_vat_pkey PRIMARY KEY (id);


--
-- Name: tend tend_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tend
    ADD CONSTRAINT tend_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: tend tend_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tend
    ADD CONSTRAINT tend_pkey PRIMARY KEY (id);


--
-- Name: tick tick_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tick
    ADD CONSTRAINT tick_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: tick tick_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tick
    ADD CONSTRAINT tick_pkey PRIMARY KEY (id);


--
-- Name: urns urns_ilk_id_identifier_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.urns
    ADD CONSTRAINT urns_ilk_id_identifier_key UNIQUE (ilk_id, identifier);


--
-- Name: urns urns_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.urns
    ADD CONSTRAINT urns_pkey PRIMARY KEY (id);


--
-- Name: vat_dai vat_dai_block_number_block_hash_guy_dai_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_dai
    ADD CONSTRAINT vat_dai_block_number_block_hash_guy_dai_key UNIQUE (block_number, block_hash, guy, dai);


--
-- Name: vat_dai vat_dai_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_dai
    ADD CONSTRAINT vat_dai_pkey PRIMARY KEY (id);


--
-- Name: vat_debt vat_debt_block_number_block_hash_debt_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_debt
    ADD CONSTRAINT vat_debt_block_number_block_hash_debt_key UNIQUE (block_number, block_hash, debt);


--
-- Name: vat_debt vat_debt_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_debt
    ADD CONSTRAINT vat_debt_pkey PRIMARY KEY (id);


--
-- Name: vat_file_debt_ceiling vat_file_debt_ceiling_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_debt_ceiling
    ADD CONSTRAINT vat_file_debt_ceiling_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_file_debt_ceiling vat_file_debt_ceiling_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_debt_ceiling
    ADD CONSTRAINT vat_file_debt_ceiling_pkey PRIMARY KEY (id);


--
-- Name: vat_file_ilk vat_file_ilk_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_ilk
    ADD CONSTRAINT vat_file_ilk_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_file_ilk vat_file_ilk_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_ilk
    ADD CONSTRAINT vat_file_ilk_pkey PRIMARY KEY (id);


--
-- Name: vat_flux vat_flux_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_flux
    ADD CONSTRAINT vat_flux_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_flux vat_flux_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_flux
    ADD CONSTRAINT vat_flux_pkey PRIMARY KEY (id);


--
-- Name: vat_fold vat_fold_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold
    ADD CONSTRAINT vat_fold_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_fold vat_fold_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold
    ADD CONSTRAINT vat_fold_pkey PRIMARY KEY (id);


--
-- Name: vat_fork vat_fork_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fork
    ADD CONSTRAINT vat_fork_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_fork vat_fork_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fork
    ADD CONSTRAINT vat_fork_pkey PRIMARY KEY (id);


--
-- Name: vat_frob vat_frob_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob
    ADD CONSTRAINT vat_frob_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_frob vat_frob_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob
    ADD CONSTRAINT vat_frob_pkey PRIMARY KEY (id);


--
-- Name: vat_gem vat_gem_block_number_block_hash_ilk_id_guy_gem_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_gem
    ADD CONSTRAINT vat_gem_block_number_block_hash_ilk_id_guy_gem_key UNIQUE (block_number, block_hash, ilk_id, guy, gem);


--
-- Name: vat_gem vat_gem_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_gem
    ADD CONSTRAINT vat_gem_pkey PRIMARY KEY (id);


--
-- Name: vat_grab vat_grab_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_grab
    ADD CONSTRAINT vat_grab_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_grab vat_grab_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_grab
    ADD CONSTRAINT vat_grab_pkey PRIMARY KEY (id);


--
-- Name: vat_heal vat_heal_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_heal
    ADD CONSTRAINT vat_heal_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_heal vat_heal_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_heal
    ADD CONSTRAINT vat_heal_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_art vat_ilk_art_block_number_block_hash_ilk_id_art_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_art
    ADD CONSTRAINT vat_ilk_art_block_number_block_hash_ilk_id_art_key UNIQUE (block_number, block_hash, ilk_id, art);


--
-- Name: vat_ilk_art vat_ilk_art_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_art
    ADD CONSTRAINT vat_ilk_art_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_dust vat_ilk_dust_block_number_block_hash_ilk_id_dust_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_dust
    ADD CONSTRAINT vat_ilk_dust_block_number_block_hash_ilk_id_dust_key UNIQUE (block_number, block_hash, ilk_id, dust);


--
-- Name: vat_ilk_dust vat_ilk_dust_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_dust
    ADD CONSTRAINT vat_ilk_dust_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_line vat_ilk_line_block_number_block_hash_ilk_id_line_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_line
    ADD CONSTRAINT vat_ilk_line_block_number_block_hash_ilk_id_line_key UNIQUE (block_number, block_hash, ilk_id, line);


--
-- Name: vat_ilk_line vat_ilk_line_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_line
    ADD CONSTRAINT vat_ilk_line_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_rate vat_ilk_rate_block_number_block_hash_ilk_id_rate_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_rate
    ADD CONSTRAINT vat_ilk_rate_block_number_block_hash_ilk_id_rate_key UNIQUE (block_number, block_hash, ilk_id, rate);


--
-- Name: vat_ilk_rate vat_ilk_rate_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_rate
    ADD CONSTRAINT vat_ilk_rate_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_spot vat_ilk_spot_block_number_block_hash_ilk_id_spot_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_spot
    ADD CONSTRAINT vat_ilk_spot_block_number_block_hash_ilk_id_spot_key UNIQUE (block_number, block_hash, ilk_id, spot);


--
-- Name: vat_ilk_spot vat_ilk_spot_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_spot
    ADD CONSTRAINT vat_ilk_spot_pkey PRIMARY KEY (id);


--
-- Name: vat_init vat_init_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_init
    ADD CONSTRAINT vat_init_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_init vat_init_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_init
    ADD CONSTRAINT vat_init_pkey PRIMARY KEY (id);


--
-- Name: vat_line vat_line_block_number_block_hash_line_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_line
    ADD CONSTRAINT vat_line_block_number_block_hash_line_key UNIQUE (block_number, block_hash, line);


--
-- Name: vat_line vat_line_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_line
    ADD CONSTRAINT vat_line_pkey PRIMARY KEY (id);


--
-- Name: vat_live vat_live_block_number_block_hash_live_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_live
    ADD CONSTRAINT vat_live_block_number_block_hash_live_key UNIQUE (block_number, block_hash, live);


--
-- Name: vat_live vat_live_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_live
    ADD CONSTRAINT vat_live_pkey PRIMARY KEY (id);


--
-- Name: vat_move vat_move_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_move
    ADD CONSTRAINT vat_move_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_move vat_move_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_move
    ADD CONSTRAINT vat_move_pkey PRIMARY KEY (id);


--
-- Name: vat_sin vat_sin_block_number_block_hash_guy_sin_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_sin
    ADD CONSTRAINT vat_sin_block_number_block_hash_guy_sin_key UNIQUE (block_number, block_hash, guy, sin);


--
-- Name: vat_sin vat_sin_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_sin
    ADD CONSTRAINT vat_sin_pkey PRIMARY KEY (id);


--
-- Name: vat_slip vat_slip_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_slip
    ADD CONSTRAINT vat_slip_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_slip vat_slip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_slip
    ADD CONSTRAINT vat_slip_pkey PRIMARY KEY (id);


--
-- Name: vat_suck vat_suck_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_suck
    ADD CONSTRAINT vat_suck_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_suck vat_suck_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_suck
    ADD CONSTRAINT vat_suck_pkey PRIMARY KEY (id);


--
-- Name: vat_urn_art vat_urn_art_block_number_block_hash_urn_id_art_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_art
    ADD CONSTRAINT vat_urn_art_block_number_block_hash_urn_id_art_key UNIQUE (block_number, block_hash, urn_id, art);


--
-- Name: vat_urn_art vat_urn_art_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_art
    ADD CONSTRAINT vat_urn_art_pkey PRIMARY KEY (id);


--
-- Name: vat_urn_ink vat_urn_ink_block_number_block_hash_urn_id_ink_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_ink
    ADD CONSTRAINT vat_urn_ink_block_number_block_hash_urn_id_ink_key UNIQUE (block_number, block_hash, urn_id, ink);


--
-- Name: vat_urn_ink vat_urn_ink_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_ink
    ADD CONSTRAINT vat_urn_ink_pkey PRIMARY KEY (id);


--
-- Name: vat_vice vat_vice_block_number_block_hash_vice_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_vice
    ADD CONSTRAINT vat_vice_block_number_block_hash_vice_key UNIQUE (block_number, block_hash, vice);


--
-- Name: vat_vice vat_vice_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_vice
    ADD CONSTRAINT vat_vice_pkey PRIMARY KEY (id);


--
-- Name: vow_ash vow_ash_block_number_block_hash_ash_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_ash
    ADD CONSTRAINT vow_ash_block_number_block_hash_ash_key UNIQUE (block_number, block_hash, ash);


--
-- Name: vow_ash vow_ash_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_ash
    ADD CONSTRAINT vow_ash_pkey PRIMARY KEY (id);


--
-- Name: vow_bump vow_bump_block_number_block_hash_bump_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_bump
    ADD CONSTRAINT vow_bump_block_number_block_hash_bump_key UNIQUE (block_number, block_hash, bump);


--
-- Name: vow_bump vow_bump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_bump
    ADD CONSTRAINT vow_bump_pkey PRIMARY KEY (id);


--
-- Name: vow_dump vow_dump_block_number_block_hash_dump_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_dump
    ADD CONSTRAINT vow_dump_block_number_block_hash_dump_key UNIQUE (block_number, block_hash, dump);


--
-- Name: vow_dump vow_dump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_dump
    ADD CONSTRAINT vow_dump_pkey PRIMARY KEY (id);


--
-- Name: vow_fess vow_fess_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_fess
    ADD CONSTRAINT vow_fess_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vow_fess vow_fess_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_fess
    ADD CONSTRAINT vow_fess_pkey PRIMARY KEY (id);


--
-- Name: vow_file vow_file_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file
    ADD CONSTRAINT vow_file_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vow_file vow_file_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file
    ADD CONSTRAINT vow_file_pkey PRIMARY KEY (id);


--
-- Name: vow_flapper vow_flapper_block_number_block_hash_flapper_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flapper
    ADD CONSTRAINT vow_flapper_block_number_block_hash_flapper_key UNIQUE (block_number, block_hash, flapper);


--
-- Name: vow_flapper vow_flapper_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flapper
    ADD CONSTRAINT vow_flapper_pkey PRIMARY KEY (id);


--
-- Name: vow_flog vow_flog_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flog
    ADD CONSTRAINT vow_flog_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vow_flog vow_flog_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flog
    ADD CONSTRAINT vow_flog_pkey PRIMARY KEY (id);


--
-- Name: vow_flopper vow_flopper_block_number_block_hash_flopper_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flopper
    ADD CONSTRAINT vow_flopper_block_number_block_hash_flopper_key UNIQUE (block_number, block_hash, flopper);


--
-- Name: vow_flopper vow_flopper_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flopper
    ADD CONSTRAINT vow_flopper_pkey PRIMARY KEY (id);


--
-- Name: vow_hump vow_hump_block_number_block_hash_hump_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_hump
    ADD CONSTRAINT vow_hump_block_number_block_hash_hump_key UNIQUE (block_number, block_hash, hump);


--
-- Name: vow_hump vow_hump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_hump
    ADD CONSTRAINT vow_hump_pkey PRIMARY KEY (id);


--
-- Name: vow_sin_integer vow_sin_integer_block_number_block_hash_sin_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_integer
    ADD CONSTRAINT vow_sin_integer_block_number_block_hash_sin_key UNIQUE (block_number, block_hash, sin);


--
-- Name: vow_sin_integer vow_sin_integer_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_integer
    ADD CONSTRAINT vow_sin_integer_pkey PRIMARY KEY (id);


--
-- Name: vow_sin_mapping vow_sin_mapping_block_number_block_hash_era_tab_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_mapping
    ADD CONSTRAINT vow_sin_mapping_block_number_block_hash_era_tab_key UNIQUE (block_number, block_hash, era, tab);


--
-- Name: vow_sin_mapping vow_sin_mapping_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_mapping
    ADD CONSTRAINT vow_sin_mapping_pkey PRIMARY KEY (id);


--
-- Name: vow_sump vow_sump_block_number_block_hash_sump_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sump
    ADD CONSTRAINT vow_sump_block_number_block_hash_sump_key UNIQUE (block_number, block_hash, sump);


--
-- Name: vow_sump vow_sump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sump
    ADD CONSTRAINT vow_sump_pkey PRIMARY KEY (id);


--
-- Name: vow_vat vow_vat_block_number_block_hash_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_vat
    ADD CONSTRAINT vow_vat_block_number_block_hash_vat_key UNIQUE (block_number, block_hash, vat);


--
-- Name: vow_vat vow_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_vat
    ADD CONSTRAINT vow_vat_pkey PRIMARY KEY (id);


--
-- Name: vow_wait vow_wait_block_number_block_hash_wait_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_wait
    ADD CONSTRAINT vow_wait_block_number_block_hash_wait_key UNIQUE (block_number, block_hash, wait);


--
-- Name: vow_wait vow_wait_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_wait
    ADD CONSTRAINT vow_wait_pkey PRIMARY KEY (id);


--
-- Name: yank yank_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.yank
    ADD CONSTRAINT yank_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: yank yank_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.yank
    ADD CONSTRAINT yank_pkey PRIMARY KEY (id);


--
-- Name: addresses addresses_address_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_address_key UNIQUE (address);


--
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- Name: blocks blocks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT blocks_pkey PRIMARY KEY (id);


--
-- Name: checked_headers checked_headers_header_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers
    ADD CONSTRAINT checked_headers_header_id_key UNIQUE (header_id);


--
-- Name: checked_headers checked_headers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers
    ADD CONSTRAINT checked_headers_pkey PRIMARY KEY (id);


--
-- Name: blocks eth_node_id_block_number_uc; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT eth_node_id_block_number_uc UNIQUE (number, eth_node_id);


--
-- Name: eth_nodes eth_node_uc; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes
    ADD CONSTRAINT eth_node_uc UNIQUE (genesis_block, network_id, eth_node_id);


--
-- Name: full_sync_logs full_sync_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_logs
    ADD CONSTRAINT full_sync_logs_pkey PRIMARY KEY (id);


--
-- Name: full_sync_receipts full_sync_receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_receipts
    ADD CONSTRAINT full_sync_receipts_pkey PRIMARY KEY (id);


--
-- Name: full_sync_transactions full_sync_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_transactions
    ADD CONSTRAINT full_sync_transactions_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: header_sync_logs header_sync_logs_header_id_tx_index_log_index_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_logs
    ADD CONSTRAINT header_sync_logs_header_id_tx_index_log_index_key UNIQUE (header_id, tx_index, log_index);


--
-- Name: header_sync_logs header_sync_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_logs
    ADD CONSTRAINT header_sync_logs_pkey PRIMARY KEY (id);


--
-- Name: header_sync_receipts header_sync_receipts_header_id_transaction_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_receipts
    ADD CONSTRAINT header_sync_receipts_header_id_transaction_id_key UNIQUE (header_id, transaction_id);


--
-- Name: header_sync_receipts header_sync_receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_receipts
    ADD CONSTRAINT header_sync_receipts_pkey PRIMARY KEY (id);


--
-- Name: header_sync_transactions header_sync_transactions_header_id_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_transactions
    ADD CONSTRAINT header_sync_transactions_header_id_hash_key UNIQUE (header_id, hash);


--
-- Name: header_sync_transactions header_sync_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_transactions
    ADD CONSTRAINT header_sync_transactions_pkey PRIMARY KEY (id);


--
-- Name: headers headers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_pkey PRIMARY KEY (id);


--
-- Name: log_filters name_uc; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.log_filters
    ADD CONSTRAINT name_uc UNIQUE (name);


--
-- Name: eth_nodes nodes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes
    ADD CONSTRAINT nodes_pkey PRIMARY KEY (id);


--
-- Name: queued_storage queued_storage_block_height_block_hash_contract_storage_key_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage
    ADD CONSTRAINT queued_storage_block_height_block_hash_contract_storage_key_key UNIQUE (block_height, block_hash, contract, storage_key, storage_value);


--
-- Name: queued_storage queued_storage_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage
    ADD CONSTRAINT queued_storage_pkey PRIMARY KEY (id);


--
-- Name: uncles uncles_block_id_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles
    ADD CONSTRAINT uncles_block_id_hash_key UNIQUE (block_id, hash);


--
-- Name: uncles uncles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles
    ADD CONSTRAINT uncles_pkey PRIMARY KEY (id);


--
-- Name: watched_contracts watched_contracts_contract_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_contracts
    ADD CONSTRAINT watched_contracts_contract_hash_key UNIQUE (contract_hash);


--
-- Name: watched_contracts watched_contracts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_contracts
    ADD CONSTRAINT watched_contracts_pkey PRIMARY KEY (contract_id);


--
-- Name: watched_logs watched_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_logs
    ADD CONSTRAINT watched_logs_pkey PRIMARY KEY (id);


--
-- Name: bite_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX bite_header_index ON maker.bite USING btree (header_id);


--
-- Name: bite_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX bite_urn_index ON maker.bite USING btree (urn_id);


--
-- Name: cat_file_chop_lump_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_chop_lump_header_index ON maker.cat_file_chop_lump USING btree (header_id);


--
-- Name: cat_file_chop_lump_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_chop_lump_ilk_index ON maker.cat_file_chop_lump USING btree (ilk_id);


--
-- Name: cat_file_flip_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_flip_header_index ON maker.cat_file_flip USING btree (header_id);


--
-- Name: cat_file_flip_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_flip_ilk_index ON maker.cat_file_flip USING btree (ilk_id);


--
-- Name: cat_file_vow_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_vow_header_index ON maker.cat_file_vow USING btree (header_id);


--
-- Name: cat_ilk_chop_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_ilk_chop_block_number_index ON maker.cat_ilk_chop USING btree (block_number);


--
-- Name: cat_ilk_chop_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_ilk_chop_ilk_index ON maker.cat_ilk_chop USING btree (ilk_id);


--
-- Name: cat_ilk_flip_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_ilk_flip_block_number_index ON maker.cat_ilk_flip USING btree (block_number);


--
-- Name: cat_ilk_flip_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_ilk_flip_ilk_index ON maker.cat_ilk_flip USING btree (ilk_id);


--
-- Name: cat_ilk_lump_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_ilk_lump_block_number_index ON maker.cat_ilk_lump USING btree (block_number);


--
-- Name: cat_ilk_lump_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_ilk_lump_ilk_index ON maker.cat_ilk_lump USING btree (ilk_id);


--
-- Name: cdp_manager_cdpi_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_cdpi_block_number_index ON maker.cdp_manager_cdpi USING btree (block_number);


--
-- Name: cdp_manager_cdpi_cdpi_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_cdpi_cdpi_index ON maker.cdp_manager_cdpi USING btree (cdpi);


--
-- Name: cdp_manager_ilks_cdpi_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_ilks_cdpi_index ON maker.cdp_manager_ilks USING btree (cdpi);


--
-- Name: cdp_manager_ilks_ilk_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_ilks_ilk_id_index ON maker.cdp_manager_ilks USING btree (ilk_id);


--
-- Name: cdp_manager_owns_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_owns_block_number_index ON maker.cdp_manager_owns USING btree (block_number);


--
-- Name: cdp_manager_owns_cdpi_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_owns_cdpi_index ON maker.cdp_manager_owns USING btree (cdpi);


--
-- Name: cdp_manager_owns_owner_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_owns_owner_index ON maker.cdp_manager_owns USING btree (owner);


--
-- Name: cdp_manager_urns_cdpi_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_urns_cdpi_index ON maker.cdp_manager_urns USING btree (cdpi);


--
-- Name: cdp_manager_urns_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_urns_urn_index ON maker.cdp_manager_urns USING btree (urn);


--
-- Name: deal_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deal_address_id_index ON maker.deal USING btree (address_id);


--
-- Name: deal_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deal_bid_id_index ON maker.deal USING btree (bid_id);


--
-- Name: deal_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deal_header_index ON maker.deal USING btree (header_id);


--
-- Name: dent_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX dent_header_index ON maker.dent USING btree (header_id);


--
-- Name: flap_bid_bid_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_bid_address_id_index ON maker.flap_bid_bid USING btree (address_id);


--
-- Name: flap_bid_bid_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_bid_bid_id_index ON maker.flap_bid_bid USING btree (bid_id);


--
-- Name: flap_bid_bid_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_bid_block_number_index ON maker.flap_bid_bid USING btree (block_number);


--
-- Name: flap_bid_end_bid_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_end_bid_address_id_index ON maker.flap_bid_end USING btree (address_id);


--
-- Name: flap_bid_end_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_end_bid_id_index ON maker.flap_bid_end USING btree (bid_id);


--
-- Name: flap_bid_end_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_end_block_number_index ON maker.flap_bid_end USING btree (block_number);


--
-- Name: flap_bid_guy_bid_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_guy_bid_address_id_index ON maker.flap_bid_guy USING btree (address_id);


--
-- Name: flap_bid_guy_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_guy_bid_id_index ON maker.flap_bid_guy USING btree (bid_id);


--
-- Name: flap_bid_guy_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_guy_block_number_index ON maker.flap_bid_guy USING btree (block_number);


--
-- Name: flap_bid_lot_bid_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_lot_bid_address_id_index ON maker.flap_bid_lot USING btree (address_id);


--
-- Name: flap_bid_lot_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_lot_bid_id_index ON maker.flap_bid_lot USING btree (bid_id);


--
-- Name: flap_bid_lot_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_lot_block_number_index ON maker.flap_bid_lot USING btree (block_number);


--
-- Name: flap_bid_tic_bid_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_tic_bid_address_id_index ON maker.flap_bid_tic USING btree (address_id);


--
-- Name: flap_bid_tic_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_tic_bid_id_index ON maker.flap_bid_tic USING btree (bid_id);


--
-- Name: flap_bid_tic_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_tic_block_number_index ON maker.flap_bid_tic USING btree (block_number);


--
-- Name: flap_kick_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_kick_header_index ON maker.flap_kick USING btree (header_id);


--
-- Name: flap_kicks_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_kicks_address_id_index ON maker.flap_kicks USING btree (address_id);


--
-- Name: flap_kicks_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_kicks_block_number_index ON maker.flap_kicks USING btree (block_number);


--
-- Name: flap_kicks_kicks_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_kicks_kicks_index ON maker.flap_kicks USING btree (kicks);


--
-- Name: flip_bid_bid_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_bid_address_id_index ON maker.flip_bid_bid USING btree (address_id);


--
-- Name: flip_bid_bid_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_bid_bid_id_index ON maker.flip_bid_bid USING btree (bid_id);


--
-- Name: flip_bid_bid_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_bid_block_number_index ON maker.flip_bid_bid USING btree (block_number);


--
-- Name: flip_bid_end_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_end_address_id_index ON maker.flip_bid_end USING btree (address_id);


--
-- Name: flip_bid_end_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_end_bid_id_index ON maker.flip_bid_end USING btree (bid_id);


--
-- Name: flip_bid_end_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_end_block_number_index ON maker.flip_bid_end USING btree (block_number);


--
-- Name: flip_bid_gal_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_gal_address_id_index ON maker.flip_bid_gal USING btree (address_id);


--
-- Name: flip_bid_gal_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_gal_bid_id_index ON maker.flip_bid_gal USING btree (bid_id);


--
-- Name: flip_bid_gal_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_gal_block_number_index ON maker.flip_bid_gal USING btree (block_number);


--
-- Name: flip_bid_guy_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_guy_address_id_index ON maker.flip_bid_guy USING btree (address_id);


--
-- Name: flip_bid_guy_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_guy_bid_id_index ON maker.flip_bid_guy USING btree (bid_id);


--
-- Name: flip_bid_guy_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_guy_block_number_index ON maker.flip_bid_guy USING btree (block_number);


--
-- Name: flip_bid_lot_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_lot_address_id_index ON maker.flip_bid_lot USING btree (address_id);


--
-- Name: flip_bid_lot_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_lot_bid_id_index ON maker.flip_bid_lot USING btree (bid_id);


--
-- Name: flip_bid_lot_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_lot_block_number_index ON maker.flip_bid_lot USING btree (block_number);


--
-- Name: flip_bid_tab_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_tab_address_id_index ON maker.flip_bid_tab USING btree (address_id);


--
-- Name: flip_bid_tab_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_tab_bid_id_index ON maker.flip_bid_tab USING btree (bid_id);


--
-- Name: flip_bid_tab_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_tab_block_number_index ON maker.flip_bid_tab USING btree (block_number);


--
-- Name: flip_bid_tic_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_tic_address_id_index ON maker.flip_bid_tic USING btree (address_id);


--
-- Name: flip_bid_tic_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_tic_bid_id_index ON maker.flip_bid_tic USING btree (bid_id);


--
-- Name: flip_bid_tic_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_tic_block_number_index ON maker.flip_bid_tic USING btree (block_number);


--
-- Name: flip_bid_usr_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_usr_address_id_index ON maker.flip_bid_usr USING btree (address_id);


--
-- Name: flip_bid_usr_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_usr_bid_id_index ON maker.flip_bid_usr USING btree (bid_id);


--
-- Name: flip_bid_usr_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_usr_block_number_index ON maker.flip_bid_usr USING btree (block_number);


--
-- Name: flip_ilk_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_ilk_block_number_index ON maker.flip_ilk USING btree (block_number);


--
-- Name: flip_ilk_ilk_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_ilk_ilk_id_index ON maker.flip_ilk USING btree (ilk_id);


--
-- Name: flip_kick_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_kick_address_id_index ON maker.flip_kick USING btree (address_id);


--
-- Name: flip_kick_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_kick_bid_id_index ON maker.flip_kick USING btree (bid_id);


--
-- Name: flip_kick_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_kick_header_index ON maker.flip_kick USING btree (header_id);


--
-- Name: flip_kicks_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_kicks_address_id_index ON maker.flip_kicks USING btree (address_id);


--
-- Name: flip_kicks_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_kicks_block_number_index ON maker.flip_kicks USING btree (block_number);


--
-- Name: flip_kicks_kicks_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_kicks_kicks_index ON maker.flip_kicks USING btree (kicks);


--
-- Name: flop_bid_bid_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_bid_address_id_index ON maker.flop_bid_bid USING btree (address_id);


--
-- Name: flop_bid_bid_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_bid_bid_id_index ON maker.flop_bid_bid USING btree (bid_id);


--
-- Name: flop_bid_bid_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_bid_block_number_index ON maker.flop_bid_bid USING btree (block_number);


--
-- Name: flop_bid_end_bid_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_end_bid_address_id_index ON maker.flop_bid_end USING btree (address_id);


--
-- Name: flop_bid_end_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_end_bid_id_index ON maker.flop_bid_end USING btree (bid_id);


--
-- Name: flop_bid_end_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_end_block_number_index ON maker.flop_bid_end USING btree (block_number);


--
-- Name: flop_bid_guy_bid_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_guy_bid_address_id_index ON maker.flop_bid_guy USING btree (address_id);


--
-- Name: flop_bid_guy_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_guy_bid_id_index ON maker.flop_bid_guy USING btree (bid_id);


--
-- Name: flop_bid_guy_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_guy_block_number_index ON maker.flop_bid_guy USING btree (block_number);


--
-- Name: flop_bid_lot_bid_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_lot_bid_address_id_index ON maker.flop_bid_lot USING btree (address_id);


--
-- Name: flop_bid_lot_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_lot_bid_id_index ON maker.flop_bid_lot USING btree (bid_id);


--
-- Name: flop_bid_lot_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_lot_block_number_index ON maker.flop_bid_lot USING btree (block_number);


--
-- Name: flop_bid_tic_bid_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_tic_bid_address_id_index ON maker.flop_bid_tic USING btree (address_id);


--
-- Name: flop_bid_tic_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_tic_bid_id_index ON maker.flop_bid_tic USING btree (bid_id);


--
-- Name: flop_bid_tic_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_tic_block_number_index ON maker.flop_bid_tic USING btree (block_number);


--
-- Name: flop_kick_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_kick_header_index ON maker.flop_kick USING btree (header_id);


--
-- Name: flop_kicks_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_kicks_address_id_index ON maker.flop_kicks USING btree (address_id);


--
-- Name: flop_kicks_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_kicks_block_number_index ON maker.flop_kicks USING btree (block_number);


--
-- Name: flop_kicks_kicks_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_kicks_kicks_index ON maker.flop_kicks USING btree (kicks);


--
-- Name: jug_drip_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_drip_header_index ON maker.jug_drip USING btree (header_id);


--
-- Name: jug_drip_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_drip_ilk_index ON maker.jug_drip USING btree (ilk_id);


--
-- Name: jug_file_base_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_base_header_index ON maker.jug_file_base USING btree (header_id);


--
-- Name: jug_file_ilk_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_ilk_header_index ON maker.jug_file_ilk USING btree (header_id);


--
-- Name: jug_file_ilk_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_ilk_ilk_index ON maker.jug_file_ilk USING btree (ilk_id);


--
-- Name: jug_file_vow_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_vow_header_index ON maker.jug_file_vow USING btree (header_id);


--
-- Name: jug_ilk_duty_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_ilk_duty_block_number_index ON maker.jug_ilk_duty USING btree (block_number);


--
-- Name: jug_ilk_duty_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_ilk_duty_ilk_index ON maker.jug_ilk_duty USING btree (ilk_id);


--
-- Name: jug_ilk_rho_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_ilk_rho_block_number_index ON maker.jug_ilk_rho USING btree (block_number);


--
-- Name: jug_ilk_rho_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_ilk_rho_ilk_index ON maker.jug_ilk_rho USING btree (ilk_id);


--
-- Name: jug_init_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_init_header_index ON maker.jug_init USING btree (header_id);


--
-- Name: jug_init_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_init_ilk_index ON maker.jug_init USING btree (ilk_id);


--
-- Name: spot_file_mat_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_mat_header_index ON maker.spot_file_mat USING btree (header_id);


--
-- Name: spot_file_mat_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_mat_ilk_index ON maker.spot_file_mat USING btree (ilk_id);


--
-- Name: spot_file_pip_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_pip_header_index ON maker.spot_file_pip USING btree (header_id);


--
-- Name: spot_file_pip_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_pip_ilk_index ON maker.spot_file_pip USING btree (ilk_id);


--
-- Name: spot_ilk_mat_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_ilk_mat_block_number_index ON maker.spot_ilk_mat USING btree (block_number);


--
-- Name: spot_ilk_mat_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_ilk_mat_ilk_index ON maker.spot_ilk_mat USING btree (ilk_id);


--
-- Name: spot_ilk_pip_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_ilk_pip_block_number_index ON maker.spot_ilk_pip USING btree (block_number);


--
-- Name: spot_ilk_pip_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_ilk_pip_ilk_index ON maker.spot_ilk_pip USING btree (ilk_id);


--
-- Name: spot_poke_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_poke_header_index ON maker.spot_poke USING btree (header_id);


--
-- Name: spot_poke_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_poke_ilk_index ON maker.spot_poke USING btree (ilk_id);


--
-- Name: tend_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX tend_header_index ON maker.tend USING btree (header_id);


--
-- Name: tick_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX tick_bid_id_index ON maker.tick USING btree (bid_id);


--
-- Name: tick_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX tick_header_index ON maker.tick USING btree (header_id);


--
-- Name: urn_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX urn_ilk_index ON maker.urns USING btree (ilk_id);


--
-- Name: vat_file_debt_ceiling_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_file_debt_ceiling_header_index ON maker.vat_file_debt_ceiling USING btree (header_id);


--
-- Name: vat_file_ilk_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_file_ilk_header_index ON maker.vat_file_ilk USING btree (header_id);


--
-- Name: vat_file_ilk_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_file_ilk_ilk_index ON maker.vat_file_ilk USING btree (ilk_id);


--
-- Name: vat_flux_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_flux_header_index ON maker.vat_flux USING btree (header_id);


--
-- Name: vat_flux_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_flux_ilk_index ON maker.vat_flux USING btree (ilk_id);


--
-- Name: vat_fold_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_fold_header_index ON maker.vat_fold USING btree (header_id);


--
-- Name: vat_fold_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_fold_urn_index ON maker.vat_fold USING btree (urn_id);


--
-- Name: vat_fork_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_fork_header_index ON maker.vat_fork USING btree (header_id);


--
-- Name: vat_fork_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_fork_ilk_index ON maker.vat_fork USING btree (ilk_id);


--
-- Name: vat_frob_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_frob_header_index ON maker.vat_frob USING btree (header_id);


--
-- Name: vat_frob_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_frob_urn_index ON maker.vat_frob USING btree (urn_id);


--
-- Name: vat_gem_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_gem_ilk_index ON maker.vat_gem USING btree (ilk_id);


--
-- Name: vat_grab_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_grab_header_index ON maker.vat_grab USING btree (header_id);


--
-- Name: vat_grab_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_grab_urn_index ON maker.vat_grab USING btree (urn_id);


--
-- Name: vat_heal_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_heal_header_index ON maker.vat_heal USING btree (header_id);


--
-- Name: vat_ilk_art_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_art_block_number_index ON maker.vat_ilk_art USING btree (block_number);


--
-- Name: vat_ilk_art_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_art_ilk_index ON maker.vat_ilk_art USING btree (ilk_id);


--
-- Name: vat_ilk_dust_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_dust_block_number_index ON maker.vat_ilk_dust USING btree (block_number);


--
-- Name: vat_ilk_dust_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_dust_ilk_index ON maker.vat_ilk_dust USING btree (ilk_id);


--
-- Name: vat_ilk_line_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_line_block_number_index ON maker.vat_ilk_line USING btree (block_number);


--
-- Name: vat_ilk_line_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_line_ilk_index ON maker.vat_ilk_line USING btree (ilk_id);


--
-- Name: vat_ilk_rate_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_rate_block_number_index ON maker.vat_ilk_rate USING btree (block_number);


--
-- Name: vat_ilk_rate_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_rate_ilk_index ON maker.vat_ilk_rate USING btree (ilk_id);


--
-- Name: vat_ilk_spot_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_spot_block_number_index ON maker.vat_ilk_spot USING btree (block_number);


--
-- Name: vat_ilk_spot_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_spot_ilk_index ON maker.vat_ilk_spot USING btree (ilk_id);


--
-- Name: vat_init_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_init_header_index ON maker.vat_init USING btree (header_id);


--
-- Name: vat_init_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_init_ilk_index ON maker.vat_init USING btree (ilk_id);


--
-- Name: vat_move_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_move_header_index ON maker.vat_move USING btree (header_id);


--
-- Name: vat_slip_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_slip_header_index ON maker.vat_slip USING btree (header_id);


--
-- Name: vat_slip_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_slip_ilk_index ON maker.vat_slip USING btree (ilk_id);


--
-- Name: vat_suck_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_suck_header_index ON maker.vat_suck USING btree (header_id);


--
-- Name: vat_urn_art_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_urn_art_block_number_index ON maker.vat_urn_art USING btree (block_number);


--
-- Name: vat_urn_art_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_urn_art_urn_index ON maker.vat_urn_art USING btree (urn_id);


--
-- Name: vat_urn_ink_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_urn_ink_block_number_index ON maker.vat_urn_ink USING btree (block_number);


--
-- Name: vat_urn_ink_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_urn_ink_urn_index ON maker.vat_urn_ink USING btree (urn_id);


--
-- Name: vow_fess_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_fess_header_index ON maker.vow_fess USING btree (header_id);


--
-- Name: vow_file_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_file_header_index ON maker.vow_file USING btree (header_id);


--
-- Name: vow_flog_era_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_flog_era_index ON maker.vow_flog USING btree (era);


--
-- Name: vow_flog_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_flog_header_index ON maker.vow_flog USING btree (header_id);


--
-- Name: vow_sin_mapping_block_number_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_sin_mapping_block_number_index ON maker.vow_sin_mapping USING btree (block_number);


--
-- Name: vow_sin_mapping_era_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_sin_mapping_era_index ON maker.vow_sin_mapping USING btree (era);


--
-- Name: yank_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX yank_bid_id_index ON maker.yank USING btree (bid_id);


--
-- Name: yank_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX yank_header_index ON maker.yank USING btree (header_id);


--
-- Name: block_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX block_id_index ON public.full_sync_transactions USING btree (block_id);


--
-- Name: headers_block_number; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX headers_block_number ON public.headers USING btree (block_number);


--
-- Name: node_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX node_id_index ON public.blocks USING btree (eth_node_id);


--
-- Name: number_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX number_index ON public.blocks USING btree (number);


--
-- Name: tx_from_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tx_from_index ON public.full_sync_transactions USING btree (tx_from);


--
-- Name: tx_to_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tx_to_index ON public.full_sync_transactions USING btree (tx_to);


--
-- Name: flap_bid_bid flap_bid_bid; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_bid_bid AFTER INSERT OR UPDATE ON maker.flap_bid_bid FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flap_bid();


--
-- Name: flap_bid_end flap_bid_end; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_bid_end AFTER INSERT OR UPDATE ON maker.flap_bid_end FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flap_end();


--
-- Name: flap_bid_guy flap_bid_guy; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_bid_guy AFTER INSERT OR UPDATE ON maker.flap_bid_guy FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flap_guy();


--
-- Name: flap_bid_lot flap_bid_lot; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_bid_lot AFTER INSERT OR UPDATE ON maker.flap_bid_lot FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flap_lot();


--
-- Name: flap_bid_tic flap_bid_tic; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_bid_tic AFTER INSERT OR UPDATE ON maker.flap_bid_tic FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flap_tic();


--
-- Name: flap_kick flap_created_trigger; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_created_trigger AFTER INSERT ON maker.flap_kick FOR EACH ROW EXECUTE PROCEDURE maker.flap_created();


--
-- Name: flip_bid_bid flip_bid_bid; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_bid_bid AFTER INSERT OR UPDATE ON maker.flip_bid_bid FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flip_bid();


--
-- Name: flip_bid_end flip_bid_end; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_bid_end AFTER INSERT OR UPDATE ON maker.flip_bid_end FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flip_end();


--
-- Name: flip_bid_gal flip_bid_gal; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_bid_gal AFTER INSERT OR UPDATE ON maker.flip_bid_gal FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flip_gal();


--
-- Name: flip_bid_guy flip_bid_guy; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_bid_guy AFTER INSERT OR UPDATE ON maker.flip_bid_guy FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flip_guy();


--
-- Name: flip_bid_lot flip_bid_lot; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_bid_lot AFTER INSERT OR UPDATE ON maker.flip_bid_lot FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flip_lot();


--
-- Name: flip_bid_tab flip_bid_tab; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_bid_tab AFTER INSERT OR UPDATE ON maker.flip_bid_tab FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flip_tab();


--
-- Name: flip_bid_tic flip_bid_tic; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_bid_tic AFTER INSERT OR UPDATE ON maker.flip_bid_tic FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flip_tic();


--
-- Name: flip_kick flip_created_trigger; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_created_trigger AFTER INSERT ON maker.flip_kick FOR EACH ROW EXECUTE PROCEDURE maker.flip_created();


--
-- Name: flop_bid_bid flop_bid_bid; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_bid_bid AFTER INSERT OR UPDATE ON maker.flop_bid_bid FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flop_bid();


--
-- Name: flop_bid_end flop_bid_end; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_bid_end AFTER INSERT OR UPDATE ON maker.flop_bid_end FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flop_end();


--
-- Name: flop_bid_guy flop_bid_guy; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_bid_guy AFTER INSERT OR UPDATE ON maker.flop_bid_guy FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flop_guy();


--
-- Name: flop_bid_lot flop_bid_lot; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_bid_lot AFTER INSERT OR UPDATE ON maker.flop_bid_lot FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flop_lot();


--
-- Name: flop_bid_tic flop_bid_tic; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_bid_tic AFTER INSERT OR UPDATE ON maker.flop_bid_tic FOR EACH ROW EXECUTE PROCEDURE maker.insert_updated_flop_tic();


--
-- Name: flop_kick flop_created_trigger; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_created_trigger AFTER INSERT ON maker.flop_kick FOR EACH ROW EXECUTE PROCEDURE maker.flop_created();


--
-- Name: vat_ilk_art ilk_art; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_art AFTER INSERT OR UPDATE ON maker.vat_ilk_art FOR EACH ROW EXECUTE PROCEDURE maker.insert_ilk_art();


--
-- Name: cat_ilk_chop ilk_chop; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_chop AFTER INSERT OR UPDATE ON maker.cat_ilk_chop FOR EACH ROW EXECUTE PROCEDURE maker.insert_ilk_chop();


--
-- Name: vat_ilk_dust ilk_dust; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_dust AFTER INSERT OR UPDATE ON maker.vat_ilk_dust FOR EACH ROW EXECUTE PROCEDURE maker.insert_ilk_dust();


--
-- Name: jug_ilk_duty ilk_duty; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_duty AFTER INSERT OR UPDATE ON maker.jug_ilk_duty FOR EACH ROW EXECUTE PROCEDURE maker.insert_ilk_duty();


--
-- Name: cat_ilk_flip ilk_flip; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_flip AFTER INSERT OR UPDATE ON maker.cat_ilk_flip FOR EACH ROW EXECUTE PROCEDURE maker.insert_ilk_flip();


--
-- Name: vat_ilk_line ilk_line; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_line AFTER INSERT OR UPDATE ON maker.vat_ilk_line FOR EACH ROW EXECUTE PROCEDURE maker.insert_ilk_line();


--
-- Name: cat_ilk_lump ilk_lump; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_lump AFTER INSERT OR UPDATE ON maker.cat_ilk_lump FOR EACH ROW EXECUTE PROCEDURE maker.insert_ilk_lump();


--
-- Name: spot_ilk_mat ilk_mat; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_mat AFTER INSERT OR UPDATE ON maker.spot_ilk_mat FOR EACH ROW EXECUTE PROCEDURE maker.insert_ilk_mat();


--
-- Name: spot_ilk_pip ilk_pip; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_pip AFTER INSERT OR UPDATE ON maker.spot_ilk_pip FOR EACH ROW EXECUTE PROCEDURE maker.insert_ilk_pip();


--
-- Name: vat_ilk_rate ilk_rate; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_rate AFTER INSERT OR UPDATE ON maker.vat_ilk_rate FOR EACH ROW EXECUTE PROCEDURE maker.insert_ilk_rate();


--
-- Name: jug_ilk_rho ilk_rho; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_rho AFTER INSERT OR UPDATE ON maker.jug_ilk_rho FOR EACH ROW EXECUTE PROCEDURE maker.insert_ilk_rho();


--
-- Name: vat_ilk_spot ilk_spot; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_spot AFTER INSERT OR UPDATE ON maker.vat_ilk_spot FOR EACH ROW EXECUTE PROCEDURE maker.insert_ilk_spot();


--
-- Name: cdp_manager_cdpi managed_cdp_cdpi; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER managed_cdp_cdpi AFTER INSERT OR UPDATE ON maker.cdp_manager_cdpi FOR EACH ROW EXECUTE PROCEDURE maker.insert_cdp_created();


--
-- Name: cdp_manager_ilks managed_cdp_ilk; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER managed_cdp_ilk AFTER INSERT OR UPDATE ON maker.cdp_manager_ilks FOR EACH ROW EXECUTE PROCEDURE maker.insert_cdp_ilk_identifier();


--
-- Name: cdp_manager_urns managed_cdp_urn; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER managed_cdp_urn AFTER INSERT OR UPDATE ON maker.cdp_manager_urns FOR EACH ROW EXECUTE PROCEDURE maker.insert_cdp_urn_identifier();


--
-- Name: cdp_manager_owns managed_cdp_usr; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER managed_cdp_usr AFTER INSERT OR UPDATE ON maker.cdp_manager_owns FOR EACH ROW EXECUTE PROCEDURE maker.insert_cdp_usr();


--
-- Name: bite bite_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite
    ADD CONSTRAINT bite_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: bite bite_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite
    ADD CONSTRAINT bite_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: bite bite_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite
    ADD CONSTRAINT bite_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id) ON DELETE CASCADE;


--
-- Name: cat_file_chop_lump cat_file_chop_lump_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_chop_lump
    ADD CONSTRAINT cat_file_chop_lump_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cat_file_chop_lump cat_file_chop_lump_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_chop_lump
    ADD CONSTRAINT cat_file_chop_lump_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: cat_file_chop_lump cat_file_chop_lump_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_chop_lump
    ADD CONSTRAINT cat_file_chop_lump_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: cat_file_flip cat_file_flip_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_flip
    ADD CONSTRAINT cat_file_flip_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cat_file_flip cat_file_flip_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_flip
    ADD CONSTRAINT cat_file_flip_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: cat_file_flip cat_file_flip_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_flip
    ADD CONSTRAINT cat_file_flip_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: cat_file_vow cat_file_vow_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_vow
    ADD CONSTRAINT cat_file_vow_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cat_file_vow cat_file_vow_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_vow
    ADD CONSTRAINT cat_file_vow_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: cat_ilk_chop cat_ilk_chop_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_chop
    ADD CONSTRAINT cat_ilk_chop_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: cat_ilk_flip cat_ilk_flip_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_flip
    ADD CONSTRAINT cat_ilk_flip_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: cat_ilk_lump cat_ilk_lump_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_lump
    ADD CONSTRAINT cat_ilk_lump_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_ilks cdp_manager_ilks_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_ilks
    ADD CONSTRAINT cdp_manager_ilks_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: deal deal_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deal
    ADD CONSTRAINT deal_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: deal deal_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deal
    ADD CONSTRAINT deal_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: deal deal_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deal
    ADD CONSTRAINT deal_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: dent dent_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.dent
    ADD CONSTRAINT dent_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: dent dent_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.dent
    ADD CONSTRAINT dent_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: dent dent_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.dent
    ADD CONSTRAINT dent_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: flap flap_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap
    ADD CONSTRAINT flap_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_beg flap_beg_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_beg
    ADD CONSTRAINT flap_beg_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_bid_bid flap_bid_bid_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_bid
    ADD CONSTRAINT flap_bid_bid_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_bid_end flap_bid_end_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_end
    ADD CONSTRAINT flap_bid_end_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_bid_guy flap_bid_guy_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_guy
    ADD CONSTRAINT flap_bid_guy_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_bid_lot flap_bid_lot_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_lot
    ADD CONSTRAINT flap_bid_lot_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_bid_tic flap_bid_tic_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_tic
    ADD CONSTRAINT flap_bid_tic_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_gem flap_gem_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_gem
    ADD CONSTRAINT flap_gem_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_kick flap_kick_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kick
    ADD CONSTRAINT flap_kick_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_kick flap_kick_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kick
    ADD CONSTRAINT flap_kick_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flap_kick flap_kick_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kick
    ADD CONSTRAINT flap_kick_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: flap_kicks flap_kicks_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kicks
    ADD CONSTRAINT flap_kicks_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_live flap_live_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_live
    ADD CONSTRAINT flap_live_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_tau flap_tau_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_tau
    ADD CONSTRAINT flap_tau_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_ttl flap_ttl_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_ttl
    ADD CONSTRAINT flap_ttl_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_vat flap_vat_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_vat
    ADD CONSTRAINT flap_vat_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip flip_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip
    ADD CONSTRAINT flip_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_beg flip_beg_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_beg
    ADD CONSTRAINT flip_beg_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_bid flip_bid_bid_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_bid
    ADD CONSTRAINT flip_bid_bid_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_end flip_bid_end_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_end
    ADD CONSTRAINT flip_bid_end_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_gal flip_bid_gal_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_gal
    ADD CONSTRAINT flip_bid_gal_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_guy flip_bid_guy_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_guy
    ADD CONSTRAINT flip_bid_guy_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_lot flip_bid_lot_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_lot
    ADD CONSTRAINT flip_bid_lot_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_tab flip_bid_tab_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tab
    ADD CONSTRAINT flip_bid_tab_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_tic flip_bid_tic_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tic
    ADD CONSTRAINT flip_bid_tic_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_usr flip_bid_usr_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_usr
    ADD CONSTRAINT flip_bid_usr_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_ilk flip_ilk_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ilk
    ADD CONSTRAINT flip_ilk_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_ilk flip_ilk_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ilk
    ADD CONSTRAINT flip_ilk_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: flip_kick flip_kick_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kick
    ADD CONSTRAINT flip_kick_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_kick flip_kick_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kick
    ADD CONSTRAINT flip_kick_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_kick flip_kick_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kick
    ADD CONSTRAINT flip_kick_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: flip_kicks flip_kicks_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kicks
    ADD CONSTRAINT flip_kicks_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_tau flip_tau_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_tau
    ADD CONSTRAINT flip_tau_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_ttl flip_ttl_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ttl
    ADD CONSTRAINT flip_ttl_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_vat flip_vat_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_vat
    ADD CONSTRAINT flip_vat_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop flop_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop
    ADD CONSTRAINT flop_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_beg flop_beg_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_beg
    ADD CONSTRAINT flop_beg_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_bid_bid flop_bid_bid_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_bid
    ADD CONSTRAINT flop_bid_bid_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_bid_end flop_bid_end_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_end
    ADD CONSTRAINT flop_bid_end_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_bid_guy flop_bid_guy_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_guy
    ADD CONSTRAINT flop_bid_guy_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_bid_lot flop_bid_lot_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_lot
    ADD CONSTRAINT flop_bid_lot_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_bid_tic flop_bid_tic_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_tic
    ADD CONSTRAINT flop_bid_tic_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_gem flop_gem_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_gem
    ADD CONSTRAINT flop_gem_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_kick flop_kick_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kick
    ADD CONSTRAINT flop_kick_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_kick flop_kick_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kick
    ADD CONSTRAINT flop_kick_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_kick flop_kick_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kick
    ADD CONSTRAINT flop_kick_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: flop_kicks flop_kicks_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kicks
    ADD CONSTRAINT flop_kicks_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_live flop_live_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_live
    ADD CONSTRAINT flop_live_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_pad flop_pad_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_pad
    ADD CONSTRAINT flop_pad_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_tau flop_tau_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_tau
    ADD CONSTRAINT flop_tau_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_ttl flop_ttl_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_ttl
    ADD CONSTRAINT flop_ttl_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_vat flop_vat_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vat
    ADD CONSTRAINT flop_vat_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: jug_drip jug_drip_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_drip
    ADD CONSTRAINT jug_drip_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_drip jug_drip_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_drip
    ADD CONSTRAINT jug_drip_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: jug_drip jug_drip_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_drip
    ADD CONSTRAINT jug_drip_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: jug_file_base jug_file_base_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_base
    ADD CONSTRAINT jug_file_base_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_file_base jug_file_base_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_base
    ADD CONSTRAINT jug_file_base_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: jug_file_ilk jug_file_ilk_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_ilk
    ADD CONSTRAINT jug_file_ilk_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_file_ilk jug_file_ilk_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_ilk
    ADD CONSTRAINT jug_file_ilk_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: jug_file_ilk jug_file_ilk_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_ilk
    ADD CONSTRAINT jug_file_ilk_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: jug_file_vow jug_file_vow_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow
    ADD CONSTRAINT jug_file_vow_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_file_vow jug_file_vow_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow
    ADD CONSTRAINT jug_file_vow_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: jug_ilk_duty jug_ilk_duty_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_duty
    ADD CONSTRAINT jug_ilk_duty_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: jug_ilk_rho jug_ilk_rho_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_rho
    ADD CONSTRAINT jug_ilk_rho_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: jug_init jug_init_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_init
    ADD CONSTRAINT jug_init_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_init jug_init_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_init
    ADD CONSTRAINT jug_init_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: jug_init jug_init_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_init
    ADD CONSTRAINT jug_init_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: new_cdp new_cdp_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.new_cdp
    ADD CONSTRAINT new_cdp_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: new_cdp new_cdp_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.new_cdp
    ADD CONSTRAINT new_cdp_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: spot_file_mat spot_file_mat_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_mat
    ADD CONSTRAINT spot_file_mat_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: spot_file_mat spot_file_mat_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_mat
    ADD CONSTRAINT spot_file_mat_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: spot_file_mat spot_file_mat_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_mat
    ADD CONSTRAINT spot_file_mat_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: spot_file_pip spot_file_pip_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_pip
    ADD CONSTRAINT spot_file_pip_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: spot_file_pip spot_file_pip_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_pip
    ADD CONSTRAINT spot_file_pip_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: spot_file_pip spot_file_pip_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_pip
    ADD CONSTRAINT spot_file_pip_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: spot_ilk_mat spot_ilk_mat_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_mat
    ADD CONSTRAINT spot_ilk_mat_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: spot_ilk_pip spot_ilk_pip_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_pip
    ADD CONSTRAINT spot_ilk_pip_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: spot_poke spot_poke_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_poke
    ADD CONSTRAINT spot_poke_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: spot_poke spot_poke_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_poke
    ADD CONSTRAINT spot_poke_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: spot_poke spot_poke_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_poke
    ADD CONSTRAINT spot_poke_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: tend tend_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tend
    ADD CONSTRAINT tend_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: tend tend_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tend
    ADD CONSTRAINT tend_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: tend tend_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tend
    ADD CONSTRAINT tend_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: tick tick_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tick
    ADD CONSTRAINT tick_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: tick tick_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tick
    ADD CONSTRAINT tick_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: tick tick_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tick
    ADD CONSTRAINT tick_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: urns urns_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.urns
    ADD CONSTRAINT urns_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_file_debt_ceiling vat_file_debt_ceiling_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_debt_ceiling
    ADD CONSTRAINT vat_file_debt_ceiling_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_file_debt_ceiling vat_file_debt_ceiling_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_debt_ceiling
    ADD CONSTRAINT vat_file_debt_ceiling_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vat_file_ilk vat_file_ilk_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_ilk
    ADD CONSTRAINT vat_file_ilk_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_file_ilk vat_file_ilk_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_ilk
    ADD CONSTRAINT vat_file_ilk_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_file_ilk vat_file_ilk_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_ilk
    ADD CONSTRAINT vat_file_ilk_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vat_flux vat_flux_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_flux
    ADD CONSTRAINT vat_flux_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_flux vat_flux_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_flux
    ADD CONSTRAINT vat_flux_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_flux vat_flux_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_flux
    ADD CONSTRAINT vat_flux_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vat_fold vat_fold_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold
    ADD CONSTRAINT vat_fold_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_fold vat_fold_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold
    ADD CONSTRAINT vat_fold_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vat_fold vat_fold_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold
    ADD CONSTRAINT vat_fold_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id) ON DELETE CASCADE;


--
-- Name: vat_fork vat_fork_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fork
    ADD CONSTRAINT vat_fork_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_fork vat_fork_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fork
    ADD CONSTRAINT vat_fork_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_fork vat_fork_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fork
    ADD CONSTRAINT vat_fork_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vat_frob vat_frob_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob
    ADD CONSTRAINT vat_frob_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_frob vat_frob_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob
    ADD CONSTRAINT vat_frob_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vat_frob vat_frob_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob
    ADD CONSTRAINT vat_frob_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id) ON DELETE CASCADE;


--
-- Name: vat_gem vat_gem_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_gem
    ADD CONSTRAINT vat_gem_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_grab vat_grab_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_grab
    ADD CONSTRAINT vat_grab_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_grab vat_grab_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_grab
    ADD CONSTRAINT vat_grab_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vat_grab vat_grab_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_grab
    ADD CONSTRAINT vat_grab_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id) ON DELETE CASCADE;


--
-- Name: vat_heal vat_heal_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_heal
    ADD CONSTRAINT vat_heal_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_heal vat_heal_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_heal
    ADD CONSTRAINT vat_heal_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_art vat_ilk_art_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_art
    ADD CONSTRAINT vat_ilk_art_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_dust vat_ilk_dust_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_dust
    ADD CONSTRAINT vat_ilk_dust_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_line vat_ilk_line_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_line
    ADD CONSTRAINT vat_ilk_line_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_rate vat_ilk_rate_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_rate
    ADD CONSTRAINT vat_ilk_rate_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_spot vat_ilk_spot_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_spot
    ADD CONSTRAINT vat_ilk_spot_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_init vat_init_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_init
    ADD CONSTRAINT vat_init_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_init vat_init_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_init
    ADD CONSTRAINT vat_init_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_init vat_init_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_init
    ADD CONSTRAINT vat_init_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vat_move vat_move_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_move
    ADD CONSTRAINT vat_move_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_move vat_move_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_move
    ADD CONSTRAINT vat_move_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vat_slip vat_slip_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_slip
    ADD CONSTRAINT vat_slip_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_slip vat_slip_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_slip
    ADD CONSTRAINT vat_slip_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_slip vat_slip_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_slip
    ADD CONSTRAINT vat_slip_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vat_suck vat_suck_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_suck
    ADD CONSTRAINT vat_suck_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_suck vat_suck_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_suck
    ADD CONSTRAINT vat_suck_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vat_urn_art vat_urn_art_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_art
    ADD CONSTRAINT vat_urn_art_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id) ON DELETE CASCADE;


--
-- Name: vat_urn_ink vat_urn_ink_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_ink
    ADD CONSTRAINT vat_urn_ink_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id) ON DELETE CASCADE;


--
-- Name: vow_fess vow_fess_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_fess
    ADD CONSTRAINT vow_fess_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_fess vow_fess_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_fess
    ADD CONSTRAINT vow_fess_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vow_file vow_file_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file
    ADD CONSTRAINT vow_file_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_file vow_file_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file
    ADD CONSTRAINT vow_file_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: vow_flog vow_flog_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flog
    ADD CONSTRAINT vow_flog_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_flog vow_flog_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flog
    ADD CONSTRAINT vow_flog_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: yank yank_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.yank
    ADD CONSTRAINT yank_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: yank yank_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.yank
    ADD CONSTRAINT yank_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: yank yank_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.yank
    ADD CONSTRAINT yank_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.header_sync_logs(id) ON DELETE CASCADE;


--
-- Name: full_sync_receipts blocks_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_receipts
    ADD CONSTRAINT blocks_fk FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: checked_headers checked_headers_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers
    ADD CONSTRAINT checked_headers_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: full_sync_receipts full_sync_receipts_contract_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_receipts
    ADD CONSTRAINT full_sync_receipts_contract_address_id_fkey FOREIGN KEY (contract_address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: full_sync_transactions full_sync_transactions_block_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_transactions
    ADD CONSTRAINT full_sync_transactions_block_id_fkey FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: header_sync_logs header_sync_logs_address_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_logs
    ADD CONSTRAINT header_sync_logs_address_fkey FOREIGN KEY (address) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: header_sync_logs header_sync_logs_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_logs
    ADD CONSTRAINT header_sync_logs_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: header_sync_receipts header_sync_receipts_contract_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_receipts
    ADD CONSTRAINT header_sync_receipts_contract_address_id_fkey FOREIGN KEY (contract_address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: header_sync_receipts header_sync_receipts_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_receipts
    ADD CONSTRAINT header_sync_receipts_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: header_sync_receipts header_sync_receipts_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_receipts
    ADD CONSTRAINT header_sync_receipts_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES public.header_sync_transactions(id) ON DELETE CASCADE;


--
-- Name: header_sync_transactions header_sync_transactions_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_transactions
    ADD CONSTRAINT header_sync_transactions_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: headers headers_eth_node_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_eth_node_id_fkey FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: blocks node_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT node_fk FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: full_sync_logs receipts_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_logs
    ADD CONSTRAINT receipts_fk FOREIGN KEY (receipt_id) REFERENCES public.full_sync_receipts(id) ON DELETE CASCADE;


--
-- Name: uncles uncles_block_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles
    ADD CONSTRAINT uncles_block_id_fkey FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: uncles uncles_eth_node_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles
    ADD CONSTRAINT uncles_eth_node_id_fkey FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

