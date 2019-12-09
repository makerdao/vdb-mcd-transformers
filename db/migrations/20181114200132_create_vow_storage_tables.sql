-- +goose Up
CREATE TABLE maker.vow_vat
(
    id        SERIAL PRIMARY KEY,
    diff_id   INTEGER NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vat       TEXT,
    UNIQUE (diff_id, header_id, vat)
);

CREATE TABLE maker.vow_flapper
(
    id        SERIAL PRIMARY KEY,
    diff_id   INTEGER NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    flapper   TEXT,
    UNIQUE (diff_id, header_id, flapper)
);

CREATE TABLE maker.vow_flopper
(
    id        SERIAL PRIMARY KEY,
    diff_id   INTEGER NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    flopper   TEXT,
    UNIQUE (diff_id, header_id, flopper)
);

CREATE TABLE maker.vow_sin_integer
(
    id        SERIAL PRIMARY KEY,
    diff_id   INTEGER NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    sin       numeric,
    UNIQUE (diff_id, header_id, sin)
);

CREATE TABLE maker.vow_sin_mapping
(
    id        SERIAL PRIMARY KEY,
    diff_id   INTEGER NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    era       numeric,
    tab       numeric,
    UNIQUE (diff_id, header_id, era, tab)
);

CREATE INDEX vow_sin_mapping_era_index
    ON maker.vow_sin_mapping (era);

CREATE TABLE maker.vow_ash
(
    id        SERIAL PRIMARY KEY,
    diff_id   INTEGER NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ash       numeric,
    UNIQUE (diff_id, header_id, ash)
);

CREATE TABLE maker.vow_wait
(
    id        SERIAL PRIMARY KEY,
    diff_id   INTEGER NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    wait      numeric,
    UNIQUE (diff_id, header_id, wait)
);

CREATE TABLE maker.vow_dump
(
    id        SERIAL PRIMARY KEY,
    diff_id   INTEGER NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    dump      NUMERIC,
    UNIQUE (diff_id, header_id, dump)
);

CREATE TABLE maker.vow_sump
(
    id        SERIAL PRIMARY KEY,
    diff_id   INTEGER NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    sump      numeric,
    UNIQUE (diff_id, header_id, sump)
);

CREATE TABLE maker.vow_bump
(
    id        SERIAL PRIMARY KEY,
    diff_id   INTEGER NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    bump      numeric,
    UNIQUE (diff_id, header_id, bump)
);

CREATE TABLE maker.vow_hump
(
    id        SERIAL PRIMARY KEY,
    diff_id   INTEGER NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    hump      numeric,
    UNIQUE (diff_id, header_id, hump)
);

-- +goose Down
DROP INDEX maker.vow_sin_mapping_era_index;

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
