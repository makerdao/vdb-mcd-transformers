-- +goose Up
CREATE TABLE maker.median_val
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    val       NUMERIC NOT NULL,
    UNIQUE (header_id, address_id, val)
);

CREATE INDEX median_val_diff_id_index
    ON maker.median_val (diff_id);
CREATE INDEX median_val_header_id_index
    ON maker.median_val (header_id);
CREATE INDEX median_val_address_id_index
    ON maker.median_val (address_id);

CREATE TABLE maker.median_bar
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bar       NUMERIC NOT NULL,
    UNIQUE (header_id, address_id, bar)
);

CREATE INDEX median_bar_diff_id_index
    ON maker.median_bar (diff_id);
CREATE INDEX median_bar_header_id_index
    ON maker.median_bar (header_id);
CREATE INDEX median_bar_address_id_index
    ON maker.median_bar (address_id);

CREATE TABLE maker.median_age
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    age       NUMERIC NOT NULL,
    UNIQUE (header_id, address_id, age)
);

CREATE INDEX median_age_diff_id_index
    ON maker.median_age (diff_id);
CREATE INDEX median_age_header_id_index
    ON maker.median_age (header_id);
CREATE INDEX median_age_address_id_index
    ON maker.median_age (address_id);

CREATE TABLE maker.median_orcl
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    a          INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    orcl        INTEGER NOT NULL,
    UNIQUE (header_id, address_id, a, orcl)
);

CREATE INDEX median_orcl_diff_id_index
    ON maker.median_orcl (diff_id);
CREATE INDEX median_orcl_header_id_index
    ON maker.median_orcl (header_id);
CREATE INDEX median_orcl_address_id_index
    ON maker.median_orcl (address_id);
CREATE INDEX median_orcl_a_index
    ON maker.median_orcl (a);

CREATE TABLE maker.median_bud
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    a          INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bud        INTEGER NOT NULL,
    UNIQUE (header_id, address_id, a, bud)
);

CREATE INDEX median_bud_header_id_index
    ON maker.median_bud (header_id);
CREATE INDEX median_bud_address_index
    ON maker.median_bud (address_id);
CREATE INDEX median_bud_a_index
    ON maker.median_bud (a);

CREATE TABLE maker.median_slot
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    slot_id    INTEGER NOT NULL,
    slot       INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, address_id, slot_id, slot)
);

CREATE INDEX median_slot_header_id_index
    ON maker.median_slot (header_id);
CREATE INDEX median_slot_address_index
    ON maker.median_slot (address_id);
CREATE INDEX median_slot_id_index
    ON maker.median_slot (slot_id);

-- +goose Down

DROP TABLE maker.median_val;
DROP TABLE maker.median_bar;
DROP TABLE maker.median_age;
DROP TABLE maker.median_orcl;
DROP TABLE maker.median_bud;
DROP TABLE maker.median_slot;
