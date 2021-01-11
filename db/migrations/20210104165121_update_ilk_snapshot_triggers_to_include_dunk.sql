-- +goose Up
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated, dunk)
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
            diff_block_timestamp,
            ilk_dunk_before_block(new_diff.ilk_id, new_diff.header_id))
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET art = new_diff.art;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated, dunk)
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
            diff_block_timestamp,
            ilk_dunk_before_block(new_diff.ilk_id, new_diff.header_id))
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET dust = new_diff.dust;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated, dunk)
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
            diff_block_timestamp,
            ilk_dunk_before_block(new_diff.ilk_id, new_diff.header_id))
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET line = new_diff.line;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated, dunk)
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
            diff_block_timestamp,
            ilk_dunk_before_block(new_diff.ilk_id, new_diff.header_id))
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET rate = new_diff.rate;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated, dunk)
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
            diff_block_timestamp,
            ilk_dunk_before_block(new_diff.ilk_id, new_diff.header_id))
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET spot = new_diff.spot;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated, dunk)
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
            diff_block_timestamp,
            ilk_dunk_before_block(new_diff.ilk_id, new_diff.header_id))
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET chop = new_diff.chop;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated, dunk)
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
            diff_block_timestamp,
            ilk_dunk_before_block(new_diff.ilk_id, new_diff.header_id))
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET lump = new_diff.lump;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated, dunk)
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
            diff_block_timestamp,
            ilk_dunk_before_block(new_diff.ilk_id, new_diff.header_id))
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET flip = new_diff.flip;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated, dunk)
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
            diff_block_timestamp,
            ilk_dunk_before_block(new_diff.ilk_id, new_diff.header_id))
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET rho = new_diff.rho;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated, dunk)
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
            diff_block_timestamp,
            ilk_dunk_before_block(new_diff.ilk_id, new_diff.header_id))
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET duty = new_diff.duty;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated, dunk)
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
            diff_block_timestamp,
            ilk_dunk_before_block(new_diff.ilk_id, new_diff.header_id))
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET pip = new_diff.pip;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
                           duty, pip, mat, created, updated, dunk)
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
            diff_block_timestamp,
            ilk_dunk_before_block(new_diff.ilk_id, new_diff.header_id))
    ON CONFLICT (ilk_identifier, block_number)
        DO UPDATE SET mat = new_diff.mat;
    RETURN new_diff;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd


-- +goose Down
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
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
    INTO api.ilk_snapshot (ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho,
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
