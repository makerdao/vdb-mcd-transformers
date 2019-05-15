-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.bite_event AS (
  ilk_name     TEXT,
  -- ilk object
  urn_guy       TEXT,
  -- urn object
  ink          NUMERIC,
  art          NUMERIC,
  tab          NUMERIC,
  block_height BIGINT,
  tx_idx       INTEGER
  -- tx
);

COMMENT ON COLUMN api.bite_event.block_height IS E'@omit';
COMMENT ON COLUMN api.bite_event.tx_idx IS E'@omit';

CREATE FUNCTION api.all_bites(ilk_name TEXT)
  RETURNS SETOF api.bite_event AS
$$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = $1)

  SELECT $1 AS ilk_name, guy AS urn_guy, ink, art, tab, block_number AS block_height, tx_idx
  FROM maker.bite
  LEFT JOIN maker.urns ON bite.urn_id = urns.id
  LEFT JOIN headers ON bite.header_id = headers.id
  WHERE urns.ilk_id = (SELECT id FROM ilk)
  ORDER BY guy, block_number DESC
$$ LANGUAGE sql STABLE;


CREATE FUNCTION api.urn_bites(ilk_name TEXT, urn TEXT)
  RETURNS SETOF api.bite_event AS
$body$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = $1),
    urn AS (
      SELECT id FROM maker.urns
      WHERE ilk_id = (SELECT id FROM ilk)
        AND guy = $2
    )

  SELECT $1 AS ilk_name, $2 AS urn_guy, ink, art, tab, block_number AS block_height, tx_idx
  FROM maker.bite LEFT JOIN headers ON bite.header_id = headers.id
  WHERE bite.urn_id = (SELECT id FROM urn)
  ORDER BY block_number DESC
$body$
  LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.urn_bites(TEXT, TEXT);
DROP FUNCTION api.all_bites(TEXT);
DROP TYPE api.bite_event CASCADE;
