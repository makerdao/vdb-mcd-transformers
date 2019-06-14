-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.spot_ilk_pip
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    pip          TEXT,
    UNIQUE (block_number, block_hash, ilk_id, pip)
);

CREATE TABLE maker.spot_ilk_mat
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    mat          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, ilk_id, mat)
);

CREATE TABLE maker.spot_vat
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    vat          TEXT,
    UNIQUE (block_number, block_hash, vat)
);

CREATE TABLE maker.spot_par
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    par          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, par)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE maker.spot_par;
DROP TABLE maker.spot_vat;
DROP TABLE maker.spot_ilk_mat;
DROP TABLE maker.spot_ilk_pip;
