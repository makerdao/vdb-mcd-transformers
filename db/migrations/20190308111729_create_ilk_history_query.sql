-- +goose Up
-- +goose StatementBegin

-- Function returning the history of a given ilk as of the given block height
CREATE FUNCTION api.all_ilk_states(block_height BIGINT, ilk_name TEXT)
  RETURNS SETOF api.ilk_state AS $$
DECLARE
  r api.relevant_block;
BEGIN
  FOR r IN SELECT get_ilk_blocks_before.block_height FROM api.get_ilk_blocks_before($1, $2)
  LOOP
    RETURN QUERY
    SELECT * FROM api.get_ilk(r.block_height, $2);
  END LOOP;
END;
$$
LANGUAGE plpgsql STABLE STRICT;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.all_ilk_states(BIGINT, TEXT);
