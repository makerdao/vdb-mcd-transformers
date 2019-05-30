-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION api.all_urn_states(ilk_identifier TEXT, urn_guy TEXT, block_height BIGINT DEFAULT api.max_block())
    RETURNS SETOF api.urn_state AS
$$
DECLARE
    blocks  BIGINT[];
    i       BIGINT;
    _ilk_id NUMERIC;
    _urn_id NUMERIC;
BEGIN
    SELECT id
    FROM maker.ilks
    WHERE ilks.identifier = ilk_identifier INTO _ilk_id;
    SELECT id
    FROM maker.urns
    WHERE urns.guy = urn_guy
      AND urns.ilk_id = _ilk_id INTO _urn_id;

    blocks := ARRAY(
            SELECT block_number
            FROM (SELECT block_number
                  FROM maker.vat_urn_ink
                  WHERE vat_urn_ink.urn_id = _urn_id
                    AND block_number <= all_urn_states.block_height
                  UNION
                  SELECT block_number
                  FROM maker.vat_urn_art
                  WHERE vat_urn_art.urn_id = _urn_id
                    AND block_number <= all_urn_states.block_height) inks_and_arts
            ORDER BY block_number DESC
        );

    FOREACH i IN ARRAY blocks
        LOOP
            RETURN QUERY
                SELECT * FROM api.get_urn(ilk_identifier, urn_guy, i);
        END LOOP;
END;
$$
    LANGUAGE plpgsql
    STABLE
    STRICT;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.all_urn_states(TEXT, TEXT, BIGINT);
