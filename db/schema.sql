--
-- PostgreSQL database dump
--

-- Dumped from database version 11.3
-- Dumped by pg_dump version 11.3

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
-- Name: postgraphile_watch; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA postgraphile_watch;


--
-- Name: bite_event; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.bite_event AS (
	ilk_name text,
	urn_id text,
	ink numeric,
	art numeric,
	tab numeric,
	block_height bigint,
	tx_idx integer
);


--
-- Name: COLUMN bite_event.block_height; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.bite_event.block_height IS '@omit';


--
-- Name: COLUMN bite_event.tx_idx; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.bite_event.tx_idx IS '@omit';


--
-- Name: era; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.era AS (
	epoch bigint,
	iso timestamp without time zone
);


--
-- Name: file_event; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.file_event AS (
	id text,
	ilk_name text,
	what text,
	data text,
	block_height bigint,
	tx_idx integer
);


--
-- Name: COLUMN file_event.block_height; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.file_event.block_height IS '@omit';


--
-- Name: COLUMN file_event.tx_idx; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.file_event.tx_idx IS '@omit';


--
-- Name: frob_event; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.frob_event AS (
	ilk_name text,
	urn_id text,
	dink numeric,
	dart numeric,
	block_height bigint,
	tx_idx integer
);


--
-- Name: COLUMN frob_event.block_height; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.frob_event.block_height IS '@omit';


--
-- Name: COLUMN frob_event.tx_idx; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.frob_event.tx_idx IS '@omit';


--
-- Name: ilk_state; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.ilk_state AS (
	ilk_id integer,
	ilk_name text,
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
	created timestamp without time zone,
	updated timestamp without time zone
);


--
-- Name: log_value; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.log_value AS (
	val numeric,
	block_number bigint,
	tx_idx integer,
	contract_address text
);


--
-- Name: COLUMN log_value.tx_idx; Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON COLUMN api.log_value.tx_idx IS '@omit';


--
-- Name: relevant_block; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.relevant_block AS (
	block_height bigint,
	block_hash text,
	ilk_id integer
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
-- Name: urn_state; Type: TYPE; Schema: api; Owner: -
--

CREATE TYPE api.urn_state AS (
	urn_id text,
	ilk_name text,
	block_height bigint,
	ink numeric,
	art numeric,
	ratio numeric,
	safe boolean,
	created timestamp without time zone,
	updated timestamp without time zone
);


--
-- Name: address_files(text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.address_files(address text) RETURNS SETOF api.file_event
    LANGUAGE sql STABLE
    AS $_$
  WITH
    lowerAddress AS (SELECT lower($1))

-- ilk files
  SELECT cat_file_chop_lump.raw_log::json->>'address' AS id, ilks.name AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.cat_file_chop_lump
  LEFT JOIN maker.ilks ON cat_file_chop_lump.ilk_id = ilks.id
  LEFT JOIN headers    ON cat_file_chop_lump.header_id = headers.id
  WHERE lower(cat_file_chop_lump.raw_log::json->>'address') = (SELECT * FROM lowerAddress)
  UNION
  SELECT cat_file_flip.raw_log::json->>'address' AS id, ilks.name AS ilk_name, what, flip AS data, block_number AS block_height, tx_idx
  FROM maker.cat_file_flip
  LEFT JOIN maker.ilks ON cat_file_flip.ilk_id = ilks.id
  LEFT JOIN headers ON cat_file_flip.header_id = headers.id
  WHERE lower(cat_file_flip.raw_log::json->>'address') = (SELECT * FROM lowerAddress)
  UNION
  SELECT jug_file_ilk.raw_log::json->>'address' AS id, ilks.name AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.jug_file_ilk
  LEFT JOIN maker.ilks ON jug_file_ilk.ilk_id = ilks.id
  LEFT JOIN headers ON jug_file_ilk.header_id = headers.id
  WHERE lower(jug_file_ilk.raw_log::json->>'address') = (SELECT * FROM lowerAddress)
  UNION
  SELECT vat_file_ilk.raw_log::json->>'address' AS id, ilks.name AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.vat_file_ilk
  LEFT JOIN maker.ilks ON vat_file_ilk.ilk_id = ilks.id
  LEFT JOIN headers ON vat_file_ilk.header_id = headers.id
  WHERE lower(vat_file_ilk.raw_log::json->>'address') = (SELECT * FROM lowerAddress)

-- contract files
  UNION
  SELECT cat_file_vow.raw_log::json->>'address' AS id, NULL AS ilk_name, what, data, block_number AS block_height, tx_idx
  FROM maker.cat_file_vow
  LEFT JOIN headers ON cat_file_vow.header_id = headers.id
  WHERE lower(cat_file_vow.raw_log::json->>'address') = (SELECT * FROM lowerAddress)
  UNION
  SELECT jug_file_base.raw_log::json->>'address' AS id, NULL AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.jug_file_base
  LEFT JOIN headers ON jug_file_base.header_id = headers.id
  WHERE lower(jug_file_base.raw_log::json->>'address') = (SELECT * FROM lowerAddress)
  UNION
  SELECT jug_file_vow.raw_log::json->>'address' AS id, NULL AS ilk_name, what, data, block_number AS block_height, tx_idx
  FROM maker.jug_file_vow
  LEFT JOIN headers ON jug_file_vow.header_id = headers.id
  WHERE lower(jug_file_vow.raw_log::json->>'address') = (SELECT * FROM lowerAddress)
  UNION
  SELECT vat_file_debt_ceiling.raw_log::json->>'address' AS id, NULL AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.vat_file_debt_ceiling
  LEFT JOIN headers on vat_file_debt_ceiling.header_id = headers.id
  WHERE lower(vat_file_debt_ceiling.raw_log::json->>'address') = (SELECT * FROM lowerAddress)

  ORDER BY block_height DESC
$_$;


--
-- Name: all_bites(text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_bites(ilk_name text) RETURNS SETOF api.bite_event
    LANGUAGE sql STABLE
    AS $_$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = $1)

  SELECT $1 AS ilk_name, guy AS urn_id, ink, art, tab, block_number AS block_height, tx_idx
  FROM maker.bite
  LEFT JOIN maker.urns ON bite.urn_id = urns.id
  LEFT JOIN headers    ON bite.header_id = headers.id
  WHERE urns.ilk_id = (SELECT id FROM ilk)
  ORDER BY guy, block_number DESC
$_$;


--
-- Name: all_frobs(text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_frobs(ilk_name text) RETURNS SETOF api.frob_event
    LANGUAGE sql STABLE
    AS $_$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = $1)

  SELECT $1 AS ilk_name, guy AS urn_id, dink, dart, block_number AS block_height, tx_idx
  FROM maker.vat_frob
  LEFT JOIN maker.urns ON vat_frob.urn_id = urns.id
  LEFT JOIN headers    ON vat_frob.header_id = headers.id
  WHERE urns.ilk_id = (SELECT id FROM ilk)
  ORDER BY guy, block_number DESC
$_$;


--
-- Name: all_ilk_states(bigint, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_ilk_states(block_height bigint, ilk_id integer) RETURNS SETOF api.ilk_state
    LANGUAGE plpgsql STABLE
    AS $_$
DECLARE
  r api.relevant_block;
BEGIN
  FOR r IN SELECT * FROM api.get_ilk_blocks_before($1, $2)
  LOOP
    RETURN QUERY
    SELECT * FROM api.get_ilk(r.block_height, $2::INTEGER);
  END LOOP;
END;
$_$;


--
-- Name: all_ilks(bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_ilks(block_height bigint) RETURNS SETOF api.ilk_state
    LANGUAGE sql STABLE
    AS $_$
WITH rates AS (
  SELECT DISTINCT ON (ilk_id) rate, ilk_id, block_hash
  FROM maker.vat_ilk_rate
  WHERE block_number <= $1
  ORDER BY ilk_id, block_number DESC
), arts AS (
  SELECT DISTINCT ON (ilk_id) art, ilk_id, block_hash
  FROM maker.vat_ilk_art
  WHERE block_number <= $1
  ORDER BY ilk_id, block_number DESC
), spots AS (
  SELECT DISTINCT ON (ilk_id) spot, ilk_id, block_hash
  FROM maker.vat_ilk_spot
  WHERE block_number <= $1
  ORDER BY ilk_id, block_number DESC
), lines AS (
  SELECT DISTINCT ON (ilk_id) line, ilk_id, block_hash
  FROM maker.vat_ilk_line
  WHERE block_number <= $1
  ORDER BY ilk_id, block_number DESC
), dusts AS (
  SELECT DISTINCT ON (ilk_id) dust, ilk_id, block_hash
  FROM maker.vat_ilk_dust
  WHERE block_number <= $1
  ORDER BY ilk_id, block_number DESC
), chops AS (
  SELECT DISTINCT ON (ilk_id) chop, ilk_id, block_hash
  FROM maker.cat_ilk_chop
  WHERE block_number <= $1
  ORDER BY ilk_id, block_number DESC
), lumps AS (
  SELECT DISTINCT ON (ilk_id) lump, ilk_id, block_hash
  FROM maker.cat_ilk_lump
  WHERE block_number <= $1
  ORDER BY ilk_id, block_number DESC
), flips AS (
  SELECT DISTINCT ON (ilk_id) flip, ilk_id, block_hash
  FROM maker.cat_ilk_flip
  WHERE block_number <= $1
  ORDER BY ilk_id, block_number DESC
), rhos AS (
  SELECT DISTINCT ON (ilk_id) rho, ilk_id, block_hash
  FROM maker.jug_ilk_rho
  WHERE block_number <= $1
  ORDER BY ilk_id, block_number DESC
), duties AS (
  SELECT DISTINCT ON (ilk_id) duty, ilk_id, block_hash
  FROM maker.jug_ilk_duty
  WHERE block_number <= $1
  ORDER BY ilk_id, block_number DESC
)
  SELECT
    ilks.id,
    ilks.name,
    $1 block_height,
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
    (
      SELECT TIMESTAMP 'epoch' + h.block_timestamp * INTERVAL '1 second'
      FROM api.get_ilk_blocks_before($1, ilks.id) b
      JOIN headers h on h.block_number = b.block_height
      ORDER BY h.block_number DESC
      LIMIT 1
    ),
    (
      SELECT TIMESTAMP 'epoch' + h.block_timestamp * INTERVAL '1 second'
      FROM api.get_ilk_blocks_before($1, ilks.id) b
      JOIN headers h on h.block_number = b.block_height
      ORDER BY h.block_number ASC
      LIMIT 1
    )
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
    duties.duty is not null
  )
$_$;


--
-- Name: all_urn_states(text, text, bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_urn_states(ilk_name text, urn text, block_height bigint) RETURNS SETOF api.urn_state
    LANGUAGE plpgsql STABLE
    AS $_$
DECLARE
  blocks BIGINT[];
  i BIGINT;
  _ilk_id NUMERIC;
  _urn_id NUMERIC;
BEGIN
  SELECT id FROM maker.ilks WHERE ilks.name = $1 INTO _ilk_id;
  SELECT id FROM maker.urns WHERE urns.guy = $2 AND urns.ilk_id = _ilk_id INTO _urn_id;

  blocks := ARRAY(
    SELECT block_number
    FROM (
       SELECT block_number
       FROM maker.vat_urn_ink
       WHERE vat_urn_ink.urn_id = _urn_id
         AND block_number <= $3
       UNION
       SELECT block_number
       FROM maker.vat_urn_art
       WHERE vat_urn_art.urn_id = _urn_id
         AND block_number <= $3
     ) inks_and_arts
    ORDER BY block_number DESC
  );

  FOREACH i IN ARRAY blocks
    LOOP
      RETURN QUERY
        SELECT * FROM api.get_urn(ilk_name, urn, i);
    END LOOP;
END;
$_$;


--
-- Name: all_urns(bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.all_urns(block_height bigint) RETURNS SETOF api.urn_state
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
    SELECT urn_id, (SELECT TIMESTAMP 'epoch' + block_timestamp * INTERVAL '1 second') AS datetime
    FROM
      (
        SELECT DISTINCT ON (urn_id) urn_id, block_hash FROM maker.vat_urn_ink
        ORDER BY urn_id, block_number ASC
      ) earliest_blocks
        LEFT JOIN public.headers ON hash = block_hash
  ),

  updated AS (
    SELECT DISTINCT ON (urn_id) urn_id, (SELECT TIMESTAMP 'epoch' + block_timestamp * INTERVAL '1 second') AS datetime
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

SELECT urns.guy, ilks.name, $1, inks.ink, arts.art, ratios.ratio,
       COALESCE(safe.safe, arts.art = 0), created.datetime, updated.datetime
FROM inks
  LEFT JOIN arts       ON arts.urn_id = inks.urn_id
  LEFT JOIN urns       ON arts.urn_id = urns.urn_id
  LEFT JOIN ratios     ON ratios.guy = urns.guy
  LEFT JOIN safe       ON safe.guy = ratios.guy
  LEFT JOIN created    ON created.urn_id = urns.urn_id
  LEFT JOIN updated    ON updated.urn_id = urns.urn_id
  LEFT JOIN maker.ilks ON ilks.id = urns.ilk_id
$_$;


--
-- Name: bite_event_ilk(api.bite_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bite_event_ilk(event api.bite_event) RETURNS SETOF api.ilk_state
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM api.get_ilk(
     event.block_height,
     (SELECT id FROM maker.ilks WHERE name = event.ilk_name))
$$;


--
-- Name: bite_event_tx(api.bite_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bite_event_tx(event api.bite_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
  SELECT txs.hash, txs.tx_index, headers.block_number AS block_height, headers.hash, tx_from, tx_to
  FROM public.header_sync_transactions txs
         LEFT JOIN headers ON txs.header_id = headers.id
  WHERE block_number <= event.block_height AND txs.tx_index = event.tx_idx
  ORDER BY block_height DESC
$$;


--
-- Name: bite_event_urn(api.bite_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bite_event_urn(event api.bite_event) RETURNS SETOF api.urn_state
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM api.get_urn(event.ilk_name, event.urn_id, event.block_height)
$$;


--
-- Name: file_event_ilk(api.file_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.file_event_ilk(event api.file_event) RETURNS SETOF api.ilk_state
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM api.get_ilk(
    event.block_height,
    (SELECT id FROM maker.ilks WHERE name = event.ilk_name)
  )
$$;


--
-- Name: file_event_tx(api.file_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.file_event_tx(event api.file_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
  SELECT txs.hash, txs.tx_index, headers.block_number AS block_height, headers.hash, tx_from, tx_to
  FROM public.header_sync_transactions txs
  LEFT JOIN headers ON txs.header_id = headers.id
  WHERE block_number <= event.block_height AND txs.tx_index = event.tx_idx
  ORDER BY block_height DESC
  LIMIT 1
$$;


--
-- Name: frob_event_ilk(api.frob_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.frob_event_ilk(event api.frob_event) RETURNS SETOF api.ilk_state
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM api.get_ilk(
    event.block_height,
    (SELECT id FROM maker.ilks WHERE name = event.ilk_name))
$$;


--
-- Name: frob_event_tx(api.frob_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.frob_event_tx(event api.frob_event) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
  SELECT txs.hash, txs.tx_index, headers.block_number AS block_height, headers.hash, tx_from, tx_to
  FROM public.header_sync_transactions txs
  LEFT JOIN headers ON txs.header_id = headers.id
  WHERE block_number <= event.block_height AND txs.tx_index = event.tx_idx
  ORDER BY block_height DESC
  LIMIT 1 -- Should always be true anyway?
$$;


--
-- Name: frob_event_urn(api.frob_event); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.frob_event_urn(event api.frob_event) RETURNS SETOF api.urn_state
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM api.get_urn(event.ilk_name, event.urn_id, event.block_height)
$$;


--
-- Name: get_ilk(bigint, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_ilk(block_height bigint, _ilk_id integer) RETURNS api.ilk_state
    LANGUAGE sql STABLE
    AS $_$
WITH rates AS (
    SELECT
      rate,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_rate
    WHERE ilk_id = _ilk_id
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), arts AS (
    SELECT
      art,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_art
    WHERE ilk_id = _ilk_id
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), spots AS (
    SELECT
      spot,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_spot
    WHERE ilk_id = _ilk_id
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), lines AS (
    SELECT
      line,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_line
    WHERE ilk_id = _ilk_id
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), dusts AS (
    SELECT
      dust,
      ilk_id,
      block_hash
    FROM maker.vat_ilk_dust
    WHERE ilk_id = _ilk_id
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), chops AS (
    SELECT
      chop,
      ilk_id,
      block_hash
    FROM maker.cat_ilk_chop
    WHERE ilk_id = _ilk_id
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), lumps AS (
    SELECT
      lump,
      ilk_id,
      block_hash
    FROM maker.cat_ilk_lump
    WHERE ilk_id = _ilk_id
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), flips AS (
    SELECT
      flip,
      ilk_id,
      block_hash
    FROM maker.cat_ilk_flip
    WHERE ilk_id = _ilk_id
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), rhos AS (
    SELECT
      rho,
      ilk_id,
      block_hash
    FROM maker.jug_ilk_rho
    WHERE ilk_id = _ilk_id
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), duties AS (
    SELECT
      duty,
      ilk_id,
      block_hash
    FROM maker.jug_ilk_duty
    WHERE ilk_id = _ilk_id
          AND block_number <= $1
    ORDER BY ilk_id, block_number DESC
    LIMIT 1
), relevant_blocks AS (
  SELECT * FROM api.get_ilk_blocks_before($1, $2)
), created AS (
    SELECT DISTINCT ON (relevant_blocks.ilk_id, relevant_blocks.block_height)
      relevant_blocks.block_height,
      relevant_blocks.block_hash,
      relevant_blocks.ilk_id,
      (SELECT TIMESTAMP 'epoch' + headers.block_timestamp * INTERVAL '1 second') as datetime
    FROM relevant_blocks
      LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
    ORDER BY relevant_blocks.block_height ASC
    LIMIT 1
), updated AS (
    SELECT DISTINCT ON (relevant_blocks.ilk_id, relevant_blocks.block_height)
      relevant_blocks.block_height,
      relevant_blocks.block_hash,
      relevant_blocks.ilk_id,
      (SELECT TIMESTAMP 'epoch' + headers.block_timestamp * INTERVAL '1 second') as datetime
    FROM relevant_blocks
      LEFT JOIN public.headers AS headers on headers.hash = relevant_blocks.block_hash
    ORDER BY relevant_blocks.block_height DESC
    LIMIT 1
)

SELECT
  ilks.id,
  ilks.name,
  $1 block_height,
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
  duties.duty is not null
)
$_$;


--
-- Name: get_ilk_blocks_before(bigint, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_ilk_blocks_before(block_height bigint, ilk_id integer) RETURNS SETOF api.relevant_block
    LANGUAGE sql STABLE
    AS $_$
SELECT
  block_number AS block_height,
  block_hash,
  ilk_id
FROM maker.vat_ilk_rate
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number AS block_height,
  block_hash,
  ilk_id
FROM maker.vat_ilk_art
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number AS block_height,
  block_hash,
  ilk_id
FROM maker.vat_ilk_spot
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number AS block_height,
  block_hash,
  ilk_id
FROM maker.vat_ilk_line
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number AS block_height,
  block_hash,
  ilk_id
FROM maker.vat_ilk_dust
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number AS block_height,
  block_hash,
  ilk_id
FROM maker.cat_ilk_chop
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number AS block_height,
  block_hash,
  ilk_id
FROM maker.cat_ilk_lump
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number AS block_height,
  block_hash,
  ilk_id
FROM maker.cat_ilk_flip
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number AS block_height,
  block_hash,
  ilk_id
FROM maker.jug_ilk_rho
WHERE block_number <= $1
      AND ilk_id = $2
UNION
SELECT
  block_number AS block_height,
  block_hash,
  ilk_id
FROM maker.jug_ilk_duty
WHERE block_number <= $1
      AND ilk_id = $2
$_$;


--
-- Name: FUNCTION get_ilk_blocks_before(block_height bigint, ilk_id integer); Type: COMMENT; Schema: api; Owner: -
--

COMMENT ON FUNCTION api.get_ilk_blocks_before(block_height bigint, ilk_id integer) IS '@omit';


--
-- Name: get_urn(text, text, bigint); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_urn(ilk text, urn text, block_height bigint) RETURNS api.urn_state
    LANGUAGE sql STABLE
    AS $_$
WITH
  urn AS (
    SELECT urns.id AS urn_id, ilks.id AS ilk_id, ilks.ilk, urns.guy
    FROM maker.urns urns
    LEFT JOIN maker.ilks ilks
    ON urns.ilk_id = ilks.id
    WHERE ilks.name = $1 AND urns.guy = $2
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
    WHERE ilk_id = (SELECT ilk_id FROM urn) AND block_number <= block_height
    ORDER BY ilk_id, block_number DESC
  ),

  spot AS ( -- Get latest price update for ilk. Problematic from update frequency, slow query?
    SELECT DISTINCT ON (ilk_id) ilk_id, spot, block_number
    FROM maker.vat_ilk_spot
    WHERE ilk_id = (SELECT ilk_id FROM urn) AND block_number <= block_height
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
    SELECT urn_id, (SELECT TIMESTAMP 'epoch' + block_timestamp * INTERVAL '1 second') AS datetime
    FROM
      (
        SELECT DISTINCT ON (urn_id) urn_id, block_hash FROM maker.vat_urn_ink
        WHERE urn_id = (SELECT urn_id from urn where guy = $2)
        ORDER BY urn_id, block_number ASC
      ) earliest_blocks
        LEFT JOIN public.headers ON hash = block_hash
  ),

  updated AS (
    SELECT DISTINCT ON (urn_id) urn_id, (SELECT TIMESTAMP 'epoch' + block_timestamp * INTERVAL '1 second') AS datetime
    FROM
      (
        SELECT urn_id, block_number FROM ink
        UNION
        SELECT urn_id, block_number FROM art
      ) last_blocks
        LEFT JOIN public.headers ON headers.block_number = last_blocks.block_number
    ORDER BY urn_id, block_timestamp DESC
  )

SELECT $2 AS urn_id, $1 AS ilk_name, $3 AS block_height, ink.ink, art.art, ratio.ratio,
       COALESCE(safe.safe, art.art = 0), created.datetime, updated.datetime
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
-- Name: ilk_files(text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_files(ilk_name text) RETURNS SETOF api.file_event
    LANGUAGE sql STABLE
    AS $_$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = $1)

  SELECT cat_file_chop_lump.raw_log::json->>'address' AS id, $1 AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.cat_file_chop_lump
  LEFT JOIN headers ON cat_file_chop_lump.header_id = headers.id
  WHERE cat_file_chop_lump.ilk_id = (SELECT id FROM ilk)
  UNION
  SELECT cat_file_flip.raw_log::json->>'address' AS id, $1 AS ilk_name, what, flip AS data, block_number AS block_height, tx_idx
  FROM maker.cat_file_flip
  LEFT JOIN headers ON cat_file_flip.header_id = headers.id
  WHERE cat_file_flip.ilk_id = (SELECT id FROM ilk)
  UNION
  SELECT jug_file_ilk.raw_log::json->>'address' AS id, $1 AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.jug_file_ilk
  LEFT JOIN headers ON jug_file_ilk.header_id = headers.id
  WHERE jug_file_ilk.ilk_id = (SELECT id FROM ilk)
  UNION
  SELECT vat_file_ilk.raw_log::json->>'address' AS id, $1 AS ilk_name, what, data::text, block_number AS block_height, tx_idx
  FROM maker.vat_file_ilk
  LEFT JOIN headers ON vat_file_ilk.header_id = headers.id
  WHERE vat_file_ilk.ilk_id = (SELECT id FROM ilk)
  ORDER BY block_height DESC
$_$;


--
-- Name: ilk_state_files(api.ilk_state); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_state_files(state api.ilk_state) RETURNS SETOF api.file_event
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM api.ilk_files(state.ilk_name)
  WHERE block_height <= state.block_height
$$;


--
-- Name: ilk_state_frobs(api.ilk_state); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.ilk_state_frobs(state api.ilk_state) RETURNS SETOF api.frob_event
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM api.all_frobs(state.ilk_name)
  WHERE block_height <= state.block_height
$$;


--
-- Name: log_value_tx(api.log_value); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.log_value_tx(priceupdate api.log_value) RETURNS api.tx
    LANGUAGE sql STABLE
    AS $$
  SELECT txs.hash, txs.tx_index, headers.block_number, headers.hash, txs.tx_from, txs.tx_to
    FROM maker.pip_log_value plv
    LEFT JOIN public.header_sync_transactions txs ON plv.header_id = txs.header_id
    LEFT JOIN headers ON plv.header_id = headers.id
    WHERE headers.block_number = priceUpdate.block_number AND priceUpdate.tx_idx = txs.tx_index
    ORDER BY headers.block_number DESC
$$;


--
-- Name: log_values(integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.log_values(begintime integer, endtime integer) RETURNS SETOF api.log_value
    LANGUAGE sql STABLE
    AS $_$
  SELECT val, pip_log_value.block_number, tx_idx, contract_address FROM maker.pip_log_value
    LEFT JOIN public.headers ON pip_log_value.header_id = headers.id
    WHERE block_timestamp BETWEEN $1 AND $2
$_$;


--
-- Name: tx_era(api.tx); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.tx_era(tx api.tx) RETURNS api.era
    LANGUAGE sql STABLE
    AS $$
SELECT block_timestamp::BIGINT AS "epoch", (SELECT TIMESTAMP 'epoch' + block_timestamp * INTERVAL '1 second') AS iso
  FROM headers WHERE block_number = tx.block_height
$$;


--
-- Name: urn_frobs(text, text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.urn_frobs(ilk_name text, urn text) RETURNS SETOF api.frob_event
    LANGUAGE sql STABLE
    AS $_$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = $1),
    urn AS (
      SELECT id FROM maker.urns
      WHERE ilk_id = (SELECT id FROM ilk)
        AND guy = $2
    )

  SELECT $1 AS ilk_name, $2 AS urn_id, dink, dart, block_number AS block_height, tx_idx
  FROM maker.vat_frob LEFT JOIN headers ON vat_frob.header_id = headers.id
  WHERE vat_frob.urn_id = (SELECT id FROM urn)
  ORDER BY block_number DESC
$_$;


--
-- Name: urn_state_frobs(api.urn_state); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.urn_state_frobs(state api.urn_state) RETURNS SETOF api.frob_event
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM api.urn_frobs(state.ilk_name, state.urn_id)
  WHERE block_height <= state.block_height
$$;


--
-- Name: urn_state_ilk(api.urn_state); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.urn_state_ilk(state api.urn_state) RETURNS api.ilk_state
    LANGUAGE sql STABLE
    AS $$
  SELECT * FROM api.get_ilk(
    state.block_height,
    (SELECT id FROM maker.ilks WHERE name = state.ilk_name)
  )
$$;


--
-- Name: notify_watchers_ddl(); Type: FUNCTION; Schema: postgraphile_watch; Owner: -
--

CREATE FUNCTION postgraphile_watch.notify_watchers_ddl() RETURNS event_trigger
    LANGUAGE plpgsql
    AS $$
begin
  perform pg_notify(
    'postgraphile_watch',
    json_build_object(
      'type',
      'ddl',
      'payload',
      (select json_agg(json_build_object('schema', schema_name, 'command', command_tag)) from pg_event_trigger_ddl_commands() as x)
    )::text
  );
end;
$$;


--
-- Name: notify_watchers_drop(); Type: FUNCTION; Schema: postgraphile_watch; Owner: -
--

CREATE FUNCTION postgraphile_watch.notify_watchers_drop() RETURNS event_trigger
    LANGUAGE plpgsql
    AS $$
begin
  perform pg_notify(
    'postgraphile_watch',
    json_build_object(
      'type',
      'drop',
      'payload',
      (select json_agg(distinct x.schema_name) from pg_event_trigger_dropped_objects() as x)
    )::text
  );
end;
$$;


--
-- Name: notify_pip_log_value(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.notify_pip_log_value() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  PERFORM pg_notify(
    CAST('postgraphile:pip_log_value' AS text),
    json_build_object('__node__', json_build_array('pip_log_value', NEW.id))::text
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
-- Name: TABLE bite; Type: COMMENT; Schema: maker; Owner: -
--

COMMENT ON TABLE maker.bite IS '@name raw_bites';


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
    ilk text NOT NULL,
    name text NOT NULL
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
-- Name: jug_file_base; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.jug_file_base (
    id integer NOT NULL,
    header_id integer NOT NULL,
    what text,
    data numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
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
-- Name: pip_log_value; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE maker.pip_log_value (
    id integer NOT NULL,
    block_number bigint NOT NULL,
    header_id integer NOT NULL,
    contract_address text,
    val numeric,
    log_idx integer NOT NULL,
    tx_idx integer NOT NULL,
    raw_log jsonb
);


--
-- Name: pip_log_value_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE maker.pip_log_value_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pip_log_value_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE maker.pip_log_value_id_seq OWNED BY maker.pip_log_value.id;


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
    wad numeric,
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
    usr text,
    wad numeric,
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
    "timestamp" numeric,
    sin numeric
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
    flip_kick_checked integer DEFAULT 0 NOT NULL,
    vat_frob_checked integer DEFAULT 0 NOT NULL,
    pip_log_value_checked integer DEFAULT 0 NOT NULL,
    tend_checked integer DEFAULT 0 NOT NULL,
    bite_checked integer DEFAULT 0 NOT NULL,
    dent_checked integer DEFAULT 0 NOT NULL,
    vat_file_debt_ceiling_checked integer DEFAULT 0 NOT NULL,
    vat_file_ilk_checked integer DEFAULT 0 NOT NULL,
    vat_init_checked integer DEFAULT 0 NOT NULL,
    jug_file_base_checked integer DEFAULT 0 NOT NULL,
    jug_file_ilk_checked integer DEFAULT 0 NOT NULL,
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
-- Name: full_sync_receipts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.full_sync_receipts (
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
-- Name: header_sync_receipts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.header_sync_receipts (
    id integer NOT NULL,
    transaction_id integer NOT NULL,
    header_id integer NOT NULL,
    contract_address character varying(42),
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
-- Name: jug_vat id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vat ALTER COLUMN id SET DEFAULT nextval('maker.jug_vat_id_seq'::regclass);


--
-- Name: jug_vow id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_vow ALTER COLUMN id SET DEFAULT nextval('maker.jug_vow_id_seq'::regclass);


--
-- Name: pip_log_value id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pip_log_value ALTER COLUMN id SET DEFAULT nextval('maker.pip_log_value_id_seq'::regclass);


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
-- Name: logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs ALTER COLUMN id SET DEFAULT nextval('public.logs_id_seq'::regclass);


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
-- Name: cat_flip_ilk cat_flip_ilk_block_number_block_hash_flip_ilk_id_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_ilk
    ADD CONSTRAINT cat_flip_ilk_block_number_block_hash_flip_ilk_id_key UNIQUE (block_number, block_hash, flip, ilk_id);


--
-- Name: cat_flip_ilk cat_flip_ilk_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_ilk
    ADD CONSTRAINT cat_flip_ilk_pkey PRIMARY KEY (id);


--
-- Name: cat_flip_ink cat_flip_ink_block_number_block_hash_flip_ink_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_ink
    ADD CONSTRAINT cat_flip_ink_block_number_block_hash_flip_ink_key UNIQUE (block_number, block_hash, flip, ink);


--
-- Name: cat_flip_ink cat_flip_ink_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_ink
    ADD CONSTRAINT cat_flip_ink_pkey PRIMARY KEY (id);


--
-- Name: cat_flip_tab cat_flip_tab_block_number_block_hash_flip_tab_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_tab
    ADD CONSTRAINT cat_flip_tab_block_number_block_hash_flip_tab_key UNIQUE (block_number, block_hash, flip, tab);


--
-- Name: cat_flip_tab cat_flip_tab_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_tab
    ADD CONSTRAINT cat_flip_tab_pkey PRIMARY KEY (id);


--
-- Name: cat_flip_urn cat_flip_urn_block_number_block_hash_flip_urn_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_urn
    ADD CONSTRAINT cat_flip_urn_block_number_block_hash_flip_urn_key UNIQUE (block_number, block_hash, flip, urn);


--
-- Name: cat_flip_urn cat_flip_urn_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_flip_urn
    ADD CONSTRAINT cat_flip_urn_pkey PRIMARY KEY (id);


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
-- Name: cat_nflip cat_nflip_block_number_block_hash_nflip_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_nflip
    ADD CONSTRAINT cat_nflip_block_number_block_hash_nflip_key UNIQUE (block_number, block_hash, nflip);


--
-- Name: cat_nflip cat_nflip_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.cat_nflip
    ADD CONSTRAINT cat_nflip_pkey PRIMARY KEY (id);


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
-- Name: ilks ilks_name_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.ilks
    ADD CONSTRAINT ilks_name_key UNIQUE (name);


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
-- Name: jug_file_base jug_file_base_header_id_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_base
    ADD CONSTRAINT jug_file_base_header_id_tx_idx_log_idx_key UNIQUE (header_id, tx_idx, log_idx);


--
-- Name: jug_file_base jug_file_base_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_base
    ADD CONSTRAINT jug_file_base_pkey PRIMARY KEY (id);


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
-- Name: pip_log_value pip_log_value_header_id_contract_address_tx_idx_log_idx_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pip_log_value
    ADD CONSTRAINT pip_log_value_header_id_contract_address_tx_idx_log_idx_key UNIQUE (header_id, contract_address, tx_idx, log_idx);


--
-- Name: pip_log_value pip_log_value_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pip_log_value
    ADD CONSTRAINT pip_log_value_pkey PRIMARY KEY (id);


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
-- Name: vow_cow vow_cow_block_number_block_hash_cow_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_cow
    ADD CONSTRAINT vow_cow_block_number_block_hash_cow_key UNIQUE (block_number, block_hash, cow);


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
-- Name: vow_row vow_row_block_number_block_hash_row_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_row
    ADD CONSTRAINT vow_row_block_number_block_hash_row_key UNIQUE (block_number, block_hash, "row");


--
-- Name: vow_row vow_row_pkey; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_row
    ADD CONSTRAINT vow_row_pkey PRIMARY KEY (id);


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
-- Name: vow_sin_mapping vow_sin_mapping_block_number_block_hash_timestamp_sin_key; Type: CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.vow_sin_mapping
    ADD CONSTRAINT vow_sin_mapping_block_number_block_hash_timestamp_sin_key UNIQUE (block_number, block_hash, "timestamp", sin);


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
-- Name: pip_log_value notify_pip_log_value; Type: TRIGGER; Schema: maker; Owner: -
--

CREATE TRIGGER notify_pip_log_value AFTER INSERT ON maker.pip_log_value FOR EACH ROW EXECUTE PROCEDURE public.notify_pip_log_value();


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
-- Name: jug_file_base jug_file_base_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_base
    ADD CONSTRAINT jug_file_base_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
-- Name: jug_file_vow jug_file_vow_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_file_vow
    ADD CONSTRAINT jug_file_vow_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: jug_ilk_duty jug_ilk_duty_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_duty
    ADD CONSTRAINT jug_ilk_duty_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: jug_ilk_rho jug_ilk_rho_ilk_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.jug_ilk_rho
    ADD CONSTRAINT jug_ilk_rho_ilk_id_fkey FOREIGN KEY (ilk_id) REFERENCES maker.ilks(id);


--
-- Name: pip_log_value pip_log_value_header_id_fkey; Type: FK CONSTRAINT; Schema: maker; Owner: -
--

ALTER TABLE ONLY maker.pip_log_value
    ADD CONSTRAINT pip_log_value_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


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
-- Name: full_sync_transactions full_sync_transactions_block_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_transactions
    ADD CONSTRAINT full_sync_transactions_block_id_fkey FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


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
-- Name: logs receipts_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs
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
-- Name: postgraphile_watch_ddl; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER postgraphile_watch_ddl ON ddl_command_end
         WHEN TAG IN ('ALTER AGGREGATE', 'ALTER DOMAIN', 'ALTER EXTENSION', 'ALTER FOREIGN TABLE', 'ALTER FUNCTION', 'ALTER POLICY', 'ALTER SCHEMA', 'ALTER TABLE', 'ALTER TYPE', 'ALTER VIEW', 'COMMENT', 'CREATE AGGREGATE', 'CREATE DOMAIN', 'CREATE EXTENSION', 'CREATE FOREIGN TABLE', 'CREATE FUNCTION', 'CREATE INDEX', 'CREATE POLICY', 'CREATE RULE', 'CREATE SCHEMA', 'CREATE TABLE', 'CREATE TABLE AS', 'CREATE VIEW', 'DROP AGGREGATE', 'DROP DOMAIN', 'DROP EXTENSION', 'DROP FOREIGN TABLE', 'DROP FUNCTION', 'DROP INDEX', 'DROP OWNED', 'DROP POLICY', 'DROP RULE', 'DROP SCHEMA', 'DROP TABLE', 'DROP TYPE', 'DROP VIEW', 'GRANT', 'REVOKE', 'SELECT INTO')
   EXECUTE PROCEDURE postgraphile_watch.notify_watchers_ddl();


--
-- Name: postgraphile_watch_drop; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER postgraphile_watch_drop ON sql_drop
   EXECUTE PROCEDURE postgraphile_watch.notify_watchers_drop();


--
-- PostgreSQL database dump complete
--

