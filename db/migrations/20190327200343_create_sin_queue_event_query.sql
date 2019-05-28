-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TYPE api.sin_act AS ENUM (
  'flog',
  'fess'
);

CREATE TYPE api.sin_queue_event AS (
  era          NUMERIC,
  act          api.sin_act,
  block_height BIGINT,
  tx_idx       INTEGER
  -- tx
);

COMMENT ON COLUMN api.sin_queue_event.block_height IS E'@omit';
COMMENT ON COLUMN api.sin_queue_event.tx_idx IS E'@omit';

CREATE FUNCTION api.all_sin_queue_events(era NUMERIC)
  RETURNS SETOF api.sin_queue_event AS
$$
  SELECT block_timestamp AS era, 'fess'::api.sin_act AS act, block_number AS block_height, tx_idx
  FROM maker.vow_fess
  LEFT JOIN headers ON vow_fess.header_id = headers.id
  WHERE block_timestamp = $1
  UNION
  SELECT era, 'flog'::api.sin_act AS act, block_number AS block_height, tx_idx
  FROM maker.vow_flog
  LEFT JOIN headers ON vow_flog.header_id = headers.id
  where vow_flog.era = $1
  ORDER BY block_height DESC
$$ LANGUAGE sql STABLE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.all_sin_queue_events(NUMERIC);
DROP TYPE api.sin_queue_event;
DROP TYPE api.sin_act;
