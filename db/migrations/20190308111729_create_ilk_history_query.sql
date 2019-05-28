-- +goose Up
-- +goose StatementBegin

-- Function returning the history of a given ilk as of the given block height
CREATE FUNCTION api.all_ilk_states(ilk_identifier TEXT, block_height BIGINT DEFAULT api.max_block())
  RETURNS SETOF api.ilk_state AS $$
DECLARE
  r api.relevant_block;
BEGIN
  FOR r IN SELECT get_ilk_blocks_before.block_height FROM api.get_ilk_blocks_before($2, ilk_identifier)
  LOOP
    RETURN QUERY
    SELECT * FROM api.get_ilk(ilk_identifier, r.block_height);
  END LOOP;
END;
$$
LANGUAGE plpgsql STABLE STRICT;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.all_ilk_states(TEXT, BIGINT);
