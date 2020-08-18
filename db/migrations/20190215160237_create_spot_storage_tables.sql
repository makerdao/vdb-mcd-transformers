-- +goose Up
CREATE TABLE maker.spot_ilk_pip
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    pip       TEXT,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, ilk_id, pip)
);

CREATE INDEX spot_ilk_pip_header_id_index
    ON maker.spot_ilk_pip (header_id);
CREATE INDEX spot_ilk_pip_ilk_index
    ON maker.spot_ilk_pip (ilk_id);

CREATE TABLE maker.spot_ilk_mat
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    mat       NUMERIC NOT NULL,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, ilk_id, mat)
);

CREATE INDEX spot_ilk_mat_header_id_index
    ON maker.spot_ilk_mat (header_id);
CREATE INDEX spot_ilk_mat_ilk_index
    ON maker.spot_ilk_mat (ilk_id);

CREATE TABLE maker.spot_vat
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    vat       TEXT,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, vat)
);

CREATE INDEX spot_vat_header_id_index
    ON maker.spot_vat (header_id);

CREATE TABLE maker.spot_par
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    par       NUMERIC NOT NULL,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, par)
);

CREATE INDEX spot_par_header_id_index
    ON maker.spot_par (header_id);

CREATE TABLE maker.spot_live
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    live      NUMERIC NOT NULL,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, live)
);

CREATE INDEX spot_live_header_id_index
    ON maker.spot_live (header_id);

-- +goose Down
DROP INDEX maker.spot_ilk_pip_header_id_index;
DROP INDEX maker.spot_ilk_pip_ilk_index;
DROP INDEX maker.spot_ilk_mat_header_id_index;
DROP INDEX maker.spot_ilk_mat_ilk_index;
DROP INDEX maker.spot_vat_header_id_index;
DROP INDEX maker.spot_par_header_id_index;
DROP INDEX maker.spot_live_header_id_index;

DROP TABLE maker.spot_par;
DROP TABLE maker.spot_vat;
DROP TABLE maker.spot_ilk_mat;
DROP TABLE maker.spot_ilk_pip;
DROP TABLE maker.spot_live;