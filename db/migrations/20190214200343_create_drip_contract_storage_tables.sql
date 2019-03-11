-- +goose Up
CREATE TABLE maker.drip_ilk_rho (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk           INTEGER NOT NULL REFERENCES maker.ilks (id), 
  rho           NUMERIC NOT NULL
);

CREATE TABLE maker.drip_ilk_tax (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk           INTEGER NOT NULL REFERENCES maker.ilks (id), 
  tax           NUMERIC NOT NULL
);

CREATE TABLE maker.drip_vat (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  vat           TEXT
);

CREATE TABLE maker.drip_vow (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  vow           TEXT
);

CREATE TABLE maker.drip_repo (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  repo          TEXT
);

-- +goose Down
DROP TABLE maker.drip_ilk_rho;
DROP TABLE maker.drip_ilk_tax;
DROP TABLE maker.drip_vat;
DROP TABLE maker.drip_vow;
DROP TABLE maker.drip_repo;
