-- +goose Up
CREATE TABLE maker.jug_ilk_rho(
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk_id        INTEGER NOT NULL REFERENCES maker.ilks (id),
  rho           NUMERIC NOT NULL
);

CREATE TABLE maker.jug_ilk_duty(
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk_id        INTEGER NOT NULL REFERENCES maker.ilks (id),
  duty          NUMERIC NOT NULL
);

CREATE TABLE maker.jug_vat(
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  vat           TEXT
);

CREATE TABLE maker.jug_vow(
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  vow           TEXT
);

CREATE TABLE maker.jug_base (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  base          TEXT
);

-- +goose Down
DROP TABLE maker.jug_ilk_rho;
DROP TABLE maker.jug_ilk_duty;
DROP TABLE maker.jug_vat;
DROP TABLE maker.jug_vow;
DROP TABLE maker.jug_base;
