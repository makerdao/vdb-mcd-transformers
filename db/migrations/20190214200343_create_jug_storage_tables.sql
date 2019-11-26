-- +goose Up
CREATE TABLE maker.jug_ilk_rho
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    rho       NUMERIC NOT NULL,
    UNIQUE (header_id, ilk_id, rho)
);

CREATE INDEX jug_ilk_rho_header_id_index
    ON maker.jug_ilk_rho (header_id);
CREATE INDEX jug_ilk_rho_ilk_index
    ON maker.jug_ilk_rho (ilk_id);

CREATE TABLE maker.jug_ilk_duty
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    duty      NUMERIC NOT NULL,
    UNIQUE (header_id, ilk_id, duty)
);

CREATE INDEX jug_ilk_duty_header_id_index
    ON maker.jug_ilk_duty (header_id);
CREATE INDEX jug_ilk_duty_ilk_index
    ON maker.jug_ilk_duty (ilk_id);

CREATE TABLE maker.jug_vat
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vat       TEXT,
    UNIQUE (header_id, vat)
);

CREATE TABLE maker.jug_vow
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vow       TEXT,
    UNIQUE (header_id, vow)
);

CREATE TABLE maker.jug_base
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    base      TEXT,
    UNIQUE (header_id, base)
);

-- +goose Down
DROP INDEX maker.jug_ilk_rho_header_id_index;
DROP INDEX maker.jug_ilk_rho_ilk_index;
DROP INDEX maker.jug_ilk_duty_header_id_index;
DROP INDEX maker.jug_ilk_duty_ilk_index;

DROP TABLE maker.jug_ilk_rho;
DROP TABLE maker.jug_ilk_duty;
DROP TABLE maker.jug_vat;
DROP TABLE maker.jug_vow;
DROP TABLE maker.jug_base;
