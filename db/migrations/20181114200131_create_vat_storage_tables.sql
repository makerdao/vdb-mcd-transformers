-- +goose Up
CREATE TABLE maker.vat_debt
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    debt         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, debt)
);

CREATE TABLE maker.vat_vice
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    vice         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, vice)
);

CREATE TABLE maker.vat_ilk_art
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    art          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, ilk_id, art)
);

CREATE INDEX vat_ilk_art_block_number_index
    ON maker.vat_ilk_art (block_number);

CREATE INDEX vat_ilk_art_ilk_index
    ON maker.vat_ilk_art (ilk_id);

CREATE TABLE maker.vat_ilk_dust
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    dust         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, ilk_id, dust)
);

CREATE INDEX vat_ilk_dust_block_number_index
    ON maker.vat_ilk_dust (block_number);

CREATE INDEX vat_ilk_dust_ilk_index
    ON maker.vat_ilk_dust (ilk_id);

CREATE TABLE maker.vat_ilk_line
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    line         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, ilk_id, line)
);

CREATE INDEX vat_ilk_line_block_number_index
    ON maker.vat_ilk_line (block_number);

CREATE INDEX vat_ilk_line_ilk_index
    ON maker.vat_ilk_line (ilk_id);

CREATE TABLE maker.vat_ilk_spot
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    spot         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, ilk_id, spot)
);

CREATE INDEX vat_ilk_spot_block_number_index
    ON maker.vat_ilk_spot (block_number);

CREATE INDEX vat_ilk_spot_ilk_index
    ON maker.vat_ilk_spot (ilk_id);

CREATE TABLE maker.vat_ilk_rate
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    rate         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, ilk_id, rate)
);

CREATE INDEX vat_ilk_rate_block_number_index
    ON maker.vat_ilk_rate (block_number);

CREATE INDEX vat_ilk_rate_ilk_index
    ON maker.vat_ilk_rate (ilk_id);

CREATE TABLE maker.vat_urn_art
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    urn_id       INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
    art          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, urn_id, art)
);

CREATE INDEX vat_urn_art_block_number_index
    ON maker.vat_urn_art (block_number);

CREATE INDEX vat_urn_art_urn_index
    ON maker.vat_urn_art (urn_id);

CREATE TABLE maker.vat_urn_ink
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    urn_id       INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
    ink          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, urn_id, ink)
);

CREATE INDEX vat_urn_ink_block_number_index
    ON maker.vat_urn_ink (block_number);

CREATE INDEX vat_urn_ink_urn_index
    ON maker.vat_urn_ink (urn_id);

CREATE TABLE maker.vat_gem
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    ilk_id       INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    guy          TEXT,
    gem          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, ilk_id, guy, gem)
);

CREATE INDEX vat_gem_ilk_index
    ON maker.vat_gem (ilk_id);

CREATE TABLE maker.vat_dai
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    guy          TEXT,
    dai          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, guy, dai)
);

CREATE TABLE maker.vat_sin
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    guy          TEXT,
    sin          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, guy, sin)
);

CREATE TABLE maker.vat_line
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    line         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, line)
);

CREATE TABLE maker.vat_live
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    live         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, live)
);

-- +goose Down
DROP INDEX maker.vat_ilk_art_block_number_index;
DROP INDEX maker.vat_ilk_art_ilk_index;
DROP INDEX maker.vat_ilk_dust_block_number_index;
DROP INDEX maker.vat_ilk_dust_ilk_index;
DROP INDEX maker.vat_ilk_line_block_number_index;
DROP INDEX maker.vat_ilk_line_ilk_index;
DROP INDEX maker.vat_ilk_spot_block_number_index;
DROP INDEX maker.vat_ilk_spot_ilk_index;
DROP INDEX maker.vat_ilk_rate_block_number_index;
DROP INDEX maker.vat_ilk_rate_ilk_index;
DROP INDEX maker.vat_urn_art_block_number_index;
DROP INDEX maker.vat_urn_art_urn_index;
DROP INDEX maker.vat_urn_ink_block_number_index;
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
