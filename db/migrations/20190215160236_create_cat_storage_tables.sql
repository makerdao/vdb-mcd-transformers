-- +goose Up
CREATE TABLE maker.cat_live
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    live      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, live)
);

CREATE INDEX cat_live_header_id_index
    ON maker.cat_live (header_id);
COMMENT ON TABLE maker.cat_live
    IS E'@omit';

CREATE TABLE maker.cat_vat
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vat       TEXT,
    UNIQUE (diff_id, header_id, vat)
);

CREATE INDEX cat_vat_header_id_index
    ON maker.cat_vat (header_id);
COMMENT ON TABLE maker.cat_vat
    IS E'@omit';

CREATE TABLE maker.cat_vow
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vow       TEXT,
    UNIQUE (diff_id, header_id, vow)
);

CREATE INDEX cat_vow_header_id_index
    ON maker.cat_vow (header_id);
COMMENT ON TABLE maker.cat_vow
    IS E'@omit';

CREATE TABLE maker.cat_ilk_flip
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    flip      TEXT,
    UNIQUE (diff_id, header_id, ilk_id, flip)
);

CREATE INDEX cat_ilk_flip_header_id_index
    ON maker.cat_ilk_flip (header_id);
COMMENT ON TABLE maker.cat_ilk_flip
    IS E'@omit';
CREATE INDEX cat_ilk_flip_ilk_index
    ON maker.cat_ilk_flip (ilk_id);

CREATE TABLE maker.cat_ilk_chop
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    chop      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, ilk_id, chop)
);

CREATE INDEX cat_ilk_chop_header_id_index
    ON maker.cat_ilk_chop (header_id);
COMMENT ON TABLE maker.cat_ilk_chop
    IS E'@omit';
CREATE INDEX cat_ilk_chop_ilk_index
    ON maker.cat_ilk_chop (ilk_id);

CREATE TABLE maker.cat_ilk_lump
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    lump      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, ilk_id, lump)
);

CREATE INDEX cat_ilk_lump_header_id_index
    ON maker.cat_ilk_lump (header_id);
COMMENT ON TABLE maker.cat_ilk_lump
    IS E'@omit';
CREATE INDEX cat_ilk_lump_ilk_index
    ON maker.cat_ilk_lump (ilk_id);


-- +goose Down
DROP INDEX maker.cat_live_header_id_index;
DROP INDEX maker.cat_vat_header_id_index;
DROP INDEX maker.cat_vow_header_id_index;
DROP INDEX maker.cat_ilk_flip_header_id_index;
DROP INDEX maker.cat_ilk_flip_ilk_index;
DROP INDEX maker.cat_ilk_chop_header_id_index;
DROP INDEX maker.cat_ilk_chop_ilk_index;
DROP INDEX maker.cat_ilk_lump_header_id_index;
DROP INDEX maker.cat_ilk_lump_ilk_index;

DROP TABLE maker.cat_live;
DROP TABLE maker.cat_vat;
DROP TABLE maker.cat_vow;
DROP TABLE maker.cat_ilk_flip;
DROP TABLE maker.cat_ilk_chop;
DROP TABLE maker.cat_ilk_lump;
