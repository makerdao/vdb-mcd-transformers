-- +goose Up

CREATE TYPE api.min_or_max AS ENUM (
	'min',
	'max'
);

CREATE TYPE api.storage_transformations AS (
	min_or_max 		api.min_or_max,
	address 		BYTEA,
	block_hash 		BYTEA,
	block_height 	BIGINT,
	from_backfill 	BOOLEAN,
	status 			public.diff_status,
	storage_key 	BYTEA,
	storage_value 	BYTEA,
	created 		TIMESTAMP WITHOUT TIME ZONE,
	updated 		TIMESTAMP WITHOUT TIME ZONE
);

CREATE OR REPLACE FUNCTION api.get_storage_transformations_for_status(status TEXT)
	RETURNS SETOF api.storage_transformations
	LANGUAGE sql
	STABLE
AS
$$
(SELECT 'max'::api.min_or_max as min_or_max, address, block_hash, block_height, from_backfill, status, storage_key, storage_value, created, updated
		FROM public.storage_diff
 		WHERE status = LOWER(get_storage_transformations_for_status.status)::public.diff_status
		ORDER BY block_height DESC
		LIMIT 1)
UNION
(SELECT 'min'::api.min_or_max as min_or_max, address, block_hash, block_height, from_backfill, status, storage_key, storage_value, created, updated
		FROM public.storage_diff
 		WHERE status = LOWER(get_storage_transformations_for_status.status)::public.diff_status
 		ORDER BY block_height ASC
		LIMIT 1)
$$;

-- +goose Down

DROP FUNCTION api.get_storage_transformations_for_status(status TEXT);
DROP TYPE api.storage_transformations;
DROP TYPE api.min_or_max;
