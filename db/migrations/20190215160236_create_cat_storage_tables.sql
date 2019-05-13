-- +goose Up
CREATE TABLE maker.cat_nflip (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  nflip NUMERIC NOT NULL,
  UNIQUE (block_number, block_hash, nflip)
);

CREATE TABLE maker.cat_live (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  live NUMERIC NOT NULL,
  UNIQUE (block_number, block_hash, live)
);

CREATE TABLE maker.cat_vat (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  vat TEXT,
  UNIQUE (block_number, block_hash, vat)
);

CREATE TABLE maker.cat_vow (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  vow TEXT,
  UNIQUE (block_number, block_hash, vow)
);

CREATE TABLE maker.cat_ilk_flip (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  ilk_id INTEGER NOT NULL REFERENCES maker.ilks (id),
  flip TEXT,
  UNIQUE (block_number, block_hash, ilk_id, flip)
);

CREATE TABLE maker.cat_ilk_chop (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  ilk_id INTEGER NOT NULL REFERENCES maker.ilks (id),
  chop NUMERIC NOT NULL,
  UNIQUE (block_number, block_hash, ilk_id, chop)
);

CREATE TABLE maker.cat_ilk_lump (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  ilk_id INTEGER NOT NULL REFERENCES maker.ilks (id),
  lump NUMERIC NOT NULL,
  UNIQUE (block_number, block_hash, ilk_id, lump)
);

CREATE TABLE maker.cat_flip_ilk (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  flip NUMERIC NOT NULL,
  ilk_id INTEGER NOT NULL REFERENCES maker.ilks (id),
  UNIQUE (block_number, block_hash, flip, ilk_id)
);

CREATE TABLE maker.cat_flip_urn (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  flip NUMERIC NOT NULL,
  urn TEXT,
  UNIQUE (block_number, block_hash, flip, urn)
);

CREATE TABLE maker.cat_flip_ink (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  flip NUMERIC NOT NULL,
  ink NUMERIC NOT NULL,
  UNIQUE (block_number, block_hash, flip, ink)
);

CREATE TABLE maker.cat_flip_tab (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  flip NUMERIC NOT NULL,
  tab NUMERIC NOT NULL,
  UNIQUE (block_number, block_hash, flip, tab)
);


-- +goose Down
DROP TABLE maker.cat_nflip;
DROP TABLE maker.cat_live;
DROP TABLE maker.cat_vat;
DROP TABLE maker.cat_vow;
DROP TABLE maker.cat_ilk_flip;
DROP TABLE maker.cat_ilk_chop;
DROP TABLE maker.cat_ilk_lump;
DROP TABLE maker.cat_flip_ilk;
DROP TABLE maker.cat_flip_urn;
DROP TABLE maker.cat_flip_ink;
DROP TABLE maker.cat_flip_tab;
