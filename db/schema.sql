--
-- PostgreSQL database dump
--

-- Dumped from database version 10.5
-- Dumped by pg_dump version 10.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: maker; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA maker;


--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: era; Type: TYPE; Schema: maker; Owner: -
--

CREATE TYPE maker.era AS (
	epoch bigint,
	iso timestamp without time zone
);


--
-- Name: frob_event; Type: TYPE; Schema: maker; Owner: -
--

CREATE TYPE maker.frob_event AS (
	ilkid text,
	urnid text,
	dink numeric,
	dart numeric,
	block_number bigint
);


--
-- Name: ilk_state; Type: TYPE; Schema: maker; Owner: -
--

CREATE TYPE maker.ilk_state AS (
	ilk_id integer,
	ilk text,
	block_number bigint,
	rate numeric,
	art numeric,
	spot numeric,
	line numeric,
	dust numeric,
	chop numeric,
	lump numeric,
	flip text,
	rho numeric,
	tax numeric,
	created numeric,
	updated numeric
);


--
-- Name: relevant_block; Type: TYPE; Schema: maker; Owner: -
--

CREATE TYPE maker.relevant_block AS (
	block_number bigint,
	block_hash text,
	ilk_id integer
);


--
-- Name: tx; Type: TYPE; Schema: maker; Owner: -
--

CREATE TYPE maker.tx AS (
	transaction_hash text,
	transaction_index integer,
	block_number bigint,
	block_hash text,
	tx_from text,
	tx_to text
);


--
-- Name: urn_state; Type: TYPE; Schema: maker; Owner: -
--

CREATE TYPE maker.urn_state AS (
	urnid text,
	ilkid text,
	blockheight bigint,
	ink numeric,
	art numeric,
	ratio numeric,
	safe boolean,
	created numeric,
	updated numeric
);


--
-- Name: all_frobs(text); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.all_frobs(ilk text) RETURNS SETOF maker.frob_event
    LANGUAGE sql STABLE
    AS $_$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.ilk = $1)

  SELECT $1 AS ilkId, guy AS urnId, dink, dart, block_number
  FROM maker.vat_frob
  LEFT JOIN maker.urns ON vat_frob.urn_id = urns.id
  LEFT JOIN headers    ON vat_frob.header_id = headers.id
  WHERE urns.ilk_id = (SELECT id FROM ilk)
  ORDER BY guy, block_number DESC
$_$;


--
-- Name: frob_event_ilk(maker.frob_event); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.frob_event_ilk(event maker.frob_event) RETURNS SETOF maker.ilk_state
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM maker.get_ilk_at_block_number(
    event.block_number,
    (SELECT id FROM maker.ilks WHERE ilk = event.ilkid))
$$;


--
-- Name: frob_event_tx(maker.frob_event); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.frob_event_tx(event maker.frob_event) RETURNS maker.tx
    LANGUAGE sql STABLE
    AS $$
  SELECT txs.hash, txs.tx_index, block_number, headers.hash, tx_from, tx_to
  FROM public.light_sync_transactions txs
  LEFT JOIN headers ON txs.header_id = headers.id
  WHERE block_number <= event.block_number
  ORDER BY block_number DESC
  LIMIT 1 -- Should always be true anyway?
$$;


--
-- Name: frob_event_urn(maker.frob_event); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.frob_event_urn(event maker.frob_event) RETURNS SETOF maker.urn_state
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM maker.get_urn_state_at_block(event.ilkid, event.urnid, event.block_number)
$$;


--
-- Name: get_all_urn_states_at_block(bigint); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.get_all_urn_states_at_block(block_height bigint) RETURNS SETOF maker.urn_state
    LANGUAGE sql STABLE
    AS $_$
WITH
  urns AS (
    SELECT urns.id AS urn_id, ilks.id AS ilk_id, ilks.ilk, urns.guy
    FROM maker.urns urns
    LEFT JOIN maker.ilks ilks
    ON urns.ilk_id = ilks.id
  ),

  inks AS ( -- Latest ink for each urn
    SELECT DISTINCT ON (urn_id) urn_id, ink, block_number
    FROM maker.vat_urn_ink
    WHERE block_number <= block_height
    ORDER BY urn_id, block_number DESC
  ),

  arts AS ( -- Latest art for each urn
    SELECT DISTINCT ON (urn_id) urn_id, art, block_number
    FROM maker.vat_urn_art
    WHERE block_number <= block_height
    ORDER BY urn_id, block_number DESC
  ),

  rates AS ( -- Latest rate for each ilk
    SELECT DISTINCT ON (ilk_id) ilk_id, rate, block_number
    FROM maker.vat_ilk_rate
    WHERE block_number <= block_height
    ORDER BY ilk_id, block_number DESC
  ),

  spots AS ( -- Get latest price update for ilk. Problematic from update frequency, slow query?
    SELECT DISTINCT ON (ilk_id) ilk_id, spot, block_number
    FROM maker.vat_ilk_spot
    WHERE block_number <= block_height
    ORDER BY ilk_id, block_number DESC
  ),

  ratio_data AS (
    SELECT urns.ilk, urns.guy, inks.ink, spots.spot, arts.art, rates.rate
    FROM inks
      JOIN urns ON inks.urn_id = urns.urn_id
      JOIN arts ON arts.urn_id = inks.urn_id
      JOIN spots ON spots.ilk_id = urns.ilk_id
      JOIN rates ON rates.ilk_id = spots.ilk_id
  ),

  ratios AS (
    SELECT ilk, guy, ((1.0 * ink * spot) / NULLIF(art * rate, 0)) AS ratio FROM ratio_data
  ),

  safe AS (
    SELECT ilk, guy, (ratio >= 1) AS safe FROM ratios
  ),

  created AS (
    SELECT urn_id, block_timestamp AS created
    FROM
      (
        SELECT DISTINCT ON (urn_id) urn_id, block_hash FROM maker.vat_urn_ink
        ORDER BY urn_id, block_number ASC
      ) earliest_blocks
        LEFT JOIN public.headers ON hash = block_hash
  ),

  updated AS (
    SELECT DISTINCT ON (urn_id) urn_id, headers.block_timestamp AS updated
    FROM
      (
        (SELECT DISTINCT ON (urn_id) urn_id, block_hash FROM maker.vat_urn_ink
         WHERE block_number <= block_height
         ORDER BY urn_id, block_number DESC)
        UNION
        (SELECT DISTINCT ON (urn_id) urn_id, block_hash FROM maker.vat_urn_art
         WHERE block_number <= block_height
         ORDER BY urn_id, block_number DESC)
      ) last_blocks
        LEFT JOIN public.headers ON headers.hash = last_blocks.block_hash
    ORDER BY urn_id, headers.block_timestamp DESC
  )

SELECT urns.guy, urns.ilk, $1, inks.ink, arts.art, ratios.ratio,
       COALESCE(safe.safe, arts.art = 0), created.created, updated.updated
FROM inks
  LEFT JOIN arts     ON arts.urn_id = inks.urn_id
  LEFT JOIN urns     ON arts.urn_id = urns.urn_id
  LEFT JOIN ratios   ON ratios.guy = urns.guy
  LEFT JOIN safe     ON safe.guy = ratios.guy
  LEFT JOIN created  ON created.urn_id = urns.urn_id
  LEFT JOIN updated  ON updated.urn_id = urns.urn_id
  -- Add collections of frob and bite events?
$_$;


--
-- Name: get_ilk_at_block_number(bigint, integer); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.get_ilk_at_block_number(block_number bigint, ilkid integer) RETURNS maker.ilk_state
    LANGUAGE sql STABLE
    AS $_$
WITH rates AS (
    SELECT
      rate,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_rate
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), arts AS (
    SELECT
      art,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_art
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), spots AS (
    SELECT
      spot,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_spot
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), lines AS (
    SELECT
      line,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_line
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), dusts AS (
    SELECT
      dust,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_dust
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), chops AS (
    SELECT
      chop,
      ilk_id,
      block_hash
    FROM maker.cat_ilk_chop
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), lumps AS (
    SELECT
      lump,
      ilk_id,
      block_hash
    FROM maker.cat_ilk_lump
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), flips AS (
    SELECT
      flip,
      ilk_id,
      block_hash
    FROM maker.cat_ilk_flip
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), rhos AS (
    SELECT
      rho,
      ilk_id,
      block_hash
    FROM maker.jug_ilk_rho
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), taxes AS (
    SELECT
      tax,
      ilk_id,
      block_hash
    FROM maker.jug_ilk_tax
    WHERE ilk_id = ilkId
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), relevant_blocks AS (
  SELECT * FROM maker.get_ilk_blocks_before($1, $2)
), created AS (
    SELECT DISTINCT ON (relevant_blocks.ilk_id, relevant_blocks.block_number)
      relevant_blocks.block_number,
      relevant_blocks.block_hash,
      relevant_blocks.ilk_id,
      headers.block_timestamp
    FROM relevant_blocks
      LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
    ORDER BY block_number ASC
    LIMIT 1
), updated AS (
    SELECT DISTINCT ON (relevant_blocks.ilk_id, relevant_blocks.block_number)
      relevant_blocks.block_number,
      relevant_blocks.block_hash,
      relevant_blocks.ilk_id,
      headers.block_timestamp
    FROM relevant_blocks
      LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
    ORDER BY block_number DESC
    LIMIT 1
)

SELECT
  ilks.id,
  ilks.ilk,
  $1 block_number,
  rates.rate,
  arts.art,
  spots.spot,
  lines.line,
  dusts.dust,
  chops.chop,
  lumps.lump,
  flips.flip,
  rhos.rho,
  taxes.tax,
  created.block_timestamp,
  updated.block_timestamp
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
  LEFT JOIN taxes ON taxes.ilk_id = ilks.id
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
  taxes.tax is not null
)
$_$;


--
-- Name: get_ilk_blocks_before(bigint, integer); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.get_ilk_blocks_before(block_number bigint, ilk_id integer) RETURNS SETOF maker.relevant_block
    LANGUAGE sql STABLE
    AS $_$
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.vat_ilk_rate
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.vat_ilk_art
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.vat_ilk_spot
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.vat_ilk_line
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.vat_ilk_dust
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.cat_ilk_chop
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.cat_ilk_lump
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.cat_ilk_flip
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.jug_ilk_rho
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number,
  block_hash,
  ilk_id
FROM maker.jug_ilk_tax
WHERE block_number <= $1
      AND ilk_id = $2
$_$;


--
-- Name: get_ilk_history_before_block(bigint, integer); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.get_ilk_history_before_block(block_number bigint, ilk_id integer) RETURNS SETOF maker.ilk_state
    LANGUAGE plpgsql STABLE
    AS $_$
DECLARE
  r maker.relevant_block;
BEGIN
  FOR r IN SELECT * FROM maker.get_ilk_blocks_before($1, $2)
  LOOP
    RETURN QUERY
    SELECT * FROM maker.get_ilk_at_block_number(r.block_number, $2::integer);
  END LOOP;
END;
$_$;


--
-- Name: get_urn_history_before_block(text, text, bigint); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.get_urn_history_before_block(ilk text, urn text, block_height bigint) RETURNS SETOF maker.urn_state
    LANGUAGE plpgsql STABLE
    AS $_$
DECLARE
  blocks BIGINT[];
  i BIGINT;
  ilkId NUMERIC;
  urnId NUMERIC;
BEGIN
  SELECT id FROM maker.ilks WHERE ilks.ilk = $1 INTO ilkId;
  SELECT id FROM maker.urns WHERE urns.guy = $2 AND urns.ilk_id = ilkID INTO urnId;

  blocks := ARRAY(
    SELECT block_number
    FROM (
       SELECT block_number
       FROM maker.vat_urn_ink
       WHERE vat_urn_ink.urn_id = urnId
         AND block_number <= $3
       UNION
       SELECT block_number
       FROM maker.vat_urn_art
       WHERE vat_urn_art.urn_id = urnId
         AND block_number <= $3
     ) inks_and_arts
    ORDER BY block_number DESC
  );

  FOREACH i IN ARRAY blocks
    LOOP
      RETURN QUERY
        SELECT * FROM maker.get_urn_state_at_block(ilk, urn, i);
    END LOOP;
END;
$_$;


--
-- Name: get_urn_state_at_block(text, text, bigint); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.get_urn_state_at_block(ilk text, urn text, block_height bigint) RETURNS maker.urn_state
    LANGUAGE sql STABLE
    AS $_$
WITH
  urn AS (
    SELECT urns.id AS urn_id, ilks.id AS ilk_id, ilks.ilk, urns.guy
    FROM maker.urns urns
    LEFT JOIN maker.ilks ilks
    ON urns.ilk_id = ilks.id
    WHERE ilks.ilk = $1 AND urns.guy = $2
  ),

  ink AS ( -- Latest ink
    SELECT DISTINCT ON (urn_id) urn_id, ink, block_number
    FROM maker.vat_urn_ink
    WHERE urn_id = (SELECT urn_id from urn where guy = $2) AND block_number <= block_height
    ORDER BY urn_id, block_number DESC
  ),

  art AS ( -- Latest art
    SELECT DISTINCT ON (urn_id) urn_id, art, block_number
    FROM maker.vat_urn_art
    WHERE urn_id = (SELECT urn_id from urn where guy = $2) AND block_number <= block_height
    ORDER BY urn_id, block_number DESC
  ),

  rate AS ( -- Latest rate for ilk
    SELECT DISTINCT ON (ilk_id) ilk_id, rate, block_number
    FROM maker.vat_ilk_rate
    WHERE ilk_id = (SELECT ilk_id from urn where ilk = $1) AND block_number <= block_height
    ORDER BY ilk_id, block_number DESC
  ),

  spot AS ( -- Get latest price update for ilk. Problematic from update frequency, slow query?
    SELECT DISTINCT ON (ilk_id) ilk_id, spot, block_number
    FROM maker.vat_ilk_spot
    WHERE ilk_id = (SELECT ilk_id from urn where ilk = $1) AND block_number <= block_height
    ORDER BY ilk_id, block_number DESC
  ),

  ratio_data AS (
    SELECT urn.ilk, urn.guy, ink, spot, art, rate
    FROM ink
      JOIN urn ON ink.urn_id = urn.urn_id
      JOIN art ON art.urn_id = ink.urn_id
      JOIN spot ON spot.ilk_id = urn.ilk_id
      JOIN rate ON rate.ilk_id = spot.ilk_id
  ),

  ratio AS (
    SELECT ilk, guy, ((1.0 * ink * spot) / NULLIF(art * rate, 0)) AS ratio
    FROM ratio_data
  ),

  safe AS (
    SELECT ilk, guy, (ratio >= 1) AS safe FROM ratio
  ),

  created AS (
    SELECT urn_id, block_timestamp AS created
    FROM
      (
        SELECT DISTINCT ON (urn_id) urn_id, block_hash FROM maker.vat_urn_ink
        WHERE urn_id = (SELECT urn_id from urn where guy = $2)
        ORDER BY urn_id, block_number ASC
      ) earliest_blocks
        LEFT JOIN public.headers ON hash = block_hash
  ),

  updated AS (
    SELECT DISTINCT ON (urn_id) urn_id, headers.block_timestamp AS updated
    FROM
      (
        SELECT urn_id, block_number FROM ink
        UNION
        SELECT urn_id, block_number FROM art
      ) last_blocks
        LEFT JOIN public.headers ON headers.block_number = last_blocks.block_number
    ORDER BY urn_id, headers.block_timestamp DESC
  )

SELECT $2 AS urnId, $1 AS ilkId, $3 AS block_height, ink.ink, art.art, ratio.ratio,
       COALESCE(safe.safe, art.art = 0), created.created, updated.updated
FROM ink
  LEFT JOIN art     ON art.urn_id = ink.urn_id
  LEFT JOIN urn     ON urn.urn_id = ink.urn_id
  LEFT JOIN ratio   ON ratio.ilk = urn.ilk   AND ratio.guy = urn.guy
  LEFT JOIN safe    ON safe.ilk = ratio.ilk  AND safe.guy = ratio.guy
  LEFT JOIN created ON created.urn_id = art.urn_id
  LEFT JOIN updated ON updated.urn_id = art.urn_id
  -- Add collections of frob and bite events?
WHERE ink.urn_id IS NOT NULL
$_$;


--
-- Name: ilk_state_frobs(maker.ilk_state); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.ilk_state_frobs(state maker.ilk_state) RETURNS SETOF maker.frob_event
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM maker.all_frobs(state.ilk)
  WHERE block_number <= state.block_number
$$;


--
-- Name: tx_era(maker.tx); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.tx_era(tx maker.tx) RETURNS maker.era
    LANGUAGE sql STABLE
    AS $$
SELECT block_timestamp::BIGINT AS "epoch", (SELECT TIMESTAMP 'epoch' + block_timestamp * INTERVAL '1 second') AS iso
  FROM headers WHERE block_number = tx.block_number
$$;


--
-- Name: urn_frobs(text, text); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.urn_frobs(ilk text, urn text) RETURNS SETOF maker.frob_event
    LANGUAGE sql STABLE
    AS $_$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.ilk = $1),
    urn AS (
      SELECT id FROM maker.urns
      WHERE ilk_id = (SELECT id FROM ilk)
        AND guy = $2
    )

  SELECT $1 AS ilkId, $2 AS urnId, dink, dart, block_number
  FROM maker.vat_frob LEFT JOIN headers ON vat_frob.header_id = headers.id
  WHERE vat_frob.urn_id = (SELECT id FROM urn)
  ORDER BY block_number DESC
$_$;


--
-- Name: urn_state_frobs(maker.urn_state); Type: FUNCTION; Schema: maker; Owner: -
--

CREATE FUNCTION maker.urn_state_frobs(state maker.urn_state) RETURNS SETOF maker.frob_event
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM maker.urn_frobs(state.ilkid, state.urnid)
  WHERE block_number <= state.blockheight
$$;


--
-- Name: notify_pricefeed(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.notify_pricefeed() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  PERFORM pg_notify(
    CAST('postgraphile:price_feed' AS text),
    json_build_object('__node__', json_build_array('price_feeds', NEW.id))::text
  );
  RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: bite; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.bite (
    id integer NOT NULL,
    header_id integer NOT NULL,
    urn_id integer NOT NULL,
    ink numeric,
    art numeric,
    tab numeric,
    flip numeric,
    tx_idx integer NOT NULL,
    log_idx integer NOT NULL,
    raw_log jsonb
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
    ilk_id integer NOT NULL,
    what text,
    data numeric,
    tx_idx integer NOT NULL,
    log_idx integer NOT NULL,
    raw_log jsonb
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
    ilk_id integer NOT NULL,
    what text,
    flip text,
    tx_idx integer NOT NULL,
    log_idx integer NOT NULL,
    raw_log jsonb
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
    what text,
    data text,
    tx_idx integer NOT NULL,
    log_idx integer NOT NULL,
    raw_log jsonb
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
-- Name: cat_flip_ilk; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_flip_ilk (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    flip numeric NOT NULL,
    ilk_id integer NOT NULL
);


--
-- Name: cat_flip_ilk_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_flip_ilk_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_flip_ilk_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_flip_ilk_id_seq OWNED BY maker.cat_flip_ilk.id;


--
-- Name: cat_flip_ink; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_flip_ink (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    flip numeric NOT NULL,
    ink numeric NOT NULL
);


--
-- Name: cat_flip_ink_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_flip_ink_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_flip_ink_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_flip_ink_id_seq OWNED BY maker.cat_flip_ink.id;


--
-- Name: cat_flip_tab; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_flip_tab (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    flip numeric NOT NULL,
    tab numeric NOT NULL
);


--
-- Name: cat_flip_tab_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_flip_tab_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_flip_tab_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_flip_tab_id_seq OWNED BY maker.cat_flip_tab.id;


--
-- Name: cat_flip_urn; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_flip_urn (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    flip numeric NOT NULL,
    urn text
);


--
-- Name: cat_flip_urn_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_flip_urn_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_flip_urn_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_flip_urn_id_seq OWNED BY maker.cat_flip_urn.id;


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
-- Name: cat_nflip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_nflip (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    nflip numeric NOT NULL
);


--
-- Name: cat_nflip_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_nflip_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_nflip_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_nflip_id_seq OWNED BY maker.cat_nflip.id;


--
-- Name: cat_pit; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.cat_pit (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    pit text
);


--
-- Name: cat_pit_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.cat_pit_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cat_pit_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.cat_pit_id_seq OWNED BY maker.cat_pit.id;


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
-- Name: deal; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.deal (
    id integer NOT NULL,
    header_id integer NOT NULL,
    bid_id numeric NOT NULL,
    contract_address character varying,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
    bid_id numeric NOT NULL,
    lot numeric,
    bid numeric,
    guy bytea,
    tic numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
-- Name: flap_kick; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flap_kick (
    id integer NOT NULL,
    header_id integer NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric NOT NULL,
    bid numeric NOT NULL,
    gal text,
    "end" timestamp with time zone,
    tx_idx integer NOT NULL,
    log_idx integer NOT NULL,
    raw_log jsonb
);


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
-- Name: flip_kick; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flip_kick (
    id integer NOT NULL,
    header_id integer NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric,
    bid numeric,
    gal text,
    "end" timestamp with time zone,
    urn text,
    tab numeric,
    tx_idx integer NOT NULL,
    log_idx integer NOT NULL,
    raw_log jsonb
);


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
-- Name: flop_kick; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.flop_kick (
    id integer NOT NULL,
    header_id integer NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric NOT NULL,
    bid numeric NOT NULL,
    gal text,
    "end" timestamp with time zone,
    tx_idx integer NOT NULL,
    log_idx integer NOT NULL,
    raw_log jsonb
);


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
-- Name: ilks; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.ilks (
    id integer NOT NULL,
    ilk text
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
-- Name: jug_drip; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_drip (
    id integer NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
-- Name: jug_file_ilk; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_file_ilk (
    id integer NOT NULL,
    header_id integer NOT NULL,
    ilk_id integer NOT NULL,
    what text,
    data numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
-- Name: jug_file_repo; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_file_repo (
    id integer NOT NULL,
    header_id integer NOT NULL,
    what text,
    data numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
);


--
-- Name: jug_file_repo_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_file_repo_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_file_repo_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_file_repo_id_seq OWNED BY maker.jug_file_repo.id;


--
-- Name: jug_file_vow; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_file_vow (
    id integer NOT NULL,
    header_id integer NOT NULL,
    what text,
    data text,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
-- Name: jug_ilk_tax; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_ilk_tax (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    ilk_id integer NOT NULL,
    tax numeric NOT NULL
);


--
-- Name: jug_ilk_tax_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_ilk_tax_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_ilk_tax_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_ilk_tax_id_seq OWNED BY maker.jug_ilk_tax.id;


--
-- Name: jug_repo; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_repo (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    repo text
);


--
-- Name: jug_repo_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.jug_repo_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: jug_repo_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.jug_repo_id_seq OWNED BY maker.jug_repo.id;


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
-- Name: price_feeds; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.price_feeds (
    id integer NOT NULL,
    block_number bigint NOT NULL,
    header_id integer NOT NULL,
    medianizer_address text,
    age numeric,
    usd_value numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
);


--
-- Name: price_feeds_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.price_feeds_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: price_feeds_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.price_feeds_id_seq OWNED BY maker.price_feeds.id;


--
-- Name: tend; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.tend (
    id integer NOT NULL,
    header_id integer NOT NULL,
    bid_id numeric NOT NULL,
    lot numeric,
    bid numeric,
    guy text,
    tic numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
-- Name: urns; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.urns (
    id integer NOT NULL,
    ilk_id integer NOT NULL,
    guy text
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
    what text,
    data numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
    ilk_id integer NOT NULL,
    what text,
    data numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
    ilk_id integer NOT NULL,
    src text,
    dst text,
    rad numeric,
    tx_idx integer NOT NULL,
    log_idx integer NOT NULL,
    raw_log jsonb
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
    urn_id integer NOT NULL,
    rate numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
-- Name: vat_frob; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vat_frob (
    id integer NOT NULL,
    header_id integer NOT NULL,
    urn_id integer NOT NULL,
    v text,
    w text,
    dink numeric,
    dart numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
    urn_id integer NOT NULL,
    v text,
    w text,
    dink numeric,
    dart numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
    urn text,
    v text,
    rad numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
    ilk_id integer NOT NULL,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
    src text NOT NULL,
    dst text NOT NULL,
    rad numeric NOT NULL,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
    ilk_id integer NOT NULL,
    guy text,
    rad numeric,
    tx_idx integer NOT NULL,
    log_idx integer NOT NULL,
    raw_log jsonb
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
-- Name: vow_cow; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_cow (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    cow text
);


--
-- Name: vow_cow_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_cow_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_cow_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_cow_id_seq OWNED BY maker.vow_cow.id;


--
-- Name: vow_fess; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_fess (
    id integer NOT NULL,
    header_id integer NOT NULL,
    tab numeric NOT NULL,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
-- Name: vow_flog; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_flog (
    id integer NOT NULL,
    header_id integer NOT NULL,
    era integer NOT NULL,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
-- Name: vow_row; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_row (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    "row" text
);


--
-- Name: vow_row_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_row_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_row_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_row_id_seq OWNED BY maker.vow_row.id;


--
-- Name: vow_sin; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.vow_sin (
    id integer NOT NULL,
    block_number bigint,
    block_hash text,
    sin numeric
);


--
-- Name: vow_sin_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.vow_sin_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vow_sin_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.vow_sin_id_seq OWNED BY maker.vow_sin.id;


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
-- Name: logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.logs (
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
 SELECT max(logs.block_number) AS max_block,
    min(logs.block_number) AS min_block
   FROM public.logs;


--
-- Name: blocks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.blocks (
    id integer NOT NULL,
    difficulty bigint,
    extra_data character varying,
    gaslimit bigint,
    gasused bigint,
    hash character varying(66),
    miner character varying(42),
    nonce character varying(20),
    number bigint,
    parenthash character varying(66),
    reward double precision,
    uncles_reward double precision,
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
    flip_kick_checked integer DEFAULT 0 NOT NULL,
    vat_frob_checked integer DEFAULT 0 NOT NULL,
    price_feeds_checked integer DEFAULT 0 NOT NULL,
    tend_checked integer DEFAULT 0 NOT NULL,
    bite_checked integer DEFAULT 0 NOT NULL,
    dent_checked integer DEFAULT 0 NOT NULL,
    vat_file_debt_ceiling_checked integer DEFAULT 0 NOT NULL,
    vat_file_ilk_checked integer DEFAULT 0 NOT NULL,
    vat_init_checked integer DEFAULT 0 NOT NULL,
    jug_file_ilk_checked integer DEFAULT 0 NOT NULL,
    jug_file_repo_checked integer DEFAULT 0 NOT NULL,
    jug_file_vow_checked integer DEFAULT 0 NOT NULL,
    deal_checked integer DEFAULT 0 NOT NULL,
    jug_drip_checked integer DEFAULT 0 NOT NULL,
    cat_file_chop_lump_checked integer DEFAULT 0 NOT NULL,
    cat_file_flip_checked integer DEFAULT 0 NOT NULL,
    cat_file_vow_checked integer DEFAULT 0 NOT NULL,
    flop_kick_checked integer DEFAULT 0 NOT NULL,
    vat_move_checked integer DEFAULT 0 NOT NULL,
    vat_fold_checked integer DEFAULT 0 NOT NULL,
    vat_heal_checked integer DEFAULT 0 NOT NULL,
    vat_grab_checked integer DEFAULT 0 NOT NULL,
    vat_flux_checked integer DEFAULT 0 NOT NULL,
    vat_slip_checked integer DEFAULT 0 NOT NULL,
    vow_flog_checked integer DEFAULT 0 NOT NULL,
    flap_kick_checked integer DEFAULT 0 NOT NULL,
    vow_fess_checked integer DEFAULT 0 NOT NULL
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
-- Name: full_sync_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.full_sync_transactions (
    id integer NOT NULL,
    block_id integer NOT NULL,
    gaslimit numeric,
    gasprice numeric,
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
-- Name: headers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.headers (
    id integer NOT NULL,
    hash character varying(66),
    block_number bigint,
    raw jsonb,
    block_timestamp numeric,
    eth_node_id integer,
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
-- Name: light_sync_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.light_sync_transactions (
    id integer NOT NULL,
    header_id integer NOT NULL,
    hash text,
    gaslimit numeric,
    gasprice numeric,
    input_data bytea,
    nonce numeric,
    raw bytea,
    tx_from text,
    tx_index integer,
    tx_to text,
    value numeric
);


--
-- Name: light_sync_transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.light_sync_transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: light_sync_transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.light_sync_transactions_id_seq OWNED BY public.light_sync_transactions.id;


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
-- Name: logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.logs_id_seq OWNED BY public.logs.id;


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
-- Name: receipts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.receipts (
    id integer NOT NULL,
    contract_address character varying(42),
    cumulative_gas_used numeric,
    gas_used numeric,
    state_root character varying(66),
    status integer,
    tx_hash character varying(66),
    block_id integer NOT NULL
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
    logs.id,
    logs.block_number,
    logs.address,
    logs.tx_hash,
    logs.index,
    logs.topic0,
    logs.topic1,
    logs.topic2,
    logs.topic3,
    logs.data,
    logs.receipt_id
   FROM ((public.log_filters
     CROSS JOIN public.block_stats)
     JOIN public.logs ON ((((logs.address)::text = (log_filters.address)::text) AND (logs.block_number >= COALESCE(log_filters.from_block, block_stats.min_block)) AND (logs.block_number <= COALESCE(log_filters.to_block, block_stats.max_block)))))
  WHERE ((((log_filters.topic0)::text = (logs.topic0)::text) OR (log_filters.topic0 IS NULL)) AND (((log_filters.topic1)::text = (logs.topic1)::text) OR (log_filters.topic1 IS NULL)) AND (((log_filters.topic2)::text = (logs.topic2)::text) OR (log_filters.topic2 IS NULL)) AND (((log_filters.topic3)::text = (logs.topic3)::text) OR (log_filters.topic3 IS NULL)));


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
-- Name: cat_flip_ilk id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_ilk ALTER COLUMN id SET DEFAULT nextval('maker.cat_flip_ilk_id_seq'::regclass);


--
-- Name: cat_flip_ink id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_ink ALTER COLUMN id SET DEFAULT nextval('maker.cat_flip_ink_id_seq'::regclass);


--
-- Name: cat_flip_tab id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_tab ALTER COLUMN id SET DEFAULT nextval('maker.cat_flip_tab_id_seq'::regclass);


--
-- Name: cat_flip_urn id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_urn ALTER COLUMN id SET DEFAULT nextval('maker.cat_flip_urn_id_seq'::regclass);


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
-- Name: cat_nflip id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_nflip ALTER COLUMN id SET DEFAULT nextval('maker.cat_nflip_id_seq'::regclass);


--
-- Name: cat_pit id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_pit ALTER COLUMN id SET DEFAULT nextval('maker.cat_pit_id_seq'::regclass);


--
-- Name: cat_vat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vat ALTER COLUMN id SET DEFAULT nextval('maker.cat_vat_id_seq'::regclass);


--
-- Name: cat_vow id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vow ALTER COLUMN id SET DEFAULT nextval('maker.cat_vow_id_seq'::regclass);


--
-- Name: deal id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deal ALTER COLUMN id SET DEFAULT nextval('maker.deal_id_seq'::regclass);


--
-- Name: dent id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.dent ALTER COLUMN id SET DEFAULT nextval('maker.dent_id_seq'::regclass);


--
-- Name: flap_kick id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kick ALTER COLUMN id SET DEFAULT nextval('maker.flap_kick_id_seq'::regclass);


--
-- Name: flip_kick id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kick ALTER COLUMN id SET DEFAULT nextval('maker.flip_kick_id_seq'::regclass);


--
-- Name: flop_kick id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kick ALTER COLUMN id SET DEFAULT nextval('maker.flop_kick_id_seq'::regclass);


--
-- Name: ilks id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.ilks ALTER COLUMN id SET DEFAULT nextval('maker.ilks_id_seq'::regclass);


--
-- Name: jug_drip id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_drip ALTER COLUMN id SET DEFAULT nextval('maker.jug_drip_id_seq'::regclass);


--
-- Name: jug_file_ilk id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_ilk ALTER COLUMN id SET DEFAULT nextval('maker.jug_file_ilk_id_seq'::regclass);


--
-- Name: jug_file_repo id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_repo ALTER COLUMN id SET DEFAULT nextval('maker.jug_file_repo_id_seq'::regclass);


--
-- Name: jug_file_vow id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow ALTER COLUMN id SET DEFAULT nextval('maker.jug_file_vow_id_seq'::regclass);


--
-- Name: jug_ilk_rho id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_rho ALTER COLUMN id SET DEFAULT nextval('maker.jug_ilk_rho_id_seq'::regclass);


--
-- Name: jug_ilk_tax id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_tax ALTER COLUMN id SET DEFAULT nextval('maker.jug_ilk_tax_id_seq'::regclass);


--
-- Name: jug_repo id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_repo ALTER COLUMN id SET DEFAULT nextval('maker.jug_repo_id_seq'::regclass);


--
-- Name: jug_vat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vat ALTER COLUMN id SET DEFAULT nextval('maker.jug_vat_id_seq'::regclass);


--
-- Name: jug_vow id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vow ALTER COLUMN id SET DEFAULT nextval('maker.jug_vow_id_seq'::regclass);


--
-- Name: price_feeds id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.price_feeds ALTER COLUMN id SET DEFAULT nextval('maker.price_feeds_id_seq'::regclass);


--
-- Name: tend id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tend ALTER COLUMN id SET DEFAULT nextval('maker.tend_id_seq'::regclass);


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
-- Name: vow_cow id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_cow ALTER COLUMN id SET DEFAULT nextval('maker.vow_cow_id_seq'::regclass);


--
-- Name: vow_fess id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_fess ALTER COLUMN id SET DEFAULT nextval('maker.vow_fess_id_seq'::regclass);


--
-- Name: vow_flog id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flog ALTER COLUMN id SET DEFAULT nextval('maker.vow_flog_id_seq'::regclass);


--
-- Name: vow_hump id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_hump ALTER COLUMN id SET DEFAULT nextval('maker.vow_hump_id_seq'::regclass);


--
-- Name: vow_row id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_row ALTER COLUMN id SET DEFAULT nextval('maker.vow_row_id_seq'::regclass);


--
-- Name: vow_sin id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin ALTER COLUMN id SET DEFAULT nextval('maker.vow_sin_id_seq'::regclass);


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
-- Name: full_sync_transactions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_transactions ALTER COLUMN id SET DEFAULT nextval('public.full_sync_transactions_id_seq'::regclass);


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: headers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers ALTER COLUMN id SET DEFAULT nextval('public.headers_id_seq'::regclass);


--
-- Name: light_sync_transactions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_transactions ALTER COLUMN id SET DEFAULT nextval('public.light_sync_transactions_id_seq'::regclass);


--
-- Name: log_filters id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.log_filters ALTER COLUMN id SET DEFAULT nextval('public.log_filters_id_seq'::regclass);


--
-- Name: logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs ALTER COLUMN id SET DEFAULT nextval('public.logs_id_seq'::regclass);


--
-- Name: queued_storage id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage ALTER COLUMN id SET DEFAULT nextval('public.queued_storage_id_seq'::regclass);


--
-- Name: receipts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipts ALTER COLUMN id SET DEFAULT nextval('public.receipts_id_seq'::regclass);


--
-- Name: watched_contracts contract_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_contracts ALTER COLUMN contract_id SET DEFAULT nextval('public.watched_contracts_contract_id_seq'::regclass);


--
-- Name: bite bite_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite
    ADD CONSTRAINT bite_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: bite bite_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite
    ADD CONSTRAINT bite_pkey PRIMARY KEY (id);


--
-- Name: cat_file_chop_lump cat_file_chop_lump_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_chop_lump
    ADD CONSTRAINT cat_file_chop_lump_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: cat_file_chop_lump cat_file_chop_lump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_chop_lump
    ADD CONSTRAINT cat_file_chop_lump_pkey PRIMARY KEY (id);


--
-- Name: cat_file_flip cat_file_flip_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_flip
    ADD CONSTRAINT cat_file_flip_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: cat_file_flip cat_file_flip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_flip
    ADD CONSTRAINT cat_file_flip_pkey PRIMARY KEY (id);


--
-- Name: cat_file_vow cat_file_vow_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_vow
    ADD CONSTRAINT cat_file_vow_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: cat_file_vow cat_file_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_vow
    ADD CONSTRAINT cat_file_vow_pkey PRIMARY KEY (id);


--
-- Name: cat_flip_ilk cat_flip_ilk_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_ilk
    ADD CONSTRAINT cat_flip_ilk_pkey PRIMARY KEY (id);


--
-- Name: cat_flip_ink cat_flip_ink_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_ink
    ADD CONSTRAINT cat_flip_ink_pkey PRIMARY KEY (id);


--
-- Name: cat_flip_tab cat_flip_tab_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_tab
    ADD CONSTRAINT cat_flip_tab_pkey PRIMARY KEY (id);


--
-- Name: cat_flip_urn cat_flip_urn_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_urn
    ADD CONSTRAINT cat_flip_urn_pkey PRIMARY KEY (id);


--
-- Name: cat_ilk_chop cat_ilk_chop_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_chop
    ADD CONSTRAINT cat_ilk_chop_pkey PRIMARY KEY (id);


--
-- Name: cat_ilk_flip cat_ilk_flip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_flip
    ADD CONSTRAINT cat_ilk_flip_pkey PRIMARY KEY (id);


--
-- Name: cat_ilk_lump cat_ilk_lump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_lump
    ADD CONSTRAINT cat_ilk_lump_pkey PRIMARY KEY (id);


--
-- Name: cat_live cat_live_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_live
    ADD CONSTRAINT cat_live_pkey PRIMARY KEY (id);


--
-- Name: cat_nflip cat_nflip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_nflip
    ADD CONSTRAINT cat_nflip_pkey PRIMARY KEY (id);


--
-- Name: cat_pit cat_pit_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_pit
    ADD CONSTRAINT cat_pit_pkey PRIMARY KEY (id);


--
-- Name: cat_vat cat_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vat
    ADD CONSTRAINT cat_vat_pkey PRIMARY KEY (id);


--
-- Name: cat_vow cat_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_vow
    ADD CONSTRAINT cat_vow_pkey PRIMARY KEY (id);


--
-- Name: deal deal_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deal
    ADD CONSTRAINT deal_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: deal deal_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deal
    ADD CONSTRAINT deal_pkey PRIMARY KEY (id);


--
-- Name: dent dent_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.dent
    ADD CONSTRAINT dent_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: dent dent_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.dent
    ADD CONSTRAINT dent_pkey PRIMARY KEY (id);


--
-- Name: flap_kick flap_kick_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kick
    ADD CONSTRAINT flap_kick_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: flap_kick flap_kick_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kick
    ADD CONSTRAINT flap_kick_pkey PRIMARY KEY (id);


--
-- Name: flip_kick flip_kick_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kick
    ADD CONSTRAINT flip_kick_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: flip_kick flip_kick_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kick
    ADD CONSTRAINT flip_kick_pkey PRIMARY KEY (id);


--
-- Name: flop_kick flop_kick_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kick
    ADD CONSTRAINT flop_kick_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: flop_kick flop_kick_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kick
    ADD CONSTRAINT flop_kick_pkey PRIMARY KEY (id);


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
-- Name: jug_drip jug_drip_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_drip
    ADD CONSTRAINT jug_drip_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: jug_drip jug_drip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_drip
    ADD CONSTRAINT jug_drip_pkey PRIMARY KEY (id);


--
-- Name: jug_file_ilk jug_file_ilk_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_ilk
    ADD CONSTRAINT jug_file_ilk_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: jug_file_ilk jug_file_ilk_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_ilk
    ADD CONSTRAINT jug_file_ilk_pkey PRIMARY KEY (id);


--
-- Name: jug_file_repo jug_file_repo_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_repo
    ADD CONSTRAINT jug_file_repo_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: jug_file_repo jug_file_repo_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_repo
    ADD CONSTRAINT jug_file_repo_pkey PRIMARY KEY (id);


--
-- Name: jug_file_vow jug_file_vow_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow
    ADD CONSTRAINT jug_file_vow_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: jug_file_vow jug_file_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow
    ADD CONSTRAINT jug_file_vow_pkey PRIMARY KEY (id);


--
-- Name: jug_ilk_rho jug_ilk_rho_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_rho
    ADD CONSTRAINT jug_ilk_rho_pkey PRIMARY KEY (id);


--
-- Name: jug_ilk_tax jug_ilk_tax_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_tax
    ADD CONSTRAINT jug_ilk_tax_pkey PRIMARY KEY (id);


--
-- Name: jug_repo jug_repo_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_repo
    ADD CONSTRAINT jug_repo_pkey PRIMARY KEY (id);


--
-- Name: jug_vat jug_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vat
    ADD CONSTRAINT jug_vat_pkey PRIMARY KEY (id);


--
-- Name: jug_vow jug_vow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vow
    ADD CONSTRAINT jug_vow_pkey PRIMARY KEY (id);


--
-- Name: price_feeds price_feeds_header_id_medianizer_address_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.price_feeds
    ADD CONSTRAINT price_feeds_header_id_medianizer_address_tx_idx_log_idx_key UNIQUE (header_id, medianizer_address, tx_idx, log_idx);


--
-- Name: price_feeds price_feeds_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.price_feeds
    ADD CONSTRAINT price_feeds_pkey PRIMARY KEY (id);


--
-- Name: tend tend_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tend
    ADD CONSTRAINT tend_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: tend tend_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tend
    ADD CONSTRAINT tend_pkey PRIMARY KEY (id);


--
-- Name: urns urns_ilk_id_guy_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.urns
    ADD CONSTRAINT urns_ilk_id_guy_key UNIQUE (ilk_id, guy);


--
-- Name: urns urns_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.urns
    ADD CONSTRAINT urns_pkey PRIMARY KEY (id);


--
-- Name: vat_dai vat_dai_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_dai
    ADD CONSTRAINT vat_dai_pkey PRIMARY KEY (id);


--
-- Name: vat_debt vat_debt_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_debt
    ADD CONSTRAINT vat_debt_pkey PRIMARY KEY (id);


--
-- Name: vat_file_debt_ceiling vat_file_debt_ceiling_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_debt_ceiling
    ADD CONSTRAINT vat_file_debt_ceiling_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: vat_file_debt_ceiling vat_file_debt_ceiling_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_debt_ceiling
    ADD CONSTRAINT vat_file_debt_ceiling_pkey PRIMARY KEY (id);


--
-- Name: vat_file_ilk vat_file_ilk_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_ilk
    ADD CONSTRAINT vat_file_ilk_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: vat_file_ilk vat_file_ilk_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_ilk
    ADD CONSTRAINT vat_file_ilk_pkey PRIMARY KEY (id);


--
-- Name: vat_flux vat_flux_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_flux
    ADD CONSTRAINT vat_flux_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: vat_flux vat_flux_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_flux
    ADD CONSTRAINT vat_flux_pkey PRIMARY KEY (id);


--
-- Name: vat_fold vat_fold_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold
    ADD CONSTRAINT vat_fold_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: vat_fold vat_fold_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold
    ADD CONSTRAINT vat_fold_pkey PRIMARY KEY (id);


--
-- Name: vat_frob vat_frob_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob
    ADD CONSTRAINT vat_frob_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: vat_frob vat_frob_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob
    ADD CONSTRAINT vat_frob_pkey PRIMARY KEY (id);


--
-- Name: vat_gem vat_gem_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_gem
    ADD CONSTRAINT vat_gem_pkey PRIMARY KEY (id);


--
-- Name: vat_grab vat_grab_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_grab
    ADD CONSTRAINT vat_grab_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: vat_grab vat_grab_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_grab
    ADD CONSTRAINT vat_grab_pkey PRIMARY KEY (id);


--
-- Name: vat_heal vat_heal_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_heal
    ADD CONSTRAINT vat_heal_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: vat_heal vat_heal_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_heal
    ADD CONSTRAINT vat_heal_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_art vat_ilk_art_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_art
    ADD CONSTRAINT vat_ilk_art_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_dust vat_ilk_dust_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_dust
    ADD CONSTRAINT vat_ilk_dust_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_line vat_ilk_line_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_line
    ADD CONSTRAINT vat_ilk_line_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_rate vat_ilk_rate_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_rate
    ADD CONSTRAINT vat_ilk_rate_pkey PRIMARY KEY (id);


--
-- Name: vat_ilk_spot vat_ilk_spot_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_spot
    ADD CONSTRAINT vat_ilk_spot_pkey PRIMARY KEY (id);


--
-- Name: vat_init vat_init_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_init
    ADD CONSTRAINT vat_init_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: vat_init vat_init_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_init
    ADD CONSTRAINT vat_init_pkey PRIMARY KEY (id);


--
-- Name: vat_line vat_line_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_line
    ADD CONSTRAINT vat_line_pkey PRIMARY KEY (id);


--
-- Name: vat_live vat_live_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_live
    ADD CONSTRAINT vat_live_pkey PRIMARY KEY (id);


--
-- Name: vat_move vat_move_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_move
    ADD CONSTRAINT vat_move_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: vat_move vat_move_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_move
    ADD CONSTRAINT vat_move_pkey PRIMARY KEY (id);


--
-- Name: vat_sin vat_sin_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_sin
    ADD CONSTRAINT vat_sin_pkey PRIMARY KEY (id);


--
-- Name: vat_slip vat_slip_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_slip
    ADD CONSTRAINT vat_slip_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: vat_slip vat_slip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_slip
    ADD CONSTRAINT vat_slip_pkey PRIMARY KEY (id);


--
-- Name: vat_urn_art vat_urn_art_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_art
    ADD CONSTRAINT vat_urn_art_pkey PRIMARY KEY (id);


--
-- Name: vat_urn_ink vat_urn_ink_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_ink
    ADD CONSTRAINT vat_urn_ink_pkey PRIMARY KEY (id);


--
-- Name: vat_vice vat_vice_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_vice
    ADD CONSTRAINT vat_vice_pkey PRIMARY KEY (id);


--
-- Name: vow_ash vow_ash_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_ash
    ADD CONSTRAINT vow_ash_pkey PRIMARY KEY (id);


--
-- Name: vow_bump vow_bump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_bump
    ADD CONSTRAINT vow_bump_pkey PRIMARY KEY (id);


--
-- Name: vow_cow vow_cow_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_cow
    ADD CONSTRAINT vow_cow_pkey PRIMARY KEY (id);


--
-- Name: vow_fess vow_fess_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_fess
    ADD CONSTRAINT vow_fess_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: vow_fess vow_fess_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_fess
    ADD CONSTRAINT vow_fess_pkey PRIMARY KEY (id);


--
-- Name: vow_flog vow_flog_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flog
    ADD CONSTRAINT vow_flog_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: vow_flog vow_flog_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flog
    ADD CONSTRAINT vow_flog_pkey PRIMARY KEY (id);


--
-- Name: vow_hump vow_hump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_hump
    ADD CONSTRAINT vow_hump_pkey PRIMARY KEY (id);


--
-- Name: vow_row vow_row_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_row
    ADD CONSTRAINT vow_row_pkey PRIMARY KEY (id);


--
-- Name: vow_sin vow_sin_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin
    ADD CONSTRAINT vow_sin_pkey PRIMARY KEY (id);


--
-- Name: vow_sump vow_sump_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sump
    ADD CONSTRAINT vow_sump_pkey PRIMARY KEY (id);


--
-- Name: vow_vat vow_vat_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_vat
    ADD CONSTRAINT vow_vat_pkey PRIMARY KEY (id);


--
-- Name: vow_wait vow_wait_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_wait
    ADD CONSTRAINT vow_wait_pkey PRIMARY KEY (id);


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
-- Name: headers headers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_pkey PRIMARY KEY (id);


--
-- Name: light_sync_transactions light_sync_transactions_header_id_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_transactions
    ADD CONSTRAINT light_sync_transactions_header_id_hash_key UNIQUE (header_id, hash);


--
-- Name: light_sync_transactions light_sync_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_transactions
    ADD CONSTRAINT light_sync_transactions_pkey PRIMARY KEY (id);


--
-- Name: logs logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT logs_pkey PRIMARY KEY (id);


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
-- Name: queued_storage queued_storage_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage
    ADD CONSTRAINT queued_storage_pkey PRIMARY KEY (id);


--
-- Name: receipts receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_pkey PRIMARY KEY (id);


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
-- Name: price_feeds notify_pricefeeds; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER notify_pricefeeds AFTER INSERT ON maker.price_feeds FOR EACH ROW EXECUTE PROCEDURE public.notify_pricefeed();


--
-- Name: bite bite_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite
    ADD CONSTRAINT bite_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: bite bite_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.bite
    ADD CONSTRAINT bite_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id);


--
-- Name: cat_file_chop_lump cat_file_chop_lump_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_chop_lump
    ADD CONSTRAINT cat_file_chop_lump_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cat_file_chop_lump cat_file_chop_lump_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_chop_lump
    ADD CONSTRAINT cat_file_chop_lump_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: cat_file_flip cat_file_flip_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_flip
    ADD CONSTRAINT cat_file_flip_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cat_file_flip cat_file_flip_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_flip
    ADD CONSTRAINT cat_file_flip_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: cat_file_vow cat_file_vow_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_file_vow
    ADD CONSTRAINT cat_file_vow_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: cat_flip_ilk cat_flip_ilk_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_ilk
    ADD CONSTRAINT cat_flip_ilk_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: cat_ilk_chop cat_ilk_chop_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_chop
    ADD CONSTRAINT cat_ilk_chop_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: cat_ilk_flip cat_ilk_flip_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_flip
    ADD CONSTRAINT cat_ilk_flip_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: cat_ilk_lump cat_ilk_lump_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_ilk_lump
    ADD CONSTRAINT cat_ilk_lump_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: deal deal_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.deal
    ADD CONSTRAINT deal_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: dent dent_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.dent
    ADD CONSTRAINT dent_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flap_kick flap_kick_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flap_kick
    ADD CONSTRAINT flap_kick_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flip_kick flip_kick_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flip_kick
    ADD CONSTRAINT flip_kick_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: flop_kick flop_kick_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.flop_kick
    ADD CONSTRAINT flop_kick_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_drip jug_drip_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_drip
    ADD CONSTRAINT jug_drip_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_drip jug_drip_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_drip
    ADD CONSTRAINT jug_drip_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: jug_file_ilk jug_file_ilk_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_ilk
    ADD CONSTRAINT jug_file_ilk_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_file_ilk jug_file_ilk_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_ilk
    ADD CONSTRAINT jug_file_ilk_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: jug_file_repo jug_file_repo_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_repo
    ADD CONSTRAINT jug_file_repo_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_file_vow jug_file_vow_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow
    ADD CONSTRAINT jug_file_vow_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_ilk_rho jug_ilk_rho_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_rho
    ADD CONSTRAINT jug_ilk_rho_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: jug_ilk_tax jug_ilk_tax_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_tax
    ADD CONSTRAINT jug_ilk_tax_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: price_feeds price_feeds_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.price_feeds
    ADD CONSTRAINT price_feeds_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: tend tend_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.tend
    ADD CONSTRAINT tend_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: urns urns_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.urns
    ADD CONSTRAINT urns_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: vat_file_debt_ceiling vat_file_debt_ceiling_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_debt_ceiling
    ADD CONSTRAINT vat_file_debt_ceiling_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_file_ilk vat_file_ilk_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_ilk
    ADD CONSTRAINT vat_file_ilk_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_file_ilk vat_file_ilk_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_file_ilk
    ADD CONSTRAINT vat_file_ilk_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: vat_flux vat_flux_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_flux
    ADD CONSTRAINT vat_flux_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_flux vat_flux_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_flux
    ADD CONSTRAINT vat_flux_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: vat_fold vat_fold_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold
    ADD CONSTRAINT vat_fold_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_fold vat_fold_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_fold
    ADD CONSTRAINT vat_fold_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id);


--
-- Name: vat_frob vat_frob_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob
    ADD CONSTRAINT vat_frob_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_frob vat_frob_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_frob
    ADD CONSTRAINT vat_frob_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id);


--
-- Name: vat_gem vat_gem_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_gem
    ADD CONSTRAINT vat_gem_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: vat_grab vat_grab_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_grab
    ADD CONSTRAINT vat_grab_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_grab vat_grab_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_grab
    ADD CONSTRAINT vat_grab_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id);


--
-- Name: vat_heal vat_heal_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_heal
    ADD CONSTRAINT vat_heal_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_ilk_art vat_ilk_art_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_art
    ADD CONSTRAINT vat_ilk_art_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: vat_ilk_dust vat_ilk_dust_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_dust
    ADD CONSTRAINT vat_ilk_dust_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: vat_ilk_line vat_ilk_line_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_line
    ADD CONSTRAINT vat_ilk_line_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: vat_ilk_rate vat_ilk_rate_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_rate
    ADD CONSTRAINT vat_ilk_rate_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: vat_ilk_spot vat_ilk_spot_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_ilk_spot
    ADD CONSTRAINT vat_ilk_spot_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: vat_init vat_init_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_init
    ADD CONSTRAINT vat_init_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_init vat_init_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_init
    ADD CONSTRAINT vat_init_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: vat_move vat_move_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_move
    ADD CONSTRAINT vat_move_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_slip vat_slip_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_slip
    ADD CONSTRAINT vat_slip_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vat_slip vat_slip_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_slip
    ADD CONSTRAINT vat_slip_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: vat_urn_art vat_urn_art_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_art
    ADD CONSTRAINT vat_urn_art_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id);


--
-- Name: vat_urn_ink vat_urn_ink_urn_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vat_urn_ink
    ADD CONSTRAINT vat_urn_ink_urn_id_fkey FOREIGN KEY (urn_id) REFERENCES maker.urns(id);


--
-- Name: vow_fess vow_fess_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_fess
    ADD CONSTRAINT vow_fess_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: vow_flog vow_flog_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_flog
    ADD CONSTRAINT vow_flog_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: receipts blocks_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT blocks_fk FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: checked_headers checked_headers_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers
    ADD CONSTRAINT checked_headers_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: headers eth_nodes_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT eth_nodes_fk FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: full_sync_transactions full_sync_transactions_block_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_transactions
    ADD CONSTRAINT full_sync_transactions_block_id_fkey FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: light_sync_transactions light_sync_transactions_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.light_sync_transactions
    ADD CONSTRAINT light_sync_transactions_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: blocks node_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT node_fk FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: logs receipts_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT receipts_fk FOREIGN KEY (receipt_id) REFERENCES public.receipts(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

