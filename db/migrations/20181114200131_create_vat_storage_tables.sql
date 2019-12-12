-- +goose Up
CREATE TABLE maker.vat_debt
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    debt      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, debt)
);

CREATE TABLE maker.vat_vice
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vice      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, vice)
);

CREATE TABLE maker.vat_ilk_art
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    art       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, ilk_id, art)
);

CREATE INDEX vat_ilk_art_header_id_index
    ON maker.vat_ilk_art (header_id);
CREATE INDEX vat_ilk_art_ilk_index
    ON maker.vat_ilk_art (ilk_id);

CREATE TABLE maker.vat_ilk_dust
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    dust      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, ilk_id, dust)
);

CREATE INDEX vat_ilk_dust_header_id_index
    ON maker.vat_ilk_dust (header_id);
CREATE INDEX vat_ilk_dust_ilk_index
    ON maker.vat_ilk_dust (ilk_id);

CREATE TABLE maker.vat_ilk_line
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    line      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, ilk_id, line)
);

CREATE INDEX vat_ilk_line_header_id_index
    ON maker.vat_ilk_line (header_id);
CREATE INDEX vat_ilk_line_ilk_index
    ON maker.vat_ilk_line (ilk_id);

CREATE TABLE maker.vat_ilk_spot
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    spot      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, ilk_id, spot)
);

CREATE INDEX vat_ilk_spot_header_id_index
    ON maker.vat_ilk_spot (header_id);
CREATE INDEX vat_ilk_spot_ilk_index
    ON maker.vat_ilk_spot (ilk_id);

CREATE TABLE maker.vat_ilk_rate
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    rate      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, ilk_id, rate)
);

CREATE INDEX vat_ilk_rate_header_id_index
    ON maker.vat_ilk_rate (header_id);
CREATE INDEX vat_ilk_rate_ilk_index
    ON maker.vat_ilk_rate (ilk_id);

CREATE TABLE maker.vat_urn_art
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    urn_id    INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
    art       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, urn_id, art)
);

CREATE INDEX vat_urn_art_header_id_index
    ON maker.vat_urn_art (header_id);
CREATE INDEX vat_urn_art_urn_index
    ON maker.vat_urn_art (urn_id);

CREATE TABLE maker.vat_urn_ink
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    urn_id    INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
    ink       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, urn_id, ink)
);

CREATE INDEX vat_urn_ink_header_id_index
    ON maker.vat_urn_ink (header_id);
CREATE INDEX vat_urn_ink_urn_index
    ON maker.vat_urn_ink (urn_id);

CREATE TABLE maker.vat_gem
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    guy       TEXT,
    gem       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, ilk_id, guy, gem)
);

CREATE INDEX vat_gem_ilk_index
    ON maker.vat_gem (ilk_id);

CREATE TABLE maker.vat_dai
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    guy       TEXT,
    dai       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, guy, dai)
);

CREATE TABLE maker.vat_sin
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    guy       TEXT,
    sin       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, guy, sin)
);

CREATE TABLE maker.vat_line
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    line      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, line)
);

CREATE TABLE maker.vat_live
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    live      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, live)
);

-- +goose Down
DROP INDEX maker.vat_ilk_art_header_id_index;
DROP INDEX maker.vat_ilk_art_ilk_index;
DROP INDEX maker.vat_ilk_dust_header_id_index;
DROP INDEX maker.vat_ilk_dust_ilk_index;
DROP INDEX maker.vat_ilk_line_header_id_index;
DROP INDEX maker.vat_ilk_line_ilk_index;
DROP INDEX maker.vat_ilk_spot_header_id_index;
DROP INDEX maker.vat_ilk_spot_ilk_index;
DROP INDEX maker.vat_ilk_rate_header_id_index;
DROP INDEX maker.vat_ilk_rate_ilk_index;
DROP INDEX maker.vat_urn_art_header_id_index;
DROP INDEX maker.vat_urn_art_urn_index;
DROP INDEX maker.vat_urn_ink_header_id_index;
DROP INDEX maker.vat_urn_ink_urn_index;
DROP INDEX maker.vat_gem_ilk_index;

DROP TABLE maker.vat_debt;
DROP TABLE maker.vat_vice;
DROP TABLE maker.vat_ilk_art;
DROP TABLE maker.vat_ilk_dust;
DROP TABLE maker.vat_ilk_line;
DROP TABLE maker.vat_ilk_spot;
DROP TABLE maker.vat_ilk_rate;
DROP TABLE maker.vat_urn_art;
DROP TABLE maker.vat_urn_ink;
DROP TABLE maker.vat_gem;
DROP TABLE maker.vat_dai;
DROP TABLE maker.vat_sin;
DROP TABLE maker.vat_line;
DROP TABLE maker.vat_live;
