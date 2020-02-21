-- +goose Up
CREATE TABLE maker.vow_vat
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vat       TEXT,
    UNIQUE (diff_id, header_id, vat)
);

CREATE INDEX vow_vat_header_id_index
    ON maker.vow_vat (header_id);

COMMENT ON TABLE maker.vow_vat
    IS E'Value of the Vow contract\'s vat variable as of a block header.';

CREATE TABLE maker.vow_flapper
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    flapper   TEXT,
    UNIQUE (diff_id, header_id, flapper)
);

CREATE INDEX vow_flapper_header_id_index
    ON maker.vow_flapper (header_id);

COMMENT ON TABLE maker.vow_flapper
    IS E'Value of the Vow contract\'s flapper variable as of a block header.';

CREATE TABLE maker.vow_flopper
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    flopper   TEXT,
    UNIQUE (diff_id, header_id, flopper)
);

CREATE INDEX vow_flopper_header_id_index
    ON maker.vow_flopper (header_id);

COMMENT ON TABLE maker.vow_flopper
    IS E'Value of the Vow contract\'s flopper variable as of a block header.';

CREATE TABLE maker.vow_sin_integer
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    sin       numeric,
    UNIQUE (diff_id, header_id, sin)
);

CREATE INDEX vow_sin_integer_header_id_index
    ON maker.vow_sin_integer (header_id);

COMMENT ON TABLE maker.vow_sin_integer
    IS E'Value of the Vow contract\'s Sin variable as of a block header.';

CREATE TABLE maker.vow_sin_mapping
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    era       numeric,
    tab       numeric,
    UNIQUE (diff_id, header_id, era, tab)
);

CREATE INDEX vow_sin_mapping_header_id_index
    ON maker.vow_sin_mapping (header_id);
CREATE INDEX vow_sin_mapping_era_index
    ON maker.vow_sin_mapping (era);

COMMENT ON TABLE maker.vow_sin_mapping
    IS E'Value of an entry in the Vow contract\'s sin mapping as of a block header.';

CREATE TABLE maker.vow_ash
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ash       numeric,
    UNIQUE (diff_id, header_id, ash)
);

CREATE INDEX vow_ash_header_id_index
    ON maker.vow_ash (header_id);

COMMENT ON TABLE maker.vow_ash
    IS E'Value of the Vow contract\'s Ash variable as of a block header.';

CREATE TABLE maker.vow_wait
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    wait      numeric,
    UNIQUE (diff_id, header_id, wait)
);

CREATE INDEX vow_wait_header_id_index
    ON maker.vow_wait (header_id);

COMMENT ON TABLE maker.vow_wait
    IS E'Value of the Vow contract\'s wait variable as of a block header.';

CREATE TABLE maker.vow_dump
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    dump      NUMERIC,
    UNIQUE (diff_id, header_id, dump)
);

CREATE INDEX vow_dump_header_id_index
    ON maker.vow_dump (header_id);

COMMENT ON TABLE maker.vow_dump
    IS E'Value of the Vow contract\'s dump variable as of a block header.';

CREATE TABLE maker.vow_sump
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    sump      numeric,
    UNIQUE (diff_id, header_id, sump)
);

CREATE INDEX vow_sump_header_id_index
    ON maker.vow_sump (header_id);

COMMENT ON TABLE maker.vow_sump
    IS E'Value of the Vow contract\'s sump variable as of a block header.';

CREATE TABLE maker.vow_bump
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    bump      numeric,
    UNIQUE (diff_id, header_id, bump)
);

CREATE INDEX vow_bump_header_id_index
    ON maker.vow_bump (header_id);

COMMENT ON TABLE maker.vow_bump
    IS E'Value of the Vow contract\'s bump variable as of a block header.';

CREATE TABLE maker.vow_hump
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    hump      numeric,
    UNIQUE (diff_id, header_id, hump)
);

CREATE INDEX vow_hump_header_id_index
    ON maker.vow_hump (header_id);

COMMENT ON TABLE maker.vow_hump
    IS E'Value of the Vow contract\'s hump variable as of a block header.';

-- +goose Down
DROP INDEX maker.vow_vat_header_id_index;
DROP INDEX maker.vow_flapper_header_id_index;
DROP INDEX maker.vow_flopper_header_id_index;
DROP INDEX maker.vow_sin_integer_header_id_index;
DROP INDEX maker.vow_sin_mapping_header_id_index;
DROP INDEX maker.vow_sin_mapping_era_index;
DROP INDEX maker.vow_ash_header_id_index;
DROP INDEX maker.vow_wait_header_id_index;
DROP INDEX maker.vow_dump_header_id_index;
DROP INDEX maker.vow_sump_header_id_index;
DROP INDEX maker.vow_bump_header_id_index;
DROP INDEX maker.vow_hump_header_id_index;

DROP TABLE maker.vow_vat;
DROP TABLE maker.vow_flapper;
DROP TABLE maker.vow_flopper;
DROP TABLE maker.vow_sin_integer;
DROP TABLE maker.vow_sin_mapping;
DROP TABLE maker.vow_ash;
DROP TABLE maker.vow_wait;
DROP TABLE maker.vow_dump;
DROP TABLE maker.vow_sump;
DROP TABLE maker.vow_bump;
DROP TABLE maker.vow_hump;
