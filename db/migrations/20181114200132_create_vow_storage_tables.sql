-- +goose Up
CREATE TABLE maker.vow_vat (
  id           SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash   TEXT,
  vat          TEXT,
  UNIQUE (block_number, block_hash, vat)
);

CREATE TABLE maker.vow_cow (
  id           SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash   TEXT,
  cow          TEXT,
  UNIQUE (block_number, block_hash, cow)
);

CREATE TABLE maker.vow_row (
  id           SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash   TEXT,
  row          TEXT,
  UNIQUE (block_number, block_hash, row)
);

CREATE TABLE maker.vow_sin_integer (
  id           SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash   TEXT,
  sin          numeric,
  UNIQUE (block_number, block_hash, sin)
);

CREATE TABLE maker.vow_sin_mapping (
  id           SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash   TEXT,
  era          numeric,
  tab          numeric,
  UNIQUE (block_number, block_hash, era, tab)
);

CREATE TABLE maker.vow_ash (
  id           SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash   TEXT,
  ash          numeric,
  UNIQUE (block_number, block_hash, ash)
);

CREATE TABLE maker.vow_wait (
  id           SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash   TEXT,
  wait         numeric,
  UNIQUE (block_number, block_hash, wait)
);

CREATE TABLE maker.vow_sump (
  id           SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash   TEXT,
  sump         numeric,
  UNIQUE (block_number, block_hash, sump)
);

CREATE TABLE maker.vow_bump (
  id           SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash   TEXT,
  bump         numeric,
  UNIQUE (block_number, block_hash, bump)
);

CREATE TABLE maker.vow_hump (
  id           SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash   TEXT,
  hump         numeric,
  UNIQUE (block_number, block_hash, hump)
);

-- +goose Down
DROP TABLE maker.vow_vat;
DROP TABLE maker.vow_cow;
DROP TABLE maker.vow_row;
DROP TABLE maker.vow_sin_integer;
DROP TABLE maker.vow_sin_mapping;
DROP TABLE maker.vow_ash;
DROP TABLE maker.vow_wait;
DROP TABLE maker.vow_sump;
DROP TABLE maker.vow_bump;
DROP TABLE maker.vow_hump;
