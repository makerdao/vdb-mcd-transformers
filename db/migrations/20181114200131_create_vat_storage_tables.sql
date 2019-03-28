-- +goose Up
CREATE TABLE maker.vat_debt (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  debt          NUMERIC NOT NULL
);

CREATE TABLE maker.vat_vice (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  vice          NUMERIC NOT NULL
);

CREATE TABLE maker.vat_ilk_art (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk_id        INTEGER NOT NULL REFERENCES maker.ilks (id),
  art           NUMERIC NOT NULL
);

CREATE TABLE maker.vat_ilk_dust (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk_id        INTEGER NOT NULL REFERENCES maker.ilks (id),
  dust          NUMERIC NOT NULL
);

CREATE TABLE maker.vat_ilk_line (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk_id        INTEGER NOT NULL REFERENCES maker.ilks (id),
  line          NUMERIC NOT NULL
);

CREATE TABLE maker.vat_ilk_spot (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk_id        INTEGER NOT NULL REFERENCES maker.ilks (id),
  spot          NUMERIC NOT NULL
);

-- TODO: remove this once the ilk query no longer depends on it
CREATE TABLE maker.vat_ilk_ink (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk_id        INTEGER NOT NULL REFERENCES maker.ilks (id),
  ink           NUMERIC NOT NULL
);

CREATE TABLE maker.vat_ilk_rate (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk_id        INTEGER NOT NULL REFERENCES maker.ilks (id),
  rate          NUMERIC NOT NULL
);

-- TODO: remove this once the ilk query no longer depends on it
CREATE TABLE maker.vat_ilk_take (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk_id        INTEGER NOT NULL REFERENCES maker.ilks (id),
  take          NUMERIC NOT NULL
);

CREATE TABLE maker.vat_urn_art (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  urn_id        INTEGER NOT NULL REFERENCES maker.urns (id),
  art           NUMERIC NOT NULL
);

CREATE TABLE maker.vat_urn_ink (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  urn_id        INTEGER NOT NULL REFERENCES maker.urns (id),
  ink           NUMERIC NOT NULL
);

CREATE TABLE maker.vat_gem (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk_id        INTEGER NOT NULL REFERENCES maker.ilks (id),
  guy           TEXT,
  gem           NUMERIC NOT NULL
);

CREATE TABLE maker.vat_dai (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  guy           TEXT,
  dai           NUMERIC NOT NULL
);

CREATE TABLE maker.vat_sin (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  guy           TEXT,
  sin           NUMERIC NOT NULL
);

CREATE TABLE maker.vat_line (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  line          NUMERIC NOT NULL
);

CREATE TABLE maker.vat_live (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  live          NUMERIC NOT NULL
);

-- +goose Down
DROP TABLE maker.vat_debt;
DROP TABLE maker.vat_vice;
DROP TABLE maker.vat_ilk_art;
DROP TABLE maker.vat_ilk_dust;
DROP TABLE maker.vat_ilk_line;
DROP TABLE maker.vat_ilk_spot;
DROP TABLE maker.vat_ilk_ink; -- TODO: remove this once the ilk query no longer depends on it
DROP TABLE maker.vat_ilk_rate;
DROP TABLE maker.vat_ilk_take; -- TODO: remove this once the ilk query no longer depends on it
DROP TABLE maker.vat_urn_art;
DROP TABLE maker.vat_urn_ink;
DROP TABLE maker.vat_gem;
DROP TABLE maker.vat_dai;
DROP TABLE maker.vat_sin;
DROP TABLE maker.vat_line;
DROP TABLE maker.vat_live;
