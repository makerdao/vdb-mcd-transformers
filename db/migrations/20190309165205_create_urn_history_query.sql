-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.get_urn_history_before_block(ilk TEXT, urn TEXT, block_height NUMERIC)
  RETURNS SETOF maker.urn_state AS $$
DECLARE
  i NUMERIC;
  ilkId NUMERIC;
  urnId NUMERIC;
BEGIN
  SELECT id FROM maker.ilks WHERE ilks.ilk = $1 INTO ilkId;
  SELECT id FROM maker.urns WHERE urns.guy = $2 AND urns.ilk = ilkID INTO urnId;

  CREATE TEMP TABLE updated ON COMMIT DROP AS
  SELECT block_number FROM (
    SELECT block_number FROM maker.vat_urn_ink
    WHERE vat_urn_ink.urn_id = urnId AND block_number <= $3
    UNION
    SELECT block_number FROM maker.vat_urn_art
    WHERE vat_urn_art.urn_id = urnId AND block_number <= $3
  ) inks_and_arts
  ORDER BY block_number DESC;

  FOR i IN SELECT block_number FROM updated
    LOOP
      RETURN QUERY
        SELECT * FROM maker.get_urn_state_at_block(ilk, urn, i);
    END LOOP;
END;
$$
  LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION  maker.get_urn_history_before_block(TEXT, TEXT, NUMERIC);
