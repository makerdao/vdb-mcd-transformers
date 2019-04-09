-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.get_urn_history_before_block(ilk TEXT, urn TEXT, block_height BIGINT)
  RETURNS SETOF maker.urn_state AS $$
DECLARE
  blocks BIGINT[];
  i BIGINT;
  ilkId NUMERIC;
  urnId NUMERIC;
BEGIN
  SELECT id FROM maker.ilks WHERE ilks.ilk = $1 INTO ilkId;
  SELECT id FROM maker.urns WHERE urns.guy = $2 AND urns.ilk_id = ilkID INTO urnId;

  blocks := ARRAY(
    SELECT block_number
    FROM (
       SELECT block_number
       FROM maker.vat_urn_ink
       WHERE vat_urn_ink.urn_id = urnId
         AND block_number <= $3
       UNION
       SELECT block_number
       FROM maker.vat_urn_art
       WHERE vat_urn_art.urn_id = urnId
         AND block_number <= $3
     ) inks_and_arts
    ORDER BY block_number DESC
  );

  FOREACH i IN ARRAY blocks
    LOOP
      RETURN QUERY
        SELECT * FROM maker.get_urn_state_at_block(ilk, urn, i);
    END LOOP;
END;
$$
LANGUAGE plpgsql
STABLE;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION  maker.get_urn_history_before_block(TEXT, TEXT, BIGINT);
