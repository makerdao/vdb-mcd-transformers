-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.cdp_manager_vat
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    vat       TEXT,
    UNIQUE (diff_id, header_id, vat)
);

CREATE TABLE maker.cdp_manager_cdpi
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    cdpi      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, cdpi)
);

CREATE INDEX cdp_manager_cdpi_cdpi_index
    ON maker.cdp_manager_cdpi (cdpi);

CREATE TABLE maker.cdp_manager_urns
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    cdpi      NUMERIC NOT NULL,
    urn       TEXT,
    UNIQUE (diff_id, header_id, cdpi, urn)
);

CREATE INDEX cdp_manager_urns_urn_index
    ON maker.cdp_manager_urns (urn);
CREATE INDEX cdp_manager_urns_cdpi_index
    ON maker.cdp_manager_urns (cdpi);

CREATE TABLE maker.cdp_manager_list_prev
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    cdpi      NUMERIC NOT NULL,
    prev      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, cdpi, prev)
);

CREATE TABLE maker.cdp_manager_list_next
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    cdpi      NUMERIC NOT NULL,
    next      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, cdpi, next)
);

CREATE TABLE maker.cdp_manager_owns
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    cdpi      NUMERIC NOT NULL,
    owner     TEXT,
    UNIQUE (diff_id, header_id, cdpi, owner)
);

CREATE INDEX cdp_manager_owns_cdpi_index
    ON maker.cdp_manager_owns (cdpi);
CREATE INDEX cdp_manager_owns_owner_index
    ON maker.cdp_manager_owns (owner);

CREATE TABLE maker.cdp_manager_ilks
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    cdpi      NUMERIC NOT NULL,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, cdpi, ilk_id)
);

CREATE INDEX cdp_manager_ilks_cdpi_index
    ON maker.cdp_manager_ilks (cdpi);
CREATE INDEX cdp_manager_ilks_ilk_id_index
    ON maker.cdp_manager_ilks (ilk_id);

CREATE TABLE maker.cdp_manager_first
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    owner     TEXT,
    first     NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, owner, first)
);

CREATE TABLE maker.cdp_manager_last
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    owner     TEXT,
    last      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, owner, last)
);

CREATE TABLE maker.cdp_manager_count
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    owner     TEXT,
    count     NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, owner, count)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.cdp_manager_ilks_cdpi_index;
DROP INDEX maker.cdp_manager_ilks_ilk_id_index;
DROP INDEX maker.cdp_manager_owns_cdpi_index;
DROP INDEX maker.cdp_manager_owns_owner_index;
DROP INDEX maker.cdp_manager_urns_urn_index;
DROP INDEX maker.cdp_manager_urns_cdpi_index;
DROP INDEX maker.cdp_manager_cdpi_cdpi_index;

DROP TABLE maker.cdp_manager_cdpi;
DROP TABLE maker.cdp_manager_vat;
DROP TABLE maker.cdp_manager_urns;
DROP TABLE maker.cdp_manager_list_prev;
DROP TABLE maker.cdp_manager_list_next;
DROP TABLE maker.cdp_manager_owns;
DROP TABLE maker.cdp_manager_ilks;
DROP TABLE maker.cdp_manager_first;
DROP TABLE maker.cdp_manager_last;
DROP TABLE maker.cdp_manager_count;
