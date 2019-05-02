-- +goose Up
-- +goose StatementBegin

-- Function returning the history of a given ilk as of the given block height
CREATE OR REPLACE FUNCTION maker.all_ilk_states(block_height BIGINT, ilk_id INT)
  RETURNS SETOF maker.ilk_state AS $$
DECLARE
  r maker.relevant_block;
BEGIN
  FOR r IN SELECT * FROM maker.get_ilk_blocks_before($1, $2)
  LOOP
    RETURN QUERY
    SELECT * FROM maker.get_ilk(r.block_height, $2::INTEGER);
  END LOOP;
END;
$$
LANGUAGE plpgsql
STABLE SECURITY DEFINER;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION maker.all_ilk_states(BIGINT, INT);
