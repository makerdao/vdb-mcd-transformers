-- +goose Up
-- +goose StatementBegin
create or replace function maker.all_ilk_states(block_height bigint, ilk_id int)
  returns setof maker.ilk_state as $$
DECLARE
  r maker.relevant_block;
BEGIN
  FOR r IN SELECT * FROM maker.get_ilk_blocks_before($1, $2)
  LOOP
    RETURN QUERY
    SELECT * FROM maker.get_ilk(r.block_number, $2::integer);
  END LOOP;
END;
$$
LANGUAGE plpgsql
STABLE SECURITY DEFINER;
-- +goose StatementEnd

-- +goose Down
drop function  maker.all_ilk_states(block_height bigint, ilk_id int);
