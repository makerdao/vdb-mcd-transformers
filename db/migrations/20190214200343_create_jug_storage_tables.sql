-- +goose Up
CREATE TABLE maker.jug_ilk_rho
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    rho       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, ilk_id, rho)
);

CREATE INDEX jug_ilk_rho_header_id_index
    ON maker.jug_ilk_rho (header_id);
CREATE INDEX jug_ilk_rho_ilk_index
    ON maker.jug_ilk_rho (ilk_id);

COMMENT ON TABLE maker.jug_ilk_rho
    IS E'Value of an Ilk\'s rho field on the Jug contract as of a block header.';

CREATE TABLE maker.jug_ilk_duty
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    duty      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, ilk_id, duty)
);

CREATE INDEX jug_ilk_duty_header_id_index
    ON maker.jug_ilk_duty (header_id);
CREATE INDEX jug_ilk_duty_ilk_index
    ON maker.jug_ilk_duty (ilk_id);

COMMENT ON TABLE maker.jug_ilk_duty
    IS E'Value of an Ilk\'s duty field on the Jug contract as of a block header.';

CREATE TABLE maker.jug_vat
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vat       TEXT,
    UNIQUE (diff_id, header_id, vat)
);

CREATE INDEX jug_vat_header_id_index
    ON maker.jug_vat (header_id);

COMMENT ON TABLE maker.jug_vat
    IS E'Value of the Jug contract\'s vat variable as of a block header.';

CREATE TABLE maker.jug_vow
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vow       TEXT,
    UNIQUE (diff_id, header_id, vow)
);

CREATE INDEX jug_vow_header_id_index
    ON maker.jug_vow (header_id);

COMMENT ON TABLE maker.jug_vow
    IS E'Value of the Jug contract\'s vow variable as of a block header.';

CREATE TABLE maker.jug_base
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    base      TEXT,
    UNIQUE (diff_id, header_id, base)
);

CREATE INDEX jug_base_header_id_index
    ON maker.jug_base (header_id);

COMMENT ON TABLE maker.jug_base
    IS E'Value of the Jug contract\'s base variable as of a block header.';

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
