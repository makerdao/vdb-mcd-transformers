-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.pot_user_pie
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    "user"    INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    pie       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, "user", pie)
);

CREATE INDEX pot_user_pie_header_id_index
    ON maker.pot_user_pie (header_id);
CREATE INDEX pot_user_pie_user_index
    ON maker.pot_user_pie ("user");

COMMENT ON TABLE maker.pot_user_pie
    IS E'Value of an entry in the Pot contract\'s pie mapping as of a block header.';

CREATE TABLE maker.pot_pie
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    pie       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, pie)
);

CREATE INDEX pot_pie_header_id_index
    ON maker.pot_pie (header_id);

COMMENT ON TABLE maker.pot_pie
    IS E'Value of the Pot contract\'s Pie variable as of a block header.';

CREATE TABLE maker.pot_dsr
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    dsr       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, dsr)
);

CREATE INDEX pot_dsr_header_id_index
    ON maker.pot_dsr (header_id);

COMMENT ON TABLE maker.pot_dsr
    IS E'Value of the Pot contract\'s dsr variable as of a block header.';

CREATE TABLE maker.pot_chi
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    chi       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, chi)
);

CREATE INDEX pot_chi_header_id_index
    ON maker.pot_chi (header_id);

COMMENT ON TABLE maker.pot_chi
    IS E'Value of the Pot contract\'s chi variable as of a block header.';

CREATE TABLE maker.pot_vat
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vat       INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, vat)
);

CREATE INDEX pot_vat_header_id_index
    ON maker.pot_vat (header_id);
CREATE INDEX pot_vat_vat_index
    ON maker.pot_vat (vat);

COMMENT ON TABLE maker.pot_vat
    IS E'Value of the Pot contract\'s vat variable as of a block header.';

CREATE TABLE maker.pot_vow
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vow       INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, vow)
);

CREATE INDEX pot_vow_header_id_index
    ON maker.pot_vow (header_id);
CREATE INDEX pot_vow_vow_index
    ON maker.pot_vow (vow);

COMMENT ON TABLE maker.pot_vow
    IS E'Value of the Pot contract\'s vow variable as of a block header.';

CREATE TABLE maker.pot_rho
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    rho       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, rho)
);

CREATE INDEX pot_rho_header_id_index
    ON maker.pot_rho (header_id);

COMMENT ON TABLE maker.pot_rho
    IS E'Value of the Pot contract\'s rho variable as of a block header.';

CREATE TABLE maker.pot_live
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    live      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, live)
);

CREATE INDEX pot_live_header_id_index
    ON maker.pot_live (header_id);

COMMENT ON TABLE maker.pot_live
    IS E'Value of the Pot contract\'s live variable as of a block header.';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.pot_live_header_id_index;
DROP INDEX maker.pot_rho_header_id_index;
DROP INDEX maker.pot_vow_vow_index;
DROP INDEX maker.pot_vow_header_id_index;
DROP INDEX maker.pot_vat_vat_index;
DROP INDEX maker.pot_vat_header_id_index;
DROP INDEX maker.pot_chi_header_id_index;
DROP INDEX maker.pot_dsr_header_id_index;
DROP INDEX maker.pot_pie_header_id_index;
DROP INDEX maker.pot_user_pie_user_index;
DROP INDEX maker.pot_user_pie_header_id_index;

DROP TABLE maker.pot_live;
DROP TABLE maker.pot_rho;
DROP TABLE maker.pot_vow;
DROP TABLE maker.pot_vat;
DROP TABLE maker.pot_chi;
DROP TABLE maker.pot_dsr;
DROP TABLE maker.pot_pie;
DROP TABLE maker.pot_user_pie;
