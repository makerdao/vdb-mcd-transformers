-- +goose Up
CREATE TABLE api.historical_ilk_state
(
    ilk_identifier TEXT,
    block_number   BIGINT,
    rate           NUMERIC   DEFAULT NULL,
    art            NUMERIC   DEFAULT NULL,
    spot           NUMERIC   DEFAULT NULL,
    line           NUMERIC   DEFAULT NULL,
    dust           NUMERIC   DEFAULT NULL,
    chop           NUMERIC   DEFAULT NULL,
    lump           NUMERIC   DEFAULT NULL,
    flip           TEXT      DEFAULT NULL,
    rho            NUMERIC   DEFAULT NULL,
    duty           NUMERIC   DEFAULT NULL,
    pip            TEXT      DEFAULT NULL,
    mat            NUMERIC   DEFAULT NULL,
    created        TIMESTAMP DEFAULT NULL,
    updated        TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (ilk_identifier, block_number)
);

COMMENT ON COLUMN api.historical_ilk_state.ilk_identifier IS '@name id';


CREATE FUNCTION ilk_rate_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT rate
FROM maker.vat_ilk_rate
         LEFT JOIN public.headers ON vat_ilk_rate.header_id = headers.id
WHERE vat_ilk_rate.ilk_id = ilk_rate_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_rate_before_block
    IS E'@omit';


CREATE FUNCTION ilk_art_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT art
FROM maker.vat_ilk_art
         LEFT JOIN public.headers ON vat_ilk_art.header_id = headers.id
WHERE vat_ilk_art.ilk_id = ilk_art_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_art_before_block
    IS E'@omit';


CREATE FUNCTION ilk_spot_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT spot
FROM maker.vat_ilk_spot
         LEFT JOIN public.headers ON vat_ilk_spot.header_id = headers.id
WHERE vat_ilk_spot.ilk_id = ilk_spot_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_spot_before_block
    IS E'@omit';


CREATE FUNCTION ilk_line_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT line
FROM maker.vat_ilk_line
         LEFT JOIN public.headers ON vat_ilk_line.header_id = headers.id
WHERE vat_ilk_line.ilk_id = ilk_line_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_line_before_block
    IS E'@omit';


CREATE FUNCTION ilk_dust_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT dust
FROM maker.vat_ilk_dust
         LEFT JOIN public.headers ON vat_ilk_dust.header_id = headers.id
WHERE vat_ilk_dust.ilk_id = ilk_dust_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_dust_before_block
    IS E'@omit';


CREATE FUNCTION ilk_chop_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT chop
FROM maker.cat_ilk_chop
         LEFT JOIN public.headers ON cat_ilk_chop.header_id = headers.id
WHERE cat_ilk_chop.ilk_id = ilk_chop_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_chop_before_block
    IS E'@omit';


CREATE FUNCTION ilk_lump_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT lump
FROM maker.cat_ilk_lump
         LEFT JOIN public.headers ON cat_ilk_lump.header_id = headers.id
WHERE cat_ilk_lump.ilk_id = ilk_lump_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_lump_before_block
    IS E'@omit';


CREATE FUNCTION ilk_flip_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS TEXT AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT flip
FROM maker.cat_ilk_flip
         LEFT JOIN public.headers ON cat_ilk_flip.header_id = headers.id
WHERE cat_ilk_flip.ilk_id = ilk_flip_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_flip_before_block
    IS E'@omit';


CREATE FUNCTION ilk_rho_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$

WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT rho
FROM maker.jug_ilk_rho
         LEFT JOIN public.headers ON jug_ilk_rho.header_id = headers.id
WHERE jug_ilk_rho.ilk_id = ilk_rho_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_rho_before_block
    IS E'@omit';


CREATE FUNCTION ilk_duty_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT duty
FROM maker.jug_ilk_duty
         LEFT JOIN public.headers ON jug_ilk_duty.header_id = headers.id
WHERE jug_ilk_duty.ilk_id = ilk_duty_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_duty_before_block
    IS E'@omit';


CREATE FUNCTION ilk_pip_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS TEXT AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT pip
FROM maker.spot_ilk_pip
         LEFT JOIN public.headers ON spot_ilk_pip.header_id = headers.id
WHERE spot_ilk_pip.ilk_id = ilk_pip_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_pip_before_block
    IS E'@omit';


CREATE FUNCTION ilk_mat_before_block(ilk_id INTEGER, header_id INTEGER) RETURNS NUMERIC AS
$$
WITH passed_block_number AS (
    SELECT block_number
    FROM public.headers
    WHERE id = header_id
)

SELECT mat
FROM maker.spot_ilk_mat
         LEFT JOIN public.headers ON spot_ilk_mat.header_id = headers.id
WHERE spot_ilk_mat.ilk_id = ilk_mat_before_block.ilk_id
  AND headers.block_number < (SELECT block_number FROM passed_block_number)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_mat_before_block
    IS E'@omit';


CREATE FUNCTION ilk_time_created(ilk_id INTEGER) RETURNS TIMESTAMP AS
$$
SELECT api.epoch_to_datetime(MIN(block_timestamp))
FROM public.headers
         LEFT JOIN maker.vat_init ON vat_init.header_id = headers.id
WHERE vat_init.ilk_id = ilk_time_created.ilk_id
$$
    LANGUAGE sql;

COMMENT ON FUNCTION ilk_time_created
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.delete_redundant_ilk_state(ilk_id INTEGER, header_id INTEGER) RETURNS api.historical_ilk_state
AS
$$
DECLARE
    associated_ilk         TEXT                     := (
        SELECT identifier
        FROM maker.ilks
        WHERE id = delete_redundant_ilk_state.ilk_id);
    ilk_state_block_number BIGINT                   := (
        SELECT block_number
        FROM public.headers
        WHERE id = header_id);
    ilk_state              api.historical_ilk_state := (
        SELECT (ilk_identifier, historical_ilk_state.block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                duty, pip, mat, created, updated)
        FROM api.historical_ilk_state
        WHERE ilk_identifier = associated_ilk
          AND historical_ilk_state.block_number = ilk_state_block_number);
    prev_ilk_state         api.historical_ilk_state := (
        SELECT (ilk_identifier, historical_ilk_state.block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                duty, pip, mat, created, updated)
        FROM api.historical_ilk_state
        WHERE ilk_identifier = associated_ilk
          AND historical_ilk_state.block_number < ilk_state_block_number
        ORDER BY historical_ilk_state.block_number DESC
        LIMIT 1);
BEGIN
    DELETE
    FROM api.historical_ilk_state
    WHERE historical_ilk_state.ilk_identifier = associated_ilk
      AND historical_ilk_state.block_number = ilk_state_block_number
      AND (ilk_state.rate IS NULL OR ilk_state.rate = prev_ilk_state.rate)
      AND (ilk_state.art IS NULL OR ilk_state.art = prev_ilk_state.art)
      AND (ilk_state.spot IS NULL OR ilk_state.spot = prev_ilk_state.spot)
      AND (ilk_state.line IS NULL OR ilk_state.line = prev_ilk_state.line)
      AND (ilk_state.dust IS NULL OR ilk_state.dust = prev_ilk_state.dust)
      AND (ilk_state.chop IS NULL OR ilk_state.chop = prev_ilk_state.chop)
      AND (ilk_state.lump IS NULL OR ilk_state.lump = prev_ilk_state.lump)
      AND (ilk_state.flip IS NULL OR ilk_state.flip = prev_ilk_state.flip)
      AND (ilk_state.rho IS NULL OR ilk_state.rho = prev_ilk_state.rho)
      AND (ilk_state.duty IS NULL OR ilk_state.duty = prev_ilk_state.duty)
      AND (ilk_state.pip IS NULL OR ilk_state.pip = prev_ilk_state.pip)
      AND (ilk_state.mat IS NULL OR ilk_state.mat = prev_ilk_state.mat);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.delete_redundant_ilk_state
    IS E'@omit';


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_rate(new_diff maker.vat_ilk_rate) RETURNS maker.vat_ilk_rate
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
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.historical_ilk_state (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                                   duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            new_diff.rate,
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
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET rate = new_diff.rate;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_rate
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_rates_until_next_diff(start_at_diff maker.vat_ilk_rate, new_rate NUMERIC) RETURNS maker.vat_ilk_rate
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
    next_rate_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_ilk_rate
                 LEFT JOIN public.headers ON vat_ilk_rate.header_id = headers.id
        WHERE vat_ilk_rate.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.historical_ilk_state
    SET rate = new_rate
    WHERE historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.block_number >= diff_block_number
      AND (next_rate_diff_block IS NULL
        OR historical_ilk_state.block_number < next_rate_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_rates_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_rates() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_rate(NEW);
        PERFORM maker.update_rates_until_next_diff(NEW, NEW.rate);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_rates_until_next_diff(OLD, ilk_rate_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_state(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_rate
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.vat_ilk_rate
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_rates();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_art(new_diff maker.vat_ilk_art) RETURNS maker.vat_ilk_art
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
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.historical_ilk_state (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                                   duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.art,
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
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET art = new_diff.art;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_art
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_arts_until_next_diff(start_at_diff maker.vat_ilk_art, new_art NUMERIC) RETURNS maker.vat_ilk_art
AS
$$
DECLARE
    diff_ilk_identifier TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_art_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_ilk_art
                 LEFT JOIN public.headers ON vat_ilk_art.header_id = headers.id
        WHERE vat_ilk_art.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.historical_ilk_state
    SET art = new_art
    WHERE historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.block_number >= diff_block_number
      AND (next_art_diff_block IS NULL
        OR historical_ilk_state.block_number < next_art_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_arts_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_arts() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_art(NEW);
        PERFORM maker.update_arts_until_next_diff(NEW, NEW.art);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_arts_until_next_diff(OLD, ilk_art_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_state(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_art
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.vat_ilk_art
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_arts();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_spot(new_diff maker.vat_ilk_spot) RETURNS maker.vat_ilk_spot
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
    INTO api.historical_ilk_state (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                                   duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.spot,
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET spot = new_diff.spot;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_spot
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_spots_until_next_diff(start_at_diff maker.vat_ilk_spot, new_spot NUMERIC) RETURNS maker.vat_ilk_spot
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
    next_spot_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_ilk_spot
                 LEFT JOIN public.headers ON vat_ilk_spot.header_id = headers.id
        WHERE vat_ilk_spot.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.historical_ilk_state
    SET spot = new_spot
    WHERE historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.block_number >= diff_block_number
      AND (next_spot_diff_block IS NULL
        OR historical_ilk_state.block_number < next_spot_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_spots_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_spots() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_spot(NEW);
        PERFORM maker.update_spots_until_next_diff(NEW, NEW.spot);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_spots_until_next_diff(OLD, ilk_spot_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_state(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_spot
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.vat_ilk_spot
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_spots();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_line(new_diff maker.vat_ilk_line) RETURNS maker.vat_ilk_line
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
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.historical_ilk_state (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                                   duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.line,
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET line = new_diff.line;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_line
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_lines_until_next_diff(start_at_diff maker.vat_ilk_line, new_line NUMERIC) RETURNS maker.vat_ilk_line
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
    next_line_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_ilk_line
                 LEFT JOIN public.headers ON vat_ilk_line.header_id = headers.id
        WHERE vat_ilk_line.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.historical_ilk_state
    SET line = new_line
    WHERE historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.block_number >= diff_block_number
      AND (next_line_diff_block IS NULL
        OR historical_ilk_state.block_number < next_line_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_lines_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_lines() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_line(NEW);
        PERFORM maker.update_lines_until_next_diff(NEW, NEW.line);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_lines_until_next_diff(OLD, ilk_line_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_state(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_line
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.vat_ilk_line
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_lines();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_dust(new_diff maker.vat_ilk_dust) RETURNS maker.vat_ilk_dust
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
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.historical_ilk_state (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                                   duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.dust,
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET dust = new_diff.dust;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_dust
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_dusts_until_next_diff(start_at_diff maker.vat_ilk_dust, new_dust NUMERIC) RETURNS maker.vat_ilk_dust
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
    next_dust_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.vat_ilk_dust
                 LEFT JOIN public.headers ON vat_ilk_dust.header_id = headers.id
        WHERE vat_ilk_dust.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.historical_ilk_state
    SET dust = new_dust
    WHERE historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.block_number >= diff_block_number
      AND (next_dust_diff_block IS NULL
        OR historical_ilk_state.block_number < next_dust_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_dusts_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_dusts() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_dust(NEW);
        PERFORM maker.update_dusts_until_next_diff(NEW, NEW.dust);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_dusts_until_next_diff(OLD, ilk_dust_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_state(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_dust
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.vat_ilk_dust
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_dusts();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_chop(new_diff maker.cat_ilk_chop) RETURNS maker.cat_ilk_chop
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
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.historical_ilk_state (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                                   duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.chop,
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET chop = new_diff.chop;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_chop
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_chops_until_next_diff(start_at_diff maker.cat_ilk_chop, new_chop NUMERIC) RETURNS maker.cat_ilk_chop
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
    next_chop_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.cat_ilk_chop
                 LEFT JOIN public.headers ON cat_ilk_chop.header_id = headers.id
        WHERE cat_ilk_chop.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.historical_ilk_state
    SET chop = new_chop
    WHERE historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.block_number >= diff_block_number
      AND (next_chop_diff_block IS NULL
        OR historical_ilk_state.block_number < next_chop_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_chops_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_chops() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_chop(NEW);
        PERFORM maker.update_chops_until_next_diff(NEW, NEW.chop);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_chops_until_next_diff(OLD, ilk_chop_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_state(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_chop
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.cat_ilk_chop
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_chops();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_lump(new_diff maker.cat_ilk_lump) RETURNS maker.cat_ilk_lump
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
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.historical_ilk_state (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                                   duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.lump,
            ilk_flip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET lump = new_diff.lump;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_lump
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_lumps_until_next_diff(start_at_diff maker.cat_ilk_lump, new_lump NUMERIC) RETURNS maker.cat_ilk_lump
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
    next_lump_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.cat_ilk_lump
                 LEFT JOIN public.headers ON cat_ilk_lump.header_id = headers.id
        WHERE cat_ilk_lump.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.historical_ilk_state
    SET lump = new_lump
    WHERE historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.block_number >= diff_block_number
      AND (next_lump_diff_block IS NULL
        OR historical_ilk_state.block_number < next_lump_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_lumps_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_lumps() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_lump(NEW);
        PERFORM maker.update_lumps_until_next_diff(NEW, NEW.lump);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_lumps_until_next_diff(OLD, ilk_lump_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_state(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_lump
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.cat_ilk_lump
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_lumps();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_flip(new_diff maker.cat_ilk_flip) RETURNS maker.cat_ilk_flip
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
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.historical_ilk_state (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                                   duty, pip, mat, created, updated)
    VALUES (diff_ilk_identifier,
            diff_block_number,
            ilk_rate_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_art_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_spot_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_line_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_dust_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_chop_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_lump_before_block(new_diff.ilk_id, new_diff.header_id),
            new_diff.flip,
            ilk_rho_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET flip = new_diff.flip;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_flip
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_flips_until_next_diff(start_at_diff maker.cat_ilk_flip, new_flip TEXT) RETURNS maker.cat_ilk_flip
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
    next_flip_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.cat_ilk_flip
                 LEFT JOIN public.headers ON cat_ilk_flip.header_id = headers.id
        WHERE cat_ilk_flip.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.historical_ilk_state
    SET flip = new_flip
    WHERE historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.block_number >= diff_block_number
      AND (next_flip_diff_block IS NULL
        OR historical_ilk_state.block_number < next_flip_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_flips_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_flips() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_flip(NEW);
        PERFORM maker.update_flips_until_next_diff(NEW, NEW.flip);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_flips_until_next_diff(OLD, ilk_flip_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_state(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_flip
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.cat_ilk_flip
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_flips();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_rho(new_diff maker.jug_ilk_rho) RETURNS maker.jug_ilk_rho
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
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.historical_ilk_state (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                                   duty, pip, mat, created, updated)
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
            new_diff.rho,
            ilk_duty_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET rho = new_diff.rho;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_rho
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_rhos_until_next_diff(start_at_diff maker.jug_ilk_rho, new_rho NUMERIC) RETURNS maker.jug_ilk_rho
AS
$$
DECLARE
    diff_ilk_identifier TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_rho_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.jug_ilk_rho
                 LEFT JOIN public.headers ON jug_ilk_rho.header_id = headers.id
        WHERE jug_ilk_rho.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.historical_ilk_state
    SET rho = new_rho
    WHERE historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.block_number >= diff_block_number
      AND (next_rho_diff_block IS NULL
        OR historical_ilk_state.block_number < next_rho_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_rhos_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_rhos() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_rho(NEW);
        PERFORM maker.update_rhos_until_next_diff(NEW, NEW.rho);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_rhos_until_next_diff(OLD, ilk_rho_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_state(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_rho
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.jug_ilk_rho
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_rhos();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_duty(new_diff maker.jug_ilk_duty) RETURNS maker.jug_ilk_duty
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
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.historical_ilk_state (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                                   duty, pip, mat, created, updated)
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
            new_diff.duty,
            ilk_pip_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET duty = new_diff.duty;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_duty
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_duties_until_next_diff(start_at_diff maker.jug_ilk_duty, new_duty NUMERIC) RETURNS maker.jug_ilk_duty
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
    next_duty_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.jug_ilk_duty
                 LEFT JOIN public.headers ON jug_ilk_duty.header_id = headers.id
        WHERE jug_ilk_duty.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.historical_ilk_state
    SET duty = new_duty
    WHERE historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.block_number >= diff_block_number
      AND (next_duty_diff_block IS NULL
        OR historical_ilk_state.block_number < next_duty_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_duties_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_duties() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_duty(NEW);
        PERFORM maker.update_duties_until_next_diff(NEW, NEW.duty);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_duties_until_next_diff(OLD, ilk_duty_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_state(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_duty
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.jug_ilk_duty
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_duties();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_pip(new_diff maker.spot_ilk_pip) RETURNS maker.spot_ilk_pip
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
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.historical_ilk_state (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                                   duty, pip, mat, created, updated)
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
            new_diff.pip,
            ilk_mat_before_block(new_diff.ilk_id, new_diff.header_id),
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET pip = new_diff.pip;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_pip
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_pips_until_next_diff(start_at_diff maker.spot_ilk_pip, new_pip TEXT) RETURNS maker.spot_ilk_pip
AS
$$
DECLARE
    diff_ilk_identifier TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_pip_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.spot_ilk_pip
                 LEFT JOIN public.headers ON spot_ilk_pip.header_id = headers.id
        WHERE spot_ilk_pip.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.historical_ilk_state
    SET pip = new_pip
    WHERE historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.block_number >= diff_block_number
      AND (next_pip_diff_block IS NULL
        OR historical_ilk_state.block_number < next_pip_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_pips_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_pips() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_pip(NEW);
        PERFORM maker.update_pips_until_next_diff(NEW, NEW.pip);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_pips_until_next_diff(OLD, ilk_pip_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_state(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_pip
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.spot_ilk_pip
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_pips();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_mat(new_diff maker.spot_ilk_mat) RETURNS maker.spot_ilk_mat
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
        WHERE id = new_diff.header_id);
    diff_block_number    NUMERIC   := (
        SELECT block_number
        FROM public.headers
        WHERE id = new_diff.header_id);
BEGIN
    INSERT
    INTO api.historical_ilk_state (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                                   duty, pip, mat, created, updated)
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
            new_diff.mat,
            ilk_time_created(new_diff.ilk_id),
            diff_block_timestamp)
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET mat = new_diff.mat;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_mat
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_mats_until_next_diff(start_at_diff maker.spot_ilk_mat, new_mat NUMERIC) RETURNS maker.spot_ilk_mat
AS
$$
DECLARE
    diff_ilk_identifier TEXT   := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = start_at_diff.ilk_id);
    diff_block_number   BIGINT := (
        SELECT block_number
        FROM public.headers
        WHERE id = start_at_diff.header_id);
    next_mat_diff_block BIGINT := (
        SELECT MIN(block_number)
        FROM maker.spot_ilk_mat
                 LEFT JOIN public.headers ON spot_ilk_mat.header_id = headers.id
        WHERE spot_ilk_mat.ilk_id = start_at_diff.ilk_id
          AND block_number > diff_block_number);
BEGIN
    UPDATE api.historical_ilk_state
    SET mat = new_mat
    WHERE historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.block_number >= diff_block_number
      AND (next_mat_diff_block IS NULL
        OR historical_ilk_state.block_number < next_mat_diff_block);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_mats_until_next_diff
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_ilk_mats() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_mat(NEW);
        PERFORM maker.update_mats_until_next_diff(NEW, NEW.mat);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.update_mats_until_next_diff(OLD, ilk_mat_before_block(OLD.ilk_id, OLD.header_id));
        PERFORM maker.delete_redundant_ilk_state(OLD.ilk_id, OLD.header_id);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER ilk_mat
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.spot_ilk_mat
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_ilk_mats();


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.insert_new_time_created(new_event maker.vat_init) RETURNS maker.vat_init
AS
$$
DECLARE
    diff_ilk_identifier TEXT      := (
        SELECT identifier
        FROM maker.ilks
        WHERE ilks.id = new_event.ilk_id);
    diff_timestamp      TIMESTAMP := (
        SELECT api.epoch_to_datetime(block_timestamp)
        FROM public.headers
        WHERE headers.id = new_event.header_id);
BEGIN
    UPDATE api.historical_ilk_state
    SET created = diff_timestamp
    FROM public.headers
    WHERE headers.block_number = historical_ilk_state.block_number
      AND historical_ilk_state.ilk_identifier = diff_ilk_identifier
      AND historical_ilk_state.created IS NULL;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.insert_new_time_created
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.clear_time_created(old_event maker.vat_init) RETURNS maker.vat_init
AS
$$
BEGIN
    UPDATE api.historical_ilk_state
    SET created = ilk_time_created(old_event.ilk_id)
    FROM maker.ilks
    WHERE ilks.identifier = historical_ilk_state.ilk_identifier
      AND ilks.id = old_event.ilk_id;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.clear_time_created
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_time_created() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP IN ('INSERT', 'UPDATE')) THEN
        PERFORM maker.insert_new_time_created(NEW);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.clear_time_created(OLD);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_time_created
    IS E'@omit';

CREATE TRIGGER ilk_init
    AFTER INSERT OR UPDATE OR DELETE
    ON maker.vat_init
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_time_created();


-- +goose Down
DROP TRIGGER ilk_init ON maker.vat_init;
DROP TRIGGER ilk_mat ON maker.spot_ilk_mat;
DROP TRIGGER ilk_pip ON maker.spot_ilk_pip;
DROP TRIGGER ilk_duty ON maker.jug_ilk_duty;
DROP TRIGGER ilk_rho ON maker.jug_ilk_rho;
DROP TRIGGER ilk_flip ON maker.cat_ilk_flip;
DROP TRIGGER ilk_lump ON maker.cat_ilk_lump;
DROP TRIGGER ilk_chop ON maker.cat_ilk_chop;
DROP TRIGGER ilk_dust ON maker.vat_ilk_dust;
DROP TRIGGER ilk_line ON maker.vat_ilk_line;
DROP TRIGGER ilk_spot ON maker.vat_ilk_spot;
DROP TRIGGER ilk_art ON maker.vat_ilk_art;
DROP TRIGGER ilk_rate ON maker.vat_ilk_rate;

DROP FUNCTION maker.update_time_created();
DROP FUNCTION maker.update_ilk_mats();
DROP FUNCTION maker.update_ilk_pips();
DROP FUNCTION maker.update_ilk_duties();
DROP FUNCTION maker.update_ilk_rhos();
DROP FUNCTION maker.update_ilk_flips();
DROP FUNCTION maker.update_ilk_lumps();
DROP FUNCTION maker.update_ilk_chops();
DROP FUNCTION maker.update_ilk_dusts();
DROP FUNCTION maker.update_ilk_lines();
DROP FUNCTION maker.update_ilk_spots();
DROP FUNCTION maker.update_ilk_arts();
DROP FUNCTION maker.update_ilk_rates();

DROP FUNCTION maker.clear_time_created(maker.vat_init);
DROP FUNCTION maker.update_mats_until_next_diff(maker.spot_ilk_mat, NUMERIC);
DROP FUNCTION maker.update_pips_until_next_diff(maker.spot_ilk_pip, TEXT);
DROP FUNCTION maker.update_duties_until_next_diff(maker.jug_ilk_duty, NUMERIC);
DROP FUNCTION maker.update_rhos_until_next_diff(maker.jug_ilk_rho, NUMERIC);
DROP FUNCTION maker.update_flips_until_next_diff(maker.cat_ilk_flip, TEXT);
DROP FUNCTION maker.update_lumps_until_next_diff(maker.cat_ilk_lump, NUMERIC);
DROP FUNCTION maker.update_chops_until_next_diff(maker.cat_ilk_chop, NUMERIC);
DROP FUNCTION maker.update_dusts_until_next_diff(maker.vat_ilk_dust, NUMERIC);
DROP FUNCTION maker.update_lines_until_next_diff(maker.vat_ilk_line, NUMERIC);
DROP FUNCTION maker.update_spots_until_next_diff(maker.vat_ilk_spot, NUMERIC);
DROP FUNCTION maker.update_arts_until_next_diff(maker.vat_ilk_art, NUMERIC);
DROP FUNCTION maker.update_rates_until_next_diff(maker.vat_ilk_rate, NUMERIC);

DROP FUNCTION maker.insert_new_time_created(maker.vat_init);
DROP FUNCTION maker.insert_new_mat(maker.spot_ilk_mat);
DROP FUNCTION maker.insert_new_pip(maker.spot_ilk_pip);
DROP FUNCTION maker.insert_new_duty(maker.jug_ilk_duty);
DROP FUNCTION maker.insert_new_rho(maker.jug_ilk_rho);
DROP FUNCTION maker.insert_new_flip(maker.cat_ilk_flip);
DROP FUNCTION maker.insert_new_lump(maker.cat_ilk_lump);
DROP FUNCTION maker.insert_new_chop(maker.cat_ilk_chop);
DROP FUNCTION maker.insert_new_dust(maker.vat_ilk_dust);
DROP FUNCTION maker.insert_new_line(maker.vat_ilk_line);
DROP FUNCTION maker.insert_new_spot(maker.vat_ilk_spot);
DROP FUNCTION maker.insert_new_art(maker.vat_ilk_art);
DROP FUNCTION maker.insert_new_rate(maker.vat_ilk_rate);

DROP FUNCTION maker.delete_redundant_ilk_state(INTEGER, INTEGER);

DROP FUNCTION ilk_time_created(INTEGER);
DROP FUNCTION ilk_mat_before_block(INTEGER, INTEGER);
DROP FUNCTION ilk_pip_before_block(INTEGER, INTEGER);
DROP FUNCTION ilk_duty_before_block(INTEGER, INTEGER);
DROP FUNCTION ilk_rho_before_block(INTEGER, INTEGER);
DROP FUNCTION ilk_flip_before_block(INTEGER, INTEGER);
DROP FUNCTION ilk_lump_before_block(INTEGER, INTEGER);
DROP FUNCTION ilk_chop_before_block(INTEGER, INTEGER);
DROP FUNCTION ilk_dust_before_block(INTEGER, INTEGER);
DROP FUNCTION ilk_line_before_block(INTEGER, INTEGER);
DROP FUNCTION ilk_spot_before_block(INTEGER, INTEGER);
DROP FUNCTION ilk_art_before_block(INTEGER, INTEGER);
DROP FUNCTION ilk_rate_before_block(INTEGER, INTEGER);

DROP TABLE api.historical_ilk_state;
