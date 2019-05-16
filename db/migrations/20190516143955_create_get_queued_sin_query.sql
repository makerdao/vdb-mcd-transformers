-- +goose Up
CREATE TYPE api.queued_sin AS (
  era     NUMERIC,
  tab     NUMERIC,
  -- events
  flogged BOOLEAN,
  created TIMESTAMP,
  updated TIMESTAMP
);

CREATE FUNCTION api.get_queued_sin(era NUMERIC)
  RETURNS SETOF api.queued_sin AS
$body$
  WITH
    created AS (
      SELECT vow_sin_mapping.timestamp AS era, vow_sin_mapping.block_number, (SELECT TIMESTAMP 'epoch' + block_timestamp * INTERVAL '1 second') AS datetime
      FROM maker.vow_sin_mapping
      LEFT JOIN public.headers ON hash = block_hash
      WHERE vow_sin_mapping.timestamp = $1
      ORDER BY vow_sin_mapping.block_number ASC
      LIMIT 1
    ),

    updated AS (
      SELECT vow_sin_mapping.timestamp AS era, vow_sin_mapping.block_number, (SELECT TIMESTAMP 'epoch' + block_timestamp * INTERVAL '1 second') AS datetime
      FROM maker.vow_sin_mapping
      LEFT JOIN public.headers ON hash = block_hash
      WHERE vow_sin_mapping.timestamp = $1
      ORDER BY vow_sin_mapping.block_number DESC
      LIMIT 1
    )

  SELECT $1 AS era, sin AS tab, (SELECT EXISTS(SELECT id FROM maker.vow_flog WHERE vow_flog.era = $1)) AS flogged, created.datetime, updated.datetime
  FROM maker.vow_sin_mapping vow_sin_mapping
  LEFT JOIN created ON created.era = vow_sin_mapping.timestamp
  LEFT JOIN updated ON updated.era = vow_sin_mapping.timestamp
  WHERE vow_sin_mapping.timestamp = $1
  ORDER BY vow_sin_mapping.block_number DESC
$body$
LANGUAGE sql STABLE;


-- +goose Down
DROP FUNCTION api.get_queued_sin(NUMERIC);
DROP TYPE api.queued_sin CASCADE;