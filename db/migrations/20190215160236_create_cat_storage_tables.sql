-- +goose Up
CREATE TABLE maker.cat_live
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    live         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, live)
);

CREATE TABLE maker.cat_vat
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    vat          TEXT,
    UNIQUE (block_number, block_hash, vat)
);

CREATE TABLE maker.cat_vow
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    vow          TEXT,
    UNIQUE (block_number, block_hash, vow)
);

CREATE TABLE maker.cat_ilk_flip
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    flip         TEXT,
    UNIQUE (block_number, block_hash, ilk_id, flip)
);

CREATE TABLE maker.cat_ilk_chop
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    chop         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, ilk_id, chop)
);

CREATE TABLE maker.cat_ilk_lump
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    lump         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, ilk_id, lump)
);


-- +goose Down
DROP TABLE maker.cat_live;
DROP TABLE maker.cat_vat;
DROP TABLE maker.cat_vow;
DROP TABLE maker.cat_ilk_flip;
DROP TABLE maker.cat_ilk_chop;
DROP TABLE maker.cat_ilk_lump;
