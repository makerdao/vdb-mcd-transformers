-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION api.all_queued_sin()
  RETURNS SETOF api.queued_sin AS $$
DECLARE
  _era NUMERIC;
BEGIN
  FOR _era IN
    SELECT DISTINCT timestamp FROM maker.vow_sin_mapping
  LOOP
    RETURN QUERY
      SELECT * FROM api.get_queued_sin(_era);
  END LOOP;
END;
$$ LANGUAGE plpgsql STABLE;
-- +goose StatementEnd


-- +goose Down
DROP FUNCTION api.all_queued_sin();