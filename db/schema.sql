--
-- PostgreSQL database dump
--

-- Dumped from database version 11.6
-- Dumped by pg_dump version 11.6

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
-- Name: flap_bid_snapshot; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.flap_bid_snapshot AS (
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
-- Name: flip_bid_snapshot; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.flip_bid_snapshot AS (
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
-- Name: flop_bid_snapshot; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.flop_bid_snapshot AS (
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
                                AND flap_bid_bid.header_id = headers.id
                  LEFT JOIN maker.flap_bid_lot
                            ON deal.bid_id = flap_bid_lot.bid_id
                                AND flap_bid_lot.header_id = headers.id
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
                                AND flap_bid_bid.header_id = headers.id
                  LEFT JOIN maker.flap_bid_lot
                            ON yank.bid_id = flap_bid_lot.bid_id
                                AND flap_bid_lot.header_id = headers.id
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
                                AND flap_bid_bid.header_id = headers.id
                  LEFT JOIN maker.flap_bid_lot
                            ON tick.bid_id = flap_bid_lot.bid_id
                                AND flap_bid_lot.header_id = headers.id
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

CREATE FUNCTION api.all_flaps(max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flap_bid_snapshot
    LANGUAGE plpgsql STABLE
    AS $$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT bid_id
            FROM maker.flap
            ORDER BY bid_id DESC
            LIMIT all_flaps.max_results
            OFFSET
            all_flaps.result_offset
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
                                AND flip_bid_bid.header_id = headers.id
                  LEFT JOIN maker.flip_bid_lot
                            ON deal.bid_id = flip_bid_lot.bid_id
                                AND flip_bid_lot.header_id = headers.id
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
                                AND flip_bid_bid.header_id = headers.id
                  LEFT JOIN maker.flip_bid_lot
                            ON yank.bid_id = flip_bid_lot.bid_id
                                AND flip_bid_lot.header_id = headers.id
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
                                AND flip_bid_bid.header_id = headers.id
                  LEFT JOIN maker.flip_bid_lot
                            ON tick.bid_id = flip_bid_lot.bid_id
                                AND flip_bid_lot.header_id = headers.id
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
LIMIT all_flip_bid_events.max_results
OFFSET
all_flip_bid_events.result_offset
$$;


--
-- Name: all_flips(text, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_flips(ilk text, max_results integer DEFAULT '-1'::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flip_bid_snapshot
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
                                AND flop_bid_bid.header_id = headers.id
                  LEFT JOIN maker.flop_bid_lot
                            ON deal.bid_id = flop_bid_lot.bid_id
                                AND flop_bid_lot.header_id = headers.id
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
                                AND flop_bid_bid.header_id = headers.id
                  LEFT JOIN maker.flop_bid_lot
                            ON yank.bid_id = flop_bid_lot.bid_id
                                AND flop_bid_lot.header_id = headers.id
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
                                AND flop_bid_bid.header_id = headers.id
                  LEFT JOIN maker.flop_bid_lot
                            ON tick.bid_id = flop_bid_lot.bid_id
                                AND flop_bid_lot.header_id = headers.id
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
LIMIT all_flop_bid_events.max_results
OFFSET
all_flop_bid_events.result_offset
$$;


--
-- Name: all_flops(integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_flops(max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flop_bid_snapshot
    LANGUAGE plpgsql STABLE
    AS $$
BEGIN
    RETURN QUERY (
        WITH bid_ids AS (
            SELECT DISTINCT bid_id
            FROM maker.flop
            ORDER BY bid_id DESC
            LIMIT all_flops.max_results
            OFFSET
            all_flops.result_offset
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
                        LEFT JOIN public.headers ON vat_ilk_rate.header_id = headers.id
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


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: ilk_snapshot; Type: TABLE; Schema: api; Owner: -
--

CREATE TABLE api.ilk_snapshot (
    ilk_identifier text NOT NULL,
    block_number bigint NOT NULL,
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
-- Name: all_ilks(bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_ilks(block_height bigint DEFAULT api.max_block()) RETURNS SETOF api.ilk_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT DISTINCT ON (ilk_identifier) *
FROM api.ilk_snapshot
WHERE block_number <= block_height
ORDER BY ilk_identifier, block_number DESC
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
-- Name: FUNCTION max_timestamp(); Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON FUNCTION api.max_timestamp() IS '@omit';


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
LIMIT all_poke_events.max_results
OFFSET
all_poke_events.result_offset
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
            LIMIT all_queued_sin.max_results
            OFFSET
            all_queued_sin.result_offset
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
-- Name: urn_snapshot; Type: TABLE; Schema: api; Owner: -
--

CREATE TABLE api.urn_snapshot (
    urn_identifier text NOT NULL,
    ilk_identifier text NOT NULL,
    block_height bigint NOT NULL,
    ink numeric,
    art numeric,
    created timestamp without time zone,
    updated timestamp without time zone NOT NULL
);


--
-- Name: all_urns(bigint, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_urns(block_height bigint DEFAULT api.max_block(), max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.urn_snapshot
    LANGUAGE sql STABLE
    AS $$
WITH distinct_urn_snapshots AS (SELECT urn_identifier, ilk_identifier, MAX(block_height) AS block_height
                                FROM api.urn_snapshot
                                WHERE block_height <= all_urns.block_height
                                GROUP BY urn_identifier, ilk_identifier)
SELECT us.urn_identifier, us.ilk_identifier, us.block_height, us.ink, coalesce(us.art, 0), us.created, us.updated
    FROM api.urn_snapshot AS us, distinct_urn_snapshots AS dus
    WHERE us.urn_identifier = dus.urn_identifier
    AND us.ilk_identifier = dus.ilk_identifier
    AND us.block_height = dus.block_height
    LIMIT all_urns.max_results
    OFFSET all_urns.result_offset
$$;


--
-- Name: bite_event_bid(api.bite_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bite_event_bid(event api.bite_event) RETURNS api.flip_bid_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_flip(event.bid_id, event.ilk_identifier, event.block_height)
$$;


--
-- Name: bite_event_ilk(api.bite_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bite_event_ilk(event api.bite_event) RETURNS api.ilk_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.ilk_snapshot i
WHERE i.ilk_identifier = event.ilk_identifier
  AND i.block_number <= event.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$;


--
-- Name: bite_event_tx(api.bite_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bite_event_tx(event api.bite_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM get_tx_data(event.block_height, event.log_id)
$$;


--
-- Name: bite_event_urn(api.bite_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bite_event_urn(event api.bite_event) RETURNS api.urn_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn(event.ilk_identifier, event.urn_identifier, event.block_height)
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
-- Name: flap_bid_event_bid(api.flap_bid_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flap_bid_event_bid(event api.flap_bid_event) RETURNS api.flap_bid_snapshot
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
-- Name: flap_bid_snapshot_bid_events(api.flap_bid_snapshot, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flap_bid_snapshot_bid_events(flap api.flap_bid_snapshot, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flap_bid_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_flap_bid_events() bids
WHERE bid_id = flap.bid_id
ORDER BY bids.block_height DESC
LIMIT flap_bid_snapshot_bid_events.max_results
OFFSET
flap_bid_snapshot_bid_events.result_offset
$$;


--
-- Name: flip_bid_event_bid(api.flip_bid_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flip_bid_event_bid(event api.flip_bid_event) RETURNS api.flip_bid_snapshot
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
-- Name: flip_bid_snapshot_bid_events(api.flip_bid_snapshot, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flip_bid_snapshot_bid_events(flip api.flip_bid_snapshot, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flip_bid_event
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
LIMIT flip_bid_snapshot_bid_events.max_results
OFFSET
flip_bid_snapshot_bid_events.result_offset
$$;


--
-- Name: flip_bid_snapshot_ilk(api.flip_bid_snapshot); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flip_bid_snapshot_ilk(flip_bid_snapshot api.flip_bid_snapshot) RETURNS api.ilk_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT i.*
FROM api.ilk_snapshot i
         LEFT JOIN maker.ilks ON ilks.identifier = i.ilk_identifier
WHERE ilks.id = flip_bid_snapshot.ilk_id
  AND i.block_number <= flip_bid_snapshot.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$;


--
-- Name: flip_bid_snapshot_urn(api.flip_bid_snapshot); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flip_bid_snapshot_urn(flip api.flip_bid_snapshot) RETURNS api.urn_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn(
     (SELECT identifier FROM maker.ilks WHERE ilks.id = flip.ilk_id),
     (SELECT identifier FROM maker.urns WHERE urns.id = flip.urn_id),
     flip.block_height)
$$;


--
-- Name: flop_bid_event_bid(api.flop_bid_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flop_bid_event_bid(event api.flop_bid_event) RETURNS api.flop_bid_snapshot
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
-- Name: flop_bid_snapshot_bid_events(api.flop_bid_snapshot, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.flop_bid_snapshot_bid_events(flop api.flop_bid_snapshot, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.flop_bid_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_flop_bid_events() bids
WHERE bid_id = flop.bid_id
ORDER BY bids.block_height DESC
LIMIT flop_bid_snapshot_bid_events.max_results
OFFSET
flop_bid_snapshot_bid_events.result_offset
$$;


--
-- Name: frob_event_ilk(api.frob_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.frob_event_ilk(event api.frob_event) RETURNS api.ilk_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.ilk_snapshot i
WHERE i.ilk_identifier = event.ilk_identifier
  AND i.block_number <= event.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$;


--
-- Name: frob_event_tx(api.frob_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.frob_event_tx(event api.frob_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM get_tx_data(event.block_height, event.log_id)
$$;


--
-- Name: frob_event_urn(api.frob_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.frob_event_urn(event api.frob_event) RETURNS api.urn_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn(event.ilk_identifier, event.urn_identifier, event.block_height)
$$;


--
-- Name: get_flap(numeric, bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_flap(bid_id numeric, block_height bigint DEFAULT api.max_block()) RETURNS api.flap_bid_snapshot
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

CREATE FUNCTION api.get_flip(bid_id numeric, ilk text, block_height bigint DEFAULT api.max_block()) RETURNS api.flip_bid_snapshot
    LANGUAGE sql STABLE STRICT
    AS $$
WITH ilk_ids AS (SELECT id FROM maker.ilks WHERE ilks.identifier = get_flip.ilk),
     -- there should only ever be 1 address for a given ilk, which is why there's a LIMIT with no ORDER BY
     address_id AS (SELECT address_id
                    FROM maker.flip_ilk
                             LEFT JOIN public.headers ON flip_ilk.header_id = headers.id
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
         WHERE flip.bid_id = get_flip.bid_id
           AND flip.address_id = (SELECT address_id FROM address_id)
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

CREATE FUNCTION api.get_flop(bid_id numeric, block_height bigint DEFAULT api.max_block()) RETURNS api.flop_bid_snapshot
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
-- Name: get_queued_sin(numeric); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_queued_sin(era numeric) RETURNS api.queued_sin
    LANGUAGE sql STABLE STRICT
    AS $$
WITH created AS (SELECT era, h.block_number, api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM maker.vow_sin_mapping
                          LEFT JOIN public.headers h ON h.id = vow_sin_mapping.header_id
                 WHERE era = get_queued_sin.era
                 ORDER BY h.block_number ASC
                 LIMIT 1),
     updated AS (SELECT era, h.block_number, api.epoch_to_datetime(block_timestamp) AS datetime
                 FROM maker.vow_sin_mapping
                          LEFT JOIN public.headers h ON h.id = vow_sin_mapping.header_id
                 WHERE era = get_queued_sin.era
                 ORDER BY h.block_number DESC
                 LIMIT 1)

SELECT get_queued_sin.era,
       tab,
       (SELECT EXISTS(SELECT id FROM maker.vow_flog WHERE vow_flog.era = get_queued_sin.era)) AS flogged,
       created.datetime,
       updated.datetime
FROM maker.vow_sin_mapping
         LEFT JOIN created ON created.era = vow_sin_mapping.era
         LEFT JOIN updated ON updated.era = vow_sin_mapping.era
         LEFT JOIN public.headers ON headers.id = vow_sin_mapping.header_id
WHERE vow_sin_mapping.era = get_queued_sin.era
ORDER BY headers.block_number DESC
$$;


--
-- Name: get_urn(text, text, bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_urn(ilk_identifier text, urn_identifier text, block_height bigint DEFAULT api.max_block()) RETURNS api.urn_snapshot
    LANGUAGE sql STABLE STRICT
    AS $$

SELECT urn_identifier, ilk_identifier, get_urn.block_height, ink, art, created, updated
    FROM api.urn_snapshot
    WHERE ilk_identifier = get_urn.ilk_identifier
    AND urn_identifier = get_urn.urn_identifier
    AND block_height <= get_urn.block_height
    ORDER BY updated DESC
    LIMIT 1
$$;


--
-- Name: ilk_file_event_ilk(api.ilk_file_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_file_event_ilk(event api.ilk_file_event) RETURNS SETOF api.ilk_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.ilk_snapshot i
WHERE i.ilk_identifier = event.ilk_identifier
  AND i.block_number <= event.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$;


--
-- Name: ilk_file_event_tx(api.ilk_file_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_file_event_tx(event api.ilk_file_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM get_tx_data(event.block_height, event.log_id)
$$;


--
-- Name: ilk_snapshot_bites(api.ilk_snapshot, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_snapshot_bites(state api.ilk_snapshot, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.bite_event
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
-- Name: ilk_snapshot_frobs(api.ilk_snapshot, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_snapshot_frobs(state api.ilk_snapshot, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.frob_event
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
-- Name: ilk_snapshot_ilk_file_events(api.ilk_snapshot, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_snapshot_ilk_file_events(state api.ilk_snapshot, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.ilk_file_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_ilk_file_events(state.ilk_identifier)
LIMIT max_results
OFFSET
result_offset
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
-- Name: managed_cdp_ilk(api.managed_cdp); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.managed_cdp_ilk(cdp api.managed_cdp) RETURNS api.ilk_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.ilk_snapshot i
WHERE i.ilk_identifier = cdp.ilk_identifier
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: managed_cdp_urn(api.managed_cdp); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.managed_cdp_urn(cdp api.managed_cdp) RETURNS api.urn_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn(cdp.ilk_identifier, cdp.urn_identifier, api.max_block())
$$;


--
-- Name: poke_event_ilk(api.poke_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.poke_event_ilk(event api.poke_event) RETURNS api.ilk_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT i.*
FROM api.ilk_snapshot i
         LEFT JOIN maker.ilks ON ilks.identifier = i.ilk_identifier
WHERE ilks.id = event.ilk_id
  AND i.block_number <= event.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$;


--
-- Name: poke_event_tx(api.poke_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.poke_event_tx(priceupdate api.poke_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM get_tx_data(priceUpdate.block_height, priceUpdate.log_id)
$$;


--
-- Name: queued_sin_sin_queue_events(api.queued_sin, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.queued_sin_sin_queue_events(state api.queued_sin, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.sin_queue_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.all_sin_queue_events(state.era)
LIMIT queued_sin_sin_queue_events.max_results
OFFSET
queued_sin_sin_queue_events.result_offset
$$;


--
-- Name: sin_queue_event_tx(api.sin_queue_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.sin_queue_event_tx(event api.sin_queue_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM get_tx_data(event.block_height, event.log_id)
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
               LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
      WHERE ilks.identifier = total_ink.ilk_identifier
        AND headers.block_number <= total_ink.block_height
      ORDER BY vat_urn_ink.urn_id, headers.block_number DESC) latest_ink_by_urn
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
-- Name: bid_event; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.bid_event (
    log_id bigint NOT NULL,
    bid_id numeric NOT NULL,
    contract_address text NOT NULL,
    act api.bid_act NOT NULL,
    lot numeric,
    bid_amount numeric,
    ilk_identifier text,
    urn_identifier text,
    block_height bigint NOT NULL
);


--
-- Name: urn_bid_events(text, text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.urn_bid_events(urn_identifier text, ilk_identifier text) RETURNS SETOF maker.bid_event
    LANGUAGE sql STABLE STRICT
    AS $$
SELECT *
FROM maker.bid_event
WHERE bid_event.ilk_identifier = urn_bid_events.ilk_identifier
  AND bid_event.urn_identifier = urn_bid_events.urn_identifier
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
                        LEFT JOIN public.headers ON vat_ilk_rate.header_id = headers.id
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
-- Name: urn_snapshot_bites(api.urn_snapshot, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.urn_snapshot_bites(state api.urn_snapshot, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.bite_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.urn_bites(state.ilk_identifier, state.urn_identifier)
WHERE block_height <= state.block_height
LIMIT urn_snapshot_bites.max_results
OFFSET
urn_snapshot_bites.result_offset
$$;


--
-- Name: urn_snapshot_frobs(api.urn_snapshot, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.urn_snapshot_frobs(state api.urn_snapshot, max_results integer DEFAULT NULL::integer, result_offset integer DEFAULT 0) RETURNS SETOF api.frob_event
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.urn_frobs(state.ilk_identifier, state.urn_identifier)
WHERE block_height <= state.block_height
LIMIT urn_snapshot_frobs.max_results
OFFSET
urn_snapshot_frobs.result_offset
$$;


--
-- Name: urn_snapshot_ilk(api.urn_snapshot); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.urn_snapshot_ilk(state api.urn_snapshot) RETURNS api.ilk_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.ilk_snapshot i
WHERE i.ilk_identifier = state.ilk_identifier
  AND i.block_number <= state.block_height
ORDER BY i.block_number DESC
LIMIT 1
$$;


--
-- Name: flip_ilk; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_ilk (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    ilk_id integer NOT NULL
);


--
-- Name: clear_bid_event_ilk(maker.flip_ilk); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.clear_bid_event_ilk(old_diff maker.flip_ilk) RETURNS void
    LANGUAGE sql
    AS $$
UPDATE maker.bid_event
SET ilk_identifier = NULL
WHERE bid_event.contract_address = (SELECT address FROM public.addresses WHERE id = old_diff.address_id)
$$;


--
-- Name: FUNCTION clear_bid_event_ilk(old_diff maker.flip_ilk); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.clear_bid_event_ilk(old_diff maker.flip_ilk) IS '@omit';


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
-- Name: clear_flap_created(maker.flap_kick); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.clear_flap_created(old_event maker.flap_kick) RETURNS void
    LANGUAGE sql
    AS $$
UPDATE maker.flap
SET created = flap_bid_time_created(old_event.address_id, old_event.bid_id)
WHERE flap.address_id = old_event.address_id
  AND flap.bid_id = old_event.bid_id
$$;


--
-- Name: FUNCTION clear_flap_created(old_event maker.flap_kick); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.clear_flap_created(old_event maker.flap_kick) IS '@omit';


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
-- Name: clear_flip_created(maker.flip_kick); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.clear_flip_created(old_event maker.flip_kick) RETURNS void
    LANGUAGE sql
    AS $$
UPDATE maker.flip
SET created = flip_bid_time_created(old_event.address_id, old_event.bid_id)
WHERE flip.address_id = old_event.address_id
  AND flip.bid_id = old_event.bid_id
$$;


--
-- Name: FUNCTION clear_flip_created(old_event maker.flip_kick); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.clear_flip_created(old_event maker.flip_kick) IS '@omit';


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
-- Name: clear_flop_created(maker.flop_kick); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.clear_flop_created(old_event maker.flop_kick) RETURNS void
    LANGUAGE sql
    AS $$
UPDATE maker.flop
SET created = flop_bid_time_created(old_event.address_id, old_event.bid_id)
WHERE flop.address_id = old_event.address_id
  AND flop.bid_id = old_event.bid_id
$$;


--
-- Name: FUNCTION clear_flop_created(old_event maker.flop_kick); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.clear_flop_created(old_event maker.flop_kick) IS '@omit';


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
-- Name: clear_time_created(maker.vat_init); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.clear_time_created(old_event maker.vat_init) RETURNS maker.vat_init
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE api.ilk_snapshot
    SET created = ilk_time_created(old_event.ilk_id)
    FROM maker.ilks
    WHERE ilks.identifier = ilk_snapshot.ilk_identifier
      AND ilks.id = old_event.ilk_id;
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION clear_time_created(old_event maker.vat_init); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.clear_time_created(old_event maker.vat_init) IS '@omit';


--
-- Name: delete_obsolete_flap(numeric, integer, integer); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.delete_obsolete_flap(bid_id numeric, address_id integer, header_id integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    flap_block      BIGINT     := (
        SELECT block_number
        FROM public.headers
        WHERE id = header_id);
    flap_state      maker.flap := (
        SELECT (flap.address_id, block_number, flap.bid_id, guy, tic, "end", lot, bid, created, updated)
        FROM maker.flap
        WHERE flap.bid_id = delete_obsolete_flap.bid_id
          AND flap.address_id = delete_obsolete_flap.address_id
          AND flap.block_number = flap_block);
    prev_flap_state maker.flap := (
        SELECT (flap.address_id, block_number, flap.bid_id, guy, tic, "end", lot, bid, created, updated)
        FROM maker.flap
        WHERE flap.bid_id = delete_obsolete_flap.bid_id
          AND flap.address_id = delete_obsolete_flap.address_id
          AND flap.block_number < flap_block
        ORDER BY flap.block_number DESC
        LIMIT 1);
BEGIN
    DELETE
    FROM maker.flap
    WHERE flap.bid_id = delete_obsolete_flap.bid_id
      AND flap.address_id = delete_obsolete_flap.address_id
      AND flap.block_number = flap_block
      AND (flap_state.guy IS NULL OR flap_state.guy = prev_flap_state.guy)
      AND (flap_state.tic IS NULL OR flap_state.tic = prev_flap_state.tic)
      AND (flap_state."end" IS NULL OR flap_state."end" = prev_flap_state."end")
      AND (flap_state.lot IS NULL OR flap_state.lot = prev_flap_state.lot)
      AND (flap_state.bid IS NULL OR flap_state.bid = prev_flap_state.bid);
END
$$;


--
-- Name: FUNCTION delete_obsolete_flap(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.delete_obsolete_flap(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: delete_obsolete_flip(numeric, integer, integer); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.delete_obsolete_flip(bid_id numeric, address_id integer, header_id integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    flip_block      BIGINT     := (
        SELECT block_number
        FROM public.headers
        WHERE id = header_id);
    flip_state      maker.flip := (
        SELECT (flip.address_id, block_number, flip.bid_id, guy, tic, "end", lot, bid, usr, gal, tab, created, updated)
        FROM maker.flip
        WHERE flip.bid_id = delete_obsolete_flip.bid_id
          AND flip.address_id = delete_obsolete_flip.address_id
          AND flip.block_number = flip_block);
    prev_flip_state maker.flip := (
        SELECT (flip.address_id, block_number, flip.bid_id, guy, tic, "end", lot, bid, usr, gal, tab, created, updated)
        FROM maker.flip
        WHERE flip.bid_id = delete_obsolete_flip.bid_id
          AND flip.address_id = delete_obsolete_flip.address_id
          AND flip.block_number < flip_block
        ORDER BY flip.block_number DESC
        LIMIT 1);
BEGIN
    DELETE
    FROM maker.flip
    WHERE flip.bid_id = delete_obsolete_flip.bid_id
      AND flip.address_id = delete_obsolete_flip.address_id
      AND flip.block_number = flip_block
      AND (flip_state.guy IS NULL OR flip_state.guy = prev_flip_state.guy)
      AND (flip_state.tic IS NULL OR flip_state.tic = prev_flip_state.tic)
      AND (flip_state."end" IS NULL OR flip_state."end" = prev_flip_state."end")
      AND (flip_state.lot IS NULL OR flip_state.lot = prev_flip_state.lot)
      AND (flip_state.bid IS NULL OR flip_state.bid = prev_flip_state.bid)
      AND (flip_state.usr IS NULL OR flip_state.usr = prev_flip_state.usr)
      AND (flip_state.gal IS NULL OR flip_state.gal = prev_flip_state.gal)
      AND (flip_state.tab IS NULL OR flip_state.tab = prev_flip_state.tab);
END
$$;


--
-- Name: FUNCTION delete_obsolete_flip(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.delete_obsolete_flip(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: delete_obsolete_flop(numeric, integer, integer); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.delete_obsolete_flop(bid_id numeric, address_id integer, header_id integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    flop_block      BIGINT     := (
        SELECT block_number
        FROM public.headers
        WHERE id = header_id);
    flop_state      maker.flop := (
        SELECT (flop.address_id, block_number, flop.bid_id, guy, tic, "end", lot, bid, created, updated)
        FROM maker.flop
        WHERE flop.bid_id = delete_obsolete_flop.bid_id
          AND flop.address_id = delete_obsolete_flop.address_id
          AND flop.block_number = flop_block);
    prev_flop_state maker.flop := (
        SELECT (flop.address_id, block_number, flop.bid_id, guy, tic, "end", lot, bid, created, updated)
        FROM maker.flop
        WHERE flop.bid_id = delete_obsolete_flop.bid_id
          AND flop.address_id = delete_obsolete_flop.address_id
          AND flop.block_number < flop_block
        ORDER BY flop.block_number DESC
        LIMIT 1);
BEGIN
    DELETE
    FROM maker.flop
    WHERE flop.bid_id = delete_obsolete_flop.bid_id
      AND flop.address_id = delete_obsolete_flop.address_id
      AND flop.block_number = flop_block
      AND (flop_state.guy IS NULL OR flop_state.guy = prev_flop_state.guy)
      AND (flop_state.tic IS NULL OR flop_state.tic = prev_flop_state.tic)
      AND (flop_state."end" IS NULL OR flop_state."end" = prev_flop_state."end")
      AND (flop_state.lot IS NULL OR flop_state.lot = prev_flop_state.lot)
      AND (flop_state.bid IS NULL OR flop_state.bid = prev_flop_state.bid);
END
$$;


--
-- Name: FUNCTION delete_obsolete_flop(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.delete_obsolete_flop(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: delete_obsolete_urn_snapshot(integer, integer); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.delete_obsolete_urn_snapshot(urn_id integer, header_id integer) RETURNS api.urn_snapshot
    LANGUAGE plpgsql
    AS $$
DECLARE
    urn_snapshot_block_number BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = header_id);
BEGIN
    DELETE
    FROM api.urn_snapshot
        USING maker.urns LEFT JOIN maker.ilks ON urns.ilk_id = ilks.id
    WHERE urn_snapshot.urn_identifier = urns.identifier
      AND urn_snapshot.ilk_identifier = ilks.identifier
      AND urns.id = urn_id
      AND urn_snapshot.block_height = urn_snapshot_block_number
      AND NOT (EXISTS(
            SELECT *
            FROM maker.vat_urn_ink
            WHERE vat_urn_ink.urn_id = delete_obsolete_urn_snapshot.urn_id
              AND vat_urn_ink.header_id = delete_obsolete_urn_snapshot.header_id))
      AND NOT (EXISTS(
            SELECT *
            FROM maker.vat_urn_art
            WHERE vat_urn_art.urn_id = delete_obsolete_urn_snapshot.urn_id
              AND vat_urn_art.header_id = delete_obsolete_urn_snapshot.header_id));
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION delete_obsolete_urn_snapshot(urn_id integer, header_id integer); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.delete_obsolete_urn_snapshot(urn_id integer, header_id integer) IS '@omit';


--
-- Name: delete_redundant_ilk_snapshot(integer, integer); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.delete_redundant_ilk_snapshot(ilk_id integer, header_id integer) RETURNS api.ilk_snapshot
    LANGUAGE plpgsql
    AS $$
DECLARE
    associated_ilk        TEXT             := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = delete_redundant_ilk_snapshot.ilk_id);
    snapshot_block_number BIGINT           := (
        SELECT block_number
        FROM public.headers
        WHERE id = header_id);
    snapshot              api.ilk_snapshot := (
        SELECT (ilk_identifier, ilk_snapshot.block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                pip, mat, created, updated)
        FROM api.ilk_snapshot
        WHERE ilk_identifier = associated_ilk
          AND ilk_snapshot.block_number = snapshot_block_number);
    prev_snapshot         api.ilk_snapshot := (
        SELECT (ilk_identifier, ilk_snapshot.block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty,
                pip, mat, created, updated)
        FROM api.ilk_snapshot
        WHERE ilk_identifier = associated_ilk
          AND ilk_snapshot.block_number < snapshot_block_number
        ORDER BY ilk_snapshot.block_number DESC
        LIMIT 1);
BEGIN
    DELETE
    FROM api.ilk_snapshot
    WHERE ilk_snapshot.ilk_identifier = associated_ilk
      AND ilk_snapshot.block_number = snapshot_block_number
      AND (snapshot.rate IS NULL OR snapshot.rate = prev_snapshot.rate)
      AND (snapshot.art IS NULL OR snapshot.art = prev_snapshot.art)
      AND (snapshot.spot IS NULL OR snapshot.spot = prev_snapshot.spot)
      AND (snapshot.line IS NULL OR snapshot.line = prev_snapshot.line)
      AND (snapshot.dust IS NULL OR snapshot.dust = prev_snapshot.dust)
      AND (snapshot.chop IS NULL OR snapshot.chop = prev_snapshot.chop)
      AND (snapshot.lump IS NULL OR snapshot.lump = prev_snapshot.lump)
      AND (snapshot.flip IS NULL OR snapshot.flip = prev_snapshot.flip)
      AND (snapshot.rho IS NULL OR snapshot.rho = prev_snapshot.rho)
      AND (snapshot.duty IS NULL OR snapshot.duty = prev_snapshot.duty)
      AND (snapshot.pip IS NULL OR snapshot.pip = prev_snapshot.pip)
      AND (snapshot.mat IS NULL OR snapshot.mat = prev_snapshot.mat);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION delete_redundant_ilk_snapshot(ilk_id integer, header_id integer); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.delete_redundant_ilk_snapshot(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: insert_bid_event(bigint, numeric, integer, integer, api.bid_act, numeric, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_bid_event(log_id bigint, bid_id numeric, address_id integer, header_id integer, act api.bid_act, lot numeric, bid_amount numeric) RETURNS void
    LANGUAGE sql
    AS $$
INSERT
INTO maker.bid_event (log_id, bid_id, contract_address, act, lot, bid_amount, ilk_identifier, urn_identifier,
                      block_height)
VALUES (insert_bid_event.log_id,
        insert_bid_event.bid_id,
        (SELECT address FROM public.addresses WHERE addresses.id = insert_bid_event.address_id),
        insert_bid_event.act,
        insert_bid_event.lot,
        insert_bid_event.bid_amount,
        (SELECT ilks.identifier
         FROM maker.flip_ilk
                  JOIN maker.ilks ON flip_ilk.ilk_id = ilks.id
                  JOIN public.headers ON flip_ilk.header_id = headers.id
         WHERE flip_ilk.address_id = insert_bid_event.address_id
         ORDER BY headers.block_number DESC
         LIMIT 1),
        (SELECT usr
         FROM maker.flip_bid_usr
                  JOIN public.headers ON flip_bid_usr.header_id = headers.id
         WHERE flip_bid_usr.bid_id = insert_bid_event.bid_id
           AND flip_bid_usr.address_id = insert_bid_event.address_id
         ORDER BY headers.block_number DESC
         LIMIT 1),
        (SELECT block_number FROM public.headers WHERE id = insert_bid_event.header_id))
$$;


--
-- Name: FUNCTION insert_bid_event(log_id bigint, bid_id numeric, address_id integer, header_id integer, act api.bid_act, lot numeric, bid_amount numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_bid_event(log_id bigint, bid_id numeric, address_id integer, header_id integer, act api.bid_act, lot numeric, bid_amount numeric) IS '@omit';


--
-- Name: insert_bid_event_ilk(maker.flip_ilk); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_bid_event_ilk(new_diff maker.flip_ilk) RETURNS void
    LANGUAGE sql
    AS $$
UPDATE maker.bid_event
SET ilk_identifier = (SELECT identifier FROM maker.ilks WHERE id = new_diff.ilk_id)
WHERE bid_event.contract_address = (SELECT address FROM public.addresses WHERE id = new_diff.address_id)
$$;


--
-- Name: FUNCTION insert_bid_event_ilk(new_diff maker.flip_ilk); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_bid_event_ilk(new_diff maker.flip_ilk) IS '@omit';


--
-- Name: flip_bid_usr; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_usr (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    usr text
);


--
-- Name: insert_bid_event_urn(maker.flip_bid_usr, text); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_bid_event_urn(diff maker.flip_bid_usr, new_usr text) RETURNS void
    LANGUAGE sql
    AS $$
UPDATE maker.bid_event
SET urn_identifier = new_usr
WHERE bid_event.bid_id = diff.bid_id
  AND bid_event.contract_address = (SELECT address FROM public.addresses WHERE id = diff.address_id)
$$;


--
-- Name: FUNCTION insert_bid_event_urn(diff maker.flip_bid_usr, new_usr text); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_bid_event_urn(diff maker.flip_bid_usr, new_usr text) IS '@omit';


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
        WHERE headers.id = NEW.header_id
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
    WITH new_block_number AS (
        SELECT block_number
        FROM public.headers
        WHERE id = NEW.header_id
    )
    INSERT
    INTO api.managed_cdp (cdpi, usr)
    VALUES (NEW.cdpi, NEW.owner)
           -- only update usr if the new owner is coming from the latest owns block we know about for the given cdpi
    ON CONFLICT (cdpi)
        DO UPDATE SET usr = NEW.owner
    WHERE (SELECT block_number FROM new_block_number) >= (
        SELECT MAX(block_number)
        FROM maker.cdp_manager_owns
                 LEFT JOIN public.headers ON cdp_manager_owns.header_id = headers.id
        WHERE cdp_manager_owns.cdpi = NEW.cdpi);
    RETURN NEW;
END
$$;


--
-- Name: insert_flap_created(maker.flap_kick); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_flap_created(new_event maker.flap_kick) RETURNS void
    LANGUAGE sql
    AS $$
UPDATE maker.flap
SET created = api.epoch_to_datetime(headers.block_timestamp)
FROM public.headers
WHERE headers.id = new_event.header_id
  AND flap.address_id = new_event.address_id
  AND flap.bid_id = new_event.bid_id
  AND flap.created IS NULL
$$;


--
-- Name: FUNCTION insert_flap_created(new_event maker.flap_kick); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_flap_created(new_event maker.flap_kick) IS '@omit';


--
-- Name: insert_flip_created(maker.flip_kick); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_flip_created(new_event maker.flip_kick) RETURNS void
    LANGUAGE sql
    AS $$
UPDATE maker.flip
SET created = api.epoch_to_datetime(headers.block_timestamp)
FROM public.headers
WHERE headers.id = new_event.header_id
  AND flip.address_id = new_event.address_id
  AND flip.bid_id = new_event.bid_id
  AND flip.created IS NULL
$$;


--
-- Name: FUNCTION insert_flip_created(new_event maker.flip_kick); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_flip_created(new_event maker.flip_kick) IS '@omit';


--
-- Name: insert_flop_created(maker.flop_kick); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_flop_created(new_event maker.flop_kick) RETURNS void
    LANGUAGE sql
    AS $$
UPDATE maker.flop
SET created = api.epoch_to_datetime(headers.block_timestamp)
FROM public.headers
WHERE headers.id = new_event.header_id
  AND flop.address_id = new_event.address_id
  AND flop.bid_id = new_event.bid_id
  AND flop.created IS NULL
$$;


--
-- Name: FUNCTION insert_flop_created(new_event maker.flop_kick); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_flop_created(new_event maker.flop_kick) IS '@omit';


--
-- Name: vat_ilk_art; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_ilk_art (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    art numeric NOT NULL
);


--
-- Name: insert_new_art(maker.vat_ilk_art); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_art(new_diff maker.vat_ilk_art) RETURNS maker.vat_ilk_art
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.art,
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET art = new_diff.art;
    RETURN new_diff;
END
$$;


--
-- Name: FUNCTION insert_new_art(new_diff maker.vat_ilk_art); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_art(new_diff maker.vat_ilk_art) IS '@omit';


--
-- Name: cat_ilk_chop; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_ilk_chop (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    chop numeric NOT NULL
);


--
-- Name: insert_new_chop(maker.cat_ilk_chop); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_chop(new_diff maker.cat_ilk_chop) RETURNS maker.cat_ilk_chop
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.chop,
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET chop = new_diff.chop;
    RETURN new_diff;
END
$$;


--
-- Name: FUNCTION insert_new_chop(new_diff maker.cat_ilk_chop); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_chop(new_diff maker.cat_ilk_chop) IS '@omit';


--
-- Name: vat_ilk_dust; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_ilk_dust (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    dust numeric NOT NULL
);


--
-- Name: insert_new_dust(maker.vat_ilk_dust); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_dust(new_diff maker.vat_ilk_dust) RETURNS maker.vat_ilk_dust
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.dust,
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET dust = new_diff.dust;
    RETURN new_diff;
END
$$;


--
-- Name: FUNCTION insert_new_dust(new_diff maker.vat_ilk_dust); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_dust(new_diff maker.vat_ilk_dust) IS '@omit';


--
-- Name: jug_ilk_duty; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_ilk_duty (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    duty numeric NOT NULL
);


--
-- Name: insert_new_duty(maker.jug_ilk_duty); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_duty(new_diff maker.jug_ilk_duty) RETURNS maker.jug_ilk_duty
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.duty,
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET duty = new_diff.duty;
    RETURN new_diff;
END
$$;


--
-- Name: FUNCTION insert_new_duty(new_diff maker.jug_ilk_duty); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_duty(new_diff maker.jug_ilk_duty) IS '@omit';


--
-- Name: flap_bid_bid; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_bid_bid (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    bid numeric NOT NULL
);


--
-- Name: insert_new_flap_bid(maker.flap_bid_bid); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flap_bid(new_diff maker.flap_bid_bid) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flap (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flap_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.bid,
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flap_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET bid = new_diff.bid
$$;


--
-- Name: FUNCTION insert_new_flap_bid(new_diff maker.flap_bid_bid); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flap_bid(new_diff maker.flap_bid_bid) IS '@omit';


--
-- Name: flap_bid_end; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_bid_end (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    "end" bigint NOT NULL
);


--
-- Name: insert_new_flap_end(maker.flap_bid_end); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flap_end(new_diff maker.flap_bid_end) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flap (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flap_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff."end",
        flap_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flap_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET "end" = new_diff."end"
$$;


--
-- Name: FUNCTION insert_new_flap_end(new_diff maker.flap_bid_end); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flap_end(new_diff maker.flap_bid_end) IS '@omit';


--
-- Name: flap_bid_guy; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_bid_guy (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    guy text NOT NULL
);


--
-- Name: insert_new_flap_guy(maker.flap_bid_guy); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flap_guy(new_diff maker.flap_bid_guy) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flap (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        new_diff.guy,
        flap_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flap_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET guy = new_diff.guy
$$;


--
-- Name: FUNCTION insert_new_flap_guy(new_diff maker.flap_bid_guy); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flap_guy(new_diff maker.flap_bid_guy) IS '@omit';


--
-- Name: flap_bid_lot; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_bid_lot (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric NOT NULL
);


--
-- Name: insert_new_flap_lot(maker.flap_bid_lot); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flap_lot(new_diff maker.flap_bid_lot) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flap (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flap_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.lot,
        flap_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flap_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET lot = new_diff.lot
$$;


--
-- Name: FUNCTION insert_new_flap_lot(new_diff maker.flap_bid_lot); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flap_lot(new_diff maker.flap_bid_lot) IS '@omit';


--
-- Name: flap_bid_tic; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_bid_tic (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    tic bigint NOT NULL
);


--
-- Name: insert_new_flap_tic(maker.flap_bid_tic); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flap_tic(new_diff maker.flap_bid_tic) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flap (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flap_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.tic,
        flap_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flap_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flap_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET tic = new_diff.tic
$$;


--
-- Name: FUNCTION insert_new_flap_tic(new_diff maker.flap_bid_tic); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flap_tic(new_diff maker.flap_bid_tic) IS '@omit';


--
-- Name: cat_ilk_flip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_ilk_flip (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    flip text
);


--
-- Name: insert_new_flip(maker.cat_ilk_flip); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flip(new_diff maker.cat_ilk_flip) RETURNS maker.cat_ilk_flip
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.flip,
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET flip = new_diff.flip;
    RETURN new_diff;
END
$$;


--
-- Name: FUNCTION insert_new_flip(new_diff maker.cat_ilk_flip); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flip(new_diff maker.cat_ilk_flip) IS '@omit';


--
-- Name: flip_bid_bid; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_bid (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    bid numeric NOT NULL
);


--
-- Name: insert_new_flip_bid(maker.flip_bid_bid); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flip_bid(new_diff maker.flip_bid_bid) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.bid,
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET bid = new_diff.bid
$$;


--
-- Name: FUNCTION insert_new_flip_bid(new_diff maker.flip_bid_bid); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flip_bid(new_diff maker.flip_bid_bid) IS '@omit';


--
-- Name: flip_bid_end; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_end (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    "end" bigint NOT NULL
);


--
-- Name: insert_new_flip_end(maker.flip_bid_end); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flip_end(new_diff maker.flip_bid_end) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff."end",
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET "end" = new_diff."end"
$$;


--
-- Name: FUNCTION insert_new_flip_end(new_diff maker.flip_bid_end); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flip_end(new_diff maker.flip_bid_end) IS '@omit';


--
-- Name: flip_bid_gal; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_gal (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    gal text
);


--
-- Name: insert_new_flip_gal(maker.flip_bid_gal); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flip_gal(new_diff maker.flip_bid_gal) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.gal,
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET gal = new_diff.gal
$$;


--
-- Name: FUNCTION insert_new_flip_gal(new_diff maker.flip_bid_gal); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flip_gal(new_diff maker.flip_bid_gal) IS '@omit';


--
-- Name: flip_bid_guy; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_guy (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    guy text
);


--
-- Name: insert_new_flip_guy(maker.flip_bid_guy); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flip_guy(new_diff maker.flip_bid_guy) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        new_diff.guy,
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET guy = new_diff.guy
$$;


--
-- Name: FUNCTION insert_new_flip_guy(new_diff maker.flip_bid_guy); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flip_guy(new_diff maker.flip_bid_guy) IS '@omit';


--
-- Name: flip_bid_lot; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_lot (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric NOT NULL
);


--
-- Name: insert_new_flip_lot(maker.flip_bid_lot); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flip_lot(new_diff maker.flip_bid_lot) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.lot,
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET lot = new_diff.lot
$$;


--
-- Name: FUNCTION insert_new_flip_lot(new_diff maker.flip_bid_lot); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flip_lot(new_diff maker.flip_bid_lot) IS '@omit';


--
-- Name: flip_bid_tab; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_tab (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    tab numeric NOT NULL
);


--
-- Name: insert_new_flip_tab(maker.flip_bid_tab); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flip_tab(new_diff maker.flip_bid_tab) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.tab,
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET tab = new_diff.tab
$$;


--
-- Name: FUNCTION insert_new_flip_tab(new_diff maker.flip_bid_tab); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flip_tab(new_diff maker.flip_bid_tab) IS '@omit';


--
-- Name: flip_bid_tic; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_bid_tic (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    tic bigint NOT NULL
);


--
-- Name: insert_new_flip_tic(maker.flip_bid_tic); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flip_tic(new_diff maker.flip_bid_tic) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.tic,
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_usr_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET tic = new_diff.tic
$$;


--
-- Name: FUNCTION insert_new_flip_tic(new_diff maker.flip_bid_tic); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flip_tic(new_diff maker.flip_bid_tic) IS '@omit';


--
-- Name: insert_new_flip_usr(maker.flip_bid_usr); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flip_usr(new_diff maker.flip_bid_usr) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flip (bid_id, address_id, block_number, guy, tic, "end", lot, bid, usr, gal, tab, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flip_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.usr,
        flip_bid_gal_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flip_bid_tab_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flip_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET usr = new_diff.usr
$$;


--
-- Name: FUNCTION insert_new_flip_usr(new_diff maker.flip_bid_usr); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flip_usr(new_diff maker.flip_bid_usr) IS '@omit';


--
-- Name: flop_bid_bid; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_bid_bid (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    bid numeric NOT NULL
);


--
-- Name: insert_new_flop_bid(maker.flop_bid_bid); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flop_bid(new_diff maker.flop_bid_bid) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flop (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flop_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.bid,
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flop_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET bid = new_diff.bid
$$;


--
-- Name: FUNCTION insert_new_flop_bid(new_diff maker.flop_bid_bid); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flop_bid(new_diff maker.flop_bid_bid) IS '@omit';


--
-- Name: flop_bid_end; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_bid_end (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    "end" bigint NOT NULL
);


--
-- Name: insert_new_flop_end(maker.flop_bid_end); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flop_end(new_diff maker.flop_bid_end) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flop (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flop_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff."end",
        flop_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flop_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET "end" = new_diff."end"
$$;


--
-- Name: FUNCTION insert_new_flop_end(new_diff maker.flop_bid_end); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flop_end(new_diff maker.flop_bid_end) IS '@omit';


--
-- Name: flop_bid_guy; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_bid_guy (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    guy text
);


--
-- Name: insert_new_flop_guy(maker.flop_bid_guy); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flop_guy(new_diff maker.flop_bid_guy) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flop (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        new_diff.guy,
        flop_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flop_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET guy = new_diff.guy
$$;


--
-- Name: FUNCTION insert_new_flop_guy(new_diff maker.flop_bid_guy); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flop_guy(new_diff maker.flop_bid_guy) IS '@omit';


--
-- Name: flop_bid_lot; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_bid_lot (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric NOT NULL
);


--
-- Name: insert_new_flop_lot(maker.flop_bid_lot); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flop_lot(new_diff maker.flop_bid_lot) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flop (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flop_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_tic_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.lot,
        flop_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flop_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET lot = new_diff.lot
$$;


--
-- Name: FUNCTION insert_new_flop_lot(new_diff maker.flop_bid_lot); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flop_lot(new_diff maker.flop_bid_lot) IS '@omit';


--
-- Name: flop_bid_tic; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_bid_tic (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bid_id numeric NOT NULL,
    tic bigint NOT NULL
);


--
-- Name: insert_new_flop_tic(maker.flop_bid_tic); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_flop_tic(new_diff maker.flop_bid_tic) RETURNS void
    LANGUAGE sql
    AS $$
WITH diff_block AS (
    SELECT block_number, block_timestamp
    FROM public.headers
    WHERE id = new_diff.header_id
)
INSERT
INTO maker.flop (bid_id, address_id, block_number, guy, tic, "end", lot, bid, updated, created)
VALUES (new_diff.bid_id,
        new_diff.address_id,
        (SELECT block_number FROM diff_block),
        flop_bid_guy_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        new_diff.tic,
        flop_bid_end_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_lot_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        flop_bid_bid_before_block(new_diff.bid_id, new_diff.address_id, new_diff.header_id),
        (SELECT api.epoch_to_datetime(block_timestamp) FROM diff_block),
        flop_bid_time_created(new_diff.address_id, new_diff.bid_id))
ON CONFLICT (block_number, bid_id, address_id) DO UPDATE SET tic = new_diff.tic
$$;


--
-- Name: FUNCTION insert_new_flop_tic(new_diff maker.flop_bid_tic); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_flop_tic(new_diff maker.flop_bid_tic) IS '@omit';


--
-- Name: vat_ilk_line; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_ilk_line (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    line numeric NOT NULL
);


--
-- Name: insert_new_line(maker.vat_ilk_line); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_line(new_diff maker.vat_ilk_line) RETURNS maker.vat_ilk_line
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.line,
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET line = new_diff.line;
    RETURN new_diff;
END
$$;


--
-- Name: FUNCTION insert_new_line(new_diff maker.vat_ilk_line); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_line(new_diff maker.vat_ilk_line) IS '@omit';


--
-- Name: cat_ilk_lump; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_ilk_lump (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    lump numeric NOT NULL
);


--
-- Name: insert_new_lump(maker.cat_ilk_lump); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_lump(new_diff maker.cat_ilk_lump) RETURNS maker.cat_ilk_lump
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.lump,
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET lump = new_diff.lump;
    RETURN new_diff;
END
$$;


--
-- Name: FUNCTION insert_new_lump(new_diff maker.cat_ilk_lump); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_lump(new_diff maker.cat_ilk_lump) IS '@omit';


--
-- Name: spot_ilk_mat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_ilk_mat (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    mat numeric NOT NULL
);


--
-- Name: insert_new_mat(maker.spot_ilk_mat); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_mat(new_diff maker.spot_ilk_mat) RETURNS maker.spot_ilk_mat
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.mat,
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET mat = new_diff.mat;
    RETURN new_diff;
END
$$;


--
-- Name: FUNCTION insert_new_mat(new_diff maker.spot_ilk_mat); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_mat(new_diff maker.spot_ilk_mat) IS '@omit';


--
-- Name: spot_ilk_pip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_ilk_pip (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    pip text
);


--
-- Name: insert_new_pip(maker.spot_ilk_pip); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_pip(new_diff maker.spot_ilk_pip) RETURNS maker.spot_ilk_pip
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.pip,
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET pip = new_diff.pip;
    RETURN new_diff;
END
$$;


--
-- Name: FUNCTION insert_new_pip(new_diff maker.spot_ilk_pip); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_pip(new_diff maker.spot_ilk_pip) IS '@omit';


--
-- Name: vat_ilk_rate; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_ilk_rate (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    rate numeric NOT NULL
);


--
-- Name: insert_new_rate(maker.vat_ilk_rate); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_rate(new_diff maker.vat_ilk_rate) RETURNS maker.vat_ilk_rate
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            new_diff.rate,
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET rate = new_diff.rate;
    RETURN new_diff;
END
$$;


--
-- Name: FUNCTION insert_new_rate(new_diff maker.vat_ilk_rate); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_rate(new_diff maker.vat_ilk_rate) IS '@omit';


--
-- Name: jug_ilk_rho; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_ilk_rho (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    rho numeric NOT NULL
);


--
-- Name: insert_new_rho(maker.jug_ilk_rho); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_rho(new_diff maker.jug_ilk_rho) RETURNS maker.jug_ilk_rho
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.rho,
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET rho = new_diff.rho;
    RETURN new_diff;
END
$$;


--
-- Name: FUNCTION insert_new_rho(new_diff maker.jug_ilk_rho); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_rho(new_diff maker.jug_ilk_rho) IS '@omit';


--
-- Name: vat_ilk_spot; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_ilk_spot (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    spot numeric NOT NULL
);


--
-- Name: insert_new_spot(maker.vat_ilk_spot); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_spot(new_diff maker.vat_ilk_spot) RETURNS maker.vat_ilk_spot
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE headers.id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE headers.id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.spot,
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET spot = new_diff.spot;
    RETURN new_diff;
END
$$;


--
-- Name: FUNCTION insert_new_spot(new_diff maker.vat_ilk_spot); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_spot(new_diff maker.vat_ilk_spot) IS '@omit';


--
-- Name: insert_new_time_created(maker.vat_init); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_new_time_created(new_event maker.vat_init) RETURNS maker.vat_init
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = new_event.ilk_id);
    diff_timestamp      TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE headers.id = new_event.header_id);
BEGIN
    UPDATE api.ilk_snapshot
    SET created = diff_timestamp
    FROM public.headers
    WHERE headers.block_number = ilk_snapshot.block_number
      AND ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.created IS NULL;
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION insert_new_time_created(new_event maker.vat_init); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_new_time_created(new_event maker.vat_init) IS '@omit';


--
-- Name: vat_urn_art; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_urn_art (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    urn_id integer NOT NULL,
    art numeric NOT NULL
);


--
-- Name: insert_urn_art(maker.vat_urn_art); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_urn_art(new_diff maker.vat_urn_art) RETURNS maker.vat_urn_art
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH urn AS (
        SELECT urns.identifier AS urn_identifier, ilks.identifier AS ilk_identifier
        FROM maker.urns
                 LEFT JOIN maker.ilks ON urns.ilk_id = ilks.id
        WHERE urns.id = new_diff.urn_id),
         new_diff_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp, block_number
             FROM public.headers
             WHERE id = new_diff.header_id)
    INSERT
    INTO api.urn_snapshot (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)
    VALUES ((SELECT urn_identifier FROM urn),
            (SELECT ilk_identifier FROM urn),
            (SELECT block_number FROM new_diff_header),
            urn_ink_before_block(new_diff.urn_id, new_diff.header_id),
            new_diff.art,
            urn_time_created(new_diff.urn_id),
            (SELECT block_timestamp FROM new_diff_header))
    ON CONFLICT (urn_identifier, ilk_identifier, block_height)
        DO UPDATE SET art = new_diff.art;
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION insert_urn_art(new_diff maker.vat_urn_art); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_urn_art(new_diff maker.vat_urn_art) IS '@omit';


--
-- Name: vat_urn_ink; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_urn_ink (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    urn_id integer NOT NULL,
    ink numeric NOT NULL
);


--
-- Name: insert_urn_ink(maker.vat_urn_ink); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.insert_urn_ink(new_diff maker.vat_urn_ink) RETURNS maker.vat_urn_ink
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH urn AS (
        SELECT urns.identifier AS urn_identifier, ilks.identifier AS ilk_identifier
        FROM maker.urns
                 LEFT JOIN maker.ilks ON urns.ilk_id = ilks.id
        WHERE urns.id = new_diff.urn_id),
         new_diff_header AS (
             SELECT api.epoch_to_datetime(block_timestamp) AS block_timestamp, block_number
             FROM public.headers
             WHERE id = new_diff.header_id)
    INSERT
    INTO api.urn_snapshot (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)
    VALUES ((SELECT urn_identifier FROM urn),
            (SELECT ilk_identifier FROM urn),
            (SELECT block_number FROM new_diff_header),
            new_diff.ink,
            urn_art_before_block(new_diff.urn_id, new_diff.header_id),
            urn_time_created(new_diff.urn_id),
            (SELECT block_timestamp FROM new_diff_header))
    ON CONFLICT (urn_identifier, ilk_identifier, block_height)
        DO UPDATE SET ink = new_diff.ink;
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION insert_urn_ink(new_diff maker.vat_urn_ink); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.insert_urn_ink(new_diff maker.vat_urn_ink) IS '@omit';


--
-- Name: update_arts_until_next_diff(maker.vat_ilk_art, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_arts_until_next_diff(start_at_diff maker.vat_ilk_art, new_art numeric) RETURNS maker.vat_ilk_art
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_art_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_ilk_art
                 LEFT JOIN public.headers ON vat_ilk_art.header_id = headers.id
        WHERE vat_ilk_art.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET art = new_art
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_art_diff_block IS NULL
        OR ilk_snapshot.block_number < next_art_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_arts_until_next_diff(start_at_diff maker.vat_ilk_art, new_art numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_arts_until_next_diff(start_at_diff maker.vat_ilk_art, new_art numeric) IS '@omit';


--
-- Name: update_bid_event_ilk(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_bid_event_ilk() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        PERFORM maker.insert_bid_event_ilk(NEW);
    ELSIF TG_OP = 'DELETE' THEN
        PERFORM maker.clear_bid_event_ilk(OLD);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_bid_event_urn(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_bid_event_urn() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        PERFORM maker.insert_bid_event_urn(NEW, NEW.usr);
    ELSIF TG_OP = 'DELETE' THEN
        PERFORM maker.insert_bid_event_urn(OLD, NULL);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_bid_kick_tend_dent_event(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_bid_kick_tend_dent_event() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    PERFORM maker.insert_bid_event(NEW.log_id, NEW.bid_id, NEW.address_id, NEW.header_id, TG_ARGV[0]::api.bid_act,
                                   NEW.lot, NEW.bid);
    RETURN NULL;
END
$$;


--
-- Name: update_bid_tick_deal_yank_event(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_bid_tick_deal_yank_event() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    PERFORM maker.insert_bid_event(NEW.log_id, NEW.bid_id, NEW.address_id, NEW.header_id, TG_ARGV[0]::api.bid_act, NULL,
                                   NULL);
    RETURN NULL;
END
$$;


--
-- Name: update_chops_until_next_diff(maker.cat_ilk_chop, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_chops_until_next_diff(start_at_diff maker.cat_ilk_chop, new_chop numeric) RETURNS maker.cat_ilk_chop
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_chop_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.cat_ilk_chop
                 LEFT JOIN public.headers ON cat_ilk_chop.header_id = headers.id
        WHERE cat_ilk_chop.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET chop = new_chop
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_chop_diff_block IS NULL
        OR ilk_snapshot.block_number < next_chop_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_chops_until_next_diff(start_at_diff maker.cat_ilk_chop, new_chop numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_chops_until_next_diff(start_at_diff maker.cat_ilk_chop, new_chop numeric) IS '@omit';


--
-- Name: update_dusts_until_next_diff(maker.vat_ilk_dust, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_dusts_until_next_diff(start_at_diff maker.vat_ilk_dust, new_dust numeric) RETURNS maker.vat_ilk_dust
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_dust_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_ilk_dust
                 LEFT JOIN public.headers ON vat_ilk_dust.header_id = headers.id
        WHERE vat_ilk_dust.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET dust = new_dust
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_dust_diff_block IS NULL
        OR ilk_snapshot.block_number < next_dust_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_dusts_until_next_diff(start_at_diff maker.vat_ilk_dust, new_dust numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_dusts_until_next_diff(start_at_diff maker.vat_ilk_dust, new_dust numeric) IS '@omit';


--
-- Name: update_duties_until_next_diff(maker.jug_ilk_duty, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_duties_until_next_diff(start_at_diff maker.jug_ilk_duty, new_duty numeric) RETURNS maker.jug_ilk_duty
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_duty_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.jug_ilk_duty
                 LEFT JOIN public.headers ON jug_ilk_duty.header_id = headers.id
        WHERE jug_ilk_duty.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET duty = new_duty
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_duty_diff_block IS NULL
        OR ilk_snapshot.block_number < next_duty_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_duties_until_next_diff(start_at_diff maker.jug_ilk_duty, new_duty numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_duties_until_next_diff(start_at_diff maker.jug_ilk_duty, new_duty numeric) IS '@omit';


--
-- Name: update_flap_bids(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flap_bids() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flap_bid(NEW);
        PERFORM maker.update_flap_bids_until_next_diff(NEW, NEW.bid);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flap_bids_until_next_diff(
                OLD,
                flap_bid_bid_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flap(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flap_bids_until_next_diff(maker.flap_bid_bid, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flap_bids_until_next_diff(start_at_diff maker.flap_bid_bid, new_bid numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_bid_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_bid
                 LEFT JOIN public.headers ON flap_bid_bid.header_id = headers.id
        WHERE flap_bid_bid.bid_id = start_at_diff.bid_id
          AND flap_bid_bid.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET bid = new_bid
    WHERE flap.bid_id = start_at_diff.bid_id
      AND flap.address_id = start_at_diff.address_id
      AND flap.block_number >= diff_block_number
      AND (next_bid_diff_block IS NULL
        OR flap.block_number < next_bid_diff_block);
END
$$;


--
-- Name: FUNCTION update_flap_bids_until_next_diff(start_at_diff maker.flap_bid_bid, new_bid numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flap_bids_until_next_diff(start_at_diff maker.flap_bid_bid, new_bid numeric) IS '@omit';


--
-- Name: update_flap_created(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flap_created() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        PERFORM maker.insert_flap_created(NEW);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.clear_flap_created(OLD);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_flap_created(); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flap_created() IS '@omit';


--
-- Name: update_flap_ends(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flap_ends() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flap_end(NEW);
        PERFORM maker.update_flap_ends_until_next_diff(NEW, NEW."end");
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flap_ends_until_next_diff(
                OLD,
                flap_bid_end_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flap(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flap_ends_until_next_diff(maker.flap_bid_end, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flap_ends_until_next_diff(start_at_diff maker.flap_bid_end, new_end numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_end_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_end
                 LEFT JOIN public.headers ON flap_bid_end.header_id = headers.id
        WHERE flap_bid_end.bid_id = start_at_diff.bid_id
          AND flap_bid_end.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET "end" = new_end
    WHERE flap.bid_id = start_at_diff.bid_id
      AND flap.address_id = start_at_diff.address_id
      AND flap.block_number >= diff_block_number
      AND (next_end_diff_block IS NULL
        OR flap.block_number < next_end_diff_block);
END
$$;


--
-- Name: FUNCTION update_flap_ends_until_next_diff(start_at_diff maker.flap_bid_end, new_end numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flap_ends_until_next_diff(start_at_diff maker.flap_bid_end, new_end numeric) IS '@omit';


--
-- Name: update_flap_guys(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flap_guys() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flap_guy(NEW);
        PERFORM maker.update_flap_guys_until_next_diff(NEW, NEW.guy);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flap_guys_until_next_diff(
                OLD,
                flap_bid_guy_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flap(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flap_guys_until_next_diff(maker.flap_bid_guy, text); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flap_guys_until_next_diff(start_at_diff maker.flap_bid_guy, new_guy text) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_guy_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_guy
                 LEFT JOIN public.headers ON flap_bid_guy.header_id = headers.id
        WHERE flap_bid_guy.bid_id = start_at_diff.bid_id
          AND flap_bid_guy.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET guy = new_guy
    WHERE flap.bid_id = start_at_diff.bid_id
      AND flap.address_id = start_at_diff.address_id
      AND flap.block_number >= diff_block_number
      AND (next_guy_diff_block IS NULL
        OR flap.block_number < next_guy_diff_block);
END
$$;


--
-- Name: FUNCTION update_flap_guys_until_next_diff(start_at_diff maker.flap_bid_guy, new_guy text); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flap_guys_until_next_diff(start_at_diff maker.flap_bid_guy, new_guy text) IS '@omit';


--
-- Name: update_flap_lots(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flap_lots() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flap_lot(NEW);
        PERFORM maker.update_flap_lots_until_next_diff(NEW, NEW.lot);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flap_lots_until_next_diff(
                OLD,
                flap_bid_lot_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flap(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flap_lots_until_next_diff(maker.flap_bid_lot, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flap_lots_until_next_diff(start_at_diff maker.flap_bid_lot, new_lot numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_lot_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_lot
                 LEFT JOIN public.headers ON flap_bid_lot.header_id = headers.id
        WHERE flap_bid_lot.bid_id = start_at_diff.bid_id
          AND flap_bid_lot.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET lot = new_lot
    WHERE flap.bid_id = start_at_diff.bid_id
      AND flap.address_id = start_at_diff.address_id
      AND flap.block_number >= diff_block_number
      AND (next_lot_diff_block IS NULL
        OR flap.block_number < next_lot_diff_block);
END
$$;


--
-- Name: FUNCTION update_flap_lots_until_next_diff(start_at_diff maker.flap_bid_lot, new_lot numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flap_lots_until_next_diff(start_at_diff maker.flap_bid_lot, new_lot numeric) IS '@omit';


--
-- Name: update_flap_tics(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flap_tics() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flap_tic(NEW);
        PERFORM maker.update_flap_tics_until_next_diff(NEW, NEW.tic);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flap_tics_until_next_diff(
                OLD,
                flap_bid_tic_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flap(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flap_tics_until_next_diff(maker.flap_bid_tic, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flap_tics_until_next_diff(start_at_diff maker.flap_bid_tic, new_tic numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_tic_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flap_bid_tic
                 LEFT JOIN public.headers ON flap_bid_tic.header_id = headers.id
        WHERE flap_bid_tic.bid_id = start_at_diff.bid_id
          AND flap_bid_tic.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flap
    SET tic = new_tic
    WHERE flap.bid_id = start_at_diff.bid_id
      AND flap.address_id = start_at_diff.address_id
      AND flap.block_number >= diff_block_number
      AND (next_tic_diff_block IS NULL
        OR flap.block_number < next_tic_diff_block);
END
$$;


--
-- Name: FUNCTION update_flap_tics_until_next_diff(start_at_diff maker.flap_bid_tic, new_tic numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flap_tics_until_next_diff(start_at_diff maker.flap_bid_tic, new_tic numeric) IS '@omit';


--
-- Name: update_flip_bids(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_bids() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flip_bid(NEW);
        PERFORM maker.update_flip_bids_until_next_diff(NEW, NEW.bid);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flip_bids_until_next_diff(
                OLD,
                flip_bid_bid_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flip(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flip_bids_until_next_diff(maker.flip_bid_bid, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_bids_until_next_diff(start_at_diff maker.flip_bid_bid, new_bid numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_bid_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_bid
                 LEFT JOIN public.headers ON flip_bid_bid.header_id = headers.id
        WHERE flip_bid_bid.bid_id = start_at_diff.bid_id
          AND flip_bid_bid.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET bid = new_bid
    WHERE flip.bid_id = start_at_diff.bid_id
      AND flip.address_id = start_at_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_bid_diff_block IS NULL
        OR flip.block_number < next_bid_diff_block);
END
$$;


--
-- Name: FUNCTION update_flip_bids_until_next_diff(start_at_diff maker.flip_bid_bid, new_bid numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flip_bids_until_next_diff(start_at_diff maker.flip_bid_bid, new_bid numeric) IS '@omit';


--
-- Name: update_flip_created(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_created() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        PERFORM maker.insert_flip_created(NEW);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.clear_flip_created(OLD);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_flip_created(); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flip_created() IS '@omit';


--
-- Name: update_flip_ends(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_ends() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flip_end(NEW);
        PERFORM maker.update_flip_ends_until_next_diff(NEW, NEW."end");
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flip_ends_until_next_diff(
                OLD,
                flip_bid_end_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flip(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flip_ends_until_next_diff(maker.flip_bid_end, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_ends_until_next_diff(start_at_diff maker.flip_bid_end, new_end numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_end_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_end
                 LEFT JOIN public.headers ON flip_bid_end.header_id = headers.id
        WHERE flip_bid_end.bid_id = start_at_diff.bid_id
          AND flip_bid_end.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET "end" = new_end
    WHERE flip.bid_id = start_at_diff.bid_id
      AND flip.address_id = start_at_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_end_diff_block IS NULL
        OR flip.block_number < next_end_diff_block);
END
$$;


--
-- Name: FUNCTION update_flip_ends_until_next_diff(start_at_diff maker.flip_bid_end, new_end numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flip_ends_until_next_diff(start_at_diff maker.flip_bid_end, new_end numeric) IS '@omit';


--
-- Name: update_flip_gals(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_gals() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flip_gal(NEW);
        PERFORM maker.update_flip_gals_until_next_diff(NEW, NEW.gal);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flip_gals_until_next_diff(
                OLD,
                flip_bid_gal_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flip(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flip_gals_until_next_diff(maker.flip_bid_gal, text); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_gals_until_next_diff(start_at_diff maker.flip_bid_gal, new_gal text) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_gal_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_gal
                 LEFT JOIN public.headers ON flip_bid_gal.header_id = headers.id
        WHERE flip_bid_gal.bid_id = start_at_diff.bid_id
          AND flip_bid_gal.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET gal = new_gal
    WHERE flip.bid_id = start_at_diff.bid_id
      AND flip.address_id = start_at_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_gal_diff_block IS NULL
        OR flip.block_number < next_gal_diff_block);
END
$$;


--
-- Name: FUNCTION update_flip_gals_until_next_diff(start_at_diff maker.flip_bid_gal, new_gal text); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flip_gals_until_next_diff(start_at_diff maker.flip_bid_gal, new_gal text) IS '@omit';


--
-- Name: update_flip_guys(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_guys() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flip_guy(NEW);
        PERFORM maker.update_flip_guys_until_next_diff(NEW, NEW.guy);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flip_guys_until_next_diff(
                OLD,
                flip_bid_guy_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flip(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flip_guys_until_next_diff(maker.flip_bid_guy, text); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_guys_until_next_diff(start_at_diff maker.flip_bid_guy, new_guy text) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_guy_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_guy
                 LEFT JOIN public.headers ON flip_bid_guy.header_id = headers.id
        WHERE flip_bid_guy.bid_id = start_at_diff.bid_id
          AND flip_bid_guy.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET guy = new_guy
    WHERE flip.bid_id = start_at_diff.bid_id
      AND flip.address_id = start_at_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_guy_diff_block IS NULL
        OR flip.block_number < next_guy_diff_block);
END
$$;


--
-- Name: FUNCTION update_flip_guys_until_next_diff(start_at_diff maker.flip_bid_guy, new_guy text); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flip_guys_until_next_diff(start_at_diff maker.flip_bid_guy, new_guy text) IS '@omit';


--
-- Name: update_flip_lots(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_lots() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flip_lot(NEW);
        PERFORM maker.update_flip_lots_until_next_diff(NEW, NEW.lot);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flip_lots_until_next_diff(
                OLD,
                flip_bid_lot_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flip(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flip_lots_until_next_diff(maker.flip_bid_lot, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_lots_until_next_diff(start_at_diff maker.flip_bid_lot, new_lot numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_lot_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_lot
                 LEFT JOIN public.headers ON flip_bid_lot.header_id = headers.id
        WHERE flip_bid_lot.bid_id = start_at_diff.bid_id
          AND flip_bid_lot.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET lot = new_lot
    WHERE flip.bid_id = start_at_diff.bid_id
      AND flip.address_id = start_at_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_lot_diff_block IS NULL
        OR flip.block_number < next_lot_diff_block);
END
$$;


--
-- Name: FUNCTION update_flip_lots_until_next_diff(start_at_diff maker.flip_bid_lot, new_lot numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flip_lots_until_next_diff(start_at_diff maker.flip_bid_lot, new_lot numeric) IS '@omit';


--
-- Name: update_flip_tabs(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_tabs() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flip_tab(NEW);
        PERFORM maker.update_flip_tabs_until_next_diff(NEW, NEW.tab);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flip_tabs_until_next_diff(
                OLD,
                flip_bid_tab_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flip(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flip_tabs_until_next_diff(maker.flip_bid_tab, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_tabs_until_next_diff(start_at_diff maker.flip_bid_tab, new_tab numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_tab_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_tab
                 LEFT JOIN public.headers ON flip_bid_tab.header_id = headers.id
        WHERE flip_bid_tab.bid_id = start_at_diff.bid_id
          AND flip_bid_tab.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET tab = new_tab
    WHERE flip.bid_id = start_at_diff.bid_id
      AND flip.address_id = start_at_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_tab_diff_block IS NULL
        OR flip.block_number < next_tab_diff_block);
END
$$;


--
-- Name: FUNCTION update_flip_tabs_until_next_diff(start_at_diff maker.flip_bid_tab, new_tab numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flip_tabs_until_next_diff(start_at_diff maker.flip_bid_tab, new_tab numeric) IS '@omit';


--
-- Name: update_flip_tics(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_tics() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flip_tic(NEW);
        PERFORM maker.update_flip_tics_until_next_diff(NEW, NEW.tic);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flip_tics_until_next_diff(
                OLD,
                flip_bid_tic_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flip(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flip_tics_until_next_diff(maker.flip_bid_tic, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_tics_until_next_diff(start_at_diff maker.flip_bid_tic, new_tic numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_tic_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_tic
                 LEFT JOIN public.headers ON flip_bid_tic.header_id = headers.id
        WHERE flip_bid_tic.bid_id = start_at_diff.bid_id
          AND flip_bid_tic.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET tic = new_tic
    WHERE flip.bid_id = start_at_diff.bid_id
      AND flip.address_id = start_at_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_tic_diff_block IS NULL
        OR flip.block_number < next_tic_diff_block);
END
$$;


--
-- Name: FUNCTION update_flip_tics_until_next_diff(start_at_diff maker.flip_bid_tic, new_tic numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flip_tics_until_next_diff(start_at_diff maker.flip_bid_tic, new_tic numeric) IS '@omit';


--
-- Name: update_flip_usrs(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_usrs() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flip_usr(NEW);
        PERFORM maker.update_flip_usrs_until_next_diff(NEW, NEW.usr);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flip_usrs_until_next_diff(
                OLD,
                flip_bid_usr_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flip(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flip_usrs_until_next_diff(maker.flip_bid_usr, text); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flip_usrs_until_next_diff(stat_at_diff maker.flip_bid_usr, new_usr text) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = stat_at_diff.header_id);
    next_usr_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flip_bid_usr
                 LEFT JOIN public.headers ON flip_bid_usr.header_id = headers.id
        WHERE flip_bid_usr.bid_id = stat_at_diff.bid_id
          AND flip_bid_usr.address_id = stat_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flip
    SET usr = new_usr
    WHERE flip.bid_id = stat_at_diff.bid_id
      AND flip.address_id = stat_at_diff.address_id
      AND flip.block_number >= diff_block_number
      AND (next_usr_diff_block IS NULL
        OR flip.block_number < next_usr_diff_block);
END
$$;


--
-- Name: FUNCTION update_flip_usrs_until_next_diff(stat_at_diff maker.flip_bid_usr, new_usr text); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flip_usrs_until_next_diff(stat_at_diff maker.flip_bid_usr, new_usr text) IS '@omit';


--
-- Name: update_flips_until_next_diff(maker.cat_ilk_flip, text); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flips_until_next_diff(start_at_diff maker.cat_ilk_flip, new_flip text) RETURNS maker.cat_ilk_flip
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_flip_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.cat_ilk_flip
                 LEFT JOIN public.headers ON cat_ilk_flip.header_id = headers.id
        WHERE cat_ilk_flip.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET flip = new_flip
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_flip_diff_block IS NULL
        OR ilk_snapshot.block_number < next_flip_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_flips_until_next_diff(start_at_diff maker.cat_ilk_flip, new_flip text); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flips_until_next_diff(start_at_diff maker.cat_ilk_flip, new_flip text) IS '@omit';


--
-- Name: update_flop_bids(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flop_bids() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flop_bid(NEW);
        PERFORM maker.update_flop_bids_until_next_diff(NEW, NEW.bid);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flop_bids_until_next_diff(
                OLD,
                flop_bid_bid_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flop(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flop_bids_until_next_diff(maker.flop_bid_bid, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flop_bids_until_next_diff(start_at_diff maker.flop_bid_bid, new_bid numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_bid_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flop_bid_bid
                 LEFT JOIN public.headers ON flop_bid_bid.header_id = headers.id
        WHERE flop_bid_bid.bid_id = start_at_diff.bid_id
          AND flop_bid_bid.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flop
    SET bid = new_bid
    WHERE flop.bid_id = start_at_diff.bid_id
      AND flop.address_id = start_at_diff.address_id
      AND flop.block_number >= diff_block_number
      AND (next_bid_diff_block IS NULL
        OR flop.block_number < next_bid_diff_block);
END
$$;


--
-- Name: FUNCTION update_flop_bids_until_next_diff(start_at_diff maker.flop_bid_bid, new_bid numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flop_bids_until_next_diff(start_at_diff maker.flop_bid_bid, new_bid numeric) IS '@omit';


--
-- Name: update_flop_created(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flop_created() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        PERFORM maker.insert_flop_created(NEW);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.clear_flop_created(OLD);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_flop_created(); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flop_created() IS '@omit';


--
-- Name: update_flop_ends(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flop_ends() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flop_end(NEW);
        PERFORM maker.update_flop_ends_until_next_diff(NEW, NEW."end");
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flop_ends_until_next_diff(
                OLD,
                flop_bid_end_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flop(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flop_ends_until_next_diff(maker.flop_bid_end, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flop_ends_until_next_diff(start_at_diff maker.flop_bid_end, new_end numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_end_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flop_bid_end
                 LEFT JOIN public.headers ON flop_bid_end.header_id = headers.id
        WHERE flop_bid_end.bid_id = start_at_diff.bid_id
          AND flop_bid_end.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flop
    SET "end" = new_end
    WHERE flop.bid_id = start_at_diff.bid_id
      AND flop.address_id = start_at_diff.address_id
      AND flop.block_number >= diff_block_number
      AND (next_end_diff_block IS NULL
        OR flop.block_number < next_end_diff_block);
END
$$;


--
-- Name: FUNCTION update_flop_ends_until_next_diff(start_at_diff maker.flop_bid_end, new_end numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flop_ends_until_next_diff(start_at_diff maker.flop_bid_end, new_end numeric) IS '@omit';


--
-- Name: update_flop_guys(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flop_guys() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flop_guy(NEW);
        PERFORM maker.update_flop_guys_until_next_diff(NEW, NEW.guy);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flop_guys_until_next_diff(
                OLD,
                flop_bid_guy_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flop(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flop_guys_until_next_diff(maker.flop_bid_guy, text); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flop_guys_until_next_diff(start_at_diff maker.flop_bid_guy, new_guy text) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_guy_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flop_bid_guy
                 LEFT JOIN public.headers ON flop_bid_guy.header_id = headers.id
        WHERE flop_bid_guy.bid_id = start_at_diff.bid_id
          AND flop_bid_guy.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flop
    SET guy = new_guy
    WHERE flop.bid_id = start_at_diff.bid_id
      AND flop.address_id = start_at_diff.address_id
      AND flop.block_number >= diff_block_number
      AND (next_guy_diff_block IS NULL
        OR flop.block_number < next_guy_diff_block);
END
$$;


--
-- Name: FUNCTION update_flop_guys_until_next_diff(start_at_diff maker.flop_bid_guy, new_guy text); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flop_guys_until_next_diff(start_at_diff maker.flop_bid_guy, new_guy text) IS '@omit';


--
-- Name: update_flop_lots(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flop_lots() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flop_lot(NEW);
        PERFORM maker.update_flop_lots_until_next_diff(NEW, NEW.lot);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flop_lots_until_next_diff(
                OLD,
                flop_bid_lot_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flop(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flop_lots_until_next_diff(maker.flop_bid_lot, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flop_lots_until_next_diff(start_at_diff maker.flop_bid_lot, new_lot numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_lot_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flop_bid_lot
                 LEFT JOIN public.headers ON flop_bid_lot.header_id = headers.id
        WHERE flop_bid_lot.bid_id = start_at_diff.bid_id
          AND flop_bid_lot.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flop
    SET lot = new_lot
    WHERE flop.bid_id = start_at_diff.bid_id
      AND flop.address_id = start_at_diff.address_id
      AND flop.block_number >= diff_block_number
      AND (next_lot_diff_block IS NULL
        OR flop.block_number < next_lot_diff_block);
END
$$;


--
-- Name: FUNCTION update_flop_lots_until_next_diff(start_at_diff maker.flop_bid_lot, new_lot numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flop_lots_until_next_diff(start_at_diff maker.flop_bid_lot, new_lot numeric) IS '@omit';


--
-- Name: update_flop_tics(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flop_tics() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flop_tic(NEW);
        PERFORM maker.update_flop_tics_until_next_diff(NEW, NEW.tic);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flop_tics_until_next_diff(
                OLD,
                flop_bid_tic_before_block(OLD.bid_id, OLD.address_id, OLD.header_id));
        PERFORM maker.delete_obsolete_flop(OLD.bid_id, OLD.address_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_flop_tics_until_next_diff(maker.flop_bid_tic, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_flop_tics_until_next_diff(start_at_diff maker.flop_bid_tic, new_tic numeric) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_tic_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.flop_bid_tic
                 LEFT JOIN public.headers ON flop_bid_tic.header_id = headers.id
        WHERE flop_bid_tic.bid_id = start_at_diff.bid_id
          AND flop_bid_tic.address_id = start_at_diff.address_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE maker.flop
    SET tic = new_tic
    WHERE flop.bid_id = start_at_diff.bid_id
      AND flop.address_id = start_at_diff.address_id
      AND flop.block_number >= diff_block_number
      AND (next_tic_diff_block IS NULL
        OR flop.block_number < next_tic_diff_block);
END
$$;


--
-- Name: FUNCTION update_flop_tics_until_next_diff(start_at_diff maker.flop_bid_tic, new_tic numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_flop_tics_until_next_diff(start_at_diff maker.flop_bid_tic, new_tic numeric) IS '@omit';


--
-- Name: update_ilk_arts(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_ilk_arts() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_art(NEW);
        PERFORM maker.update_arts_until_next_diff(NEW, NEW.art);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_arts_until_next_diff(OLD, ilk_art_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_ilk_chops(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_ilk_chops() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_chop(NEW);
        PERFORM maker.update_chops_until_next_diff(NEW, NEW.chop);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_chops_until_next_diff(OLD, ilk_chop_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_ilk_dusts(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_ilk_dusts() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_dust(NEW);
        PERFORM maker.update_dusts_until_next_diff(NEW, NEW.dust);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_dusts_until_next_diff(OLD, ilk_dust_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_ilk_duties(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_ilk_duties() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_duty(NEW);
        PERFORM maker.update_duties_until_next_diff(NEW, NEW.duty);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_duties_until_next_diff(OLD, ilk_duty_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_ilk_flips(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_ilk_flips() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flip(NEW);
        PERFORM maker.update_flips_until_next_diff(NEW, NEW.flip);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flips_until_next_diff(OLD, ilk_flip_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_ilk_lines(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_ilk_lines() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_line(NEW);
        PERFORM maker.update_lines_until_next_diff(NEW, NEW.line);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_lines_until_next_diff(OLD, ilk_line_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_ilk_lumps(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_ilk_lumps() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_lump(NEW);
        PERFORM maker.update_lumps_until_next_diff(NEW, NEW.lump);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_lumps_until_next_diff(OLD, ilk_lump_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_ilk_mats(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_ilk_mats() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_mat(NEW);
        PERFORM maker.update_mats_until_next_diff(NEW, NEW.mat);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_mats_until_next_diff(OLD, ilk_mat_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_ilk_pips(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_ilk_pips() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_pip(NEW);
        PERFORM maker.update_pips_until_next_diff(NEW, NEW.pip);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_pips_until_next_diff(OLD, ilk_pip_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_ilk_rates(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_ilk_rates() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_rate(NEW);
        PERFORM maker.update_rates_until_next_diff(NEW, NEW.rate);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_rates_until_next_diff(OLD, ilk_rate_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_ilk_rhos(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_ilk_rhos() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_rho(NEW);
        PERFORM maker.update_rhos_until_next_diff(NEW, NEW.rho);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_rhos_until_next_diff(OLD, ilk_rho_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_ilk_spots(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_ilk_spots() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_spot(NEW);
        PERFORM maker.update_spots_until_next_diff(NEW, NEW.spot);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_spots_until_next_diff(OLD, ilk_spot_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_lines_until_next_diff(maker.vat_ilk_line, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_lines_until_next_diff(start_at_diff maker.vat_ilk_line, new_line numeric) RETURNS maker.vat_ilk_line
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_line_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_ilk_line
                 LEFT JOIN public.headers ON vat_ilk_line.header_id = headers.id
        WHERE vat_ilk_line.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET line = new_line
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_line_diff_block IS NULL
        OR ilk_snapshot.block_number < next_line_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_lines_until_next_diff(start_at_diff maker.vat_ilk_line, new_line numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_lines_until_next_diff(start_at_diff maker.vat_ilk_line, new_line numeric) IS '@omit';


--
-- Name: update_lumps_until_next_diff(maker.cat_ilk_lump, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_lumps_until_next_diff(start_at_diff maker.cat_ilk_lump, new_lump numeric) RETURNS maker.cat_ilk_lump
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_lump_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.cat_ilk_lump
                 LEFT JOIN public.headers ON cat_ilk_lump.header_id = headers.id
        WHERE cat_ilk_lump.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET lump = new_lump
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_lump_diff_block IS NULL
        OR ilk_snapshot.block_number < next_lump_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_lumps_until_next_diff(start_at_diff maker.cat_ilk_lump, new_lump numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_lumps_until_next_diff(start_at_diff maker.cat_ilk_lump, new_lump numeric) IS '@omit';


--
-- Name: update_mats_until_next_diff(maker.spot_ilk_mat, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_mats_until_next_diff(start_at_diff maker.spot_ilk_mat, new_mat numeric) RETURNS maker.spot_ilk_mat
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_mat_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.spot_ilk_mat
                 LEFT JOIN public.headers ON spot_ilk_mat.header_id = headers.id
        WHERE spot_ilk_mat.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET mat = new_mat
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_mat_diff_block IS NULL
        OR ilk_snapshot.block_number < next_mat_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_mats_until_next_diff(start_at_diff maker.spot_ilk_mat, new_mat numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_mats_until_next_diff(start_at_diff maker.spot_ilk_mat, new_mat numeric) IS '@omit';


--
-- Name: update_pips_until_next_diff(maker.spot_ilk_pip, text); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_pips_until_next_diff(start_at_diff maker.spot_ilk_pip, new_pip text) RETURNS maker.spot_ilk_pip
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_pip_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.spot_ilk_pip
                 LEFT JOIN public.headers ON spot_ilk_pip.header_id = headers.id
        WHERE spot_ilk_pip.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET pip = new_pip
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_pip_diff_block IS NULL
        OR ilk_snapshot.block_number < next_pip_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_pips_until_next_diff(start_at_diff maker.spot_ilk_pip, new_pip text); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_pips_until_next_diff(start_at_diff maker.spot_ilk_pip, new_pip text) IS '@omit';


--
-- Name: update_rates_until_next_diff(maker.vat_ilk_rate, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_rates_until_next_diff(start_at_diff maker.vat_ilk_rate, new_rate numeric) RETURNS maker.vat_ilk_rate
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_rate_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_ilk_rate
                 LEFT JOIN public.headers ON vat_ilk_rate.header_id = headers.id
        WHERE vat_ilk_rate.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET rate = new_rate
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_rate_diff_block IS NULL
        OR ilk_snapshot.block_number < next_rate_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_rates_until_next_diff(start_at_diff maker.vat_ilk_rate, new_rate numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_rates_until_next_diff(start_at_diff maker.vat_ilk_rate, new_rate numeric) IS '@omit';


--
-- Name: update_rhos_until_next_diff(maker.jug_ilk_rho, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_rhos_until_next_diff(start_at_diff maker.jug_ilk_rho, new_rho numeric) RETURNS maker.jug_ilk_rho
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_rho_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.jug_ilk_rho
                 LEFT JOIN public.headers ON jug_ilk_rho.header_id = headers.id
        WHERE jug_ilk_rho.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET rho = new_rho
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_rho_diff_block IS NULL
        OR ilk_snapshot.block_number < next_rho_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_rhos_until_next_diff(start_at_diff maker.jug_ilk_rho, new_rho numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_rhos_until_next_diff(start_at_diff maker.jug_ilk_rho, new_rho numeric) IS '@omit';


--
-- Name: update_spots_until_next_diff(maker.vat_ilk_spot, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_spots_until_next_diff(start_at_diff maker.vat_ilk_spot, new_spot numeric) RETURNS maker.vat_ilk_spot
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_ilk_identifier  TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_spot_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_ilk_spot
                 LEFT JOIN public.headers ON vat_ilk_spot.header_id = headers.id
        WHERE vat_ilk_spot.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET spot = new_spot
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_spot_diff_block IS NULL
        OR ilk_snapshot.block_number < next_spot_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_spots_until_next_diff(start_at_diff maker.vat_ilk_spot, new_spot numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_spots_until_next_diff(start_at_diff maker.vat_ilk_spot, new_spot numeric) IS '@omit';


--
-- Name: update_time_created(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_time_created() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_time_created(NEW);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.clear_time_created(OLD);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_time_created(); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_time_created() IS '@omit';


--
-- Name: update_urn_arts(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_urn_arts() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_urn_art(NEW);
        PERFORM maker.update_urn_arts_until_next_diff(NEW, NEW.art);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_urn_arts_until_next_diff(OLD, urn_art_before_block(OLD.urn_id, OLD.header_id));
        PERFORM maker.delete_obsolete_urn_snapshot(OLD.urn_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_urn_arts_until_next_diff(maker.vat_urn_art, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_urn_arts_until_next_diff(start_at_diff maker.vat_urn_art, new_art numeric) RETURNS maker.vat_urn_art
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_rate_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_urn_art
                 LEFT JOIN public.headers ON vat_urn_art.header_id = headers.id
        WHERE vat_urn_art.urn_id = start_at_diff.urn_id
          AND block_number > diff_block_number);
BEGIN
    WITH urn AS (
        SELECT urns.identifier AS urn_identifier, ilks.identifier AS ilk_identifier
        FROM maker.urns
                 LEFT JOIN maker.ilks ON urns.ilk_id = ilks.id
        WHERE urns.id = start_at_diff.urn_id)
    UPDATE api.urn_snapshot
    SET art = new_art
    FROM urn
    WHERE urn_snapshot.urn_identifier = urn.urn_identifier
      AND urn_snapshot.ilk_identifier = urn.ilk_identifier
      AND urn_snapshot.block_height >= diff_block_number
      AND (next_rate_diff_block IS NULL
        OR urn_snapshot.block_height < next_rate_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_urn_arts_until_next_diff(start_at_diff maker.vat_urn_art, new_art numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_urn_arts_until_next_diff(start_at_diff maker.vat_urn_art, new_art numeric) IS '@omit';


--
-- Name: update_urn_created(integer); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_urn_created(urn_id integer) RETURNS maker.vat_urn_ink
    LANGUAGE plpgsql
    AS $$
BEGIN
    WITH utc AS (select urn_time_created(urn_id) as utc)
    UPDATE api.urn_snapshot
    SET created = (SELECT utc FROM utc)
    FROM maker.urns
             LEFT JOIN maker.ilks ON urns.ilk_id = ilks.id
    WHERE urns.identifier = urn_snapshot.urn_identifier
      AND ilks.identifier = urn_snapshot.ilk_identifier
      AND urns.id = urn_id;
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_urn_created(urn_id integer); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_urn_created(urn_id integer) IS '@omit';


--
-- Name: update_urn_inks(); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_urn_inks() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_urn_ink(NEW);
        PERFORM maker.update_urn_inks_until_next_diff(NEW, NEW.ink);
        PERFORM maker.update_urn_created(NEW.urn_id);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_urn_inks_until_next_diff(OLD, urn_ink_before_block(OLD.urn_id, OLD.header_id));
        PERFORM maker.delete_obsolete_urn_snapshot(OLD.urn_id, OLD.header_id);
        PERFORM maker.update_urn_created(OLD.urn_id);
    END IF;
    RETURN NULL;
END
$$;


--
-- Name: update_urn_inks_until_next_diff(maker.vat_urn_ink, numeric); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.update_urn_inks_until_next_diff(start_at_diff maker.vat_urn_ink, new_ink numeric) RETURNS maker.vat_urn_ink
    LANGUAGE plpgsql
    AS $$
DECLARE
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_rate_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_urn_ink
                 LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
        WHERE vat_urn_ink.urn_id = start_at_diff.urn_id
          AND block_number > diff_block_number);
BEGIN
    WITH urn AS (
        SELECT urns.identifier AS urn_identifier, ilks.identifier AS ilk_identifier
        FROM maker.urns
                 LEFT JOIN maker.ilks ON urns.ilk_id = ilks.id
        WHERE urns.id = start_at_diff.urn_id)
    UPDATE api.urn_snapshot
    SET ink = new_ink
    FROM urn
    WHERE urn_snapshot.urn_identifier = urn.urn_identifier
      AND urn_snapshot.ilk_identifier = urn.ilk_identifier
      AND urn_snapshot.block_height >= diff_block_number
      AND (next_rate_diff_block IS NULL
        OR urn_snapshot.block_height < next_rate_diff_block);
    RETURN NULL;
END
$$;


--
-- Name: FUNCTION update_urn_inks_until_next_diff(start_at_diff maker.vat_urn_ink, new_ink numeric); Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON FUNCTION maker.update_urn_inks_until_next_diff(start_at_diff maker.vat_urn_ink, new_ink numeric) IS '@omit';


--
-- Name: create_back_filled_diff(bigint, bytea, bytea, bytea, bytea, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.create_back_filled_diff(block_height bigint, block_hash bytea, hashed_address bytea, storage_key bytea, storage_value bytea, eth_node_id integer) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    last_storage_value  BYTEA := (
        SELECT storage_diff.storage_value
        FROM public.storage_diff
        WHERE storage_diff.block_height <= create_back_filled_diff.block_height
          AND storage_diff.hashed_address = create_back_filled_diff.hashed_address
          AND storage_diff.storage_key = create_back_filled_diff.storage_key
        ORDER BY storage_diff.block_height DESC
        LIMIT 1
    );
    empty_storage_value BYTEA := (
        SELECT '\x0000000000000000000000000000000000000000000000000000000000000000'::BYTEA
    );
BEGIN
    IF last_storage_value = create_back_filled_diff.storage_value THEN
        RETURN;
    END IF;

    IF last_storage_value is null and create_back_filled_diff.storage_value = empty_storage_value THEN
        RETURN;
    END IF;

    INSERT INTO public.storage_diff (block_height, block_hash, hashed_address, storage_key, storage_value,
                                     eth_node_id, from_backfill)
    VALUES (create_back_filled_diff.block_height, create_back_filled_diff.block_hash,
            create_back_filled_diff.hashed_address, create_back_filled_diff.storage_key,
            create_back_filled_diff.storage_value, create_back_filled_diff.eth_node_id, true)
    ON CONFLICT DO NOTHING;

    RETURN;
END
$$;


--
-- Name: FUNCTION create_back_filled_diff(block_height bigint, block_hash bytea, hashed_address bytea, storage_key bytea, storage_value bytea, eth_node_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.create_back_filled_diff(block_height bigint, block_hash bytea, hashed_address bytea, storage_key bytea, storage_value bytea, eth_node_id integer) IS '@omit';


--
-- Name: flap_bid_bid_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flap_bid_bid_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT bid
FROM maker.flap_bid_bid
         LEFT JOIN public.headers ON flap_bid_bid.header_id = headers.id
WHERE flap_bid_bid.bid_id = flap_bid_bid_before_block.bid_id
  AND flap_bid_bid.address_id = flap_bid_bid_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flap_bid_bid_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flap_bid_bid_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flap_bid_bid_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flap_bid_end_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flap_bid_end_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS bigint
    LANGUAGE sql
    AS $$
SELECT "end"
FROM maker.flap_bid_end
         LEFT JOIN public.headers ON flap_bid_end.header_id = headers.id
WHERE flap_bid_end.bid_id = flap_bid_end_before_block.bid_id
  AND flap_bid_end.address_id = flap_bid_end_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flap_bid_end_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flap_bid_end_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flap_bid_end_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flap_bid_guy_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flap_bid_guy_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS text
    LANGUAGE sql
    AS $$
SELECT guy
FROM maker.flap_bid_guy
         LEFT JOIN public.headers ON flap_bid_guy.header_id = headers.id
WHERE flap_bid_guy.bid_id = flap_bid_guy_before_block.bid_id
  AND flap_bid_guy.address_id = flap_bid_guy_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flap_bid_guy_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flap_bid_guy_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flap_bid_guy_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flap_bid_lot_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flap_bid_lot_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT lot
FROM maker.flap_bid_lot
         LEFT JOIN public.headers ON flap_bid_lot.header_id = headers.id
WHERE flap_bid_lot.bid_id = flap_bid_lot_before_block.bid_id
  AND flap_bid_lot.address_id = flap_bid_lot_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flap_bid_lot_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flap_bid_lot_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flap_bid_lot_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flap_bid_tic_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flap_bid_tic_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS bigint
    LANGUAGE sql
    AS $$
SELECT tic
FROM maker.flap_bid_tic
         LEFT JOIN public.headers ON flap_bid_tic.header_id = headers.id
WHERE flap_bid_tic.bid_id = flap_bid_tic_before_block.bid_id
  AND flap_bid_tic.address_id = flap_bid_tic_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flap_bid_tic_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flap_bid_tic_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flap_bid_tic_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flap_bid_time_created(integer, numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flap_bid_time_created(address_id integer, bid_id numeric) RETURNS timestamp without time zone
    LANGUAGE sql
    AS $$
SELECT api.epoch_to_datetime(MIN(block_timestamp))
FROM public.headers
         LEFT JOIN maker.flap_kick ON flap_kick.header_id = headers.id
WHERE flap_kick.address_id = flap_bid_time_created.address_id
  AND flap_kick.bid_id = flap_bid_time_created.bid_id
$$;


--
-- Name: FUNCTION flap_bid_time_created(address_id integer, bid_id numeric); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flap_bid_time_created(address_id integer, bid_id numeric) IS '@omit';


--
-- Name: flip_bid_bid_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flip_bid_bid_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT bid
FROM maker.flip_bid_bid
         LEFT JOIN public.headers ON flip_bid_bid.header_id = headers.id
WHERE flip_bid_bid.bid_id = flip_bid_bid_before_block.bid_id
  AND flip_bid_bid.address_id = flip_bid_bid_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_bid_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flip_bid_bid_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flip_bid_bid_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flip_bid_end_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flip_bid_end_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS bigint
    LANGUAGE sql
    AS $$
SELECT "end"
FROM maker.flip_bid_end
         LEFT JOIN public.headers ON flip_bid_end.header_id = headers.id
WHERE flip_bid_end.bid_id = flip_bid_end_before_block.bid_id
  AND flip_bid_end.address_id = flip_bid_end_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_end_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flip_bid_end_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flip_bid_end_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flip_bid_gal_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flip_bid_gal_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS text
    LANGUAGE sql
    AS $$
SELECT gal
FROM maker.flip_bid_gal
         LEFT JOIN public.headers ON flip_bid_gal.header_id = headers.id
WHERE flip_bid_gal.bid_id = flip_bid_gal_before_block.bid_id
  AND flip_bid_gal.address_id = flip_bid_gal_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_gal_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flip_bid_gal_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flip_bid_gal_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flip_bid_guy_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flip_bid_guy_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS text
    LANGUAGE sql
    AS $$
SELECT guy
FROM maker.flip_bid_guy
         LEFT JOIN public.headers ON flip_bid_guy.header_id = headers.id
WHERE flip_bid_guy.bid_id = flip_bid_guy_before_block.bid_id
  AND flip_bid_guy.address_id = flip_bid_guy_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_guy_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flip_bid_guy_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flip_bid_guy_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flip_bid_lot_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flip_bid_lot_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT lot
FROM maker.flip_bid_lot
         LEFT JOIN public.headers ON flip_bid_lot.header_id = headers.id
WHERE flip_bid_lot.bid_id = flip_bid_lot_before_block.bid_id
  AND flip_bid_lot.address_id = flip_bid_lot_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_lot_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flip_bid_lot_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flip_bid_lot_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flip_bid_tab_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flip_bid_tab_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT tab
FROM maker.flip_bid_tab
         LEFT JOIN public.headers ON flip_bid_tab.header_id = headers.id
WHERE flip_bid_tab.bid_id = flip_bid_tab_before_block.bid_id
  AND flip_bid_tab.address_id = flip_bid_tab_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_tab_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flip_bid_tab_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flip_bid_tab_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flip_bid_tic_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flip_bid_tic_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS bigint
    LANGUAGE sql
    AS $$
SELECT tic
FROM maker.flip_bid_tic
         LEFT JOIN public.headers ON flip_bid_tic.header_id = headers.id
WHERE flip_bid_tic.bid_id = flip_bid_tic_before_block.bid_id
  AND flip_bid_tic.address_id = flip_bid_tic_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_tic_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flip_bid_tic_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flip_bid_tic_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flip_bid_time_created(integer, numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flip_bid_time_created(address_id integer, bid_id numeric) RETURNS timestamp without time zone
    LANGUAGE sql
    AS $$
SELECT api.epoch_to_datetime(MIN(block_timestamp))
FROM public.headers
         LEFT JOIN maker.flip_kick ON flip_kick.header_id = headers.id
WHERE flip_kick.address_id = flip_bid_time_created.address_id
  AND flip_kick.bid_id = flip_bid_time_created.bid_id
$$;


--
-- Name: FUNCTION flip_bid_time_created(address_id integer, bid_id numeric); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flip_bid_time_created(address_id integer, bid_id numeric) IS '@omit';


--
-- Name: flip_bid_usr_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flip_bid_usr_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS text
    LANGUAGE sql
    AS $$
SELECT usr
FROM maker.flip_bid_usr
         LEFT JOIN public.headers ON flip_bid_usr.header_id = headers.id
WHERE flip_bid_usr.bid_id = flip_bid_usr_before_block.bid_id
  AND flip_bid_usr.address_id = flip_bid_usr_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flip_bid_usr_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flip_bid_usr_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flip_bid_usr_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flop_bid_bid_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flop_bid_bid_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT bid
FROM maker.flop_bid_bid
         LEFT JOIN public.headers ON flop_bid_bid.header_id = headers.id
WHERE flop_bid_bid.bid_id = flop_bid_bid_before_block.bid_id
  AND flop_bid_bid.address_id = flop_bid_bid_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flop_bid_bid_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flop_bid_bid_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flop_bid_bid_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flop_bid_end_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flop_bid_end_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS bigint
    LANGUAGE sql
    AS $$
SELECT "end"
FROM maker.flop_bid_end
         LEFT JOIN public.headers ON flop_bid_end.header_id = headers.id
WHERE flop_bid_end.bid_id = flop_bid_end_before_block.bid_id
  AND flop_bid_end.address_id = flop_bid_end_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flop_bid_end_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flop_bid_end_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flop_bid_end_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flop_bid_guy_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flop_bid_guy_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS text
    LANGUAGE sql
    AS $$
SELECT guy
FROM maker.flop_bid_guy
         LEFT JOIN public.headers ON flop_bid_guy.header_id = headers.id
WHERE flop_bid_guy.bid_id = flop_bid_guy_before_block.bid_id
  AND flop_bid_guy.address_id = flop_bid_guy_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flop_bid_guy_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flop_bid_guy_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flop_bid_guy_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flop_bid_lot_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flop_bid_lot_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
SELECT lot
FROM maker.flop_bid_lot
         LEFT JOIN public.headers ON flop_bid_lot.header_id = headers.id
WHERE flop_bid_lot.bid_id = flop_bid_lot_before_block.bid_id
  AND flop_bid_lot.address_id = flop_bid_lot_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flop_bid_lot_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flop_bid_lot_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flop_bid_lot_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flop_bid_tic_before_block(numeric, integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flop_bid_tic_before_block(bid_id numeric, address_id integer, header_id integer) RETURNS bigint
    LANGUAGE sql
    AS $$
SELECT tic
FROM maker.flop_bid_tic
         LEFT JOIN public.headers ON flop_bid_tic.header_id = headers.id
WHERE flop_bid_tic.bid_id = flop_bid_tic_before_block.bid_id
  AND flop_bid_tic.address_id = flop_bid_tic_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = flop_bid_tic_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION flop_bid_tic_before_block(bid_id numeric, address_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flop_bid_tic_before_block(bid_id numeric, address_id integer, header_id integer) IS '@omit';


--
-- Name: flop_bid_time_created(integer, numeric); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.flop_bid_time_created(address_id integer, bid_id numeric) RETURNS timestamp without time zone
    LANGUAGE sql
    AS $$
SELECT api.epoch_to_datetime(MIN(block_timestamp))
FROM public.headers
         LEFT JOIN maker.flop_kick ON flop_kick.header_id = headers.id
WHERE flop_kick.address_id = flop_bid_time_created.address_id
  AND flop_kick.bid_id = flop_bid_time_created.bid_id
$$;


--
-- Name: FUNCTION flop_bid_time_created(address_id integer, bid_id numeric); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.flop_bid_time_created(address_id integer, bid_id numeric) IS '@omit';


--
-- Name: get_or_create_header(bigint, character varying, jsonb, numeric, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_or_create_header(block_number bigint, hash character varying, raw jsonb, block_timestamp numeric, eth_node_id integer) RETURNS integer
    LANGUAGE plpgsql
    AS $$
DECLARE
    matching_header_id    INTEGER := (
        SELECT id
        FROM public.headers
        WHERE headers.block_number = get_or_create_header.block_number
          AND headers.hash = get_or_create_header.hash
    );
    nonmatching_header_id INTEGER := (
        SELECT id
        FROM public.headers
        WHERE headers.block_number = get_or_create_header.block_number
          AND headers.hash != get_or_create_header.hash
    );
    max_block_number      BIGINT  := (
        SELECT MAX(headers.block_number)
        FROM public.headers
    );
    inserted_header_id    INTEGER;
BEGIN
    IF matching_header_id != 0 THEN
        RETURN matching_header_id;
    END IF;

    IF nonmatching_header_id != 0 AND block_number <= max_block_number - 15 THEN
        RETURN nonmatching_header_id;
    END IF;

    IF nonmatching_header_id != 0 AND block_number > max_block_number - 15 THEN
        DELETE FROM public.headers WHERE id = nonmatching_header_id;
    END IF;

    INSERT INTO public.headers (hash, block_number, raw, block_timestamp, eth_node_id)
    VALUES (get_or_create_header.hash, get_or_create_header.block_number, get_or_create_header.raw,
            get_or_create_header.block_timestamp, get_or_create_header.eth_node_id)
    RETURNING id INTO inserted_header_id;

    RETURN inserted_header_id;
END
$$;


--
-- Name: FUNCTION get_or_create_header(block_number bigint, hash character varying, raw jsonb, block_timestamp numeric, eth_node_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.get_or_create_header(block_number bigint, hash character varying, raw jsonb, block_timestamp numeric, eth_node_id integer) IS '@omit';


--
-- Name: get_tx_data(bigint, bigint); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_tx_data(block_height bigint, log_id bigint) RETURNS SETOF api.tx
    LANGUAGE sql STABLE
    AS $$
SELECT txs.hash, txs.tx_index, headers.block_number, headers.hash, tx_from, tx_to
FROM public.transactions txs
         LEFT JOIN public.headers ON txs.header_id = headers.id
         LEFT JOIN public.event_logs ON txs.tx_index = event_logs.tx_index
WHERE headers.block_number = block_height
  AND event_logs.id = log_id
ORDER BY block_number DESC

$$;


--
-- Name: ilk_art_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_art_before_block(ilk_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT art
FROM maker.vat_ilk_art
         LEFT JOIN public.headers ON vat_ilk_art.header_id = headers.id
WHERE vat_ilk_art.ilk_id = ilk_art_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION ilk_art_before_block(ilk_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_art_before_block(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: ilk_chop_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_chop_before_block(ilk_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT chop
FROM maker.cat_ilk_chop
         LEFT JOIN public.headers ON cat_ilk_chop.header_id = headers.id
WHERE cat_ilk_chop.ilk_id = ilk_chop_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION ilk_chop_before_block(ilk_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_chop_before_block(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: ilk_dust_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_dust_before_block(ilk_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT dust
FROM maker.vat_ilk_dust
         LEFT JOIN public.headers ON vat_ilk_dust.header_id = headers.id
WHERE vat_ilk_dust.ilk_id = ilk_dust_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION ilk_dust_before_block(ilk_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_dust_before_block(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: ilk_duty_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_duty_before_block(ilk_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT duty
FROM maker.jug_ilk_duty
         LEFT JOIN public.headers ON jug_ilk_duty.header_id = headers.id
WHERE jug_ilk_duty.ilk_id = ilk_duty_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION ilk_duty_before_block(ilk_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_duty_before_block(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: ilk_flip_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_flip_before_block(ilk_id integer, header_id integer) RETURNS text
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT flip
FROM maker.cat_ilk_flip
         LEFT JOIN public.headers ON cat_ilk_flip.header_id = headers.id
WHERE cat_ilk_flip.ilk_id = ilk_flip_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION ilk_flip_before_block(ilk_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_flip_before_block(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: ilk_line_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_line_before_block(ilk_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT line
FROM maker.vat_ilk_line
         LEFT JOIN public.headers ON vat_ilk_line.header_id = headers.id
WHERE vat_ilk_line.ilk_id = ilk_line_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION ilk_line_before_block(ilk_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_line_before_block(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: ilk_lump_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_lump_before_block(ilk_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT lump
FROM maker.cat_ilk_lump
         LEFT JOIN public.headers ON cat_ilk_lump.header_id = headers.id
WHERE cat_ilk_lump.ilk_id = ilk_lump_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION ilk_lump_before_block(ilk_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_lump_before_block(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: ilk_mat_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_mat_before_block(ilk_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT mat
FROM maker.spot_ilk_mat
         LEFT JOIN public.headers ON spot_ilk_mat.header_id = headers.id
WHERE spot_ilk_mat.ilk_id = ilk_mat_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION ilk_mat_before_block(ilk_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_mat_before_block(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: ilk_pip_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_pip_before_block(ilk_id integer, header_id integer) RETURNS text
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT pip
FROM maker.spot_ilk_pip
         LEFT JOIN public.headers ON spot_ilk_pip.header_id = headers.id
WHERE spot_ilk_pip.ilk_id = ilk_pip_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION ilk_pip_before_block(ilk_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_pip_before_block(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: ilk_rate_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_rate_before_block(ilk_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT rate
FROM maker.vat_ilk_rate
         LEFT JOIN public.headers ON vat_ilk_rate.header_id = headers.id
WHERE vat_ilk_rate.ilk_id = ilk_rate_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION ilk_rate_before_block(ilk_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_rate_before_block(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: ilk_rho_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_rho_before_block(ilk_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$

WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT rho
FROM maker.jug_ilk_rho
         LEFT JOIN public.headers ON jug_ilk_rho.header_id = headers.id
WHERE jug_ilk_rho.ilk_id = ilk_rho_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION ilk_rho_before_block(ilk_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_rho_before_block(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: ilk_spot_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_spot_before_block(ilk_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT spot
FROM maker.vat_ilk_spot
         LEFT JOIN public.headers ON vat_ilk_spot.header_id = headers.id
WHERE vat_ilk_spot.ilk_id = ilk_spot_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION ilk_spot_before_block(ilk_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_spot_before_block(ilk_id integer, header_id integer) IS '@omit';


--
-- Name: ilk_time_created(integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.ilk_time_created(ilk_id integer) RETURNS timestamp without time zone
    LANGUAGE sql
    AS $$
SELECT api.epoch_to_datetime(MIN(block_timestamp))
FROM public.headers
         LEFT JOIN maker.vat_init ON vat_init.header_id = headers.id
WHERE vat_init.ilk_id = ilk_time_created.ilk_id
$$;


--
-- Name: FUNCTION ilk_time_created(ilk_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.ilk_time_created(ilk_id integer) IS '@omit';


--
-- Name: set_header_updated(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_header_updated() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated = NOW();
    RETURN NEW;
END;
$$;


--
-- Name: urn_art_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.urn_art_before_block(urn_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id)
SELECT art
FROM maker.vat_urn_art
         LEFT JOIN public.headers ON vat_urn_art.header_id = headers.id
WHERE vat_urn_art.urn_id = urn_art_before_block.urn_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION urn_art_before_block(urn_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.urn_art_before_block(urn_id integer, header_id integer) IS '@omit';


--
-- Name: urn_ink_before_block(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.urn_ink_before_block(urn_id integer, header_id integer) RETURNS numeric
    LANGUAGE sql
    AS $$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id)
SELECT ink
FROM maker.vat_urn_ink
         LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
WHERE vat_urn_ink.urn_id = urn_ink_before_block.urn_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$;


--
-- Name: FUNCTION urn_ink_before_block(urn_id integer, header_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.urn_ink_before_block(urn_id integer, header_id integer) IS '@omit';


--
-- Name: urn_time_created(integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.urn_time_created(urn_id integer) RETURNS timestamp without time zone
    LANGUAGE sql
    AS $$
SELECT api.epoch_to_datetime(MIN(block_timestamp))
FROM maker.vat_urn_ink
         LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
WHERE vat_urn_ink.urn_id = urn_time_created.urn_id
$$;


--
-- Name: FUNCTION urn_time_created(urn_id integer); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.urn_time_created(urn_id integer) IS '@omit';


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
-- Name: auction_file; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.auction_file (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    address_id integer NOT NULL,
    msg_sender integer NOT NULL,
    what text,
    data numeric
);


--
-- Name: auction_file_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.auction_file_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: auction_file_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.auction_file_id_seq OWNED BY maker.auction_file.id;


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
    bid_id numeric,
    address_id integer NOT NULL
);


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
    msg_sender integer NOT NULL,
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
    msg_sender integer NOT NULL,
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
    msg_sender integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    address_id integer NOT NULL,
    msg_sender integer NOT NULL,
    bid_id numeric NOT NULL
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
    msg_sender integer NOT NULL,
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
-- Name: deny; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.deny (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    address_id integer NOT NULL,
    msg_sender integer NOT NULL,
    usr integer NOT NULL
);


--
-- Name: deny_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.deny_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: deny_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.deny_id_seq OWNED BY maker.deny.id;


--
-- Name: flap; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap (
    address_id integer NOT NULL,
    block_number bigint NOT NULL,
    bid_id numeric NOT NULL,
    guy text,
    tic bigint,
    "end" bigint,
    lot numeric,
    bid numeric,
    created timestamp without time zone,
    updated timestamp without time zone NOT NULL
);


--
-- Name: flap_beg; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_beg (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    kicks numeric NOT NULL
);


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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    address_id integer NOT NULL,
    block_number bigint NOT NULL,
    bid_id numeric NOT NULL,
    guy text,
    tic bigint,
    "end" bigint,
    lot numeric,
    bid numeric,
    usr text,
    gal text,
    tab numeric,
    created timestamp without time zone,
    updated timestamp without time zone NOT NULL
);


--
-- Name: flip_beg; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_beg (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    kicks numeric NOT NULL
);


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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    address_id integer NOT NULL,
    block_number bigint NOT NULL,
    bid_id numeric NOT NULL,
    guy text,
    tic bigint,
    "end" bigint,
    lot numeric,
    bid numeric,
    created timestamp without time zone,
    updated timestamp without time zone NOT NULL
);


--
-- Name: flop_beg; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_beg (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    kicks numeric NOT NULL
);


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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
-- Name: flop_vow; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_vow (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    vow text
);


--
-- Name: flop_vow_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.flop_vow_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: flop_vow_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.flop_vow_id_seq OWNED BY maker.flop_vow.id;


--
-- Name: ilks; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.ilks (
    id integer NOT NULL,
    ilk text NOT NULL,
    identifier text NOT NULL
);


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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    msg_sender integer NOT NULL,
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
    msg_sender integer NOT NULL,
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
    msg_sender integer NOT NULL,
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
    msg_sender integer NOT NULL,
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
    msg_sender integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
-- Name: log_bump; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_bump (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    offer_id numeric,
    pair character varying(66),
    maker integer NOT NULL,
    pay_gem integer NOT NULL,
    buy_gem integer NOT NULL,
    pay_amt numeric,
    buy_amt numeric,
    "timestamp" integer
);


--
-- Name: log_bump_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_bump_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_bump_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_bump_id_seq OWNED BY maker.log_bump.id;


--
-- Name: log_buy_enabled; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_buy_enabled (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    is_enabled boolean
);


--
-- Name: log_buy_enabled_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_buy_enabled_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_buy_enabled_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_buy_enabled_id_seq OWNED BY maker.log_buy_enabled.id;


--
-- Name: log_delete; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_delete (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    keeper integer NOT NULL,
    offer_id numeric
);


--
-- Name: log_delete_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_delete_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_delete_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_delete_id_seq OWNED BY maker.log_delete.id;


--
-- Name: log_insert; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_insert (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    keeper integer NOT NULL,
    offer_id numeric
);


--
-- Name: log_insert_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_insert_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_insert_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_insert_id_seq OWNED BY maker.log_insert.id;


--
-- Name: log_item_update; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_item_update (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    offer_id numeric
);


--
-- Name: log_item_update_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_item_update_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_item_update_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_item_update_id_seq OWNED BY maker.log_item_update.id;


--
-- Name: log_kill; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_kill (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    offer_id numeric,
    pair character varying(66),
    maker integer NOT NULL,
    pay_gem integer NOT NULL,
    buy_gem integer NOT NULL,
    pay_amt numeric,
    buy_amt numeric,
    "timestamp" bigint
);


--
-- Name: log_kill_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_kill_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_kill_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_kill_id_seq OWNED BY maker.log_kill.id;


--
-- Name: log_make; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_make (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    offer_id numeric,
    pair character varying(66),
    maker integer NOT NULL,
    pay_gem integer NOT NULL,
    buy_gem integer NOT NULL,
    pay_amt numeric,
    buy_amt numeric,
    "timestamp" integer
);


--
-- Name: log_make_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_make_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_make_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_make_id_seq OWNED BY maker.log_make.id;


--
-- Name: log_matching_enabled; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_matching_enabled (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    is_enabled boolean
);


--
-- Name: log_matching_enabled_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_matching_enabled_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_matching_enabled_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_matching_enabled_id_seq OWNED BY maker.log_matching_enabled.id;


--
-- Name: log_median_price; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_median_price (
    id integer NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    log_id bigint NOT NULL,
    val numeric,
    age numeric
);


--
-- Name: log_median_price_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_median_price_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_median_price_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_median_price_id_seq OWNED BY maker.log_median_price.id;


--
-- Name: log_min_sell; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_min_sell (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    pay_gem integer NOT NULL,
    min_amount numeric
);


--
-- Name: log_min_sell_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_min_sell_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_min_sell_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_min_sell_id_seq OWNED BY maker.log_min_sell.id;


--
-- Name: log_sorted_offer; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_sorted_offer (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    offer_id numeric
);


--
-- Name: log_sorted_offer_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_sorted_offer_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_sorted_offer_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_sorted_offer_id_seq OWNED BY maker.log_sorted_offer.id;


--
-- Name: log_take; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_take (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    offer_id numeric,
    pair character varying(66),
    maker integer NOT NULL,
    pay_gem integer NOT NULL,
    buy_gem integer NOT NULL,
    taker integer NOT NULL,
    take_amt numeric,
    give_amt numeric,
    "timestamp" integer
);


--
-- Name: log_take_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_take_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_take_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_take_id_seq OWNED BY maker.log_take.id;


--
-- Name: log_trade; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_trade (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    pay_gem integer NOT NULL,
    buy_gem integer NOT NULL,
    pay_amt numeric,
    buy_amt numeric
);


--
-- Name: log_trade_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_trade_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_trade_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_trade_id_seq OWNED BY maker.log_trade.id;


--
-- Name: log_unsorted_offer; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_unsorted_offer (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    offer_id numeric
);


--
-- Name: log_unsorted_offer_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_unsorted_offer_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_unsorted_offer_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_unsorted_offer_id_seq OWNED BY maker.log_unsorted_offer.id;


--
-- Name: log_value; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.log_value (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    address_id bigint NOT NULL,
    val numeric
);


--
-- Name: log_value_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.log_value_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_value_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.log_value_id_seq OWNED BY maker.log_value.id;


--
-- Name: median_age; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.median_age (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    age numeric NOT NULL
);


--
-- Name: median_age_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.median_age_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: median_age_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.median_age_id_seq OWNED BY maker.median_age.id;


--
-- Name: median_bar; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.median_bar (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    bar numeric NOT NULL
);


--
-- Name: median_bar_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.median_bar_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: median_bar_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.median_bar_id_seq OWNED BY maker.median_bar.id;


--
-- Name: median_bud; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.median_bud (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    a integer NOT NULL,
    bud integer NOT NULL
);


--
-- Name: median_bud_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.median_bud_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: median_bud_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.median_bud_id_seq OWNED BY maker.median_bud.id;


--
-- Name: median_diss_batch; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.median_diss_batch (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    msg_sender integer NOT NULL,
    a_length integer NOT NULL,
    a text[] NOT NULL
);


--
-- Name: median_diss_batch_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.median_diss_batch_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: median_diss_batch_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.median_diss_batch_id_seq OWNED BY maker.median_diss_batch.id;


--
-- Name: median_diss_single; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.median_diss_single (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    msg_sender integer NOT NULL,
    a integer NOT NULL
);


--
-- Name: median_diss_single_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.median_diss_single_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: median_diss_single_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.median_diss_single_id_seq OWNED BY maker.median_diss_single.id;


--
-- Name: median_drop; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.median_drop (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    msg_sender integer NOT NULL,
    a_length integer NOT NULL,
    a text[] NOT NULL
);


--
-- Name: median_drop_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.median_drop_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: median_drop_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.median_drop_id_seq OWNED BY maker.median_drop.id;


--
-- Name: median_kiss_batch; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.median_kiss_batch (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    msg_sender integer NOT NULL,
    a_length integer NOT NULL,
    a text[] NOT NULL
);


--
-- Name: median_kiss_batch_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.median_kiss_batch_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: median_kiss_batch_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.median_kiss_batch_id_seq OWNED BY maker.median_kiss_batch.id;


--
-- Name: median_kiss_single; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.median_kiss_single (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    msg_sender integer NOT NULL,
    a integer NOT NULL
);


--
-- Name: median_kiss_single_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.median_kiss_single_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: median_kiss_single_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.median_kiss_single_id_seq OWNED BY maker.median_kiss_single.id;


--
-- Name: median_lift; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.median_lift (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    msg_sender integer NOT NULL,
    a_length integer NOT NULL,
    a text[] NOT NULL
);


--
-- Name: median_lift_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.median_lift_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: median_lift_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.median_lift_id_seq OWNED BY maker.median_lift.id;


--
-- Name: median_orcl; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.median_orcl (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    a integer NOT NULL,
    orcl integer NOT NULL
);


--
-- Name: median_orcl_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.median_orcl_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: median_orcl_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.median_orcl_id_seq OWNED BY maker.median_orcl.id;


--
-- Name: median_slot; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.median_slot (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    slot_id integer NOT NULL,
    slot integer NOT NULL
);


--
-- Name: median_slot_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.median_slot_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: median_slot_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.median_slot_id_seq OWNED BY maker.median_slot.id;


--
-- Name: median_val; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.median_val (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    val numeric NOT NULL
);


--
-- Name: median_val_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.median_val_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: median_val_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.median_val_id_seq OWNED BY maker.median_val.id;


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
-- Name: osm_change; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.osm_change (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    msg_sender integer NOT NULL,
    src integer NOT NULL
);


--
-- Name: osm_change_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.osm_change_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: osm_change_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.osm_change_id_seq OWNED BY maker.osm_change.id;


--
-- Name: pot_cage; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_cage (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    msg_sender integer NOT NULL
);


--
-- Name: pot_cage_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_cage_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_cage_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_cage_id_seq OWNED BY maker.pot_cage.id;


--
-- Name: pot_chi; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_chi (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    chi numeric NOT NULL
);


--
-- Name: pot_chi_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_chi_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_chi_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_chi_id_seq OWNED BY maker.pot_chi.id;


--
-- Name: pot_drip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_drip (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    msg_sender integer NOT NULL
);


--
-- Name: pot_drip_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_drip_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_drip_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_drip_id_seq OWNED BY maker.pot_drip.id;


--
-- Name: pot_dsr; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_dsr (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    dsr numeric NOT NULL
);


--
-- Name: pot_dsr_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_dsr_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_dsr_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_dsr_id_seq OWNED BY maker.pot_dsr.id;


--
-- Name: pot_exit; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_exit (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    msg_sender integer NOT NULL,
    wad numeric
);


--
-- Name: pot_exit_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_exit_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_exit_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_exit_id_seq OWNED BY maker.pot_exit.id;


--
-- Name: pot_file_dsr; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_file_dsr (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    msg_sender integer NOT NULL,
    what text,
    data numeric
);


--
-- Name: pot_file_dsr_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_file_dsr_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_file_dsr_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_file_dsr_id_seq OWNED BY maker.pot_file_dsr.id;


--
-- Name: pot_file_vow; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_file_vow (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    msg_sender integer NOT NULL,
    what text,
    data text
);


--
-- Name: pot_file_vow_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_file_vow_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_file_vow_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_file_vow_id_seq OWNED BY maker.pot_file_vow.id;


--
-- Name: pot_join; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_join (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    msg_sender integer NOT NULL,
    wad numeric
);


--
-- Name: pot_join_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_join_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_join_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_join_id_seq OWNED BY maker.pot_join.id;


--
-- Name: pot_live; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_live (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    live numeric NOT NULL
);


--
-- Name: pot_live_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_live_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_live_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_live_id_seq OWNED BY maker.pot_live.id;


--
-- Name: pot_pie; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_pie (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    pie numeric NOT NULL
);


--
-- Name: pot_pie_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_pie_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_pie_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_pie_id_seq OWNED BY maker.pot_pie.id;


--
-- Name: pot_rho; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_rho (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    rho numeric NOT NULL
);


--
-- Name: pot_rho_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_rho_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_rho_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_rho_id_seq OWNED BY maker.pot_rho.id;


--
-- Name: pot_user_pie; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_user_pie (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    "user" integer NOT NULL,
    pie numeric NOT NULL
);


--
-- Name: pot_user_pie_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_user_pie_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_user_pie_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_user_pie_id_seq OWNED BY maker.pot_user_pie.id;


--
-- Name: pot_vat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_vat (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    vat integer NOT NULL
);


--
-- Name: pot_vat_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_vat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_vat_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_vat_id_seq OWNED BY maker.pot_vat.id;


--
-- Name: pot_vow; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pot_vow (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    vow integer NOT NULL
);


--
-- Name: pot_vow_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pot_vow_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pot_vow_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pot_vow_id_seq OWNED BY maker.pot_vow.id;


--
-- Name: rely; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.rely (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    address_id integer NOT NULL,
    msg_sender integer NOT NULL,
    usr integer NOT NULL
);


--
-- Name: rely_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.rely_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: rely_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.rely_id_seq OWNED BY maker.rely.id;


--
-- Name: set_min_sell; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.set_min_sell (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    pay_gem integer NOT NULL,
    msg_sender integer NOT NULL,
    dust numeric
);


--
-- Name: set_min_sell_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.set_min_sell_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: set_min_sell_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.set_min_sell_id_seq OWNED BY maker.set_min_sell.id;


--
-- Name: spot_file_mat; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_file_mat (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL,
    msg_sender integer NOT NULL,
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
-- Name: spot_file_par; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_file_par (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    msg_sender integer NOT NULL,
    what text,
    data numeric
);


--
-- Name: spot_file_par_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.spot_file_par_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: spot_file_par_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.spot_file_par_id_seq OWNED BY maker.spot_file_par.id;


--
-- Name: spot_file_pip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_file_pip (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    ilk_id integer NOT NULL,
    msg_sender integer NOT NULL,
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
-- Name: spot_live; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_live (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    live numeric NOT NULL
);


--
-- Name: spot_live_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.spot_live_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: spot_live_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.spot_live_id_seq OWNED BY maker.spot_live.id;


--
-- Name: spot_par; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.spot_par (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    address_id integer NOT NULL,
    msg_sender integer NOT NULL
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
    address_id integer NOT NULL,
    msg_sender integer NOT NULL
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
-- Name: vat_deny; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_deny (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    usr integer NOT NULL
);


--
-- Name: vat_deny_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_deny_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_deny_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_deny_id_seq OWNED BY maker.vat_deny.id;


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
    ilk_id integer NOT NULL,
    u text NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
-- Name: vat_hope; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_hope (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    usr integer NOT NULL
);


--
-- Name: vat_hope_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_hope_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_hope_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_hope_id_seq OWNED BY maker.vat_hope.id;


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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
-- Name: vat_nope; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_nope (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    usr integer NOT NULL
);


--
-- Name: vat_nope_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_nope_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_nope_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_nope_id_seq OWNED BY maker.vat_nope.id;


--
-- Name: vat_rely; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_rely (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    usr integer NOT NULL
);


--
-- Name: vat_rely_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vat_rely_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vat_rely_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vat_rely_id_seq OWNED BY maker.vat_rely.id;


--
-- Name: vat_sin; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_sin (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    msg_sender integer NOT NULL,
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
-- Name: vow_file_auction_address; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_file_auction_address (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    msg_sender integer NOT NULL,
    what text,
    data integer NOT NULL
);


--
-- Name: vow_file_auction_address_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_file_auction_address_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_file_auction_address_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_file_auction_address_id_seq OWNED BY maker.vow_file_auction_address.id;


--
-- Name: vow_file_auction_attributes; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_file_auction_attributes (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    msg_sender integer NOT NULL,
    what text,
    data numeric
);


--
-- Name: vow_file_auction_attributes_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_file_auction_attributes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_file_auction_attributes_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_file_auction_attributes_id_seq OWNED BY maker.vow_file_auction_attributes.id;


--
-- Name: vow_flapper; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_flapper (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    msg_sender integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
-- Name: vow_heal; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_heal (
    id integer NOT NULL,
    log_id bigint NOT NULL,
    header_id integer NOT NULL,
    msg_sender integer NOT NULL,
    rad numeric
);


--
-- Name: vow_heal_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_heal_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_heal_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_heal_id_seq OWNED BY maker.vow_heal.id;


--
-- Name: vow_hump; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_hump (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
-- Name: vow_live; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_live (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    live numeric NOT NULL
);


--
-- Name: vow_live_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_live_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_live_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_live_id_seq OWNED BY maker.vow_live.id;


--
-- Name: vow_sin_integer; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_sin_integer (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
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
-- Name: wards; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.wards (
    id integer NOT NULL,
    diff_id bigint NOT NULL,
    header_id integer NOT NULL,
    address_id integer NOT NULL,
    usr integer NOT NULL,
    wards integer NOT NULL
);


--
-- Name: wards_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.wards_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: wards_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.wards_id_seq OWNED BY maker.wards.id;


--
-- Name: yank; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.yank (
    id integer NOT NULL,
    header_id integer NOT NULL,
    log_id bigint NOT NULL,
    address_id integer NOT NULL,
    msg_sender integer NOT NULL,
    bid_id numeric NOT NULL
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
    address character varying(42),
    hashed_address character varying(66)
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
-- Name: checked_headers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.checked_headers (
    id integer NOT NULL,
    header_id integer NOT NULL,
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
-- Name: eth_nodes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.eth_nodes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: eth_nodes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.eth_nodes_id_seq OWNED BY public.eth_nodes.id;


--
-- Name: event_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.event_logs (
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
-- Name: event_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.event_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: event_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.event_logs_id_seq OWNED BY public.event_logs.id;


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
-- Name: headers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.headers (
    id integer NOT NULL,
    hash character varying(66) NOT NULL,
    block_number bigint NOT NULL,
    raw jsonb,
    block_timestamp numeric,
    check_count integer DEFAULT 0 NOT NULL,
    eth_node_id integer NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    updated timestamp without time zone DEFAULT now() NOT NULL
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
-- Name: receipts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.receipts (
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
-- Name: receipts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.receipts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: receipts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.receipts_id_seq OWNED BY public.receipts.id;


--
-- Name: storage_diff; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.storage_diff (
    id bigint NOT NULL,
    block_height bigint,
    block_hash bytea,
    hashed_address bytea,
    storage_key bytea,
    storage_value bytea,
    eth_node_id integer NOT NULL,
    checked boolean DEFAULT false NOT NULL,
    from_backfill boolean DEFAULT false NOT NULL
);


--
-- Name: storage_diff_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.storage_diff_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: storage_diff_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.storage_diff_id_seq OWNED BY public.storage_diff.id;


--
-- Name: transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.transactions (
    id integer NOT NULL,
    header_id integer NOT NULL,
    hash character varying(66) NOT NULL,
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
-- Name: transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.transactions_id_seq OWNED BY public.transactions.id;


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
-- Name: auction_file id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.auction_file ALTER COLUMN id SET DEFAULT nextval('maker.auction_file_id_seq'::regclass);


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
-- Name: deny id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deny ALTER COLUMN id SET DEFAULT nextval('maker.deny_id_seq'::regclass);


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
-- Name: flop_vow id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vow ALTER COLUMN id SET DEFAULT nextval('maker.flop_vow_id_seq'::regclass);


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
-- Name: log_bump id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_bump ALTER COLUMN id SET DEFAULT nextval('maker.log_bump_id_seq'::regclass);


--
-- Name: log_buy_enabled id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_buy_enabled ALTER COLUMN id SET DEFAULT nextval('maker.log_buy_enabled_id_seq'::regclass);


--
-- Name: log_delete id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_delete ALTER COLUMN id SET DEFAULT nextval('maker.log_delete_id_seq'::regclass);


--
-- Name: log_insert id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_insert ALTER COLUMN id SET DEFAULT nextval('maker.log_insert_id_seq'::regclass);


--
-- Name: log_item_update id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_item_update ALTER COLUMN id SET DEFAULT nextval('maker.log_item_update_id_seq'::regclass);


--
-- Name: log_kill id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_kill ALTER COLUMN id SET DEFAULT nextval('maker.log_kill_id_seq'::regclass);


--
-- Name: log_make id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_make ALTER COLUMN id SET DEFAULT nextval('maker.log_make_id_seq'::regclass);


--
-- Name: log_matching_enabled id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_matching_enabled ALTER COLUMN id SET DEFAULT nextval('maker.log_matching_enabled_id_seq'::regclass);


--
-- Name: log_median_price id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_median_price ALTER COLUMN id SET DEFAULT nextval('maker.log_median_price_id_seq'::regclass);


--
-- Name: log_min_sell id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_min_sell ALTER COLUMN id SET DEFAULT nextval('maker.log_min_sell_id_seq'::regclass);


--
-- Name: log_sorted_offer id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_sorted_offer ALTER COLUMN id SET DEFAULT nextval('maker.log_sorted_offer_id_seq'::regclass);


--
-- Name: log_take id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_take ALTER COLUMN id SET DEFAULT nextval('maker.log_take_id_seq'::regclass);


--
-- Name: log_trade id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_trade ALTER COLUMN id SET DEFAULT nextval('maker.log_trade_id_seq'::regclass);


--
-- Name: log_unsorted_offer id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_unsorted_offer ALTER COLUMN id SET DEFAULT nextval('maker.log_unsorted_offer_id_seq'::regclass);


--
-- Name: log_value id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_value ALTER COLUMN id SET DEFAULT nextval('maker.log_value_id_seq'::regclass);


--
-- Name: median_age id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_age ALTER COLUMN id SET DEFAULT nextval('maker.median_age_id_seq'::regclass);


--
-- Name: median_bar id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bar ALTER COLUMN id SET DEFAULT nextval('maker.median_bar_id_seq'::regclass);


--
-- Name: median_bud id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bud ALTER COLUMN id SET DEFAULT nextval('maker.median_bud_id_seq'::regclass);


--
-- Name: median_diss_batch id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_batch ALTER COLUMN id SET DEFAULT nextval('maker.median_diss_batch_id_seq'::regclass);


--
-- Name: median_diss_single id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_single ALTER COLUMN id SET DEFAULT nextval('maker.median_diss_single_id_seq'::regclass);


--
-- Name: median_drop id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_drop ALTER COLUMN id SET DEFAULT nextval('maker.median_drop_id_seq'::regclass);


--
-- Name: median_kiss_batch id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_batch ALTER COLUMN id SET DEFAULT nextval('maker.median_kiss_batch_id_seq'::regclass);


--
-- Name: median_kiss_single id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_single ALTER COLUMN id SET DEFAULT nextval('maker.median_kiss_single_id_seq'::regclass);


--
-- Name: median_lift id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_lift ALTER COLUMN id SET DEFAULT nextval('maker.median_lift_id_seq'::regclass);


--
-- Name: median_orcl id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_orcl ALTER COLUMN id SET DEFAULT nextval('maker.median_orcl_id_seq'::regclass);


--
-- Name: median_slot id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_slot ALTER COLUMN id SET DEFAULT nextval('maker.median_slot_id_seq'::regclass);


--
-- Name: median_val id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_val ALTER COLUMN id SET DEFAULT nextval('maker.median_val_id_seq'::regclass);


--
-- Name: new_cdp id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.new_cdp ALTER COLUMN id SET DEFAULT nextval('maker.new_cdp_id_seq'::regclass);


--
-- Name: osm_change id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.osm_change ALTER COLUMN id SET DEFAULT nextval('maker.osm_change_id_seq'::regclass);


--
-- Name: pot_cage id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_cage ALTER COLUMN id SET DEFAULT nextval('maker.pot_cage_id_seq'::regclass);


--
-- Name: pot_chi id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_chi ALTER COLUMN id SET DEFAULT nextval('maker.pot_chi_id_seq'::regclass);


--
-- Name: pot_drip id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_drip ALTER COLUMN id SET DEFAULT nextval('maker.pot_drip_id_seq'::regclass);


--
-- Name: pot_dsr id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_dsr ALTER COLUMN id SET DEFAULT nextval('maker.pot_dsr_id_seq'::regclass);


--
-- Name: pot_exit id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_exit ALTER COLUMN id SET DEFAULT nextval('maker.pot_exit_id_seq'::regclass);


--
-- Name: pot_file_dsr id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_file_dsr ALTER COLUMN id SET DEFAULT nextval('maker.pot_file_dsr_id_seq'::regclass);


--
-- Name: pot_file_vow id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_file_vow ALTER COLUMN id SET DEFAULT nextval('maker.pot_file_vow_id_seq'::regclass);


--
-- Name: pot_join id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_join ALTER COLUMN id SET DEFAULT nextval('maker.pot_join_id_seq'::regclass);


--
-- Name: pot_live id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_live ALTER COLUMN id SET DEFAULT nextval('maker.pot_live_id_seq'::regclass);


--
-- Name: pot_pie id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_pie ALTER COLUMN id SET DEFAULT nextval('maker.pot_pie_id_seq'::regclass);


--
-- Name: pot_rho id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_rho ALTER COLUMN id SET DEFAULT nextval('maker.pot_rho_id_seq'::regclass);


--
-- Name: pot_user_pie id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_user_pie ALTER COLUMN id SET DEFAULT nextval('maker.pot_user_pie_id_seq'::regclass);


--
-- Name: pot_vat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_vat ALTER COLUMN id SET DEFAULT nextval('maker.pot_vat_id_seq'::regclass);


--
-- Name: pot_vow id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_vow ALTER COLUMN id SET DEFAULT nextval('maker.pot_vow_id_seq'::regclass);


--
-- Name: rely id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.rely ALTER COLUMN id SET DEFAULT nextval('maker.rely_id_seq'::regclass);


--
-- Name: set_min_sell id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.set_min_sell ALTER COLUMN id SET DEFAULT nextval('maker.set_min_sell_id_seq'::regclass);


--
-- Name: spot_file_mat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_mat ALTER COLUMN id SET DEFAULT nextval('maker.spot_file_mat_id_seq'::regclass);


--
-- Name: spot_file_par id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_par ALTER COLUMN id SET DEFAULT nextval('maker.spot_file_par_id_seq'::regclass);


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
-- Name: spot_live id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_live ALTER COLUMN id SET DEFAULT nextval('maker.spot_live_id_seq'::regclass);


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
-- Name: vat_deny id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_deny ALTER COLUMN id SET DEFAULT nextval('maker.vat_deny_id_seq'::regclass);


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
-- Name: vat_hope id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_hope ALTER COLUMN id SET DEFAULT nextval('maker.vat_hope_id_seq'::regclass);


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
-- Name: vat_nope id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_nope ALTER COLUMN id SET DEFAULT nextval('maker.vat_nope_id_seq'::regclass);


--
-- Name: vat_rely id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_rely ALTER COLUMN id SET DEFAULT nextval('maker.vat_rely_id_seq'::regclass);


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
-- Name: vow_file_auction_address id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_address ALTER COLUMN id SET DEFAULT nextval('maker.vow_file_auction_address_id_seq'::regclass);


--
-- Name: vow_file_auction_attributes id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_attributes ALTER COLUMN id SET DEFAULT nextval('maker.vow_file_auction_attributes_id_seq'::regclass);


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
-- Name: vow_heal id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_heal ALTER COLUMN id SET DEFAULT nextval('maker.vow_heal_id_seq'::regclass);


--
-- Name: vow_hump id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_hump ALTER COLUMN id SET DEFAULT nextval('maker.vow_hump_id_seq'::regclass);


--
-- Name: vow_live id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_live ALTER COLUMN id SET DEFAULT nextval('maker.vow_live_id_seq'::regclass);


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
-- Name: wards id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.wards ALTER COLUMN id SET DEFAULT nextval('maker.wards_id_seq'::regclass);


--
-- Name: yank id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.yank ALTER COLUMN id SET DEFAULT nextval('maker.yank_id_seq'::regclass);


--
-- Name: addresses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.addresses ALTER COLUMN id SET DEFAULT nextval('public.addresses_id_seq'::regclass);


--
-- Name: checked_headers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers ALTER COLUMN id SET DEFAULT nextval('public.checked_headers_id_seq'::regclass);


--
-- Name: eth_nodes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes ALTER COLUMN id SET DEFAULT nextval('public.eth_nodes_id_seq'::regclass);


--
-- Name: event_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.event_logs ALTER COLUMN id SET DEFAULT nextval('public.event_logs_id_seq'::regclass);


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: headers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers ALTER COLUMN id SET DEFAULT nextval('public.headers_id_seq'::regclass);


--
-- Name: receipts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipts ALTER COLUMN id SET DEFAULT nextval('public.receipts_id_seq'::regclass);


--
-- Name: storage_diff id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.storage_diff ALTER COLUMN id SET DEFAULT nextval('public.storage_diff_id_seq'::regclass);


--
-- Name: transactions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactions ALTER COLUMN id SET DEFAULT nextval('public.transactions_id_seq'::regclass);


--
-- Name: watched_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_logs ALTER COLUMN id SET DEFAULT nextval('public.watched_logs_id_seq'::regclass);


--
-- Name: ilk_snapshot ilk_snapshot_pkey; Type: CONSTRAINT; Schema: api; Owner: -
--

ALTER TABLE ONLY api.ilk_snapshot
    ADD CONSTRAINT ilk_snapshot_pkey PRIMARY KEY (ilk_identifier, block_number);


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
-- Name: urn_snapshot urn_snapshot_pkey; Type: CONSTRAINT; Schema: api; Owner: -
--

ALTER TABLE ONLY api.urn_snapshot
    ADD CONSTRAINT urn_snapshot_pkey PRIMARY KEY (urn_identifier, ilk_identifier, block_height);


--
-- Name: auction_file auction_file_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.auction_file
    ADD CONSTRAINT auction_file_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: auction_file auction_file_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.auction_file
    ADD CONSTRAINT auction_file_pkey PRIMARY KEY (id);


--
-- Name: bid_event bid_event_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bid_event
    ADD CONSTRAINT bid_event_pkey PRIMARY KEY (log_id);


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
-- Name: cat_ilk_chop cat_ilk_chop_diff_id_header_id_ilk_id_chop_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_chop
    ADD CONSTRAINT cat_ilk_chop_diff_id_header_id_ilk_id_chop_key UNIQUE (diff_id, header_id, ilk_id, chop);


--
-- Name: cat_ilk_chop cat_ilk_chop_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_chop
    ADD CONSTRAINT cat_ilk_chop_pkey PRIMARY KEY (id);


--
-- Name: cat_ilk_flip cat_ilk_flip_diff_id_header_id_ilk_id_flip_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_flip
    ADD CONSTRAINT cat_ilk_flip_diff_id_header_id_ilk_id_flip_key UNIQUE (diff_id, header_id, ilk_id, flip);


--
-- Name: cat_ilk_flip cat_ilk_flip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_flip
    ADD CONSTRAINT cat_ilk_flip_pkey PRIMARY KEY (id);


--
-- Name: cat_ilk_lump cat_ilk_lump_diff_id_header_id_ilk_id_lump_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_lump
    ADD CONSTRAINT cat_ilk_lump_diff_id_header_id_ilk_id_lump_key UNIQUE (diff_id, header_id, ilk_id, lump);


--
-- Name: cat_ilk_lump cat_ilk_lump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_lump
    ADD CONSTRAINT cat_ilk_lump_pkey PRIMARY KEY (id);


--
-- Name: cat_live cat_live_diff_id_header_id_live_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_live
    ADD CONSTRAINT cat_live_diff_id_header_id_live_key UNIQUE (diff_id, header_id, live);


--
-- Name: cat_live cat_live_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_live
    ADD CONSTRAINT cat_live_pkey PRIMARY KEY (id);


--
-- Name: cat_vat cat_vat_diff_id_header_id_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vat
    ADD CONSTRAINT cat_vat_diff_id_header_id_vat_key UNIQUE (diff_id, header_id, vat);


--
-- Name: cat_vat cat_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vat
    ADD CONSTRAINT cat_vat_pkey PRIMARY KEY (id);


--
-- Name: cat_vow cat_vow_diff_id_header_id_vow_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vow
    ADD CONSTRAINT cat_vow_diff_id_header_id_vow_key UNIQUE (diff_id, header_id, vow);


--
-- Name: cat_vow cat_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vow
    ADD CONSTRAINT cat_vow_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_cdpi cdp_manager_cdpi_diff_id_header_id_cdpi_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_cdpi
    ADD CONSTRAINT cdp_manager_cdpi_diff_id_header_id_cdpi_key UNIQUE (diff_id, header_id, cdpi);


--
-- Name: cdp_manager_cdpi cdp_manager_cdpi_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_cdpi
    ADD CONSTRAINT cdp_manager_cdpi_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_count cdp_manager_count_diff_id_header_id_owner_count_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_count
    ADD CONSTRAINT cdp_manager_count_diff_id_header_id_owner_count_key UNIQUE (diff_id, header_id, owner, count);


--
-- Name: cdp_manager_count cdp_manager_count_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_count
    ADD CONSTRAINT cdp_manager_count_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_first cdp_manager_first_diff_id_header_id_owner_first_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_first
    ADD CONSTRAINT cdp_manager_first_diff_id_header_id_owner_first_key UNIQUE (diff_id, header_id, owner, first);


--
-- Name: cdp_manager_first cdp_manager_first_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_first
    ADD CONSTRAINT cdp_manager_first_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_ilks cdp_manager_ilks_diff_id_header_id_cdpi_ilk_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_ilks
    ADD CONSTRAINT cdp_manager_ilks_diff_id_header_id_cdpi_ilk_id_key UNIQUE (diff_id, header_id, cdpi, ilk_id);


--
-- Name: cdp_manager_ilks cdp_manager_ilks_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_ilks
    ADD CONSTRAINT cdp_manager_ilks_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_last cdp_manager_last_diff_id_header_id_owner_last_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_last
    ADD CONSTRAINT cdp_manager_last_diff_id_header_id_owner_last_key UNIQUE (diff_id, header_id, owner, last);


--
-- Name: cdp_manager_last cdp_manager_last_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_last
    ADD CONSTRAINT cdp_manager_last_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_list_next cdp_manager_list_next_diff_id_header_id_cdpi_next_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_next
    ADD CONSTRAINT cdp_manager_list_next_diff_id_header_id_cdpi_next_key UNIQUE (diff_id, header_id, cdpi, next);


--
-- Name: cdp_manager_list_next cdp_manager_list_next_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_next
    ADD CONSTRAINT cdp_manager_list_next_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_list_prev cdp_manager_list_prev_diff_id_header_id_cdpi_prev_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_prev
    ADD CONSTRAINT cdp_manager_list_prev_diff_id_header_id_cdpi_prev_key UNIQUE (diff_id, header_id, cdpi, prev);


--
-- Name: cdp_manager_list_prev cdp_manager_list_prev_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_prev
    ADD CONSTRAINT cdp_manager_list_prev_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_owns cdp_manager_owns_diff_id_header_id_cdpi_owner_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_owns
    ADD CONSTRAINT cdp_manager_owns_diff_id_header_id_cdpi_owner_key UNIQUE (diff_id, header_id, cdpi, owner);


--
-- Name: cdp_manager_owns cdp_manager_owns_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_owns
    ADD CONSTRAINT cdp_manager_owns_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_urns cdp_manager_urns_diff_id_header_id_cdpi_urn_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_urns
    ADD CONSTRAINT cdp_manager_urns_diff_id_header_id_cdpi_urn_key UNIQUE (diff_id, header_id, cdpi, urn);


--
-- Name: cdp_manager_urns cdp_manager_urns_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_urns
    ADD CONSTRAINT cdp_manager_urns_pkey PRIMARY KEY (id);


--
-- Name: cdp_manager_vat cdp_manager_vat_diff_id_header_id_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_vat
    ADD CONSTRAINT cdp_manager_vat_diff_id_header_id_vat_key UNIQUE (diff_id, header_id, vat);


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
-- Name: deny deny_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deny
    ADD CONSTRAINT deny_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: deny deny_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deny
    ADD CONSTRAINT deny_pkey PRIMARY KEY (id);


--
-- Name: flap_beg flap_beg_diff_id_header_id_address_id_beg_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_beg
    ADD CONSTRAINT flap_beg_diff_id_header_id_address_id_beg_key UNIQUE (diff_id, header_id, address_id, beg);


--
-- Name: flap_beg flap_beg_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_beg
    ADD CONSTRAINT flap_beg_pkey PRIMARY KEY (id);


--
-- Name: flap_bid_bid flap_bid_bid_diff_id_header_id_address_id_bid_id_bid_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_bid
    ADD CONSTRAINT flap_bid_bid_diff_id_header_id_address_id_bid_id_bid_key UNIQUE (diff_id, header_id, address_id, bid_id, bid);


--
-- Name: flap_bid_bid flap_bid_bid_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_bid
    ADD CONSTRAINT flap_bid_bid_pkey PRIMARY KEY (id);


--
-- Name: flap_bid_end flap_bid_end_diff_id_header_id_address_id_bid_id_end_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_end
    ADD CONSTRAINT flap_bid_end_diff_id_header_id_address_id_bid_id_end_key UNIQUE (diff_id, header_id, address_id, bid_id, "end");


--
-- Name: flap_bid_end flap_bid_end_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_end
    ADD CONSTRAINT flap_bid_end_pkey PRIMARY KEY (id);


--
-- Name: flap_bid_guy flap_bid_guy_diff_id_header_id_address_id_bid_id_guy_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_guy
    ADD CONSTRAINT flap_bid_guy_diff_id_header_id_address_id_bid_id_guy_key UNIQUE (diff_id, header_id, address_id, bid_id, guy);


--
-- Name: flap_bid_guy flap_bid_guy_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_guy
    ADD CONSTRAINT flap_bid_guy_pkey PRIMARY KEY (id);


--
-- Name: flap_bid_lot flap_bid_lot_diff_id_header_id_address_id_bid_id_lot_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_lot
    ADD CONSTRAINT flap_bid_lot_diff_id_header_id_address_id_bid_id_lot_key UNIQUE (diff_id, header_id, address_id, bid_id, lot);


--
-- Name: flap_bid_lot flap_bid_lot_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_lot
    ADD CONSTRAINT flap_bid_lot_pkey PRIMARY KEY (id);


--
-- Name: flap_bid_tic flap_bid_tic_diff_id_header_id_address_id_bid_id_tic_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_tic
    ADD CONSTRAINT flap_bid_tic_diff_id_header_id_address_id_bid_id_tic_key UNIQUE (diff_id, header_id, address_id, bid_id, tic);


--
-- Name: flap_bid_tic flap_bid_tic_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_tic
    ADD CONSTRAINT flap_bid_tic_pkey PRIMARY KEY (id);


--
-- Name: flap_gem flap_gem_diff_id_header_id_address_id_gem_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_gem
    ADD CONSTRAINT flap_gem_diff_id_header_id_address_id_gem_key UNIQUE (diff_id, header_id, address_id, gem);


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
-- Name: flap_kicks flap_kicks_diff_id_header_id_address_id_kicks_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kicks
    ADD CONSTRAINT flap_kicks_diff_id_header_id_address_id_kicks_key UNIQUE (diff_id, header_id, address_id, kicks);


--
-- Name: flap_kicks flap_kicks_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kicks
    ADD CONSTRAINT flap_kicks_pkey PRIMARY KEY (id);


--
-- Name: flap_live flap_live_diff_id_header_id_address_id_live_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_live
    ADD CONSTRAINT flap_live_diff_id_header_id_address_id_live_key UNIQUE (diff_id, header_id, address_id, live);


--
-- Name: flap_live flap_live_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_live
    ADD CONSTRAINT flap_live_pkey PRIMARY KEY (id);


--
-- Name: flap flap_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap
    ADD CONSTRAINT flap_pkey PRIMARY KEY (address_id, bid_id, block_number);


--
-- Name: flap_tau flap_tau_diff_id_header_id_address_id_tau_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_tau
    ADD CONSTRAINT flap_tau_diff_id_header_id_address_id_tau_key UNIQUE (diff_id, header_id, address_id, tau);


--
-- Name: flap_tau flap_tau_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_tau
    ADD CONSTRAINT flap_tau_pkey PRIMARY KEY (id);


--
-- Name: flap_ttl flap_ttl_diff_id_header_id_address_id_ttl_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_ttl
    ADD CONSTRAINT flap_ttl_diff_id_header_id_address_id_ttl_key UNIQUE (diff_id, header_id, address_id, ttl);


--
-- Name: flap_ttl flap_ttl_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_ttl
    ADD CONSTRAINT flap_ttl_pkey PRIMARY KEY (id);


--
-- Name: flap_vat flap_vat_diff_id_header_id_address_id_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_vat
    ADD CONSTRAINT flap_vat_diff_id_header_id_address_id_vat_key UNIQUE (diff_id, header_id, address_id, vat);


--
-- Name: flap_vat flap_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_vat
    ADD CONSTRAINT flap_vat_pkey PRIMARY KEY (id);


--
-- Name: flip_beg flip_beg_diff_id_header_id_address_id_beg_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_beg
    ADD CONSTRAINT flip_beg_diff_id_header_id_address_id_beg_key UNIQUE (diff_id, header_id, address_id, beg);


--
-- Name: flip_beg flip_beg_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_beg
    ADD CONSTRAINT flip_beg_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_bid flip_bid_bid_diff_id_header_id_bid_id_address_id_bid_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_bid
    ADD CONSTRAINT flip_bid_bid_diff_id_header_id_bid_id_address_id_bid_key UNIQUE (diff_id, header_id, bid_id, address_id, bid);


--
-- Name: flip_bid_bid flip_bid_bid_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_bid
    ADD CONSTRAINT flip_bid_bid_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_end flip_bid_end_diff_id_header_id_bid_id_address_id_end_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_end
    ADD CONSTRAINT flip_bid_end_diff_id_header_id_bid_id_address_id_end_key UNIQUE (diff_id, header_id, bid_id, address_id, "end");


--
-- Name: flip_bid_end flip_bid_end_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_end
    ADD CONSTRAINT flip_bid_end_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_gal flip_bid_gal_diff_id_header_id_bid_id_address_id_gal_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_gal
    ADD CONSTRAINT flip_bid_gal_diff_id_header_id_bid_id_address_id_gal_key UNIQUE (diff_id, header_id, bid_id, address_id, gal);


--
-- Name: flip_bid_gal flip_bid_gal_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_gal
    ADD CONSTRAINT flip_bid_gal_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_guy flip_bid_guy_diff_id_header_id_bid_id_address_id_guy_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_guy
    ADD CONSTRAINT flip_bid_guy_diff_id_header_id_bid_id_address_id_guy_key UNIQUE (diff_id, header_id, bid_id, address_id, guy);


--
-- Name: flip_bid_guy flip_bid_guy_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_guy
    ADD CONSTRAINT flip_bid_guy_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_lot flip_bid_lot_diff_id_header_id_bid_id_address_id_lot_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_lot
    ADD CONSTRAINT flip_bid_lot_diff_id_header_id_bid_id_address_id_lot_key UNIQUE (diff_id, header_id, bid_id, address_id, lot);


--
-- Name: flip_bid_lot flip_bid_lot_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_lot
    ADD CONSTRAINT flip_bid_lot_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_tab flip_bid_tab_diff_id_header_id_bid_id_address_id_tab_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tab
    ADD CONSTRAINT flip_bid_tab_diff_id_header_id_bid_id_address_id_tab_key UNIQUE (diff_id, header_id, bid_id, address_id, tab);


--
-- Name: flip_bid_tab flip_bid_tab_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tab
    ADD CONSTRAINT flip_bid_tab_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_tic flip_bid_tic_diff_id_header_id_bid_id_address_id_tic_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tic
    ADD CONSTRAINT flip_bid_tic_diff_id_header_id_bid_id_address_id_tic_key UNIQUE (diff_id, header_id, bid_id, address_id, tic);


--
-- Name: flip_bid_tic flip_bid_tic_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tic
    ADD CONSTRAINT flip_bid_tic_pkey PRIMARY KEY (id);


--
-- Name: flip_bid_usr flip_bid_usr_diff_id_header_id_bid_id_address_id_usr_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_usr
    ADD CONSTRAINT flip_bid_usr_diff_id_header_id_bid_id_address_id_usr_key UNIQUE (diff_id, header_id, bid_id, address_id, usr);


--
-- Name: flip_bid_usr flip_bid_usr_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_usr
    ADD CONSTRAINT flip_bid_usr_pkey PRIMARY KEY (id);


--
-- Name: flip_ilk flip_ilk_diff_id_header_id_address_id_ilk_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ilk
    ADD CONSTRAINT flip_ilk_diff_id_header_id_address_id_ilk_id_key UNIQUE (diff_id, header_id, address_id, ilk_id);


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
-- Name: flip_kicks flip_kicks_diff_id_header_id_address_id_kicks_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kicks
    ADD CONSTRAINT flip_kicks_diff_id_header_id_address_id_kicks_key UNIQUE (diff_id, header_id, address_id, kicks);


--
-- Name: flip_kicks flip_kicks_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kicks
    ADD CONSTRAINT flip_kicks_pkey PRIMARY KEY (id);


--
-- Name: flip flip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip
    ADD CONSTRAINT flip_pkey PRIMARY KEY (address_id, bid_id, block_number);


--
-- Name: flip_tau flip_tau_diff_id_header_id_address_id_tau_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_tau
    ADD CONSTRAINT flip_tau_diff_id_header_id_address_id_tau_key UNIQUE (diff_id, header_id, address_id, tau);


--
-- Name: flip_tau flip_tau_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_tau
    ADD CONSTRAINT flip_tau_pkey PRIMARY KEY (id);


--
-- Name: flip_ttl flip_ttl_diff_id_header_id_address_id_ttl_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ttl
    ADD CONSTRAINT flip_ttl_diff_id_header_id_address_id_ttl_key UNIQUE (diff_id, header_id, address_id, ttl);


--
-- Name: flip_ttl flip_ttl_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ttl
    ADD CONSTRAINT flip_ttl_pkey PRIMARY KEY (id);


--
-- Name: flip_vat flip_vat_diff_id_header_id_address_id_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_vat
    ADD CONSTRAINT flip_vat_diff_id_header_id_address_id_vat_key UNIQUE (diff_id, header_id, address_id, vat);


--
-- Name: flip_vat flip_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_vat
    ADD CONSTRAINT flip_vat_pkey PRIMARY KEY (id);


--
-- Name: flop_beg flop_beg_diff_id_header_id_address_id_beg_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_beg
    ADD CONSTRAINT flop_beg_diff_id_header_id_address_id_beg_key UNIQUE (diff_id, header_id, address_id, beg);


--
-- Name: flop_beg flop_beg_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_beg
    ADD CONSTRAINT flop_beg_pkey PRIMARY KEY (id);


--
-- Name: flop_bid_bid flop_bid_bid_diff_id_header_id_bid_id_address_id_bid_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_bid
    ADD CONSTRAINT flop_bid_bid_diff_id_header_id_bid_id_address_id_bid_key UNIQUE (diff_id, header_id, bid_id, address_id, bid);


--
-- Name: flop_bid_bid flop_bid_bid_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_bid
    ADD CONSTRAINT flop_bid_bid_pkey PRIMARY KEY (id);


--
-- Name: flop_bid_end flop_bid_end_diff_id_header_id_bid_id_address_id_end_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_end
    ADD CONSTRAINT flop_bid_end_diff_id_header_id_bid_id_address_id_end_key UNIQUE (diff_id, header_id, bid_id, address_id, "end");


--
-- Name: flop_bid_end flop_bid_end_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_end
    ADD CONSTRAINT flop_bid_end_pkey PRIMARY KEY (id);


--
-- Name: flop_bid_guy flop_bid_guy_diff_id_header_id_bid_id_address_id_guy_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_guy
    ADD CONSTRAINT flop_bid_guy_diff_id_header_id_bid_id_address_id_guy_key UNIQUE (diff_id, header_id, bid_id, address_id, guy);


--
-- Name: flop_bid_guy flop_bid_guy_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_guy
    ADD CONSTRAINT flop_bid_guy_pkey PRIMARY KEY (id);


--
-- Name: flop_bid_lot flop_bid_lot_diff_id_header_id_bid_id_address_id_lot_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_lot
    ADD CONSTRAINT flop_bid_lot_diff_id_header_id_bid_id_address_id_lot_key UNIQUE (diff_id, header_id, bid_id, address_id, lot);


--
-- Name: flop_bid_lot flop_bid_lot_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_lot
    ADD CONSTRAINT flop_bid_lot_pkey PRIMARY KEY (id);


--
-- Name: flop_bid_tic flop_bid_tic_diff_id_header_id_bid_id_address_id_tic_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_tic
    ADD CONSTRAINT flop_bid_tic_diff_id_header_id_bid_id_address_id_tic_key UNIQUE (diff_id, header_id, bid_id, address_id, tic);


--
-- Name: flop_bid_tic flop_bid_tic_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_tic
    ADD CONSTRAINT flop_bid_tic_pkey PRIMARY KEY (id);


--
-- Name: flop_gem flop_gem_diff_id_header_id_address_id_gem_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_gem
    ADD CONSTRAINT flop_gem_diff_id_header_id_address_id_gem_key UNIQUE (diff_id, header_id, address_id, gem);


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
-- Name: flop_kicks flop_kicks_diff_id_header_id_address_id_kicks_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kicks
    ADD CONSTRAINT flop_kicks_diff_id_header_id_address_id_kicks_key UNIQUE (diff_id, header_id, address_id, kicks);


--
-- Name: flop_kicks flop_kicks_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kicks
    ADD CONSTRAINT flop_kicks_pkey PRIMARY KEY (id);


--
-- Name: flop_live flop_live_diff_id_header_id_address_id_live_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_live
    ADD CONSTRAINT flop_live_diff_id_header_id_address_id_live_key UNIQUE (diff_id, header_id, address_id, live);


--
-- Name: flop_live flop_live_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_live
    ADD CONSTRAINT flop_live_pkey PRIMARY KEY (id);


--
-- Name: flop_pad flop_pad_diff_id_header_id_address_id_pad_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_pad
    ADD CONSTRAINT flop_pad_diff_id_header_id_address_id_pad_key UNIQUE (diff_id, header_id, address_id, pad);


--
-- Name: flop_pad flop_pad_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_pad
    ADD CONSTRAINT flop_pad_pkey PRIMARY KEY (id);


--
-- Name: flop flop_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop
    ADD CONSTRAINT flop_pkey PRIMARY KEY (address_id, bid_id, block_number);


--
-- Name: flop_tau flop_tau_diff_id_header_id_address_id_tau_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_tau
    ADD CONSTRAINT flop_tau_diff_id_header_id_address_id_tau_key UNIQUE (diff_id, header_id, address_id, tau);


--
-- Name: flop_tau flop_tau_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_tau
    ADD CONSTRAINT flop_tau_pkey PRIMARY KEY (id);


--
-- Name: flop_ttl flop_ttl_diff_id_header_id_address_id_ttl_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_ttl
    ADD CONSTRAINT flop_ttl_diff_id_header_id_address_id_ttl_key UNIQUE (diff_id, header_id, address_id, ttl);


--
-- Name: flop_ttl flop_ttl_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_ttl
    ADD CONSTRAINT flop_ttl_pkey PRIMARY KEY (id);


--
-- Name: flop_vat flop_vat_diff_id_header_id_address_id_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vat
    ADD CONSTRAINT flop_vat_diff_id_header_id_address_id_vat_key UNIQUE (diff_id, header_id, address_id, vat);


--
-- Name: flop_vat flop_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vat
    ADD CONSTRAINT flop_vat_pkey PRIMARY KEY (id);


--
-- Name: flop_vow flop_vow_diff_id_header_id_address_id_vow_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vow
    ADD CONSTRAINT flop_vow_diff_id_header_id_address_id_vow_key UNIQUE (diff_id, header_id, address_id, vow);


--
-- Name: flop_vow flop_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vow
    ADD CONSTRAINT flop_vow_pkey PRIMARY KEY (id);


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
-- Name: jug_base jug_base_diff_id_header_id_base_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_base
    ADD CONSTRAINT jug_base_diff_id_header_id_base_key UNIQUE (diff_id, header_id, base);


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
-- Name: jug_ilk_duty jug_ilk_duty_diff_id_header_id_ilk_id_duty_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_duty
    ADD CONSTRAINT jug_ilk_duty_diff_id_header_id_ilk_id_duty_key UNIQUE (diff_id, header_id, ilk_id, duty);


--
-- Name: jug_ilk_duty jug_ilk_duty_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_duty
    ADD CONSTRAINT jug_ilk_duty_pkey PRIMARY KEY (id);


--
-- Name: jug_ilk_rho jug_ilk_rho_diff_id_header_id_ilk_id_rho_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_rho
    ADD CONSTRAINT jug_ilk_rho_diff_id_header_id_ilk_id_rho_key UNIQUE (diff_id, header_id, ilk_id, rho);


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
-- Name: jug_vat jug_vat_diff_id_header_id_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vat
    ADD CONSTRAINT jug_vat_diff_id_header_id_vat_key UNIQUE (diff_id, header_id, vat);


--
-- Name: jug_vat jug_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vat
    ADD CONSTRAINT jug_vat_pkey PRIMARY KEY (id);


--
-- Name: jug_vow jug_vow_diff_id_header_id_vow_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vow
    ADD CONSTRAINT jug_vow_diff_id_header_id_vow_key UNIQUE (diff_id, header_id, vow);


--
-- Name: jug_vow jug_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vow
    ADD CONSTRAINT jug_vow_pkey PRIMARY KEY (id);


--
-- Name: log_bump log_bump_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_bump
    ADD CONSTRAINT log_bump_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_bump log_bump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_bump
    ADD CONSTRAINT log_bump_pkey PRIMARY KEY (id);


--
-- Name: log_buy_enabled log_buy_enabled_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_buy_enabled
    ADD CONSTRAINT log_buy_enabled_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_buy_enabled log_buy_enabled_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_buy_enabled
    ADD CONSTRAINT log_buy_enabled_pkey PRIMARY KEY (id);


--
-- Name: log_delete log_delete_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_delete
    ADD CONSTRAINT log_delete_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_delete log_delete_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_delete
    ADD CONSTRAINT log_delete_pkey PRIMARY KEY (id);


--
-- Name: log_insert log_insert_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_insert
    ADD CONSTRAINT log_insert_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_insert log_insert_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_insert
    ADD CONSTRAINT log_insert_pkey PRIMARY KEY (id);


--
-- Name: log_item_update log_item_update_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_item_update
    ADD CONSTRAINT log_item_update_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_item_update log_item_update_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_item_update
    ADD CONSTRAINT log_item_update_pkey PRIMARY KEY (id);


--
-- Name: log_kill log_kill_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_kill
    ADD CONSTRAINT log_kill_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_kill log_kill_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_kill
    ADD CONSTRAINT log_kill_pkey PRIMARY KEY (id);


--
-- Name: log_make log_make_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_make
    ADD CONSTRAINT log_make_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_make log_make_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_make
    ADD CONSTRAINT log_make_pkey PRIMARY KEY (id);


--
-- Name: log_matching_enabled log_matching_enabled_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_matching_enabled
    ADD CONSTRAINT log_matching_enabled_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_matching_enabled log_matching_enabled_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_matching_enabled
    ADD CONSTRAINT log_matching_enabled_pkey PRIMARY KEY (id);


--
-- Name: log_median_price log_median_price_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_median_price
    ADD CONSTRAINT log_median_price_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_median_price log_median_price_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_median_price
    ADD CONSTRAINT log_median_price_pkey PRIMARY KEY (id);


--
-- Name: log_min_sell log_min_sell_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_min_sell
    ADD CONSTRAINT log_min_sell_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_min_sell log_min_sell_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_min_sell
    ADD CONSTRAINT log_min_sell_pkey PRIMARY KEY (id);


--
-- Name: log_sorted_offer log_sorted_offer_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_sorted_offer
    ADD CONSTRAINT log_sorted_offer_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_sorted_offer log_sorted_offer_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_sorted_offer
    ADD CONSTRAINT log_sorted_offer_pkey PRIMARY KEY (id);


--
-- Name: log_take log_take_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_take
    ADD CONSTRAINT log_take_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_take log_take_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_take
    ADD CONSTRAINT log_take_pkey PRIMARY KEY (id);


--
-- Name: log_trade log_trade_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_trade
    ADD CONSTRAINT log_trade_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_trade log_trade_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_trade
    ADD CONSTRAINT log_trade_pkey PRIMARY KEY (id);


--
-- Name: log_unsorted_offer log_unsorted_offer_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_unsorted_offer
    ADD CONSTRAINT log_unsorted_offer_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_unsorted_offer log_unsorted_offer_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_unsorted_offer
    ADD CONSTRAINT log_unsorted_offer_pkey PRIMARY KEY (id);


--
-- Name: log_value log_value_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_value
    ADD CONSTRAINT log_value_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: log_value log_value_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_value
    ADD CONSTRAINT log_value_pkey PRIMARY KEY (id);


--
-- Name: median_age median_age_header_id_address_id_age_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_age
    ADD CONSTRAINT median_age_header_id_address_id_age_key UNIQUE (header_id, address_id, age);


--
-- Name: median_age median_age_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_age
    ADD CONSTRAINT median_age_pkey PRIMARY KEY (id);


--
-- Name: median_bar median_bar_header_id_address_id_bar_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bar
    ADD CONSTRAINT median_bar_header_id_address_id_bar_key UNIQUE (header_id, address_id, bar);


--
-- Name: median_bar median_bar_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bar
    ADD CONSTRAINT median_bar_pkey PRIMARY KEY (id);


--
-- Name: median_bud median_bud_header_id_address_id_a_bud_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bud
    ADD CONSTRAINT median_bud_header_id_address_id_a_bud_key UNIQUE (header_id, address_id, a, bud);


--
-- Name: median_bud median_bud_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bud
    ADD CONSTRAINT median_bud_pkey PRIMARY KEY (id);


--
-- Name: median_diss_batch median_diss_batch_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_batch
    ADD CONSTRAINT median_diss_batch_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: median_diss_batch median_diss_batch_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_batch
    ADD CONSTRAINT median_diss_batch_pkey PRIMARY KEY (id);


--
-- Name: median_diss_single median_diss_single_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_single
    ADD CONSTRAINT median_diss_single_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: median_diss_single median_diss_single_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_single
    ADD CONSTRAINT median_diss_single_pkey PRIMARY KEY (id);


--
-- Name: median_drop median_drop_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_drop
    ADD CONSTRAINT median_drop_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: median_drop median_drop_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_drop
    ADD CONSTRAINT median_drop_pkey PRIMARY KEY (id);


--
-- Name: median_kiss_batch median_kiss_batch_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_batch
    ADD CONSTRAINT median_kiss_batch_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: median_kiss_batch median_kiss_batch_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_batch
    ADD CONSTRAINT median_kiss_batch_pkey PRIMARY KEY (id);


--
-- Name: median_kiss_single median_kiss_single_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_single
    ADD CONSTRAINT median_kiss_single_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: median_kiss_single median_kiss_single_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_single
    ADD CONSTRAINT median_kiss_single_pkey PRIMARY KEY (id);


--
-- Name: median_lift median_lift_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_lift
    ADD CONSTRAINT median_lift_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: median_lift median_lift_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_lift
    ADD CONSTRAINT median_lift_pkey PRIMARY KEY (id);


--
-- Name: median_orcl median_orcl_header_id_address_id_a_orcl_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_orcl
    ADD CONSTRAINT median_orcl_header_id_address_id_a_orcl_key UNIQUE (header_id, address_id, a, orcl);


--
-- Name: median_orcl median_orcl_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_orcl
    ADD CONSTRAINT median_orcl_pkey PRIMARY KEY (id);


--
-- Name: median_slot median_slot_header_id_address_id_slot_id_slot_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_slot
    ADD CONSTRAINT median_slot_header_id_address_id_slot_id_slot_key UNIQUE (header_id, address_id, slot_id, slot);


--
-- Name: median_slot median_slot_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_slot
    ADD CONSTRAINT median_slot_pkey PRIMARY KEY (id);


--
-- Name: median_val median_val_header_id_address_id_val_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_val
    ADD CONSTRAINT median_val_header_id_address_id_val_key UNIQUE (header_id, address_id, val);


--
-- Name: median_val median_val_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_val
    ADD CONSTRAINT median_val_pkey PRIMARY KEY (id);


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
-- Name: osm_change osm_change_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.osm_change
    ADD CONSTRAINT osm_change_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: osm_change osm_change_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.osm_change
    ADD CONSTRAINT osm_change_pkey PRIMARY KEY (id);


--
-- Name: pot_cage pot_cage_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_cage
    ADD CONSTRAINT pot_cage_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: pot_cage pot_cage_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_cage
    ADD CONSTRAINT pot_cage_pkey PRIMARY KEY (id);


--
-- Name: pot_chi pot_chi_diff_id_header_id_chi_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_chi
    ADD CONSTRAINT pot_chi_diff_id_header_id_chi_key UNIQUE (diff_id, header_id, chi);


--
-- Name: pot_chi pot_chi_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_chi
    ADD CONSTRAINT pot_chi_pkey PRIMARY KEY (id);


--
-- Name: pot_drip pot_drip_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_drip
    ADD CONSTRAINT pot_drip_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: pot_drip pot_drip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_drip
    ADD CONSTRAINT pot_drip_pkey PRIMARY KEY (id);


--
-- Name: pot_dsr pot_dsr_diff_id_header_id_dsr_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_dsr
    ADD CONSTRAINT pot_dsr_diff_id_header_id_dsr_key UNIQUE (diff_id, header_id, dsr);


--
-- Name: pot_dsr pot_dsr_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_dsr
    ADD CONSTRAINT pot_dsr_pkey PRIMARY KEY (id);


--
-- Name: pot_exit pot_exit_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_exit
    ADD CONSTRAINT pot_exit_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: pot_exit pot_exit_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_exit
    ADD CONSTRAINT pot_exit_pkey PRIMARY KEY (id);


--
-- Name: pot_file_dsr pot_file_dsr_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_file_dsr
    ADD CONSTRAINT pot_file_dsr_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: pot_file_dsr pot_file_dsr_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_file_dsr
    ADD CONSTRAINT pot_file_dsr_pkey PRIMARY KEY (id);


--
-- Name: pot_file_vow pot_file_vow_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_file_vow
    ADD CONSTRAINT pot_file_vow_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: pot_file_vow pot_file_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_file_vow
    ADD CONSTRAINT pot_file_vow_pkey PRIMARY KEY (id);


--
-- Name: pot_join pot_join_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_join
    ADD CONSTRAINT pot_join_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: pot_join pot_join_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_join
    ADD CONSTRAINT pot_join_pkey PRIMARY KEY (id);


--
-- Name: pot_live pot_live_diff_id_header_id_live_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_live
    ADD CONSTRAINT pot_live_diff_id_header_id_live_key UNIQUE (diff_id, header_id, live);


--
-- Name: pot_live pot_live_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_live
    ADD CONSTRAINT pot_live_pkey PRIMARY KEY (id);


--
-- Name: pot_pie pot_pie_diff_id_header_id_pie_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_pie
    ADD CONSTRAINT pot_pie_diff_id_header_id_pie_key UNIQUE (diff_id, header_id, pie);


--
-- Name: pot_pie pot_pie_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_pie
    ADD CONSTRAINT pot_pie_pkey PRIMARY KEY (id);


--
-- Name: pot_rho pot_rho_diff_id_header_id_rho_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_rho
    ADD CONSTRAINT pot_rho_diff_id_header_id_rho_key UNIQUE (diff_id, header_id, rho);


--
-- Name: pot_rho pot_rho_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_rho
    ADD CONSTRAINT pot_rho_pkey PRIMARY KEY (id);


--
-- Name: pot_user_pie pot_user_pie_diff_id_header_id_user_pie_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_user_pie
    ADD CONSTRAINT pot_user_pie_diff_id_header_id_user_pie_key UNIQUE (diff_id, header_id, "user", pie);


--
-- Name: pot_user_pie pot_user_pie_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_user_pie
    ADD CONSTRAINT pot_user_pie_pkey PRIMARY KEY (id);


--
-- Name: pot_vat pot_vat_diff_id_header_id_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_vat
    ADD CONSTRAINT pot_vat_diff_id_header_id_vat_key UNIQUE (diff_id, header_id, vat);


--
-- Name: pot_vat pot_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_vat
    ADD CONSTRAINT pot_vat_pkey PRIMARY KEY (id);


--
-- Name: pot_vow pot_vow_diff_id_header_id_vow_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_vow
    ADD CONSTRAINT pot_vow_diff_id_header_id_vow_key UNIQUE (diff_id, header_id, vow);


--
-- Name: pot_vow pot_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_vow
    ADD CONSTRAINT pot_vow_pkey PRIMARY KEY (id);


--
-- Name: rely rely_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.rely
    ADD CONSTRAINT rely_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: rely rely_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.rely
    ADD CONSTRAINT rely_pkey PRIMARY KEY (id);


--
-- Name: set_min_sell set_min_sell_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.set_min_sell
    ADD CONSTRAINT set_min_sell_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: set_min_sell set_min_sell_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.set_min_sell
    ADD CONSTRAINT set_min_sell_pkey PRIMARY KEY (id);


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
-- Name: spot_file_par spot_file_par_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_par
    ADD CONSTRAINT spot_file_par_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: spot_file_par spot_file_par_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_par
    ADD CONSTRAINT spot_file_par_pkey PRIMARY KEY (id);


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
-- Name: spot_ilk_mat spot_ilk_mat_diff_id_header_id_ilk_id_mat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_mat
    ADD CONSTRAINT spot_ilk_mat_diff_id_header_id_ilk_id_mat_key UNIQUE (diff_id, header_id, ilk_id, mat);


--
-- Name: spot_ilk_mat spot_ilk_mat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_mat
    ADD CONSTRAINT spot_ilk_mat_pkey PRIMARY KEY (id);


--
-- Name: spot_ilk_pip spot_ilk_pip_diff_id_header_id_ilk_id_pip_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_pip
    ADD CONSTRAINT spot_ilk_pip_diff_id_header_id_ilk_id_pip_key UNIQUE (diff_id, header_id, ilk_id, pip);


--
-- Name: spot_ilk_pip spot_ilk_pip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_pip
    ADD CONSTRAINT spot_ilk_pip_pkey PRIMARY KEY (id);


--
-- Name: spot_live spot_live_diff_id_header_id_live_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_live
    ADD CONSTRAINT spot_live_diff_id_header_id_live_key UNIQUE (diff_id, header_id, live);


--
-- Name: spot_live spot_live_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_live
    ADD CONSTRAINT spot_live_pkey PRIMARY KEY (id);


--
-- Name: spot_par spot_par_diff_id_header_id_par_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_par
    ADD CONSTRAINT spot_par_diff_id_header_id_par_key UNIQUE (diff_id, header_id, par);


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
-- Name: spot_vat spot_vat_diff_id_header_id_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_vat
    ADD CONSTRAINT spot_vat_diff_id_header_id_vat_key UNIQUE (diff_id, header_id, vat);


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
-- Name: vat_dai vat_dai_diff_id_header_id_guy_dai_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_dai
    ADD CONSTRAINT vat_dai_diff_id_header_id_guy_dai_key UNIQUE (diff_id, header_id, guy, dai);


--
-- Name: vat_dai vat_dai_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_dai
    ADD CONSTRAINT vat_dai_pkey PRIMARY KEY (id);


--
-- Name: vat_debt vat_debt_diff_id_header_id_debt_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_debt
    ADD CONSTRAINT vat_debt_diff_id_header_id_debt_key UNIQUE (diff_id, header_id, debt);


--
-- Name: vat_debt vat_debt_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_debt
    ADD CONSTRAINT vat_debt_pkey PRIMARY KEY (id);


--
-- Name: vat_deny vat_deny_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_deny
    ADD CONSTRAINT vat_deny_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_deny vat_deny_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_deny
    ADD CONSTRAINT vat_deny_pkey PRIMARY KEY (id);


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
-- Name: vat_gem vat_gem_diff_id_header_id_ilk_id_guy_gem_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_gem
    ADD CONSTRAINT vat_gem_diff_id_header_id_ilk_id_guy_gem_key UNIQUE (diff_id, header_id, ilk_id, guy, gem);


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
-- Name: vat_hope vat_hope_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_hope
    ADD CONSTRAINT vat_hope_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_hope vat_hope_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_hope
    ADD CONSTRAINT vat_hope_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_art vat_ilk_art_diff_id_header_id_ilk_id_art_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_art
    ADD CONSTRAINT vat_ilk_art_diff_id_header_id_ilk_id_art_key UNIQUE (diff_id, header_id, ilk_id, art);


--
-- Name: vat_ilk_art vat_ilk_art_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_art
    ADD CONSTRAINT vat_ilk_art_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_dust vat_ilk_dust_diff_id_header_id_ilk_id_dust_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_dust
    ADD CONSTRAINT vat_ilk_dust_diff_id_header_id_ilk_id_dust_key UNIQUE (diff_id, header_id, ilk_id, dust);


--
-- Name: vat_ilk_dust vat_ilk_dust_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_dust
    ADD CONSTRAINT vat_ilk_dust_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_line vat_ilk_line_diff_id_header_id_ilk_id_line_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_line
    ADD CONSTRAINT vat_ilk_line_diff_id_header_id_ilk_id_line_key UNIQUE (diff_id, header_id, ilk_id, line);


--
-- Name: vat_ilk_line vat_ilk_line_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_line
    ADD CONSTRAINT vat_ilk_line_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_rate vat_ilk_rate_diff_id_header_id_ilk_id_rate_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_rate
    ADD CONSTRAINT vat_ilk_rate_diff_id_header_id_ilk_id_rate_key UNIQUE (diff_id, header_id, ilk_id, rate);


--
-- Name: vat_ilk_rate vat_ilk_rate_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_rate
    ADD CONSTRAINT vat_ilk_rate_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_spot vat_ilk_spot_diff_id_header_id_ilk_id_spot_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_spot
    ADD CONSTRAINT vat_ilk_spot_diff_id_header_id_ilk_id_spot_key UNIQUE (diff_id, header_id, ilk_id, spot);


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
-- Name: vat_line vat_line_diff_id_header_id_line_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_line
    ADD CONSTRAINT vat_line_diff_id_header_id_line_key UNIQUE (diff_id, header_id, line);


--
-- Name: vat_line vat_line_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_line
    ADD CONSTRAINT vat_line_pkey PRIMARY KEY (id);


--
-- Name: vat_live vat_live_diff_id_header_id_live_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_live
    ADD CONSTRAINT vat_live_diff_id_header_id_live_key UNIQUE (diff_id, header_id, live);


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
-- Name: vat_nope vat_nope_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_nope
    ADD CONSTRAINT vat_nope_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_nope vat_nope_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_nope
    ADD CONSTRAINT vat_nope_pkey PRIMARY KEY (id);


--
-- Name: vat_rely vat_rely_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_rely
    ADD CONSTRAINT vat_rely_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vat_rely vat_rely_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_rely
    ADD CONSTRAINT vat_rely_pkey PRIMARY KEY (id);


--
-- Name: vat_sin vat_sin_diff_id_header_id_guy_sin_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_sin
    ADD CONSTRAINT vat_sin_diff_id_header_id_guy_sin_key UNIQUE (diff_id, header_id, guy, sin);


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
-- Name: vat_urn_art vat_urn_art_diff_id_header_id_urn_id_art_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_art
    ADD CONSTRAINT vat_urn_art_diff_id_header_id_urn_id_art_key UNIQUE (diff_id, header_id, urn_id, art);


--
-- Name: vat_urn_art vat_urn_art_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_art
    ADD CONSTRAINT vat_urn_art_pkey PRIMARY KEY (id);


--
-- Name: vat_urn_ink vat_urn_ink_diff_id_header_id_urn_id_ink_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_ink
    ADD CONSTRAINT vat_urn_ink_diff_id_header_id_urn_id_ink_key UNIQUE (diff_id, header_id, urn_id, ink);


--
-- Name: vat_urn_ink vat_urn_ink_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_ink
    ADD CONSTRAINT vat_urn_ink_pkey PRIMARY KEY (id);


--
-- Name: vat_vice vat_vice_diff_id_header_id_vice_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_vice
    ADD CONSTRAINT vat_vice_diff_id_header_id_vice_key UNIQUE (diff_id, header_id, vice);


--
-- Name: vat_vice vat_vice_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_vice
    ADD CONSTRAINT vat_vice_pkey PRIMARY KEY (id);


--
-- Name: vow_ash vow_ash_diff_id_header_id_ash_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_ash
    ADD CONSTRAINT vow_ash_diff_id_header_id_ash_key UNIQUE (diff_id, header_id, ash);


--
-- Name: vow_ash vow_ash_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_ash
    ADD CONSTRAINT vow_ash_pkey PRIMARY KEY (id);


--
-- Name: vow_bump vow_bump_diff_id_header_id_bump_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_bump
    ADD CONSTRAINT vow_bump_diff_id_header_id_bump_key UNIQUE (diff_id, header_id, bump);


--
-- Name: vow_bump vow_bump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_bump
    ADD CONSTRAINT vow_bump_pkey PRIMARY KEY (id);


--
-- Name: vow_dump vow_dump_diff_id_header_id_dump_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_dump
    ADD CONSTRAINT vow_dump_diff_id_header_id_dump_key UNIQUE (diff_id, header_id, dump);


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
-- Name: vow_file_auction_address vow_file_auction_address_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_address
    ADD CONSTRAINT vow_file_auction_address_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vow_file_auction_address vow_file_auction_address_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_address
    ADD CONSTRAINT vow_file_auction_address_pkey PRIMARY KEY (id);


--
-- Name: vow_file_auction_attributes vow_file_auction_attributes_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_attributes
    ADD CONSTRAINT vow_file_auction_attributes_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vow_file_auction_attributes vow_file_auction_attributes_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_attributes
    ADD CONSTRAINT vow_file_auction_attributes_pkey PRIMARY KEY (id);


--
-- Name: vow_flapper vow_flapper_diff_id_header_id_flapper_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flapper
    ADD CONSTRAINT vow_flapper_diff_id_header_id_flapper_key UNIQUE (diff_id, header_id, flapper);


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
-- Name: vow_flopper vow_flopper_diff_id_header_id_flopper_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flopper
    ADD CONSTRAINT vow_flopper_diff_id_header_id_flopper_key UNIQUE (diff_id, header_id, flopper);


--
-- Name: vow_flopper vow_flopper_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flopper
    ADD CONSTRAINT vow_flopper_pkey PRIMARY KEY (id);


--
-- Name: vow_heal vow_heal_header_id_log_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_heal
    ADD CONSTRAINT vow_heal_header_id_log_id_key UNIQUE (header_id, log_id);


--
-- Name: vow_heal vow_heal_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_heal
    ADD CONSTRAINT vow_heal_pkey PRIMARY KEY (id);


--
-- Name: vow_hump vow_hump_diff_id_header_id_hump_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_hump
    ADD CONSTRAINT vow_hump_diff_id_header_id_hump_key UNIQUE (diff_id, header_id, hump);


--
-- Name: vow_hump vow_hump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_hump
    ADD CONSTRAINT vow_hump_pkey PRIMARY KEY (id);


--
-- Name: vow_live vow_live_diff_id_header_id_live_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_live
    ADD CONSTRAINT vow_live_diff_id_header_id_live_key UNIQUE (diff_id, header_id, live);


--
-- Name: vow_live vow_live_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_live
    ADD CONSTRAINT vow_live_pkey PRIMARY KEY (id);


--
-- Name: vow_sin_integer vow_sin_integer_diff_id_header_id_sin_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_integer
    ADD CONSTRAINT vow_sin_integer_diff_id_header_id_sin_key UNIQUE (diff_id, header_id, sin);


--
-- Name: vow_sin_integer vow_sin_integer_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_integer
    ADD CONSTRAINT vow_sin_integer_pkey PRIMARY KEY (id);


--
-- Name: vow_sin_mapping vow_sin_mapping_diff_id_header_id_era_tab_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_mapping
    ADD CONSTRAINT vow_sin_mapping_diff_id_header_id_era_tab_key UNIQUE (diff_id, header_id, era, tab);


--
-- Name: vow_sin_mapping vow_sin_mapping_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_mapping
    ADD CONSTRAINT vow_sin_mapping_pkey PRIMARY KEY (id);


--
-- Name: vow_sump vow_sump_diff_id_header_id_sump_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sump
    ADD CONSTRAINT vow_sump_diff_id_header_id_sump_key UNIQUE (diff_id, header_id, sump);


--
-- Name: vow_sump vow_sump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sump
    ADD CONSTRAINT vow_sump_pkey PRIMARY KEY (id);


--
-- Name: vow_vat vow_vat_diff_id_header_id_vat_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_vat
    ADD CONSTRAINT vow_vat_diff_id_header_id_vat_key UNIQUE (diff_id, header_id, vat);


--
-- Name: vow_vat vow_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_vat
    ADD CONSTRAINT vow_vat_pkey PRIMARY KEY (id);


--
-- Name: vow_wait vow_wait_diff_id_header_id_wait_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_wait
    ADD CONSTRAINT vow_wait_diff_id_header_id_wait_key UNIQUE (diff_id, header_id, wait);


--
-- Name: vow_wait vow_wait_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_wait
    ADD CONSTRAINT vow_wait_pkey PRIMARY KEY (id);


--
-- Name: wards wards_diff_id_header_id_address_id_usr_wards_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.wards
    ADD CONSTRAINT wards_diff_id_header_id_address_id_usr_wards_key UNIQUE (diff_id, header_id, address_id, usr, wards);


--
-- Name: wards wards_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.wards
    ADD CONSTRAINT wards_pkey PRIMARY KEY (id);


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
-- Name: eth_nodes eth_nodes_genesis_block_network_id_eth_node_id_client_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes
    ADD CONSTRAINT eth_nodes_genesis_block_network_id_eth_node_id_client_name_key UNIQUE (genesis_block, network_id, eth_node_id, client_name);


--
-- Name: eth_nodes eth_nodes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes
    ADD CONSTRAINT eth_nodes_pkey PRIMARY KEY (id);


--
-- Name: event_logs event_logs_header_id_tx_index_log_index_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.event_logs
    ADD CONSTRAINT event_logs_header_id_tx_index_log_index_key UNIQUE (header_id, tx_index, log_index);


--
-- Name: event_logs event_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.event_logs
    ADD CONSTRAINT event_logs_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: headers headers_block_number_eth_node_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_block_number_eth_node_id_key UNIQUE (block_number, eth_node_id);


--
-- Name: headers headers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_pkey PRIMARY KEY (id);


--
-- Name: receipts receipts_header_id_transaction_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_header_id_transaction_id_key UNIQUE (header_id, transaction_id);


--
-- Name: receipts receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_pkey PRIMARY KEY (id);


--
-- Name: storage_diff storage_diff_block_height_block_hash_hashed_address_storage_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.storage_diff
    ADD CONSTRAINT storage_diff_block_height_block_hash_hashed_address_storage_key UNIQUE (block_height, block_hash, hashed_address, storage_key, storage_value);


--
-- Name: storage_diff storage_diff_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.storage_diff
    ADD CONSTRAINT storage_diff_pkey PRIMARY KEY (id);


--
-- Name: transactions transactions_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_hash_key UNIQUE (hash);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: watched_logs watched_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_logs
    ADD CONSTRAINT watched_logs_pkey PRIMARY KEY (id);


--
-- Name: auction_file_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX auction_file_address_id_index ON maker.auction_file USING btree (address_id);


--
-- Name: auction_file_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX auction_file_header_index ON maker.auction_file USING btree (header_id);


--
-- Name: auction_file_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX auction_file_log_index ON maker.auction_file USING btree (log_id);


--
-- Name: auction_file_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX auction_file_msg_sender_index ON maker.auction_file USING btree (msg_sender);


--
-- Name: bid_event_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX bid_event_index ON maker.bid_event USING btree (contract_address, bid_id);


--
-- Name: bid_event_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX bid_event_urn_index ON maker.bid_event USING btree (ilk_identifier, urn_identifier);


--
-- Name: bite_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX bite_address_index ON maker.bite USING btree (address_id);


--
-- Name: bite_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX bite_header_index ON maker.bite USING btree (header_id);


--
-- Name: bite_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX bite_log_index ON maker.bite USING btree (log_id);


--
-- Name: bite_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX bite_urn_index ON maker.bite USING btree (urn_id);


--
-- Name: cat_file_cho_lump_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_cho_lump_msg_sender_index ON maker.cat_file_chop_lump USING btree (msg_sender);


--
-- Name: cat_file_chop_lump_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_chop_lump_header_index ON maker.cat_file_chop_lump USING btree (header_id);


--
-- Name: cat_file_chop_lump_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_chop_lump_ilk_index ON maker.cat_file_chop_lump USING btree (ilk_id);


--
-- Name: cat_file_chop_lump_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_chop_lump_log_index ON maker.cat_file_chop_lump USING btree (log_id);


--
-- Name: cat_file_flip_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_flip_header_index ON maker.cat_file_flip USING btree (header_id);


--
-- Name: cat_file_flip_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_flip_ilk_index ON maker.cat_file_flip USING btree (ilk_id);


--
-- Name: cat_file_flip_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_flip_log_index ON maker.cat_file_flip USING btree (log_id);


--
-- Name: cat_file_flip_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_flip_msg_sender_index ON maker.cat_file_flip USING btree (msg_sender);


--
-- Name: cat_file_vow_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_vow_header_index ON maker.cat_file_vow USING btree (header_id);


--
-- Name: cat_file_vow_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_vow_log_index ON maker.cat_file_vow USING btree (log_id);


--
-- Name: cat_file_vow_msg_sender; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_file_vow_msg_sender ON maker.cat_file_vow USING btree (msg_sender);


--
-- Name: cat_ilk_chop_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_ilk_chop_header_id_index ON maker.cat_ilk_chop USING btree (header_id);


--
-- Name: cat_ilk_chop_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_ilk_chop_ilk_index ON maker.cat_ilk_chop USING btree (ilk_id);


--
-- Name: cat_ilk_flip_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_ilk_flip_header_id_index ON maker.cat_ilk_flip USING btree (header_id);


--
-- Name: cat_ilk_flip_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_ilk_flip_ilk_index ON maker.cat_ilk_flip USING btree (ilk_id);


--
-- Name: cat_ilk_lump_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_ilk_lump_header_id_index ON maker.cat_ilk_lump USING btree (header_id);


--
-- Name: cat_ilk_lump_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_ilk_lump_ilk_index ON maker.cat_ilk_lump USING btree (ilk_id);


--
-- Name: cat_live_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_live_header_id_index ON maker.cat_live USING btree (header_id);


--
-- Name: cat_vat_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_vat_header_id_index ON maker.cat_vat USING btree (header_id);


--
-- Name: cat_vow_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cat_vow_header_id_index ON maker.cat_vow USING btree (header_id);


--
-- Name: cdp_manager_cdpi_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_cdpi_header_id_index ON maker.cdp_manager_cdpi USING btree (header_id);


--
-- Name: cdp_manager_count_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_count_header_id_index ON maker.cdp_manager_count USING btree (header_id);


--
-- Name: cdp_manager_first_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_first_header_id_index ON maker.cdp_manager_first USING btree (header_id);


--
-- Name: cdp_manager_ilks_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_ilks_header_id_index ON maker.cdp_manager_ilks USING btree (header_id);


--
-- Name: cdp_manager_ilks_ilk_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_ilks_ilk_id_index ON maker.cdp_manager_ilks USING btree (ilk_id);


--
-- Name: cdp_manager_last_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_last_header_id_index ON maker.cdp_manager_last USING btree (header_id);


--
-- Name: cdp_manager_list_next_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_list_next_header_id_index ON maker.cdp_manager_list_next USING btree (header_id);


--
-- Name: cdp_manager_list_prev_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_list_prev_header_id_index ON maker.cdp_manager_list_prev USING btree (header_id);


--
-- Name: cdp_manager_owns_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_owns_header_id_index ON maker.cdp_manager_owns USING btree (header_id);


--
-- Name: cdp_manager_owns_owner_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_owns_owner_index ON maker.cdp_manager_owns USING btree (owner);


--
-- Name: cdp_manager_urns_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_urns_header_id_index ON maker.cdp_manager_urns USING btree (header_id);


--
-- Name: cdp_manager_urns_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_urns_urn_index ON maker.cdp_manager_urns USING btree (urn);


--
-- Name: cdp_manager_vat_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX cdp_manager_vat_header_id_index ON maker.cdp_manager_vat USING btree (header_id);


--
-- Name: deal_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deal_address_index ON maker.deal USING btree (address_id);


--
-- Name: deal_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deal_bid_id_index ON maker.deal USING btree (bid_id);


--
-- Name: deal_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deal_header_index ON maker.deal USING btree (header_id);


--
-- Name: deal_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deal_log_index ON maker.deal USING btree (log_id);


--
-- Name: deal_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deal_msg_sender_index ON maker.deal USING btree (msg_sender);


--
-- Name: dent_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX dent_address_index ON maker.dent USING btree (address_id);


--
-- Name: dent_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX dent_header_index ON maker.dent USING btree (header_id);


--
-- Name: dent_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX dent_log_index ON maker.dent USING btree (log_id);


--
-- Name: dent_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX dent_msg_sender_index ON maker.dent USING btree (msg_sender);


--
-- Name: deny_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deny_address_index ON maker.deny USING btree (address_id);


--
-- Name: deny_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deny_header_index ON maker.deny USING btree (header_id);


--
-- Name: deny_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deny_log_index ON maker.deny USING btree (log_id);


--
-- Name: deny_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deny_msg_sender_index ON maker.deny USING btree (msg_sender);


--
-- Name: deny_usr_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX deny_usr_index ON maker.deny USING btree (usr);


--
-- Name: flap_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_address_index ON maker.flap USING btree (address_id);


--
-- Name: flap_beg_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_beg_address_index ON maker.flap_beg USING btree (address_id);


--
-- Name: flap_beg_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_beg_header_id_index ON maker.flap_beg USING btree (header_id);


--
-- Name: flap_bid_bid_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_bid_address_index ON maker.flap_bid_bid USING btree (address_id);


--
-- Name: flap_bid_bid_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_bid_bid_id_index ON maker.flap_bid_bid USING btree (bid_id);


--
-- Name: flap_bid_bid_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_bid_header_id_index ON maker.flap_bid_bid USING btree (header_id);


--
-- Name: flap_bid_end_bid_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_end_bid_address_index ON maker.flap_bid_end USING btree (address_id);


--
-- Name: flap_bid_end_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_end_bid_id_index ON maker.flap_bid_end USING btree (bid_id);


--
-- Name: flap_bid_end_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_end_header_id_index ON maker.flap_bid_end USING btree (header_id);


--
-- Name: flap_bid_guy_bid_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_guy_bid_address_index ON maker.flap_bid_guy USING btree (address_id);


--
-- Name: flap_bid_guy_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_guy_bid_id_index ON maker.flap_bid_guy USING btree (bid_id);


--
-- Name: flap_bid_guy_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_guy_header_id_index ON maker.flap_bid_guy USING btree (header_id);


--
-- Name: flap_bid_lot_bid_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_lot_bid_address_index ON maker.flap_bid_lot USING btree (address_id);


--
-- Name: flap_bid_lot_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_lot_bid_id_index ON maker.flap_bid_lot USING btree (bid_id);


--
-- Name: flap_bid_lot_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_lot_header_id_index ON maker.flap_bid_lot USING btree (header_id);


--
-- Name: flap_bid_tic_bid_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_tic_bid_address_index ON maker.flap_bid_tic USING btree (address_id);


--
-- Name: flap_bid_tic_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_tic_bid_id_index ON maker.flap_bid_tic USING btree (bid_id);


--
-- Name: flap_bid_tic_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_bid_tic_header_id_index ON maker.flap_bid_tic USING btree (header_id);


--
-- Name: flap_gem_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_gem_address_index ON maker.flap_gem USING btree (address_id);


--
-- Name: flap_gem_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_gem_header_id_index ON maker.flap_gem USING btree (header_id);


--
-- Name: flap_kick_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_kick_address_index ON maker.flap_kick USING btree (address_id);


--
-- Name: flap_kick_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_kick_header_index ON maker.flap_kick USING btree (header_id);


--
-- Name: flap_kick_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_kick_log_index ON maker.flap_kick USING btree (log_id);


--
-- Name: flap_kicks_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_kicks_address_index ON maker.flap_kicks USING btree (address_id);


--
-- Name: flap_kicks_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_kicks_header_id_index ON maker.flap_kicks USING btree (header_id);


--
-- Name: flap_live_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_live_address_index ON maker.flap_live USING btree (address_id);


--
-- Name: flap_live_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_live_header_id_index ON maker.flap_live USING btree (header_id);


--
-- Name: flap_tau_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_tau_address_index ON maker.flap_tau USING btree (address_id);


--
-- Name: flap_tau_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_tau_header_id_index ON maker.flap_tau USING btree (header_id);


--
-- Name: flap_ttl_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_ttl_address_index ON maker.flap_ttl USING btree (address_id);


--
-- Name: flap_ttl_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_ttl_header_id_index ON maker.flap_ttl USING btree (header_id);


--
-- Name: flap_vat_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_vat_address_index ON maker.flap_vat USING btree (address_id);


--
-- Name: flap_vat_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flap_vat_header_id_index ON maker.flap_vat USING btree (header_id);


--
-- Name: flip_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_address_index ON maker.flip USING btree (address_id);


--
-- Name: flip_beg_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_beg_address_index ON maker.flip_beg USING btree (address_id);


--
-- Name: flip_beg_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_beg_header_id_index ON maker.flip_beg USING btree (header_id);


--
-- Name: flip_bid_bid_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_bid_address_index ON maker.flip_bid_bid USING btree (address_id);


--
-- Name: flip_bid_bid_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_bid_bid_id_index ON maker.flip_bid_bid USING btree (bid_id);


--
-- Name: flip_bid_bid_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_bid_header_id_index ON maker.flip_bid_bid USING btree (header_id);


--
-- Name: flip_bid_end_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_end_address_index ON maker.flip_bid_end USING btree (address_id);


--
-- Name: flip_bid_end_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_end_bid_id_index ON maker.flip_bid_end USING btree (bid_id);


--
-- Name: flip_bid_end_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_end_header_id_index ON maker.flip_bid_end USING btree (header_id);


--
-- Name: flip_bid_gal_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_gal_address_index ON maker.flip_bid_gal USING btree (address_id);


--
-- Name: flip_bid_gal_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_gal_bid_id_index ON maker.flip_bid_gal USING btree (bid_id);


--
-- Name: flip_bid_gal_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_gal_header_id_index ON maker.flip_bid_gal USING btree (header_id);


--
-- Name: flip_bid_guy_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_guy_address_index ON maker.flip_bid_guy USING btree (address_id);


--
-- Name: flip_bid_guy_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_guy_bid_id_index ON maker.flip_bid_guy USING btree (bid_id);


--
-- Name: flip_bid_guy_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_guy_header_id_index ON maker.flip_bid_guy USING btree (header_id);


--
-- Name: flip_bid_lot_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_lot_address_index ON maker.flip_bid_lot USING btree (address_id);


--
-- Name: flip_bid_lot_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_lot_bid_id_index ON maker.flip_bid_lot USING btree (bid_id);


--
-- Name: flip_bid_lot_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_lot_header_id_index ON maker.flip_bid_lot USING btree (header_id);


--
-- Name: flip_bid_tab_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_tab_address_index ON maker.flip_bid_tab USING btree (address_id);


--
-- Name: flip_bid_tab_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_tab_bid_id_index ON maker.flip_bid_tab USING btree (bid_id);


--
-- Name: flip_bid_tab_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_tab_header_id_index ON maker.flip_bid_tab USING btree (header_id);


--
-- Name: flip_bid_tic_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_tic_address_index ON maker.flip_bid_tic USING btree (address_id);


--
-- Name: flip_bid_tic_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_tic_bid_id_index ON maker.flip_bid_tic USING btree (bid_id);


--
-- Name: flip_bid_tic_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_tic_header_id_index ON maker.flip_bid_tic USING btree (header_id);


--
-- Name: flip_bid_usr_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_usr_address_index ON maker.flip_bid_usr USING btree (address_id);


--
-- Name: flip_bid_usr_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_usr_bid_id_index ON maker.flip_bid_usr USING btree (bid_id);


--
-- Name: flip_bid_usr_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_bid_usr_header_id_index ON maker.flip_bid_usr USING btree (header_id);


--
-- Name: flip_ilk_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_ilk_address_index ON maker.flip_ilk USING btree (address_id);


--
-- Name: flip_ilk_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_ilk_header_id_index ON maker.flip_ilk USING btree (header_id);


--
-- Name: flip_ilk_ilk_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_ilk_ilk_id_index ON maker.flip_ilk USING btree (ilk_id);


--
-- Name: flip_kick_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_kick_address_index ON maker.flip_kick USING btree (address_id);


--
-- Name: flip_kick_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_kick_bid_id_index ON maker.flip_kick USING btree (bid_id);


--
-- Name: flip_kick_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_kick_header_index ON maker.flip_kick USING btree (header_id);


--
-- Name: flip_kick_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_kick_log_index ON maker.flip_kick USING btree (log_id);


--
-- Name: flip_kicks_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_kicks_address_index ON maker.flip_kicks USING btree (address_id);


--
-- Name: flip_kicks_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_kicks_header_id_index ON maker.flip_kicks USING btree (header_id);


--
-- Name: flip_tau_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_tau_address_index ON maker.flip_tau USING btree (address_id);


--
-- Name: flip_tau_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_tau_header_id_index ON maker.flip_tau USING btree (header_id);


--
-- Name: flip_ttl_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_ttl_address_index ON maker.flip_ttl USING btree (address_id);


--
-- Name: flip_ttl_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_ttl_header_id_index ON maker.flip_ttl USING btree (header_id);


--
-- Name: flip_vat_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_vat_address_index ON maker.flip_vat USING btree (address_id);


--
-- Name: flip_vat_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flip_vat_header_id_index ON maker.flip_vat USING btree (header_id);


--
-- Name: flop_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_address_index ON maker.flop USING btree (address_id);


--
-- Name: flop_beg_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_beg_address_index ON maker.flop_beg USING btree (address_id);


--
-- Name: flop_beg_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_beg_header_id_index ON maker.flop_beg USING btree (header_id);


--
-- Name: flop_bid_bid_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_bid_address_index ON maker.flop_bid_bid USING btree (address_id);


--
-- Name: flop_bid_bid_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_bid_bid_id_index ON maker.flop_bid_bid USING btree (bid_id);


--
-- Name: flop_bid_bid_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_bid_header_id_index ON maker.flop_bid_bid USING btree (header_id);


--
-- Name: flop_bid_end_bid_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_end_bid_address_index ON maker.flop_bid_end USING btree (address_id);


--
-- Name: flop_bid_end_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_end_bid_id_index ON maker.flop_bid_end USING btree (bid_id);


--
-- Name: flop_bid_end_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_end_header_id_index ON maker.flop_bid_end USING btree (header_id);


--
-- Name: flop_bid_guy_bid_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_guy_bid_address_index ON maker.flop_bid_guy USING btree (address_id);


--
-- Name: flop_bid_guy_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_guy_bid_id_index ON maker.flop_bid_guy USING btree (bid_id);


--
-- Name: flop_bid_guy_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_guy_header_id_index ON maker.flop_bid_guy USING btree (header_id);


--
-- Name: flop_bid_lot_bid_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_lot_bid_address_index ON maker.flop_bid_lot USING btree (address_id);


--
-- Name: flop_bid_lot_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_lot_bid_id_index ON maker.flop_bid_lot USING btree (bid_id);


--
-- Name: flop_bid_lot_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_lot_header_id_index ON maker.flop_bid_lot USING btree (header_id);


--
-- Name: flop_bid_tic_bid_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_tic_bid_address_index ON maker.flop_bid_tic USING btree (address_id);


--
-- Name: flop_bid_tic_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_tic_bid_id_index ON maker.flop_bid_tic USING btree (bid_id);


--
-- Name: flop_bid_tic_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_bid_tic_header_id_index ON maker.flop_bid_tic USING btree (header_id);


--
-- Name: flop_gem_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_gem_address_index ON maker.flop_gem USING btree (address_id);


--
-- Name: flop_gem_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_gem_header_id_index ON maker.flop_gem USING btree (header_id);


--
-- Name: flop_kick_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_kick_address_index ON maker.flop_kick USING btree (address_id);


--
-- Name: flop_kick_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_kick_header_index ON maker.flop_kick USING btree (header_id);


--
-- Name: flop_kick_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_kick_log_index ON maker.flop_kick USING btree (log_id);


--
-- Name: flop_kicks_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_kicks_address_index ON maker.flop_kicks USING btree (address_id);


--
-- Name: flop_kicks_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_kicks_header_id_index ON maker.flop_kicks USING btree (header_id);


--
-- Name: flop_live_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_live_address_index ON maker.flop_live USING btree (address_id);


--
-- Name: flop_live_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_live_header_id_index ON maker.flop_live USING btree (header_id);


--
-- Name: flop_pad_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_pad_address_index ON maker.flop_pad USING btree (address_id);


--
-- Name: flop_pad_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_pad_header_id_index ON maker.flop_pad USING btree (header_id);


--
-- Name: flop_tau_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_tau_address_index ON maker.flop_tau USING btree (address_id);


--
-- Name: flop_tau_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_tau_header_id_index ON maker.flop_tau USING btree (header_id);


--
-- Name: flop_ttl_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_ttl_address_index ON maker.flop_ttl USING btree (address_id);


--
-- Name: flop_ttl_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_ttl_header_id_index ON maker.flop_ttl USING btree (header_id);


--
-- Name: flop_vat_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_vat_address_index ON maker.flop_vat USING btree (address_id);


--
-- Name: flop_vat_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_vat_header_id_index ON maker.flop_vat USING btree (header_id);


--
-- Name: flop_vow_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_vow_address_index ON maker.flop_vow USING btree (address_id);


--
-- Name: flop_vow_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX flop_vow_header_id_index ON maker.flop_vow USING btree (header_id);


--
-- Name: jog_drip_msg_sender; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jog_drip_msg_sender ON maker.jug_drip USING btree (msg_sender);


--
-- Name: jug_base_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_base_header_id_index ON maker.jug_base USING btree (header_id);


--
-- Name: jug_drip_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_drip_header_index ON maker.jug_drip USING btree (header_id);


--
-- Name: jug_drip_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_drip_ilk_index ON maker.jug_drip USING btree (ilk_id);


--
-- Name: jug_drip_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_drip_log_index ON maker.jug_drip USING btree (log_id);


--
-- Name: jug_file_base_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_base_header_index ON maker.jug_file_base USING btree (header_id);


--
-- Name: jug_file_base_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_base_log_index ON maker.jug_file_base USING btree (log_id);


--
-- Name: jug_file_base_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_base_msg_sender_index ON maker.jug_file_base USING btree (msg_sender);


--
-- Name: jug_file_ilk_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_ilk_header_index ON maker.jug_file_ilk USING btree (header_id);


--
-- Name: jug_file_ilk_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_ilk_ilk_index ON maker.jug_file_ilk USING btree (ilk_id);


--
-- Name: jug_file_ilk_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_ilk_log_index ON maker.jug_file_ilk USING btree (log_id);


--
-- Name: jug_file_ilk_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_ilk_msg_sender_index ON maker.jug_file_ilk USING btree (msg_sender);


--
-- Name: jug_file_vow_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_vow_header_index ON maker.jug_file_vow USING btree (header_id);


--
-- Name: jug_file_vow_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_vow_log_index ON maker.jug_file_vow USING btree (log_id);


--
-- Name: jug_file_vow_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_file_vow_msg_sender_index ON maker.jug_file_vow USING btree (msg_sender);


--
-- Name: jug_ilk_duty_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_ilk_duty_header_id_index ON maker.jug_ilk_duty USING btree (header_id);


--
-- Name: jug_ilk_duty_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_ilk_duty_ilk_index ON maker.jug_ilk_duty USING btree (ilk_id);


--
-- Name: jug_ilk_rho_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_ilk_rho_header_id_index ON maker.jug_ilk_rho USING btree (header_id);


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
-- Name: jug_init_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_init_log_index ON maker.jug_init USING btree (log_id);


--
-- Name: jug_init_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_init_msg_sender_index ON maker.jug_init USING btree (msg_sender);


--
-- Name: jug_vat_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_vat_header_id_index ON maker.jug_vat USING btree (header_id);


--
-- Name: jug_vow_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX jug_vow_header_id_index ON maker.jug_vow USING btree (header_id);


--
-- Name: log_bump_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_bump_address_index ON maker.log_bump USING btree (address_id);


--
-- Name: log_bump_buy_gem_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_bump_buy_gem_index ON maker.log_bump USING btree (buy_gem);


--
-- Name: log_bump_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_bump_header_index ON maker.log_bump USING btree (header_id);


--
-- Name: log_bump_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_bump_log_index ON maker.log_bump USING btree (log_id);


--
-- Name: log_bump_maker_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_bump_maker_index ON maker.log_bump USING btree (maker);


--
-- Name: log_bump_pay_gem_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_bump_pay_gem_index ON maker.log_bump USING btree (pay_gem);


--
-- Name: log_buy_enabled_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_buy_enabled_address_index ON maker.log_buy_enabled USING btree (address_id);


--
-- Name: log_buy_enabled_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_buy_enabled_header_index ON maker.log_buy_enabled USING btree (header_id);


--
-- Name: log_buy_enabled_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_buy_enabled_log_index ON maker.log_buy_enabled USING btree (log_id);


--
-- Name: log_delete_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_delete_address_index ON maker.log_delete USING btree (address_id);


--
-- Name: log_delete_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_delete_header_index ON maker.log_delete USING btree (header_id);


--
-- Name: log_delete_keeper_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_delete_keeper_index ON maker.log_delete USING btree (keeper);


--
-- Name: log_delete_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_delete_log_index ON maker.log_delete USING btree (log_id);


--
-- Name: log_insert_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_insert_address_index ON maker.log_insert USING btree (address_id);


--
-- Name: log_insert_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_insert_header_index ON maker.log_insert USING btree (header_id);


--
-- Name: log_insert_keeper_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_insert_keeper_index ON maker.log_insert USING btree (keeper);


--
-- Name: log_insert_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_insert_log_index ON maker.log_insert USING btree (log_id);


--
-- Name: log_item_update_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_item_update_address_index ON maker.log_item_update USING btree (address_id);


--
-- Name: log_item_update_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_item_update_header_index ON maker.log_item_update USING btree (header_id);


--
-- Name: log_item_update_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_item_update_log_index ON maker.log_item_update USING btree (log_id);


--
-- Name: log_kill_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_kill_address_index ON maker.log_kill USING btree (address_id);


--
-- Name: log_kill_buy_gem_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_kill_buy_gem_index ON maker.log_kill USING btree (buy_gem);


--
-- Name: log_kill_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_kill_header_index ON maker.log_kill USING btree (header_id);


--
-- Name: log_kill_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_kill_index ON maker.log_kill USING btree (log_id);


--
-- Name: log_kill_maker_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_kill_maker_index ON maker.log_kill USING btree (maker);


--
-- Name: log_kill_pay_gem_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_kill_pay_gem_index ON maker.log_kill USING btree (pay_gem);


--
-- Name: log_make_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_make_address_index ON maker.log_make USING btree (address_id);


--
-- Name: log_make_buy_gem_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_make_buy_gem_index ON maker.log_make USING btree (buy_gem);


--
-- Name: log_make_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_make_header_index ON maker.log_make USING btree (header_id);


--
-- Name: log_make_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_make_log_index ON maker.log_make USING btree (log_id);


--
-- Name: log_make_maker_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_make_maker_index ON maker.log_make USING btree (maker);


--
-- Name: log_make_pay_gem_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_make_pay_gem_index ON maker.log_make USING btree (pay_gem);


--
-- Name: log_matching_enabled_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_matching_enabled_address_index ON maker.log_matching_enabled USING btree (address_id);


--
-- Name: log_matching_enabled_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_matching_enabled_header_index ON maker.log_matching_enabled USING btree (header_id);


--
-- Name: log_matching_enabled_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_matching_enabled_log_index ON maker.log_matching_enabled USING btree (log_id);


--
-- Name: log_median_price_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_median_price_address_index ON maker.log_median_price USING btree (address_id);


--
-- Name: log_median_price_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_median_price_header_index ON maker.log_median_price USING btree (header_id);


--
-- Name: log_median_price_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_median_price_log_index ON maker.log_median_price USING btree (log_id);


--
-- Name: log_min_sell_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_min_sell_address_index ON maker.log_min_sell USING btree (address_id);


--
-- Name: log_min_sell_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_min_sell_header_index ON maker.log_min_sell USING btree (header_id);


--
-- Name: log_min_sell_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_min_sell_log_index ON maker.log_min_sell USING btree (log_id);


--
-- Name: log_min_sell_pay_gem_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_min_sell_pay_gem_index ON maker.log_min_sell USING btree (pay_gem);


--
-- Name: log_sorted_offer_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_sorted_offer_address_index ON maker.log_sorted_offer USING btree (address_id);


--
-- Name: log_sorted_offer_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_sorted_offer_header_index ON maker.log_sorted_offer USING btree (header_id);


--
-- Name: log_sorted_offer_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_sorted_offer_log_index ON maker.log_sorted_offer USING btree (log_id);


--
-- Name: log_take_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_take_address_index ON maker.log_take USING btree (address_id);


--
-- Name: log_take_buy_gem_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_take_buy_gem_index ON maker.log_take USING btree (buy_gem);


--
-- Name: log_take_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_take_header_index ON maker.log_take USING btree (header_id);


--
-- Name: log_take_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_take_log_index ON maker.log_take USING btree (log_id);


--
-- Name: log_take_maker_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_take_maker_index ON maker.log_take USING btree (maker);


--
-- Name: log_take_pay_gem_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_take_pay_gem_index ON maker.log_take USING btree (pay_gem);


--
-- Name: log_take_taker_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_take_taker_index ON maker.log_take USING btree (taker);


--
-- Name: log_trade_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_trade_address_index ON maker.log_trade USING btree (address_id);


--
-- Name: log_trade_buy_gem_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_trade_buy_gem_index ON maker.log_trade USING btree (buy_gem);


--
-- Name: log_trade_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_trade_header_index ON maker.log_trade USING btree (header_id);


--
-- Name: log_trade_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_trade_log_index ON maker.log_trade USING btree (log_id);


--
-- Name: log_trade_pay_gem_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_trade_pay_gem_index ON maker.log_trade USING btree (pay_gem);


--
-- Name: log_unsorted_offer_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_unsorted_offer_address_index ON maker.log_unsorted_offer USING btree (address_id);


--
-- Name: log_unsorted_offer_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_unsorted_offer_header_index ON maker.log_unsorted_offer USING btree (header_id);


--
-- Name: log_unsorted_offer_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_unsorted_offer_log_index ON maker.log_unsorted_offer USING btree (log_id);


--
-- Name: log_value_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_value_address_index ON maker.log_value USING btree (address_id);


--
-- Name: log_value_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_value_header_index ON maker.log_value USING btree (header_id);


--
-- Name: log_value_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX log_value_log_index ON maker.log_value USING btree (log_id);


--
-- Name: median_age_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_age_address_id_index ON maker.median_age USING btree (address_id);


--
-- Name: median_age_diff_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_age_diff_id_index ON maker.median_age USING btree (diff_id);


--
-- Name: median_age_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_age_header_id_index ON maker.median_age USING btree (header_id);


--
-- Name: median_bar_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_bar_address_id_index ON maker.median_bar USING btree (address_id);


--
-- Name: median_bar_diff_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_bar_diff_id_index ON maker.median_bar USING btree (diff_id);


--
-- Name: median_bar_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_bar_header_id_index ON maker.median_bar USING btree (header_id);


--
-- Name: median_bud_a_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_bud_a_index ON maker.median_bud USING btree (a);


--
-- Name: median_bud_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_bud_address_index ON maker.median_bud USING btree (address_id);


--
-- Name: median_bud_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_bud_header_id_index ON maker.median_bud USING btree (header_id);


--
-- Name: median_diss_batch_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_diss_batch_address_index ON maker.median_diss_batch USING btree (address_id);


--
-- Name: median_diss_batch_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_diss_batch_header_index ON maker.median_diss_batch USING btree (header_id);


--
-- Name: median_diss_batch_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_diss_batch_log_index ON maker.median_diss_batch USING btree (log_id);


--
-- Name: median_diss_batch_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_diss_batch_msg_sender_index ON maker.median_diss_batch USING btree (msg_sender);


--
-- Name: median_diss_single_a_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_diss_single_a_index ON maker.median_diss_single USING btree (a);


--
-- Name: median_diss_single_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_diss_single_address_index ON maker.median_diss_single USING btree (address_id);


--
-- Name: median_diss_single_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_diss_single_header_index ON maker.median_diss_single USING btree (header_id);


--
-- Name: median_diss_single_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_diss_single_log_index ON maker.median_diss_single USING btree (log_id);


--
-- Name: median_diss_single_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_diss_single_msg_sender_index ON maker.median_diss_single USING btree (msg_sender);


--
-- Name: median_drop_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_drop_address_index ON maker.median_drop USING btree (address_id);


--
-- Name: median_drop_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_drop_header_index ON maker.median_drop USING btree (header_id);


--
-- Name: median_drop_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_drop_log_index ON maker.median_drop USING btree (log_id);


--
-- Name: median_drop_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_drop_msg_sender_index ON maker.median_drop USING btree (msg_sender);


--
-- Name: median_kiss_batch_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_kiss_batch_address_index ON maker.median_kiss_batch USING btree (address_id);


--
-- Name: median_kiss_batch_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_kiss_batch_header_index ON maker.median_kiss_batch USING btree (header_id);


--
-- Name: median_kiss_batch_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_kiss_batch_log_index ON maker.median_kiss_batch USING btree (log_id);


--
-- Name: median_kiss_batch_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_kiss_batch_msg_sender_index ON maker.median_kiss_batch USING btree (msg_sender);


--
-- Name: median_kiss_single_a_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_kiss_single_a_index ON maker.median_kiss_single USING btree (a);


--
-- Name: median_kiss_single_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_kiss_single_address_index ON maker.median_kiss_single USING btree (address_id);


--
-- Name: median_kiss_single_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_kiss_single_header_index ON maker.median_kiss_single USING btree (header_id);


--
-- Name: median_kiss_single_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_kiss_single_log_index ON maker.median_kiss_single USING btree (log_id);


--
-- Name: median_kiss_single_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_kiss_single_msg_sender_index ON maker.median_kiss_single USING btree (msg_sender);


--
-- Name: median_lift_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_lift_address_index ON maker.median_lift USING btree (address_id);


--
-- Name: median_lift_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_lift_header_index ON maker.median_lift USING btree (header_id);


--
-- Name: median_lift_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_lift_log_index ON maker.median_lift USING btree (log_id);


--
-- Name: median_lift_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_lift_msg_sender_index ON maker.median_lift USING btree (msg_sender);


--
-- Name: median_orcl_a_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_orcl_a_index ON maker.median_orcl USING btree (a);


--
-- Name: median_orcl_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_orcl_address_id_index ON maker.median_orcl USING btree (address_id);


--
-- Name: median_orcl_diff_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_orcl_diff_id_index ON maker.median_orcl USING btree (diff_id);


--
-- Name: median_orcl_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_orcl_header_id_index ON maker.median_orcl USING btree (header_id);


--
-- Name: median_slot_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_slot_address_index ON maker.median_slot USING btree (address_id);


--
-- Name: median_slot_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_slot_header_id_index ON maker.median_slot USING btree (header_id);


--
-- Name: median_slot_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_slot_id_index ON maker.median_slot USING btree (slot_id);


--
-- Name: median_val_address_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_val_address_id_index ON maker.median_val USING btree (address_id);


--
-- Name: median_val_diff_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_val_diff_id_index ON maker.median_val USING btree (diff_id);


--
-- Name: median_val_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX median_val_header_id_index ON maker.median_val USING btree (header_id);


--
-- Name: new_cdp_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX new_cdp_log_index ON maker.new_cdp USING btree (log_id);


--
-- Name: osm_change_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX osm_change_address_index ON maker.osm_change USING btree (address_id);


--
-- Name: osm_change_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX osm_change_header_index ON maker.osm_change USING btree (header_id);


--
-- Name: osm_change_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX osm_change_log_index ON maker.osm_change USING btree (log_id);


--
-- Name: osm_change_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX osm_change_msg_sender_index ON maker.osm_change USING btree (msg_sender);


--
-- Name: osm_change_src_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX osm_change_src_index ON maker.osm_change USING btree (src);


--
-- Name: pot_cage_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_cage_header_index ON maker.pot_cage USING btree (header_id);


--
-- Name: pot_cage_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_cage_log_index ON maker.pot_cage USING btree (log_id);


--
-- Name: pot_cage_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_cage_msg_sender_index ON maker.pot_cage USING btree (msg_sender);


--
-- Name: pot_chi_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_chi_header_id_index ON maker.pot_chi USING btree (header_id);


--
-- Name: pot_drip_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_drip_header_index ON maker.pot_drip USING btree (header_id);


--
-- Name: pot_drip_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_drip_log_index ON maker.pot_drip USING btree (log_id);


--
-- Name: pot_drip_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_drip_msg_sender_index ON maker.pot_drip USING btree (msg_sender);


--
-- Name: pot_dsr_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_dsr_header_id_index ON maker.pot_dsr USING btree (header_id);


--
-- Name: pot_exit_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_exit_header_index ON maker.pot_exit USING btree (header_id);


--
-- Name: pot_exit_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_exit_log_index ON maker.pot_exit USING btree (log_id);


--
-- Name: pot_exit_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_exit_msg_sender_index ON maker.pot_exit USING btree (msg_sender);


--
-- Name: pot_file_dsr_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_file_dsr_header_index ON maker.pot_file_dsr USING btree (header_id);


--
-- Name: pot_file_dsr_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_file_dsr_log_index ON maker.pot_file_dsr USING btree (log_id);


--
-- Name: pot_file_dsr_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_file_dsr_msg_sender_index ON maker.pot_file_dsr USING btree (msg_sender);


--
-- Name: pot_file_vow_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_file_vow_header_index ON maker.pot_file_vow USING btree (header_id);


--
-- Name: pot_file_vow_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_file_vow_log_index ON maker.pot_file_vow USING btree (log_id);


--
-- Name: pot_file_vow_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_file_vow_msg_sender_index ON maker.pot_file_vow USING btree (msg_sender);


--
-- Name: pot_join_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_join_header_index ON maker.pot_join USING btree (header_id);


--
-- Name: pot_join_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_join_log_index ON maker.pot_join USING btree (log_id);


--
-- Name: pot_join_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_join_msg_sender_index ON maker.pot_join USING btree (msg_sender);


--
-- Name: pot_live_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_live_header_id_index ON maker.pot_live USING btree (header_id);


--
-- Name: pot_pie_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_pie_header_id_index ON maker.pot_pie USING btree (header_id);


--
-- Name: pot_rho_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_rho_header_id_index ON maker.pot_rho USING btree (header_id);


--
-- Name: pot_user_pie_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_user_pie_header_id_index ON maker.pot_user_pie USING btree (header_id);


--
-- Name: pot_user_pie_user_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_user_pie_user_index ON maker.pot_user_pie USING btree ("user");


--
-- Name: pot_vat_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_vat_header_id_index ON maker.pot_vat USING btree (header_id);


--
-- Name: pot_vat_vat_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_vat_vat_index ON maker.pot_vat USING btree (vat);


--
-- Name: pot_vow_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_vow_header_id_index ON maker.pot_vow USING btree (header_id);


--
-- Name: pot_vow_vow_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX pot_vow_vow_index ON maker.pot_vow USING btree (vow);


--
-- Name: rely_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX rely_address_index ON maker.rely USING btree (address_id);


--
-- Name: rely_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX rely_header_index ON maker.rely USING btree (header_id);


--
-- Name: rely_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX rely_log_index ON maker.rely USING btree (log_id);


--
-- Name: rely_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX rely_msg_sender_index ON maker.rely USING btree (msg_sender);


--
-- Name: rely_usr_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX rely_usr_index ON maker.rely USING btree (usr);


--
-- Name: set_min_sell_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX set_min_sell_address_index ON maker.set_min_sell USING btree (address_id);


--
-- Name: set_min_sell_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX set_min_sell_header_index ON maker.set_min_sell USING btree (header_id);


--
-- Name: set_min_sell_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX set_min_sell_log_index ON maker.set_min_sell USING btree (log_id);


--
-- Name: set_min_sell_msg_sender; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX set_min_sell_msg_sender ON maker.set_min_sell USING btree (msg_sender);


--
-- Name: set_min_sell_pay_gem_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX set_min_sell_pay_gem_index ON maker.set_min_sell USING btree (pay_gem);


--
-- Name: spot_file_mat_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_mat_header_index ON maker.spot_file_mat USING btree (header_id);


--
-- Name: spot_file_mat_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_mat_ilk_index ON maker.spot_file_mat USING btree (ilk_id);


--
-- Name: spot_file_mat_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_mat_log_index ON maker.spot_file_mat USING btree (log_id);


--
-- Name: spot_file_mat_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_mat_msg_sender_index ON maker.spot_file_mat USING btree (msg_sender);


--
-- Name: spot_file_par_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_par_header_index ON maker.spot_file_par USING btree (header_id);


--
-- Name: spot_file_par_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_par_log_index ON maker.spot_file_par USING btree (log_id);


--
-- Name: spot_file_par_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_par_msg_sender_index ON maker.spot_file_par USING btree (msg_sender);


--
-- Name: spot_file_pip_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_pip_header_index ON maker.spot_file_pip USING btree (header_id);


--
-- Name: spot_file_pip_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_pip_ilk_index ON maker.spot_file_pip USING btree (ilk_id);


--
-- Name: spot_file_pip_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_pip_log_index ON maker.spot_file_pip USING btree (log_id);


--
-- Name: spot_file_pip_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_file_pip_msg_sender_index ON maker.spot_file_pip USING btree (msg_sender);


--
-- Name: spot_ilk_mat_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_ilk_mat_header_id_index ON maker.spot_ilk_mat USING btree (header_id);


--
-- Name: spot_ilk_mat_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_ilk_mat_ilk_index ON maker.spot_ilk_mat USING btree (ilk_id);


--
-- Name: spot_ilk_pip_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_ilk_pip_header_id_index ON maker.spot_ilk_pip USING btree (header_id);


--
-- Name: spot_ilk_pip_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_ilk_pip_ilk_index ON maker.spot_ilk_pip USING btree (ilk_id);


--
-- Name: spot_live_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_live_header_id_index ON maker.spot_live USING btree (header_id);


--
-- Name: spot_par_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_par_header_id_index ON maker.spot_par USING btree (header_id);


--
-- Name: spot_poke_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_poke_header_index ON maker.spot_poke USING btree (header_id);


--
-- Name: spot_poke_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_poke_ilk_index ON maker.spot_poke USING btree (ilk_id);


--
-- Name: spot_poke_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_poke_log_index ON maker.spot_poke USING btree (log_id);


--
-- Name: spot_vat_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX spot_vat_header_id_index ON maker.spot_vat USING btree (header_id);


--
-- Name: tend_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX tend_address_index ON maker.tend USING btree (address_id);


--
-- Name: tend_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX tend_header_index ON maker.tend USING btree (header_id);


--
-- Name: tend_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX tend_log_index ON maker.tend USING btree (log_id);


--
-- Name: tend_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX tend_msg_sender_index ON maker.tend USING btree (msg_sender);


--
-- Name: tick_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX tick_address_index ON maker.tick USING btree (address_id);


--
-- Name: tick_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX tick_bid_id_index ON maker.tick USING btree (bid_id);


--
-- Name: tick_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX tick_header_index ON maker.tick USING btree (header_id);


--
-- Name: tick_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX tick_log_index ON maker.tick USING btree (log_id);


--
-- Name: tick_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX tick_msg_sender_index ON maker.tick USING btree (msg_sender);


--
-- Name: urn_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX urn_ilk_index ON maker.urns USING btree (ilk_id);


--
-- Name: vat_dai_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_dai_header_id_index ON maker.vat_dai USING btree (header_id);


--
-- Name: vat_debt_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_debt_header_id_index ON maker.vat_debt USING btree (header_id);


--
-- Name: vat_deny_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_deny_header_index ON maker.vat_deny USING btree (header_id);


--
-- Name: vat_deny_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_deny_log_index ON maker.vat_deny USING btree (log_id);


--
-- Name: vat_deny_usr_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_deny_usr_index ON maker.vat_deny USING btree (usr);


--
-- Name: vat_file_debt_ceiling_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_file_debt_ceiling_header_index ON maker.vat_file_debt_ceiling USING btree (header_id);


--
-- Name: vat_file_debt_ceiling_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_file_debt_ceiling_log_index ON maker.vat_file_debt_ceiling USING btree (log_id);


--
-- Name: vat_file_ilk_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_file_ilk_header_index ON maker.vat_file_ilk USING btree (header_id);


--
-- Name: vat_file_ilk_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_file_ilk_ilk_index ON maker.vat_file_ilk USING btree (ilk_id);


--
-- Name: vat_file_ilk_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_file_ilk_log_index ON maker.vat_file_ilk USING btree (log_id);


--
-- Name: vat_flux_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_flux_header_index ON maker.vat_flux USING btree (header_id);


--
-- Name: vat_flux_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_flux_ilk_index ON maker.vat_flux USING btree (ilk_id);


--
-- Name: vat_flux_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_flux_log_index ON maker.vat_flux USING btree (log_id);


--
-- Name: vat_fold_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_fold_header_index ON maker.vat_fold USING btree (header_id);


--
-- Name: vat_fold_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_fold_ilk_index ON maker.vat_fold USING btree (ilk_id);


--
-- Name: vat_fold_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_fold_log_index ON maker.vat_fold USING btree (log_id);


--
-- Name: vat_fork_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_fork_header_index ON maker.vat_fork USING btree (header_id);


--
-- Name: vat_fork_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_fork_ilk_index ON maker.vat_fork USING btree (ilk_id);


--
-- Name: vat_fork_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_fork_log_index ON maker.vat_fork USING btree (log_id);


--
-- Name: vat_frob_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_frob_header_index ON maker.vat_frob USING btree (header_id);


--
-- Name: vat_frob_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_frob_log_index ON maker.vat_frob USING btree (log_id);


--
-- Name: vat_frob_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_frob_urn_index ON maker.vat_frob USING btree (urn_id);


--
-- Name: vat_gem_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_gem_header_id_index ON maker.vat_gem USING btree (header_id);


--
-- Name: vat_gem_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_gem_ilk_index ON maker.vat_gem USING btree (ilk_id);


--
-- Name: vat_grab_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_grab_header_index ON maker.vat_grab USING btree (header_id);


--
-- Name: vat_grab_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_grab_log_index ON maker.vat_grab USING btree (log_id);


--
-- Name: vat_grab_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_grab_urn_index ON maker.vat_grab USING btree (urn_id);


--
-- Name: vat_heal_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_heal_header_index ON maker.vat_heal USING btree (header_id);


--
-- Name: vat_heal_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_heal_log_index ON maker.vat_heal USING btree (log_id);


--
-- Name: vat_hope_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_hope_header_index ON maker.vat_hope USING btree (header_id);


--
-- Name: vat_hope_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_hope_log_index ON maker.vat_hope USING btree (log_id);


--
-- Name: vat_hope_usr_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_hope_usr_index ON maker.vat_hope USING btree (usr);


--
-- Name: vat_ilk_art_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_art_header_id_index ON maker.vat_ilk_art USING btree (header_id);


--
-- Name: vat_ilk_art_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_art_ilk_index ON maker.vat_ilk_art USING btree (ilk_id);


--
-- Name: vat_ilk_dust_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_dust_header_id_index ON maker.vat_ilk_dust USING btree (header_id);


--
-- Name: vat_ilk_dust_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_dust_ilk_index ON maker.vat_ilk_dust USING btree (ilk_id);


--
-- Name: vat_ilk_line_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_line_header_id_index ON maker.vat_ilk_line USING btree (header_id);


--
-- Name: vat_ilk_line_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_line_ilk_index ON maker.vat_ilk_line USING btree (ilk_id);


--
-- Name: vat_ilk_rate_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_rate_header_id_index ON maker.vat_ilk_rate USING btree (header_id);


--
-- Name: vat_ilk_rate_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_rate_ilk_index ON maker.vat_ilk_rate USING btree (ilk_id);


--
-- Name: vat_ilk_spot_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_ilk_spot_header_id_index ON maker.vat_ilk_spot USING btree (header_id);


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
-- Name: vat_init_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_init_log_index ON maker.vat_init USING btree (log_id);


--
-- Name: vat_line_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_line_header_id_index ON maker.vat_line USING btree (header_id);


--
-- Name: vat_live_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_live_header_id_index ON maker.vat_live USING btree (header_id);


--
-- Name: vat_move_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_move_header_index ON maker.vat_move USING btree (header_id);


--
-- Name: vat_move_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_move_log_index ON maker.vat_move USING btree (log_id);


--
-- Name: vat_nope_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_nope_header_index ON maker.vat_nope USING btree (header_id);


--
-- Name: vat_nope_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_nope_log_index ON maker.vat_nope USING btree (log_id);


--
-- Name: vat_nope_usr_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_nope_usr_index ON maker.vat_nope USING btree (usr);


--
-- Name: vat_rely_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_rely_header_index ON maker.vat_rely USING btree (header_id);


--
-- Name: vat_rely_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_rely_log_index ON maker.vat_rely USING btree (log_id);


--
-- Name: vat_rely_usr_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_rely_usr_index ON maker.vat_rely USING btree (usr);


--
-- Name: vat_sin_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_sin_header_id_index ON maker.vat_sin USING btree (header_id);


--
-- Name: vat_slip_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_slip_header_index ON maker.vat_slip USING btree (header_id);


--
-- Name: vat_slip_ilk_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_slip_ilk_index ON maker.vat_slip USING btree (ilk_id);


--
-- Name: vat_slip_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_slip_log_index ON maker.vat_slip USING btree (log_id);


--
-- Name: vat_suck_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_suck_header_index ON maker.vat_suck USING btree (header_id);


--
-- Name: vat_suck_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_suck_log_index ON maker.vat_suck USING btree (log_id);


--
-- Name: vat_urn_art_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_urn_art_header_id_index ON maker.vat_urn_art USING btree (header_id);


--
-- Name: vat_urn_art_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_urn_art_urn_index ON maker.vat_urn_art USING btree (urn_id);


--
-- Name: vat_urn_ink_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_urn_ink_header_id_index ON maker.vat_urn_ink USING btree (header_id);


--
-- Name: vat_urn_ink_urn_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_urn_ink_urn_index ON maker.vat_urn_ink USING btree (urn_id);


--
-- Name: vat_vice_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vat_vice_header_id_index ON maker.vat_vice USING btree (header_id);


--
-- Name: vow_ash_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_ash_header_id_index ON maker.vow_ash USING btree (header_id);


--
-- Name: vow_bump_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_bump_header_id_index ON maker.vow_bump USING btree (header_id);


--
-- Name: vow_dump_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_dump_header_id_index ON maker.vow_dump USING btree (header_id);


--
-- Name: vow_fess_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_fess_header_index ON maker.vow_fess USING btree (header_id);


--
-- Name: vow_fess_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_fess_log_index ON maker.vow_fess USING btree (log_id);


--
-- Name: vow_fess_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_fess_msg_sender_index ON maker.vow_fess USING btree (msg_sender);


--
-- Name: vow_file_auction_address_data_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_file_auction_address_data_index ON maker.vow_file_auction_address USING btree (data);


--
-- Name: vow_file_auction_address_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_file_auction_address_header_index ON maker.vow_file_auction_address USING btree (header_id);


--
-- Name: vow_file_auction_address_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_file_auction_address_log_index ON maker.vow_file_auction_address USING btree (log_id);


--
-- Name: vow_file_auction_address_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_file_auction_address_msg_sender_index ON maker.vow_file_auction_address USING btree (msg_sender);


--
-- Name: vow_file_auction_attributes_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_file_auction_attributes_header_index ON maker.vow_file_auction_attributes USING btree (header_id);


--
-- Name: vow_file_auction_attributes_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_file_auction_attributes_log_index ON maker.vow_file_auction_attributes USING btree (log_id);


--
-- Name: vow_file_auction_attributes_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_file_auction_attributes_msg_sender_index ON maker.vow_file_auction_attributes USING btree (msg_sender);


--
-- Name: vow_flapper_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_flapper_header_id_index ON maker.vow_flapper USING btree (header_id);


--
-- Name: vow_flog_era_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_flog_era_index ON maker.vow_flog USING btree (era);


--
-- Name: vow_flog_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_flog_header_index ON maker.vow_flog USING btree (header_id);


--
-- Name: vow_flog_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_flog_log_index ON maker.vow_flog USING btree (log_id);


--
-- Name: vow_flog_msg_sender_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_flog_msg_sender_index ON maker.vow_flog USING btree (msg_sender);


--
-- Name: vow_flopper_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_flopper_header_id_index ON maker.vow_flopper USING btree (header_id);


--
-- Name: vow_heal_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_heal_header_index ON maker.vow_heal USING btree (header_id);


--
-- Name: vow_heal_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_heal_log_index ON maker.vow_heal USING btree (log_id);


--
-- Name: vow_heal_msg_sender; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_heal_msg_sender ON maker.vow_heal USING btree (msg_sender);


--
-- Name: vow_hump_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_hump_header_id_index ON maker.vow_hump USING btree (header_id);


--
-- Name: vow_live_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_live_header_id_index ON maker.vow_live USING btree (header_id);


--
-- Name: vow_sin_integer_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_sin_integer_header_id_index ON maker.vow_sin_integer USING btree (header_id);


--
-- Name: vow_sin_mapping_era_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_sin_mapping_era_index ON maker.vow_sin_mapping USING btree (era);


--
-- Name: vow_sin_mapping_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_sin_mapping_header_id_index ON maker.vow_sin_mapping USING btree (header_id);


--
-- Name: vow_sump_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_sump_header_id_index ON maker.vow_sump USING btree (header_id);


--
-- Name: vow_vat_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_vat_header_id_index ON maker.vow_vat USING btree (header_id);


--
-- Name: vow_wait_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX vow_wait_header_id_index ON maker.vow_wait USING btree (header_id);


--
-- Name: wards_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX wards_address_index ON maker.wards USING btree (address_id);


--
-- Name: wards_header_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX wards_header_id_index ON maker.wards USING btree (header_id);


--
-- Name: wards_usr_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX wards_usr_index ON maker.wards USING btree (usr);


--
-- Name: yank_address_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX yank_address_index ON maker.yank USING btree (address_id);


--
-- Name: yank_bid_id_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX yank_bid_id_index ON maker.yank USING btree (bid_id);


--
-- Name: yank_header_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX yank_header_index ON maker.yank USING btree (header_id);


--
-- Name: yank_log_index; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX yank_log_index ON maker.yank USING btree (log_id);


--
-- Name: yank_msg_sender; Type: INDEX; Schema: maker; Owner: -
--

CREATE INDEX yank_msg_sender ON maker.yank USING btree (msg_sender);


--
-- Name: event_logs_address; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX event_logs_address ON public.event_logs USING btree (address);


--
-- Name: event_logs_transaction; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX event_logs_transaction ON public.event_logs USING btree (tx_hash);


--
-- Name: event_logs_untransformed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX event_logs_untransformed ON public.event_logs USING btree (transformed) WHERE (transformed = false);


--
-- Name: headers_block_number; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX headers_block_number ON public.headers USING btree (block_number);


--
-- Name: headers_block_timestamp_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX headers_block_timestamp_index ON public.headers USING btree (block_timestamp);


--
-- Name: headers_check_count; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX headers_check_count ON public.headers USING btree (check_count);


--
-- Name: headers_eth_node; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX headers_eth_node ON public.headers USING btree (eth_node_id);


--
-- Name: receipts_contract_address; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX receipts_contract_address ON public.receipts USING btree (contract_address_id);


--
-- Name: receipts_transaction; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX receipts_transaction ON public.receipts USING btree (transaction_id);


--
-- Name: storage_diff_checked_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX storage_diff_checked_index ON public.storage_diff USING btree (checked) WHERE (checked = false);


--
-- Name: storage_diff_eth_node; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX storage_diff_eth_node ON public.storage_diff USING btree (eth_node_id);


--
-- Name: transactions_header; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX transactions_header ON public.transactions USING btree (header_id);


--
-- Name: deal deal; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER deal AFTER INSERT ON maker.deal FOR EACH ROW EXECUTE PROCEDURE maker.update_bid_tick_deal_yank_event('deal');


--
-- Name: dent dent; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER dent AFTER INSERT ON maker.dent FOR EACH ROW EXECUTE PROCEDURE maker.update_bid_kick_tend_dent_event('dent');


--
-- Name: flap_bid_bid flap_bid; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_bid AFTER INSERT OR DELETE OR UPDATE ON maker.flap_bid_bid FOR EACH ROW EXECUTE PROCEDURE maker.update_flap_bids();


--
-- Name: flap_kick flap_created_trigger; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_created_trigger AFTER INSERT OR DELETE ON maker.flap_kick FOR EACH ROW EXECUTE PROCEDURE maker.update_flap_created();


--
-- Name: flap_bid_end flap_end; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_end AFTER INSERT OR DELETE OR UPDATE ON maker.flap_bid_end FOR EACH ROW EXECUTE PROCEDURE maker.update_flap_ends();


--
-- Name: flap_bid_guy flap_guy; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_guy AFTER INSERT OR DELETE OR UPDATE ON maker.flap_bid_guy FOR EACH ROW EXECUTE PROCEDURE maker.update_flap_guys();


--
-- Name: flap_kick flap_kick; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_kick AFTER INSERT ON maker.flap_kick FOR EACH ROW EXECUTE PROCEDURE maker.update_bid_kick_tend_dent_event('kick');


--
-- Name: flap_bid_lot flap_lot; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_lot AFTER INSERT OR DELETE OR UPDATE ON maker.flap_bid_lot FOR EACH ROW EXECUTE PROCEDURE maker.update_flap_lots();


--
-- Name: flap_bid_tic flap_tic; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flap_tic AFTER INSERT OR DELETE OR UPDATE ON maker.flap_bid_tic FOR EACH ROW EXECUTE PROCEDURE maker.update_flap_tics();


--
-- Name: flip_bid_bid flip_bid; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_bid AFTER INSERT OR DELETE OR UPDATE ON maker.flip_bid_bid FOR EACH ROW EXECUTE PROCEDURE maker.update_flip_bids();


--
-- Name: flip_kick flip_created_trigger; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_created_trigger AFTER INSERT OR DELETE ON maker.flip_kick FOR EACH ROW EXECUTE PROCEDURE maker.update_flip_created();


--
-- Name: flip_bid_end flip_end; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_end AFTER INSERT OR DELETE OR UPDATE ON maker.flip_bid_end FOR EACH ROW EXECUTE PROCEDURE maker.update_flip_ends();


--
-- Name: flip_bid_gal flip_gal; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_gal AFTER INSERT OR DELETE OR UPDATE ON maker.flip_bid_gal FOR EACH ROW EXECUTE PROCEDURE maker.update_flip_gals();


--
-- Name: flip_bid_guy flip_guy; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_guy AFTER INSERT OR DELETE OR UPDATE ON maker.flip_bid_guy FOR EACH ROW EXECUTE PROCEDURE maker.update_flip_guys();


--
-- Name: flip_ilk flip_ilk; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_ilk AFTER INSERT OR DELETE ON maker.flip_ilk FOR EACH ROW EXECUTE PROCEDURE maker.update_bid_event_ilk();


--
-- Name: flip_kick flip_kick; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_kick AFTER INSERT ON maker.flip_kick FOR EACH ROW EXECUTE PROCEDURE maker.update_bid_kick_tend_dent_event('kick');


--
-- Name: flip_bid_lot flip_lot; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_lot AFTER INSERT OR DELETE OR UPDATE ON maker.flip_bid_lot FOR EACH ROW EXECUTE PROCEDURE maker.update_flip_lots();


--
-- Name: flip_bid_tab flip_tab; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_tab AFTER INSERT OR DELETE OR UPDATE ON maker.flip_bid_tab FOR EACH ROW EXECUTE PROCEDURE maker.update_flip_tabs();


--
-- Name: flip_bid_tic flip_tic; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_tic AFTER INSERT OR DELETE OR UPDATE ON maker.flip_bid_tic FOR EACH ROW EXECUTE PROCEDURE maker.update_flip_tics();


--
-- Name: flip_bid_usr flip_urn; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_urn AFTER INSERT OR DELETE ON maker.flip_bid_usr FOR EACH ROW EXECUTE PROCEDURE maker.update_bid_event_urn();


--
-- Name: flip_bid_usr flip_usr; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flip_usr AFTER INSERT OR DELETE OR UPDATE ON maker.flip_bid_usr FOR EACH ROW EXECUTE PROCEDURE maker.update_flip_usrs();


--
-- Name: flop_bid_bid flop_bid; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_bid AFTER INSERT OR DELETE OR UPDATE ON maker.flop_bid_bid FOR EACH ROW EXECUTE PROCEDURE maker.update_flop_bids();


--
-- Name: flop_kick flop_created_trigger; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_created_trigger AFTER INSERT OR DELETE ON maker.flop_kick FOR EACH ROW EXECUTE PROCEDURE maker.update_flop_created();


--
-- Name: flop_bid_end flop_end; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_end AFTER INSERT OR DELETE OR UPDATE ON maker.flop_bid_end FOR EACH ROW EXECUTE PROCEDURE maker.update_flop_ends();


--
-- Name: flop_bid_guy flop_guy; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_guy AFTER INSERT OR DELETE OR UPDATE ON maker.flop_bid_guy FOR EACH ROW EXECUTE PROCEDURE maker.update_flop_guys();


--
-- Name: flop_kick flop_kick; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_kick AFTER INSERT ON maker.flop_kick FOR EACH ROW EXECUTE PROCEDURE maker.update_bid_kick_tend_dent_event('kick');


--
-- Name: flop_bid_lot flop_lot; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_lot AFTER INSERT OR DELETE OR UPDATE ON maker.flop_bid_lot FOR EACH ROW EXECUTE PROCEDURE maker.update_flop_lots();


--
-- Name: flop_bid_tic flop_tic; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER flop_tic AFTER INSERT OR DELETE OR UPDATE ON maker.flop_bid_tic FOR EACH ROW EXECUTE PROCEDURE maker.update_flop_tics();


--
-- Name: vat_ilk_art ilk_art; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_art AFTER INSERT OR DELETE OR UPDATE ON maker.vat_ilk_art FOR EACH ROW EXECUTE PROCEDURE maker.update_ilk_arts();


--
-- Name: cat_ilk_chop ilk_chop; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_chop AFTER INSERT OR DELETE OR UPDATE ON maker.cat_ilk_chop FOR EACH ROW EXECUTE PROCEDURE maker.update_ilk_chops();


--
-- Name: vat_ilk_dust ilk_dust; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_dust AFTER INSERT OR DELETE OR UPDATE ON maker.vat_ilk_dust FOR EACH ROW EXECUTE PROCEDURE maker.update_ilk_dusts();


--
-- Name: jug_ilk_duty ilk_duty; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_duty AFTER INSERT OR DELETE OR UPDATE ON maker.jug_ilk_duty FOR EACH ROW EXECUTE PROCEDURE maker.update_ilk_duties();


--
-- Name: cat_ilk_flip ilk_flip; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_flip AFTER INSERT OR DELETE OR UPDATE ON maker.cat_ilk_flip FOR EACH ROW EXECUTE PROCEDURE maker.update_ilk_flips();


--
-- Name: vat_init ilk_init; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_init AFTER INSERT OR DELETE OR UPDATE ON maker.vat_init FOR EACH ROW EXECUTE PROCEDURE maker.update_time_created();


--
-- Name: vat_ilk_line ilk_line; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_line AFTER INSERT OR DELETE OR UPDATE ON maker.vat_ilk_line FOR EACH ROW EXECUTE PROCEDURE maker.update_ilk_lines();


--
-- Name: cat_ilk_lump ilk_lump; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_lump AFTER INSERT OR DELETE OR UPDATE ON maker.cat_ilk_lump FOR EACH ROW EXECUTE PROCEDURE maker.update_ilk_lumps();


--
-- Name: spot_ilk_mat ilk_mat; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_mat AFTER INSERT OR DELETE OR UPDATE ON maker.spot_ilk_mat FOR EACH ROW EXECUTE PROCEDURE maker.update_ilk_mats();


--
-- Name: spot_ilk_pip ilk_pip; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_pip AFTER INSERT OR DELETE OR UPDATE ON maker.spot_ilk_pip FOR EACH ROW EXECUTE PROCEDURE maker.update_ilk_pips();


--
-- Name: vat_ilk_rate ilk_rate; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_rate AFTER INSERT OR DELETE OR UPDATE ON maker.vat_ilk_rate FOR EACH ROW EXECUTE PROCEDURE maker.update_ilk_rates();


--
-- Name: jug_ilk_rho ilk_rho; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_rho AFTER INSERT OR DELETE OR UPDATE ON maker.jug_ilk_rho FOR EACH ROW EXECUTE PROCEDURE maker.update_ilk_rhos();


--
-- Name: vat_ilk_spot ilk_spot; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER ilk_spot AFTER INSERT OR DELETE OR UPDATE ON maker.vat_ilk_spot FOR EACH ROW EXECUTE PROCEDURE maker.update_ilk_spots();


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
-- Name: tend tend; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER tend AFTER INSERT ON maker.tend FOR EACH ROW EXECUTE PROCEDURE maker.update_bid_kick_tend_dent_event('tend');


--
-- Name: tick tick; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER tick AFTER INSERT ON maker.tick FOR EACH ROW EXECUTE PROCEDURE maker.update_bid_tick_deal_yank_event('tick');


--
-- Name: vat_urn_art urn_art; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER urn_art AFTER INSERT OR DELETE OR UPDATE ON maker.vat_urn_art FOR EACH ROW EXECUTE PROCEDURE maker.update_urn_arts();


--
-- Name: vat_urn_ink urn_ink; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER urn_ink AFTER INSERT OR DELETE OR UPDATE ON maker.vat_urn_ink FOR EACH ROW EXECUTE PROCEDURE maker.update_urn_inks();


--
-- Name: yank yank; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER yank AFTER INSERT ON maker.yank FOR EACH ROW EXECUTE PROCEDURE maker.update_bid_tick_deal_yank_event('yank');


--
-- Name: headers header_updated; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER header_updated BEFORE UPDATE ON public.headers FOR EACH ROW EXECUTE PROCEDURE public.set_header_updated();


--
-- Name: auction_file auction_file_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.auction_file
    ADD CONSTRAINT auction_file_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: auction_file auction_file_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.auction_file
    ADD CONSTRAINT auction_file_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: auction_file auction_file_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.auction_file
    ADD CONSTRAINT auction_file_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: auction_file auction_file_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.auction_file
    ADD CONSTRAINT auction_file_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: bid_event bid_event_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bid_event
    ADD CONSTRAINT bid_event_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: bite bite_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite
    ADD CONSTRAINT bite_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: bite bite_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite
    ADD CONSTRAINT bite_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: bite bite_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite
    ADD CONSTRAINT bite_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT cat_file_chop_lump_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: cat_file_chop_lump cat_file_chop_lump_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_chop_lump
    ADD CONSTRAINT cat_file_chop_lump_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT cat_file_flip_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: cat_file_flip cat_file_flip_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_flip
    ADD CONSTRAINT cat_file_flip_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: cat_file_vow cat_file_vow_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_vow
    ADD CONSTRAINT cat_file_vow_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cat_file_vow cat_file_vow_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_vow
    ADD CONSTRAINT cat_file_vow_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: cat_file_vow cat_file_vow_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_vow
    ADD CONSTRAINT cat_file_vow_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: cat_ilk_chop cat_ilk_chop_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_chop
    ADD CONSTRAINT cat_ilk_chop_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cat_ilk_chop cat_ilk_chop_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_chop
    ADD CONSTRAINT cat_ilk_chop_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cat_ilk_chop cat_ilk_chop_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_chop
    ADD CONSTRAINT cat_ilk_chop_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: cat_ilk_flip cat_ilk_flip_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_flip
    ADD CONSTRAINT cat_ilk_flip_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cat_ilk_flip cat_ilk_flip_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_flip
    ADD CONSTRAINT cat_ilk_flip_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cat_ilk_flip cat_ilk_flip_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_flip
    ADD CONSTRAINT cat_ilk_flip_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: cat_ilk_lump cat_ilk_lump_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_lump
    ADD CONSTRAINT cat_ilk_lump_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cat_ilk_lump cat_ilk_lump_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_lump
    ADD CONSTRAINT cat_ilk_lump_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cat_ilk_lump cat_ilk_lump_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_lump
    ADD CONSTRAINT cat_ilk_lump_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: cat_live cat_live_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_live
    ADD CONSTRAINT cat_live_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cat_live cat_live_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_live
    ADD CONSTRAINT cat_live_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cat_vat cat_vat_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vat
    ADD CONSTRAINT cat_vat_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cat_vat cat_vat_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vat
    ADD CONSTRAINT cat_vat_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cat_vow cat_vow_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vow
    ADD CONSTRAINT cat_vow_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cat_vow cat_vow_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vow
    ADD CONSTRAINT cat_vow_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_cdpi cdp_manager_cdpi_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_cdpi
    ADD CONSTRAINT cdp_manager_cdpi_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_cdpi cdp_manager_cdpi_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_cdpi
    ADD CONSTRAINT cdp_manager_cdpi_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_count cdp_manager_count_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_count
    ADD CONSTRAINT cdp_manager_count_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_count cdp_manager_count_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_count
    ADD CONSTRAINT cdp_manager_count_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_first cdp_manager_first_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_first
    ADD CONSTRAINT cdp_manager_first_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_first cdp_manager_first_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_first
    ADD CONSTRAINT cdp_manager_first_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_ilks cdp_manager_ilks_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_ilks
    ADD CONSTRAINT cdp_manager_ilks_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_ilks cdp_manager_ilks_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_ilks
    ADD CONSTRAINT cdp_manager_ilks_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_ilks cdp_manager_ilks_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_ilks
    ADD CONSTRAINT cdp_manager_ilks_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_last cdp_manager_last_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_last
    ADD CONSTRAINT cdp_manager_last_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_last cdp_manager_last_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_last
    ADD CONSTRAINT cdp_manager_last_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_list_next cdp_manager_list_next_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_next
    ADD CONSTRAINT cdp_manager_list_next_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_list_next cdp_manager_list_next_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_next
    ADD CONSTRAINT cdp_manager_list_next_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_list_prev cdp_manager_list_prev_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_prev
    ADD CONSTRAINT cdp_manager_list_prev_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_list_prev cdp_manager_list_prev_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_list_prev
    ADD CONSTRAINT cdp_manager_list_prev_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_owns cdp_manager_owns_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_owns
    ADD CONSTRAINT cdp_manager_owns_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_owns cdp_manager_owns_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_owns
    ADD CONSTRAINT cdp_manager_owns_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_urns cdp_manager_urns_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_urns
    ADD CONSTRAINT cdp_manager_urns_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_urns cdp_manager_urns_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_urns
    ADD CONSTRAINT cdp_manager_urns_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_vat cdp_manager_vat_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_vat
    ADD CONSTRAINT cdp_manager_vat_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: cdp_manager_vat cdp_manager_vat_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cdp_manager_vat
    ADD CONSTRAINT cdp_manager_vat_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT deal_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: deal deal_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deal
    ADD CONSTRAINT deal_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT dent_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: dent dent_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.dent
    ADD CONSTRAINT dent_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: deny deny_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deny
    ADD CONSTRAINT deny_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: deny deny_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deny
    ADD CONSTRAINT deny_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: deny deny_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deny
    ADD CONSTRAINT deny_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: deny deny_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deny
    ADD CONSTRAINT deny_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: deny deny_usr_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deny
    ADD CONSTRAINT deny_usr_fkey FOREIGN KEY (usr) REFERENCES public.addresses(id) ON DELETE CASCADE;


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
-- Name: flap_beg flap_beg_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_beg
    ADD CONSTRAINT flap_beg_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flap_beg flap_beg_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_beg
    ADD CONSTRAINT flap_beg_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flap_bid_bid flap_bid_bid_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_bid
    ADD CONSTRAINT flap_bid_bid_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_bid_bid flap_bid_bid_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_bid
    ADD CONSTRAINT flap_bid_bid_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flap_bid_bid flap_bid_bid_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_bid
    ADD CONSTRAINT flap_bid_bid_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flap_bid_end flap_bid_end_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_end
    ADD CONSTRAINT flap_bid_end_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_bid_end flap_bid_end_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_end
    ADD CONSTRAINT flap_bid_end_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flap_bid_end flap_bid_end_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_end
    ADD CONSTRAINT flap_bid_end_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flap_bid_guy flap_bid_guy_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_guy
    ADD CONSTRAINT flap_bid_guy_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_bid_guy flap_bid_guy_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_guy
    ADD CONSTRAINT flap_bid_guy_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flap_bid_guy flap_bid_guy_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_guy
    ADD CONSTRAINT flap_bid_guy_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flap_bid_lot flap_bid_lot_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_lot
    ADD CONSTRAINT flap_bid_lot_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_bid_lot flap_bid_lot_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_lot
    ADD CONSTRAINT flap_bid_lot_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flap_bid_lot flap_bid_lot_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_lot
    ADD CONSTRAINT flap_bid_lot_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flap_bid_tic flap_bid_tic_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_tic
    ADD CONSTRAINT flap_bid_tic_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_bid_tic flap_bid_tic_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_tic
    ADD CONSTRAINT flap_bid_tic_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flap_bid_tic flap_bid_tic_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_bid_tic
    ADD CONSTRAINT flap_bid_tic_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flap_gem flap_gem_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_gem
    ADD CONSTRAINT flap_gem_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_gem flap_gem_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_gem
    ADD CONSTRAINT flap_gem_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flap_gem flap_gem_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_gem
    ADD CONSTRAINT flap_gem_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT flap_kick_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: flap_kicks flap_kicks_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kicks
    ADD CONSTRAINT flap_kicks_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_kicks flap_kicks_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kicks
    ADD CONSTRAINT flap_kicks_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flap_kicks flap_kicks_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kicks
    ADD CONSTRAINT flap_kicks_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flap_live flap_live_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_live
    ADD CONSTRAINT flap_live_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_live flap_live_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_live
    ADD CONSTRAINT flap_live_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flap_live flap_live_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_live
    ADD CONSTRAINT flap_live_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flap_tau flap_tau_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_tau
    ADD CONSTRAINT flap_tau_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_tau flap_tau_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_tau
    ADD CONSTRAINT flap_tau_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flap_tau flap_tau_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_tau
    ADD CONSTRAINT flap_tau_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flap_ttl flap_ttl_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_ttl
    ADD CONSTRAINT flap_ttl_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_ttl flap_ttl_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_ttl
    ADD CONSTRAINT flap_ttl_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flap_ttl flap_ttl_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_ttl
    ADD CONSTRAINT flap_ttl_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flap_vat flap_vat_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_vat
    ADD CONSTRAINT flap_vat_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flap_vat flap_vat_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_vat
    ADD CONSTRAINT flap_vat_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flap_vat flap_vat_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_vat
    ADD CONSTRAINT flap_vat_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
-- Name: flip_beg flip_beg_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_beg
    ADD CONSTRAINT flip_beg_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_beg flip_beg_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_beg
    ADD CONSTRAINT flip_beg_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_bid_bid flip_bid_bid_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_bid
    ADD CONSTRAINT flip_bid_bid_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_bid flip_bid_bid_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_bid
    ADD CONSTRAINT flip_bid_bid_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_bid_bid flip_bid_bid_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_bid
    ADD CONSTRAINT flip_bid_bid_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_bid_end flip_bid_end_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_end
    ADD CONSTRAINT flip_bid_end_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_end flip_bid_end_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_end
    ADD CONSTRAINT flip_bid_end_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_bid_end flip_bid_end_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_end
    ADD CONSTRAINT flip_bid_end_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_bid_gal flip_bid_gal_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_gal
    ADD CONSTRAINT flip_bid_gal_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_gal flip_bid_gal_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_gal
    ADD CONSTRAINT flip_bid_gal_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_bid_gal flip_bid_gal_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_gal
    ADD CONSTRAINT flip_bid_gal_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_bid_guy flip_bid_guy_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_guy
    ADD CONSTRAINT flip_bid_guy_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_guy flip_bid_guy_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_guy
    ADD CONSTRAINT flip_bid_guy_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_bid_guy flip_bid_guy_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_guy
    ADD CONSTRAINT flip_bid_guy_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_bid_lot flip_bid_lot_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_lot
    ADD CONSTRAINT flip_bid_lot_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_lot flip_bid_lot_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_lot
    ADD CONSTRAINT flip_bid_lot_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_bid_lot flip_bid_lot_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_lot
    ADD CONSTRAINT flip_bid_lot_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_bid_tab flip_bid_tab_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tab
    ADD CONSTRAINT flip_bid_tab_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_tab flip_bid_tab_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tab
    ADD CONSTRAINT flip_bid_tab_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_bid_tab flip_bid_tab_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tab
    ADD CONSTRAINT flip_bid_tab_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_bid_tic flip_bid_tic_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tic
    ADD CONSTRAINT flip_bid_tic_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_tic flip_bid_tic_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tic
    ADD CONSTRAINT flip_bid_tic_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_bid_tic flip_bid_tic_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_tic
    ADD CONSTRAINT flip_bid_tic_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_bid_usr flip_bid_usr_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_usr
    ADD CONSTRAINT flip_bid_usr_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_bid_usr flip_bid_usr_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_usr
    ADD CONSTRAINT flip_bid_usr_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_bid_usr flip_bid_usr_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_bid_usr
    ADD CONSTRAINT flip_bid_usr_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_ilk flip_ilk_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ilk
    ADD CONSTRAINT flip_ilk_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_ilk flip_ilk_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ilk
    ADD CONSTRAINT flip_ilk_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_ilk flip_ilk_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ilk
    ADD CONSTRAINT flip_ilk_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT flip_kick_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: flip_kicks flip_kicks_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kicks
    ADD CONSTRAINT flip_kicks_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_kicks flip_kicks_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kicks
    ADD CONSTRAINT flip_kicks_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_kicks flip_kicks_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kicks
    ADD CONSTRAINT flip_kicks_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_tau flip_tau_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_tau
    ADD CONSTRAINT flip_tau_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_tau flip_tau_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_tau
    ADD CONSTRAINT flip_tau_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_tau flip_tau_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_tau
    ADD CONSTRAINT flip_tau_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_ttl flip_ttl_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ttl
    ADD CONSTRAINT flip_ttl_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_ttl flip_ttl_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ttl
    ADD CONSTRAINT flip_ttl_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_ttl flip_ttl_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_ttl
    ADD CONSTRAINT flip_ttl_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_vat flip_vat_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_vat
    ADD CONSTRAINT flip_vat_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flip_vat flip_vat_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_vat
    ADD CONSTRAINT flip_vat_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flip_vat flip_vat_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_vat
    ADD CONSTRAINT flip_vat_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
-- Name: flop_beg flop_beg_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_beg
    ADD CONSTRAINT flop_beg_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_beg flop_beg_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_beg
    ADD CONSTRAINT flop_beg_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_bid_bid flop_bid_bid_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_bid
    ADD CONSTRAINT flop_bid_bid_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_bid_bid flop_bid_bid_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_bid
    ADD CONSTRAINT flop_bid_bid_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_bid_bid flop_bid_bid_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_bid
    ADD CONSTRAINT flop_bid_bid_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_bid_end flop_bid_end_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_end
    ADD CONSTRAINT flop_bid_end_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_bid_end flop_bid_end_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_end
    ADD CONSTRAINT flop_bid_end_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_bid_end flop_bid_end_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_end
    ADD CONSTRAINT flop_bid_end_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_bid_guy flop_bid_guy_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_guy
    ADD CONSTRAINT flop_bid_guy_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_bid_guy flop_bid_guy_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_guy
    ADD CONSTRAINT flop_bid_guy_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_bid_guy flop_bid_guy_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_guy
    ADD CONSTRAINT flop_bid_guy_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_bid_lot flop_bid_lot_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_lot
    ADD CONSTRAINT flop_bid_lot_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_bid_lot flop_bid_lot_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_lot
    ADD CONSTRAINT flop_bid_lot_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_bid_lot flop_bid_lot_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_lot
    ADD CONSTRAINT flop_bid_lot_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_bid_tic flop_bid_tic_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_tic
    ADD CONSTRAINT flop_bid_tic_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_bid_tic flop_bid_tic_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_tic
    ADD CONSTRAINT flop_bid_tic_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_bid_tic flop_bid_tic_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_bid_tic
    ADD CONSTRAINT flop_bid_tic_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_gem flop_gem_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_gem
    ADD CONSTRAINT flop_gem_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_gem flop_gem_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_gem
    ADD CONSTRAINT flop_gem_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_gem flop_gem_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_gem
    ADD CONSTRAINT flop_gem_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT flop_kick_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: flop_kicks flop_kicks_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kicks
    ADD CONSTRAINT flop_kicks_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_kicks flop_kicks_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kicks
    ADD CONSTRAINT flop_kicks_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_kicks flop_kicks_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kicks
    ADD CONSTRAINT flop_kicks_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_live flop_live_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_live
    ADD CONSTRAINT flop_live_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_live flop_live_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_live
    ADD CONSTRAINT flop_live_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_live flop_live_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_live
    ADD CONSTRAINT flop_live_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_pad flop_pad_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_pad
    ADD CONSTRAINT flop_pad_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_pad flop_pad_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_pad
    ADD CONSTRAINT flop_pad_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_pad flop_pad_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_pad
    ADD CONSTRAINT flop_pad_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_tau flop_tau_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_tau
    ADD CONSTRAINT flop_tau_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_tau flop_tau_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_tau
    ADD CONSTRAINT flop_tau_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_tau flop_tau_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_tau
    ADD CONSTRAINT flop_tau_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_ttl flop_ttl_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_ttl
    ADD CONSTRAINT flop_ttl_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_ttl flop_ttl_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_ttl
    ADD CONSTRAINT flop_ttl_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_ttl flop_ttl_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_ttl
    ADD CONSTRAINT flop_ttl_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_vat flop_vat_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vat
    ADD CONSTRAINT flop_vat_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_vat flop_vat_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vat
    ADD CONSTRAINT flop_vat_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_vat flop_vat_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vat
    ADD CONSTRAINT flop_vat_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_vow flop_vow_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vow
    ADD CONSTRAINT flop_vow_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: flop_vow flop_vow_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vow
    ADD CONSTRAINT flop_vow_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: flop_vow flop_vow_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_vow
    ADD CONSTRAINT flop_vow_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_base jug_base_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_base
    ADD CONSTRAINT jug_base_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: jug_base jug_base_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_base
    ADD CONSTRAINT jug_base_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT jug_drip_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: jug_drip jug_drip_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_drip
    ADD CONSTRAINT jug_drip_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: jug_file_base jug_file_base_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_base
    ADD CONSTRAINT jug_file_base_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_file_base jug_file_base_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_base
    ADD CONSTRAINT jug_file_base_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: jug_file_base jug_file_base_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_base
    ADD CONSTRAINT jug_file_base_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT jug_file_ilk_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: jug_file_ilk jug_file_ilk_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_ilk
    ADD CONSTRAINT jug_file_ilk_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: jug_file_vow jug_file_vow_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow
    ADD CONSTRAINT jug_file_vow_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_file_vow jug_file_vow_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow
    ADD CONSTRAINT jug_file_vow_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: jug_file_vow jug_file_vow_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow
    ADD CONSTRAINT jug_file_vow_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: jug_ilk_duty jug_ilk_duty_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_duty
    ADD CONSTRAINT jug_ilk_duty_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: jug_ilk_duty jug_ilk_duty_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_duty
    ADD CONSTRAINT jug_ilk_duty_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_ilk_duty jug_ilk_duty_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_duty
    ADD CONSTRAINT jug_ilk_duty_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: jug_ilk_rho jug_ilk_rho_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_rho
    ADD CONSTRAINT jug_ilk_rho_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: jug_ilk_rho jug_ilk_rho_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_rho
    ADD CONSTRAINT jug_ilk_rho_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT jug_init_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: jug_init jug_init_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_init
    ADD CONSTRAINT jug_init_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: jug_vat jug_vat_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vat
    ADD CONSTRAINT jug_vat_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: jug_vat jug_vat_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vat
    ADD CONSTRAINT jug_vat_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_vow jug_vow_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vow
    ADD CONSTRAINT jug_vow_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: jug_vow jug_vow_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vow
    ADD CONSTRAINT jug_vow_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_bump log_bump_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_bump
    ADD CONSTRAINT log_bump_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_bump log_bump_buy_gem_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_bump
    ADD CONSTRAINT log_bump_buy_gem_fkey FOREIGN KEY (buy_gem) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_bump log_bump_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_bump
    ADD CONSTRAINT log_bump_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_bump log_bump_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_bump
    ADD CONSTRAINT log_bump_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_bump log_bump_maker_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_bump
    ADD CONSTRAINT log_bump_maker_fkey FOREIGN KEY (maker) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_bump log_bump_pay_gem_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_bump
    ADD CONSTRAINT log_bump_pay_gem_fkey FOREIGN KEY (pay_gem) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_buy_enabled log_buy_enabled_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_buy_enabled
    ADD CONSTRAINT log_buy_enabled_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_buy_enabled log_buy_enabled_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_buy_enabled
    ADD CONSTRAINT log_buy_enabled_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_buy_enabled log_buy_enabled_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_buy_enabled
    ADD CONSTRAINT log_buy_enabled_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_delete log_delete_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_delete
    ADD CONSTRAINT log_delete_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_delete log_delete_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_delete
    ADD CONSTRAINT log_delete_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_delete log_delete_keeper_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_delete
    ADD CONSTRAINT log_delete_keeper_fkey FOREIGN KEY (keeper) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_delete log_delete_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_delete
    ADD CONSTRAINT log_delete_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_insert log_insert_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_insert
    ADD CONSTRAINT log_insert_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_insert log_insert_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_insert
    ADD CONSTRAINT log_insert_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_insert log_insert_keeper_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_insert
    ADD CONSTRAINT log_insert_keeper_fkey FOREIGN KEY (keeper) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_insert log_insert_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_insert
    ADD CONSTRAINT log_insert_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_item_update log_item_update_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_item_update
    ADD CONSTRAINT log_item_update_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_item_update log_item_update_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_item_update
    ADD CONSTRAINT log_item_update_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_item_update log_item_update_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_item_update
    ADD CONSTRAINT log_item_update_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_kill log_kill_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_kill
    ADD CONSTRAINT log_kill_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_kill log_kill_buy_gem_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_kill
    ADD CONSTRAINT log_kill_buy_gem_fkey FOREIGN KEY (buy_gem) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_kill log_kill_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_kill
    ADD CONSTRAINT log_kill_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_kill log_kill_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_kill
    ADD CONSTRAINT log_kill_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_kill log_kill_maker_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_kill
    ADD CONSTRAINT log_kill_maker_fkey FOREIGN KEY (maker) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_kill log_kill_pay_gem_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_kill
    ADD CONSTRAINT log_kill_pay_gem_fkey FOREIGN KEY (pay_gem) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_make log_make_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_make
    ADD CONSTRAINT log_make_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_make log_make_buy_gem_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_make
    ADD CONSTRAINT log_make_buy_gem_fkey FOREIGN KEY (buy_gem) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_make log_make_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_make
    ADD CONSTRAINT log_make_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_make log_make_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_make
    ADD CONSTRAINT log_make_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_make log_make_maker_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_make
    ADD CONSTRAINT log_make_maker_fkey FOREIGN KEY (maker) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_make log_make_pay_gem_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_make
    ADD CONSTRAINT log_make_pay_gem_fkey FOREIGN KEY (pay_gem) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_matching_enabled log_matching_enabled_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_matching_enabled
    ADD CONSTRAINT log_matching_enabled_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_matching_enabled log_matching_enabled_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_matching_enabled
    ADD CONSTRAINT log_matching_enabled_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_matching_enabled log_matching_enabled_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_matching_enabled
    ADD CONSTRAINT log_matching_enabled_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_median_price log_median_price_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_median_price
    ADD CONSTRAINT log_median_price_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_median_price log_median_price_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_median_price
    ADD CONSTRAINT log_median_price_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_median_price log_median_price_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_median_price
    ADD CONSTRAINT log_median_price_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_min_sell log_min_sell_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_min_sell
    ADD CONSTRAINT log_min_sell_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_min_sell log_min_sell_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_min_sell
    ADD CONSTRAINT log_min_sell_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_min_sell log_min_sell_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_min_sell
    ADD CONSTRAINT log_min_sell_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_min_sell log_min_sell_pay_gem_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_min_sell
    ADD CONSTRAINT log_min_sell_pay_gem_fkey FOREIGN KEY (pay_gem) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_sorted_offer log_sorted_offer_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_sorted_offer
    ADD CONSTRAINT log_sorted_offer_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_sorted_offer log_sorted_offer_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_sorted_offer
    ADD CONSTRAINT log_sorted_offer_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_sorted_offer log_sorted_offer_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_sorted_offer
    ADD CONSTRAINT log_sorted_offer_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_take log_take_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_take
    ADD CONSTRAINT log_take_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_take log_take_buy_gem_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_take
    ADD CONSTRAINT log_take_buy_gem_fkey FOREIGN KEY (buy_gem) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_take log_take_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_take
    ADD CONSTRAINT log_take_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_take log_take_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_take
    ADD CONSTRAINT log_take_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_take log_take_maker_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_take
    ADD CONSTRAINT log_take_maker_fkey FOREIGN KEY (maker) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_take log_take_pay_gem_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_take
    ADD CONSTRAINT log_take_pay_gem_fkey FOREIGN KEY (pay_gem) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_take log_take_taker_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_take
    ADD CONSTRAINT log_take_taker_fkey FOREIGN KEY (taker) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_trade log_trade_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_trade
    ADD CONSTRAINT log_trade_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_trade log_trade_buy_gem_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_trade
    ADD CONSTRAINT log_trade_buy_gem_fkey FOREIGN KEY (buy_gem) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_trade log_trade_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_trade
    ADD CONSTRAINT log_trade_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_trade log_trade_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_trade
    ADD CONSTRAINT log_trade_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_trade log_trade_pay_gem_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_trade
    ADD CONSTRAINT log_trade_pay_gem_fkey FOREIGN KEY (pay_gem) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_unsorted_offer log_unsorted_offer_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_unsorted_offer
    ADD CONSTRAINT log_unsorted_offer_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_unsorted_offer log_unsorted_offer_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_unsorted_offer
    ADD CONSTRAINT log_unsorted_offer_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_unsorted_offer log_unsorted_offer_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_unsorted_offer
    ADD CONSTRAINT log_unsorted_offer_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: log_value log_value_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_value
    ADD CONSTRAINT log_value_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: log_value log_value_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_value
    ADD CONSTRAINT log_value_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: log_value log_value_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.log_value
    ADD CONSTRAINT log_value_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: median_age median_age_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_age
    ADD CONSTRAINT median_age_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_age median_age_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_age
    ADD CONSTRAINT median_age_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: median_age median_age_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_age
    ADD CONSTRAINT median_age_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: median_bar median_bar_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bar
    ADD CONSTRAINT median_bar_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_bar median_bar_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bar
    ADD CONSTRAINT median_bar_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: median_bar median_bar_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bar
    ADD CONSTRAINT median_bar_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: median_bud median_bud_a_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bud
    ADD CONSTRAINT median_bud_a_fkey FOREIGN KEY (a) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_bud median_bud_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bud
    ADD CONSTRAINT median_bud_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_bud median_bud_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bud
    ADD CONSTRAINT median_bud_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: median_bud median_bud_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_bud
    ADD CONSTRAINT median_bud_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: median_diss_batch median_diss_batch_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_batch
    ADD CONSTRAINT median_diss_batch_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_diss_batch median_diss_batch_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_batch
    ADD CONSTRAINT median_diss_batch_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: median_diss_batch median_diss_batch_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_batch
    ADD CONSTRAINT median_diss_batch_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: median_diss_batch median_diss_batch_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_batch
    ADD CONSTRAINT median_diss_batch_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_diss_single median_diss_single_a_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_single
    ADD CONSTRAINT median_diss_single_a_fkey FOREIGN KEY (a) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_diss_single median_diss_single_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_single
    ADD CONSTRAINT median_diss_single_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_diss_single median_diss_single_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_single
    ADD CONSTRAINT median_diss_single_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: median_diss_single median_diss_single_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_single
    ADD CONSTRAINT median_diss_single_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: median_diss_single median_diss_single_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_diss_single
    ADD CONSTRAINT median_diss_single_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_drop median_drop_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_drop
    ADD CONSTRAINT median_drop_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_drop median_drop_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_drop
    ADD CONSTRAINT median_drop_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: median_drop median_drop_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_drop
    ADD CONSTRAINT median_drop_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: median_drop median_drop_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_drop
    ADD CONSTRAINT median_drop_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_kiss_batch median_kiss_batch_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_batch
    ADD CONSTRAINT median_kiss_batch_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_kiss_batch median_kiss_batch_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_batch
    ADD CONSTRAINT median_kiss_batch_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: median_kiss_batch median_kiss_batch_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_batch
    ADD CONSTRAINT median_kiss_batch_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: median_kiss_batch median_kiss_batch_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_batch
    ADD CONSTRAINT median_kiss_batch_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_kiss_single median_kiss_single_a_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_single
    ADD CONSTRAINT median_kiss_single_a_fkey FOREIGN KEY (a) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_kiss_single median_kiss_single_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_single
    ADD CONSTRAINT median_kiss_single_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_kiss_single median_kiss_single_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_single
    ADD CONSTRAINT median_kiss_single_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: median_kiss_single median_kiss_single_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_single
    ADD CONSTRAINT median_kiss_single_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: median_kiss_single median_kiss_single_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_kiss_single
    ADD CONSTRAINT median_kiss_single_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_lift median_lift_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_lift
    ADD CONSTRAINT median_lift_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_lift median_lift_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_lift
    ADD CONSTRAINT median_lift_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: median_lift median_lift_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_lift
    ADD CONSTRAINT median_lift_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: median_lift median_lift_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_lift
    ADD CONSTRAINT median_lift_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_orcl median_orcl_a_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_orcl
    ADD CONSTRAINT median_orcl_a_fkey FOREIGN KEY (a) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_orcl median_orcl_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_orcl
    ADD CONSTRAINT median_orcl_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_orcl median_orcl_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_orcl
    ADD CONSTRAINT median_orcl_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: median_orcl median_orcl_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_orcl
    ADD CONSTRAINT median_orcl_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: median_slot median_slot_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_slot
    ADD CONSTRAINT median_slot_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_slot median_slot_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_slot
    ADD CONSTRAINT median_slot_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: median_slot median_slot_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_slot
    ADD CONSTRAINT median_slot_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: median_slot median_slot_slot_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_slot
    ADD CONSTRAINT median_slot_slot_fkey FOREIGN KEY (slot) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_val median_val_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_val
    ADD CONSTRAINT median_val_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: median_val median_val_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_val
    ADD CONSTRAINT median_val_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: median_val median_val_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.median_val
    ADD CONSTRAINT median_val_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: new_cdp new_cdp_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.new_cdp
    ADD CONSTRAINT new_cdp_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: new_cdp new_cdp_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.new_cdp
    ADD CONSTRAINT new_cdp_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: osm_change osm_change_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.osm_change
    ADD CONSTRAINT osm_change_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: osm_change osm_change_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.osm_change
    ADD CONSTRAINT osm_change_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: osm_change osm_change_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.osm_change
    ADD CONSTRAINT osm_change_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: osm_change osm_change_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.osm_change
    ADD CONSTRAINT osm_change_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: osm_change osm_change_src_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.osm_change
    ADD CONSTRAINT osm_change_src_fkey FOREIGN KEY (src) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: pot_cage pot_cage_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_cage
    ADD CONSTRAINT pot_cage_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_cage pot_cage_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_cage
    ADD CONSTRAINT pot_cage_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: pot_cage pot_cage_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_cage
    ADD CONSTRAINT pot_cage_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: pot_chi pot_chi_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_chi
    ADD CONSTRAINT pot_chi_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: pot_chi pot_chi_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_chi
    ADD CONSTRAINT pot_chi_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_drip pot_drip_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_drip
    ADD CONSTRAINT pot_drip_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_drip pot_drip_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_drip
    ADD CONSTRAINT pot_drip_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: pot_drip pot_drip_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_drip
    ADD CONSTRAINT pot_drip_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: pot_dsr pot_dsr_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_dsr
    ADD CONSTRAINT pot_dsr_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: pot_dsr pot_dsr_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_dsr
    ADD CONSTRAINT pot_dsr_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_exit pot_exit_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_exit
    ADD CONSTRAINT pot_exit_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_exit pot_exit_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_exit
    ADD CONSTRAINT pot_exit_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: pot_exit pot_exit_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_exit
    ADD CONSTRAINT pot_exit_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: pot_file_dsr pot_file_dsr_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_file_dsr
    ADD CONSTRAINT pot_file_dsr_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_file_dsr pot_file_dsr_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_file_dsr
    ADD CONSTRAINT pot_file_dsr_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: pot_file_dsr pot_file_dsr_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_file_dsr
    ADD CONSTRAINT pot_file_dsr_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: pot_file_vow pot_file_vow_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_file_vow
    ADD CONSTRAINT pot_file_vow_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_file_vow pot_file_vow_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_file_vow
    ADD CONSTRAINT pot_file_vow_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: pot_file_vow pot_file_vow_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_file_vow
    ADD CONSTRAINT pot_file_vow_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: pot_join pot_join_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_join
    ADD CONSTRAINT pot_join_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_join pot_join_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_join
    ADD CONSTRAINT pot_join_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: pot_join pot_join_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_join
    ADD CONSTRAINT pot_join_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: pot_live pot_live_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_live
    ADD CONSTRAINT pot_live_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: pot_live pot_live_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_live
    ADD CONSTRAINT pot_live_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_pie pot_pie_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_pie
    ADD CONSTRAINT pot_pie_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: pot_pie pot_pie_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_pie
    ADD CONSTRAINT pot_pie_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_rho pot_rho_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_rho
    ADD CONSTRAINT pot_rho_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: pot_rho pot_rho_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_rho
    ADD CONSTRAINT pot_rho_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_user_pie pot_user_pie_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_user_pie
    ADD CONSTRAINT pot_user_pie_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: pot_user_pie pot_user_pie_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_user_pie
    ADD CONSTRAINT pot_user_pie_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_user_pie pot_user_pie_user_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_user_pie
    ADD CONSTRAINT pot_user_pie_user_fkey FOREIGN KEY ("user") REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: pot_vat pot_vat_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_vat
    ADD CONSTRAINT pot_vat_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: pot_vat pot_vat_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_vat
    ADD CONSTRAINT pot_vat_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_vat pot_vat_vat_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_vat
    ADD CONSTRAINT pot_vat_vat_fkey FOREIGN KEY (vat) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: pot_vow pot_vow_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_vow
    ADD CONSTRAINT pot_vow_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: pot_vow pot_vow_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_vow
    ADD CONSTRAINT pot_vow_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: pot_vow pot_vow_vow_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pot_vow
    ADD CONSTRAINT pot_vow_vow_fkey FOREIGN KEY (vow) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: rely rely_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.rely
    ADD CONSTRAINT rely_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: rely rely_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.rely
    ADD CONSTRAINT rely_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: rely rely_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.rely
    ADD CONSTRAINT rely_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: rely rely_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.rely
    ADD CONSTRAINT rely_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: rely rely_usr_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.rely
    ADD CONSTRAINT rely_usr_fkey FOREIGN KEY (usr) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: set_min_sell set_min_sell_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.set_min_sell
    ADD CONSTRAINT set_min_sell_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: set_min_sell set_min_sell_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.set_min_sell
    ADD CONSTRAINT set_min_sell_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: set_min_sell set_min_sell_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.set_min_sell
    ADD CONSTRAINT set_min_sell_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: set_min_sell set_min_sell_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.set_min_sell
    ADD CONSTRAINT set_min_sell_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: set_min_sell set_min_sell_pay_gem_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.set_min_sell
    ADD CONSTRAINT set_min_sell_pay_gem_fkey FOREIGN KEY (pay_gem) REFERENCES public.addresses(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT spot_file_mat_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: spot_file_mat spot_file_mat_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_mat
    ADD CONSTRAINT spot_file_mat_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: spot_file_par spot_file_par_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_par
    ADD CONSTRAINT spot_file_par_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: spot_file_par spot_file_par_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_par
    ADD CONSTRAINT spot_file_par_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: spot_file_par spot_file_par_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_par
    ADD CONSTRAINT spot_file_par_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT spot_file_pip_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: spot_file_pip spot_file_pip_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_file_pip
    ADD CONSTRAINT spot_file_pip_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: spot_ilk_mat spot_ilk_mat_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_mat
    ADD CONSTRAINT spot_ilk_mat_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: spot_ilk_mat spot_ilk_mat_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_mat
    ADD CONSTRAINT spot_ilk_mat_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: spot_ilk_mat spot_ilk_mat_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_mat
    ADD CONSTRAINT spot_ilk_mat_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: spot_ilk_pip spot_ilk_pip_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_pip
    ADD CONSTRAINT spot_ilk_pip_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: spot_ilk_pip spot_ilk_pip_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_pip
    ADD CONSTRAINT spot_ilk_pip_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: spot_ilk_pip spot_ilk_pip_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_ilk_pip
    ADD CONSTRAINT spot_ilk_pip_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: spot_live spot_live_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_live
    ADD CONSTRAINT spot_live_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: spot_live spot_live_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_live
    ADD CONSTRAINT spot_live_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: spot_par spot_par_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_par
    ADD CONSTRAINT spot_par_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: spot_par spot_par_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_par
    ADD CONSTRAINT spot_par_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT spot_poke_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: spot_vat spot_vat_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_vat
    ADD CONSTRAINT spot_vat_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: spot_vat spot_vat_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.spot_vat
    ADD CONSTRAINT spot_vat_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT tend_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: tend tend_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tend
    ADD CONSTRAINT tend_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT tick_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: tick tick_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tick
    ADD CONSTRAINT tick_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: urns urns_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.urns
    ADD CONSTRAINT urns_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_dai vat_dai_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_dai
    ADD CONSTRAINT vat_dai_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_dai vat_dai_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_dai
    ADD CONSTRAINT vat_dai_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_debt vat_debt_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_debt
    ADD CONSTRAINT vat_debt_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_debt vat_debt_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_debt
    ADD CONSTRAINT vat_debt_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_deny vat_deny_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_deny
    ADD CONSTRAINT vat_deny_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_deny vat_deny_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_deny
    ADD CONSTRAINT vat_deny_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vat_deny vat_deny_usr_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_deny
    ADD CONSTRAINT vat_deny_usr_fkey FOREIGN KEY (usr) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: vat_file_debt_ceiling vat_file_debt_ceiling_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_debt_ceiling
    ADD CONSTRAINT vat_file_debt_ceiling_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_file_debt_ceiling vat_file_debt_ceiling_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_debt_ceiling
    ADD CONSTRAINT vat_file_debt_ceiling_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT vat_file_ilk_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT vat_flux_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vat_fold vat_fold_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold
    ADD CONSTRAINT vat_fold_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_fold vat_fold_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold
    ADD CONSTRAINT vat_fold_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_fold vat_fold_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold
    ADD CONSTRAINT vat_fold_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT vat_fork_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vat_frob vat_frob_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob
    ADD CONSTRAINT vat_frob_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_frob vat_frob_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob
    ADD CONSTRAINT vat_frob_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vat_frob vat_frob_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob
    ADD CONSTRAINT vat_frob_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id) ON DELETE CASCADE;


--
-- Name: vat_gem vat_gem_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_gem
    ADD CONSTRAINT vat_gem_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_gem vat_gem_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_gem
    ADD CONSTRAINT vat_gem_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT vat_grab_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT vat_heal_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vat_hope vat_hope_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_hope
    ADD CONSTRAINT vat_hope_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_hope vat_hope_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_hope
    ADD CONSTRAINT vat_hope_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vat_hope vat_hope_usr_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_hope
    ADD CONSTRAINT vat_hope_usr_fkey FOREIGN KEY (usr) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_art vat_ilk_art_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_art
    ADD CONSTRAINT vat_ilk_art_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_art vat_ilk_art_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_art
    ADD CONSTRAINT vat_ilk_art_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_art vat_ilk_art_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_art
    ADD CONSTRAINT vat_ilk_art_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_dust vat_ilk_dust_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_dust
    ADD CONSTRAINT vat_ilk_dust_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_dust vat_ilk_dust_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_dust
    ADD CONSTRAINT vat_ilk_dust_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_dust vat_ilk_dust_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_dust
    ADD CONSTRAINT vat_ilk_dust_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_line vat_ilk_line_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_line
    ADD CONSTRAINT vat_ilk_line_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_line vat_ilk_line_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_line
    ADD CONSTRAINT vat_ilk_line_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_line vat_ilk_line_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_line
    ADD CONSTRAINT vat_ilk_line_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_rate vat_ilk_rate_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_rate
    ADD CONSTRAINT vat_ilk_rate_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_rate vat_ilk_rate_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_rate
    ADD CONSTRAINT vat_ilk_rate_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_rate vat_ilk_rate_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_rate
    ADD CONSTRAINT vat_ilk_rate_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_spot vat_ilk_spot_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_spot
    ADD CONSTRAINT vat_ilk_spot_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_spot vat_ilk_spot_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_spot
    ADD CONSTRAINT vat_ilk_spot_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT vat_init_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vat_line vat_line_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_line
    ADD CONSTRAINT vat_line_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_line vat_line_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_line
    ADD CONSTRAINT vat_line_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_live vat_live_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_live
    ADD CONSTRAINT vat_live_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_live vat_live_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_live
    ADD CONSTRAINT vat_live_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_move vat_move_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_move
    ADD CONSTRAINT vat_move_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_move vat_move_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_move
    ADD CONSTRAINT vat_move_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vat_nope vat_nope_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_nope
    ADD CONSTRAINT vat_nope_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_nope vat_nope_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_nope
    ADD CONSTRAINT vat_nope_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vat_nope vat_nope_usr_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_nope
    ADD CONSTRAINT vat_nope_usr_fkey FOREIGN KEY (usr) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: vat_rely vat_rely_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_rely
    ADD CONSTRAINT vat_rely_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_rely vat_rely_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_rely
    ADD CONSTRAINT vat_rely_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vat_rely vat_rely_usr_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_rely
    ADD CONSTRAINT vat_rely_usr_fkey FOREIGN KEY (usr) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: vat_sin vat_sin_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_sin
    ADD CONSTRAINT vat_sin_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_sin vat_sin_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_sin
    ADD CONSTRAINT vat_sin_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT vat_slip_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vat_suck vat_suck_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_suck
    ADD CONSTRAINT vat_suck_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_suck vat_suck_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_suck
    ADD CONSTRAINT vat_suck_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vat_urn_art vat_urn_art_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_art
    ADD CONSTRAINT vat_urn_art_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_urn_art vat_urn_art_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_art
    ADD CONSTRAINT vat_urn_art_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_urn_art vat_urn_art_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_art
    ADD CONSTRAINT vat_urn_art_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id) ON DELETE CASCADE;


--
-- Name: vat_urn_ink vat_urn_ink_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_ink
    ADD CONSTRAINT vat_urn_ink_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_urn_ink vat_urn_ink_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_ink
    ADD CONSTRAINT vat_urn_ink_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_urn_ink vat_urn_ink_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_ink
    ADD CONSTRAINT vat_urn_ink_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id) ON DELETE CASCADE;


--
-- Name: vat_vice vat_vice_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_vice
    ADD CONSTRAINT vat_vice_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vat_vice vat_vice_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_vice
    ADD CONSTRAINT vat_vice_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_ash vow_ash_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_ash
    ADD CONSTRAINT vow_ash_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vow_ash vow_ash_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_ash
    ADD CONSTRAINT vow_ash_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_bump vow_bump_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_bump
    ADD CONSTRAINT vow_bump_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vow_bump vow_bump_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_bump
    ADD CONSTRAINT vow_bump_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_dump vow_dump_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_dump
    ADD CONSTRAINT vow_dump_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vow_dump vow_dump_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_dump
    ADD CONSTRAINT vow_dump_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_fess vow_fess_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_fess
    ADD CONSTRAINT vow_fess_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_fess vow_fess_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_fess
    ADD CONSTRAINT vow_fess_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vow_fess vow_fess_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_fess
    ADD CONSTRAINT vow_fess_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: vow_file_auction_address vow_file_auction_address_data_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_address
    ADD CONSTRAINT vow_file_auction_address_data_fkey FOREIGN KEY (data) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: vow_file_auction_address vow_file_auction_address_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_address
    ADD CONSTRAINT vow_file_auction_address_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_file_auction_address vow_file_auction_address_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_address
    ADD CONSTRAINT vow_file_auction_address_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vow_file_auction_address vow_file_auction_address_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_address
    ADD CONSTRAINT vow_file_auction_address_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: vow_file_auction_attributes vow_file_auction_attributes_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_attributes
    ADD CONSTRAINT vow_file_auction_attributes_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_file_auction_attributes vow_file_auction_attributes_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_attributes
    ADD CONSTRAINT vow_file_auction_attributes_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vow_file_auction_attributes vow_file_auction_attributes_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_file_auction_attributes
    ADD CONSTRAINT vow_file_auction_attributes_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: vow_flapper vow_flapper_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flapper
    ADD CONSTRAINT vow_flapper_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vow_flapper vow_flapper_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flapper
    ADD CONSTRAINT vow_flapper_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_flog vow_flog_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flog
    ADD CONSTRAINT vow_flog_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_flog vow_flog_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flog
    ADD CONSTRAINT vow_flog_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vow_flog vow_flog_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flog
    ADD CONSTRAINT vow_flog_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: vow_flopper vow_flopper_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flopper
    ADD CONSTRAINT vow_flopper_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vow_flopper vow_flopper_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flopper
    ADD CONSTRAINT vow_flopper_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_heal vow_heal_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_heal
    ADD CONSTRAINT vow_heal_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_heal vow_heal_log_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_heal
    ADD CONSTRAINT vow_heal_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: vow_heal vow_heal_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_heal
    ADD CONSTRAINT vow_heal_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: vow_hump vow_hump_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_hump
    ADD CONSTRAINT vow_hump_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vow_hump vow_hump_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_hump
    ADD CONSTRAINT vow_hump_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_live vow_live_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_live
    ADD CONSTRAINT vow_live_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vow_live vow_live_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_live
    ADD CONSTRAINT vow_live_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_sin_integer vow_sin_integer_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_integer
    ADD CONSTRAINT vow_sin_integer_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vow_sin_integer vow_sin_integer_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_integer
    ADD CONSTRAINT vow_sin_integer_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_sin_mapping vow_sin_mapping_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_mapping
    ADD CONSTRAINT vow_sin_mapping_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vow_sin_mapping vow_sin_mapping_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_mapping
    ADD CONSTRAINT vow_sin_mapping_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_sump vow_sump_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sump
    ADD CONSTRAINT vow_sump_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vow_sump vow_sump_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sump
    ADD CONSTRAINT vow_sump_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_vat vow_vat_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_vat
    ADD CONSTRAINT vow_vat_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vow_vat vow_vat_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_vat
    ADD CONSTRAINT vow_vat_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_wait vow_wait_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_wait
    ADD CONSTRAINT vow_wait_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: vow_wait vow_wait_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_wait
    ADD CONSTRAINT vow_wait_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: wards wards_address_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.wards
    ADD CONSTRAINT wards_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: wards wards_diff_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.wards
    ADD CONSTRAINT wards_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id) ON DELETE CASCADE;


--
-- Name: wards wards_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.wards
    ADD CONSTRAINT wards_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: wards wards_usr_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.wards
    ADD CONSTRAINT wards_usr_fkey FOREIGN KEY (usr) REFERENCES public.addresses(id) ON DELETE CASCADE;


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
    ADD CONSTRAINT yank_log_id_fkey FOREIGN KEY (log_id) REFERENCES public.event_logs(id) ON DELETE CASCADE;


--
-- Name: yank yank_msg_sender_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.yank
    ADD CONSTRAINT yank_msg_sender_fkey FOREIGN KEY (msg_sender) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: checked_headers checked_headers_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers
    ADD CONSTRAINT checked_headers_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: event_logs event_logs_address_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.event_logs
    ADD CONSTRAINT event_logs_address_fkey FOREIGN KEY (address) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: event_logs event_logs_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.event_logs
    ADD CONSTRAINT event_logs_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: event_logs event_logs_tx_hash_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.event_logs
    ADD CONSTRAINT event_logs_tx_hash_fkey FOREIGN KEY (tx_hash) REFERENCES public.transactions(hash) ON DELETE CASCADE;


--
-- Name: headers headers_eth_node_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_eth_node_id_fkey FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: receipts receipts_contract_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_contract_address_id_fkey FOREIGN KEY (contract_address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: receipts receipts_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: receipts receipts_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES public.transactions(id) ON DELETE CASCADE;


--
-- Name: storage_diff storage_diff_eth_node_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.storage_diff
    ADD CONSTRAINT storage_diff_eth_node_id_fkey FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: transactions transactions_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

