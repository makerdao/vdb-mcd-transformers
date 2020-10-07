-- +goose Up
CREATE TABLE maker.jug_ilk_rho
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    rho       NUMERIC NOT NULL,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, ilk_id, rho)
);

CREATE INDEX jug_ilk_rho_header_id_index
    ON maker.jug_ilk_rho (header_id);
CREATE INDEX jug_ilk_rho_ilk_index
    ON maker.jug_ilk_rho (ilk_id);

CREATE TABLE maker.jug_ilk_duty
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    duty      NUMERIC NOT NULL,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, ilk_id, duty)
);

CREATE INDEX jug_ilk_duty_header_id_index
    ON maker.jug_ilk_duty (header_id);
CREATE INDEX jug_ilk_duty_ilk_index
    ON maker.jug_ilk_duty (ilk_id);

CREATE TABLE maker.jug_vat
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    vat       TEXT,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, vat)
);

CREATE INDEX jug_vat_header_id_index
    ON maker.jug_vat (header_id);

CREATE TABLE maker.jug_vow
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    vow       TEXT,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, vow)
);

CREATE INDEX jug_vow_header_id_index
    ON maker.jug_vow (header_id);

CREATE TABLE maker.jug_base
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    base      TEXT,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, base)
);

CREATE INDEX jug_base_header_id_index
    ON maker.jug_base (header_id);

-- +goose Down
DROP INDEX maker.jug_ilk_rho_header_id_index;
DROP INDEX maker.jug_ilk_rho_ilk_index;
DROP INDEX maker.jug_ilk_duty_header_id_index;
DROP INDEX maker.jug_ilk_duty_ilk_index;
DROP INDEX maker.jug_vat_header_id_index;
DROP INDEX maker.jug_vow_header_id_index;
DROP INDEX maker.jug_base_header_id_index;

DROP TABLE maker.jug_ilk_rho;
DROP TABLE maker.jug_ilk_duty;
DROP TABLE maker.jug_vat;
DROP TABLE maker.jug_vow;
DROP TABLE maker.jug_base;
