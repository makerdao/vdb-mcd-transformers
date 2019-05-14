-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE maker.bite_event AS (
  ilk_name     TEXT,
  -- ilk object
  urn_id       TEXT,
  ink          NUMERIC,
  art          NUMERIC,
  tab          NUMERIC,
  block_height BIGINT,
  tx_idx       INTEGER
  -- tx
);

COMMENT ON COLUMN maker.bite_event.block_height IS E'@omit';
COMMENT ON COLUMN maker.bite_event.tx_idx IS E'@omit';

CREATE OR REPLACE FUNCTION maker.all_bites(ilk_name TEXT)
  RETURNS SETOF maker.bite_event AS
$$
  WITH
    ilk AS (SELECT id FROM maker.ilks WHERE ilks.name = $1)

  SELECT $1 AS ilk_name, guy AS urn_id, ink, art, tab, block_number AS block_height, tx_idx
  FROM maker.bite
  LEFT JOIN maker.urns ON bite.urn_id = urns.id
  LEFT JOIN headers    ON bite.header_id = headers.id
  WHERE urns.ilk_id = (SELECT id FROM ilk)
  ORDER BY guy, block_number DESC
$$ LANGUAGE sql STABLE;


-- Extend type bite_event with ilk field
CREATE OR REPLACE FUNCTION maker.bite_event_ilk(event maker.bite_event)
  RETURNS SETOF maker.ilk_state AS
$$
  SELECT * FROM maker.get_ilk(
     event.block_height,
     (SELECT id FROM maker.ilks WHERE name = event.ilk_name))
$$ LANGUAGE sql STABLE;


-- Extend type bite_event with urn field
CREATE OR REPLACE FUNCTION maker.bite_event_urn(event maker.bite_event)
  RETURNS SETOF maker.urn_state AS
$$
  SELECT * FROM maker.get_urn(event.ilk_name, event.urn_id, event.block_height)
$$ LANGUAGE sql STABLE;


-- Extend type bite_event with txs field
CREATE OR REPLACE FUNCTION maker.bite_event_tx(event maker.bite_event)
  RETURNS maker.tx AS
$$
  SELECT txs.hash, txs.tx_index, headers.block_number AS block_height, headers.hash, tx_from, tx_to
  FROM public.header_sync_transactions txs
         LEFT JOIN headers ON txs.header_id = headers.id
  WHERE block_number <= event.block_height AND txs.tx_index = event.tx_idx
  ORDER BY block_height DESC
$$ LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION maker.bite_event_tx(maker.bite_event);
DROP FUNCTION maker.bite_event_urn(maker.bite_event);
DROP FUNCTION maker.bite_event_ilk(maker.bite_event);
DROP FUNCTION maker.all_bites(TEXT);
DROP TYPE maker.bite_event CASCADE;
