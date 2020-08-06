-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- Function returning state for a single urn as of given block
CREATE FUNCTION api.get_urn(ilk_identifier text, urn_identifier text, block_height bigint DEFAULT api.max_block()) RETURNS api.urn_snapshot
    LANGUAGE sql STABLE STRICT
    AS $_$

SELECT urn_identifier, ilk_identifier, get_urn.block_height, ink, art, created, updated
    FROM api.urn_snapshot
    WHERE ilk_identifier = get_urn.ilk_identifier
    AND urn_identifier = get_urn.urn_identifier
    AND block_height <= get_urn.block_height
    ORDER BY updated DESC
    LIMIT 1
$_$;


-- +goose Down
DROP FUNCTION api.get_urn(TEXT, TEXT, BIGINT);
