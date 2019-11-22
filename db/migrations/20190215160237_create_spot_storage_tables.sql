-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.spot_ilk_pip
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    pip       TEXT,
    UNIQUE (header_id, ilk_id, pip)
);

CREATE INDEX spot_ilk_pip_ilk_index
    ON maker.spot_ilk_pip (ilk_id);

CREATE TABLE maker.spot_ilk_mat
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    mat       NUMERIC NOT NULL,
    UNIQUE (header_id, ilk_id, mat)
);

CREATE INDEX spot_ilk_mat_ilk_index
    ON maker.spot_ilk_mat (ilk_id);

CREATE TABLE maker.spot_vat
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vat       TEXT,
    UNIQUE (header_id, vat)
);

CREATE TABLE maker.spot_par
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    par       NUMERIC NOT NULL,
    UNIQUE (header_id, par)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.spot_ilk_pip_ilk_index;
DROP INDEX maker.spot_ilk_mat_ilk_index;

DROP TABLE maker.spot_par;
DROP TABLE maker.spot_vat;
DROP TABLE maker.spot_ilk_mat;
DROP TABLE maker.spot_ilk_pip;
