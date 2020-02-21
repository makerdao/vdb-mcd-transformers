-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE FUNCTION api.total_ink(ilk_identifier TEXT, block_height BIGINT DEFAULT api.max_block())
    RETURNS NUMERIC AS
$$
SELECT SUM(latest_ink_by_urn.ink)
FROM (SELECT DISTINCT ON (vat_urn_ink.urn_id) vat_urn_ink.ink
      FROM maker.ilks
               LEFT JOIN maker.urns ON urns.ilk_id = ilks.id
               LEFT JOIN maker.vat_urn_ink ON vat_urn_ink.urn_id = urns.id
               LEFT JOIN public.headers ON vat_urn_ink.header_id = headers.id
      WHERE ilks.identifier = total_ink.ilk_identifier
        AND headers.block_number <= total_ink.block_height
      ORDER BY vat_urn_ink.urn_id, headers.block_number DESC) latest_ink_by_urn
$$
    LANGUAGE sql
    STRICT
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.total_ink(TEXT, BIGINT);
