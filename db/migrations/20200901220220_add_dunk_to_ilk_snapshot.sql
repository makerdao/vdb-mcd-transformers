-- +goose Up
ALTER TABLE api.ilk_snapshot
    ADD COLUMN dunk NUMERIC DEFAULT NULL;

CREATE FUNCTION ilk_dunk_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT dunk
FROM maker.cat_ilk_dunk
         LEFT JOIN public.headers ON cat_ilk_dunk.header_id = headers.id
WHERE cat_ilk_dunk.ilk_id = ilk_dunk_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_dunk_before_block(ilk_id INTEGER, header_id INTEGER)
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_dunk(new_diff maker.cat_ilk_dunk) RETURNS maker.cat_ilk_dunk
AS
$$
DECLARE
    diff_ilk_identifier  TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = new_diff.ilk_id);
    diff_block_timestamp TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE headers.id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE headers.id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, dunk, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.dunk,
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET dunk = new_diff.dunk;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_dunk
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_dunks_until_next_diff(start_at_diff maker.cat_ilk_dunk, new_dunk NUMERIC) RETURNS maker.cat_ilk_dunk
AS
$$
DECLARE
    diff_ilk_identifier  TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number    BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_dunk_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.cat_ilk_dunk
                 LEFT JOIN public.headers ON cat_ilk_dunk.header_id = headers.id
        WHERE cat_ilk_dunk.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.ilk_snapshot
    SET dunk = new_dunk
    WHERE ilk_snapshot.ilk_identifier = diff_ilk_identifier
      AND ilk_snapshot.block_number >= diff_block_number
      AND (next_dunk_diff_block IS NULL
        OR ilk_snapshot.block_number < next_dunk_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_dunks_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_dunks() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_dunk(NEW);
        PERFORM maker.update_dunks_until_next_diff(NEW, NEW.dunk);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_dunks_until_next_diff(OLD, ilk_dunk_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_snapshot(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_dunk
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.cat_ilk_dunk
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_dunks();


-- +goose Down
ALTER TABLE api.ilk_snapshot
    DROP COLUMN dunk;

DROP TRIGGER ilk_dunk ON maker.cat_ilk_dunk;
DROP FUNCTION maker.update_ilk_dunks();
DROP FUNCTION maker.update_dunks_until_next_diff(maker.cat_ilk_dunk, NUMERIC);
DROP FUNCTION maker.insert_new_dunk(maker.cat_ilk_dunk);
DROP FUNCTION ilk_dunk_before_block(INTEGER, INTEGER);
