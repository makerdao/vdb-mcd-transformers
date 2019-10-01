-- +goose Up
CREATE TABLE maker.vow_vat
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    vat          TEXT,
    UNIQUE (block_number, block_hash, vat)
);

CREATE TABLE maker.vow_flapper
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    flapper      TEXT,
    UNIQUE (block_number, block_hash, flapper)
);

CREATE TABLE maker.vow_flopper
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    flopper      TEXT,
    UNIQUE (block_number, block_hash, flopper)
);

CREATE TABLE maker.vow_sin_integer
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    sin          numeric,
    UNIQUE (block_number, block_hash, sin)
);

CREATE TABLE maker.vow_sin_mapping
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    era          numeric,
    tab          numeric,
    UNIQUE (block_number, block_hash, era, tab)
);

CREATE INDEX vow_sin_mapping_block_number_index
    ON maker.vow_sin_mapping (block_number);

CREATE INDEX vow_sin_mapping_era_index
    ON maker.vow_sin_mapping (era);

CREATE TABLE maker.vow_ash
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ash          numeric,
    UNIQUE (block_number, block_hash, ash)
);

CREATE TABLE maker.vow_wait
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    wait         numeric,
    UNIQUE (block_number, block_hash, wait)
);

CREATE TABLE maker.vow_dump
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    dump         NUMERIC,
    UNIQUE (block_number, block_hash, dump)
);

CREATE TABLE maker.vow_sump
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    sump         numeric,
    UNIQUE (block_number, block_hash, sump)
);

CREATE TABLE maker.vow_bump
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    bump         numeric,
    UNIQUE (block_number, block_hash, bump)
);

CREATE TABLE maker.vow_hump
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    hump         numeric,
    UNIQUE (block_number, block_hash, hump)
);

-- +goose Down
DROP INDEX maker.vow_sin_mapping_block_number_index;
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
