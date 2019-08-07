-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend flip_state with ilk_state
CREATE FUNCTION api.flip_state_ilk(flip api.flip_state)
    RETURNS api.ilk_state AS
$$
SELECT *
FROM api.get_ilk((SELECT identifier FROM maker.ilks WHERE ilks.id = flip.ilk_id), flip.block_height)
$$
    LANGUAGE sql
    STABLE;


-- Extend flip_state with urn_state
CREATE FUNCTION api.flip_state_urn(flip api.flip_state)
    RETURNS SETOF api.urn_state AS
$$
SELECT *
FROM api.get_urn((SELECT identifier FROM maker.ilks WHERE ilks.id = flip.ilk_id),
                 (SELECT identifier FROM maker.urns WHERE urns.id = flip.urn_id), flip.block_height)
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.flip_state_ilk(api.flip_state);
DROP FUNCTION api.flip_state_urn(api.flip_state);
