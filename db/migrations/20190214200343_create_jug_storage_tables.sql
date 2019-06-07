-- +goose Up
CREATE TABLE maker.jug_ilk_rho
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    rho          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, ilk_id, rho)
);

CREATE TABLE maker.jug_ilk_duty
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    duty         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, ilk_id, duty)
);

CREATE TABLE maker.jug_vat
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    vat          TEXT,
    UNIQUE (block_number, block_hash, vat)
);

CREATE TABLE maker.jug_vow
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    vow          TEXT,
    UNIQUE (block_number, block_hash, vow)
);

CREATE TABLE maker.jug_base
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    base         TEXT,
    UNIQUE (block_number, block_hash, base)
);

-- +goose Down
DROP TABLE maker.jug_ilk_rho;
DROP TABLE maker.jug_ilk_duty;
DROP TABLE maker.jug_vat;
DROP TABLE maker.jug_vow;
DROP TABLE maker.jug_base;
