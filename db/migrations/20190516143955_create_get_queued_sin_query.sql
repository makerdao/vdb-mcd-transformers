-- +goose Up
CREATE TYPE api.queued_sin AS
(
    era     NUMERIC,
    tab     NUMERIC,
    -- events
    flogged BOOLEAN,
    created TIMESTAMP,
    updated TIMESTAMP
);

CREATE FUNCTION api.get_queued_sin(era NUMERIC)
    RETURNS api.queued_sin AS
$body$
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
$body$
    LANGUAGE sql
    STABLE
    STRICT;

-- +goose Down
DROP FUNCTION api.get_queued_sin(NUMERIC);
DROP TYPE api.queued_sin CASCADE;