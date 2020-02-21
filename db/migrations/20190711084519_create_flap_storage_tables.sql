-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.flap_bid_bid
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    bid        NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, bid_id, bid)
);

CREATE INDEX flap_bid_bid_header_id_index
    ON maker.flap_bid_bid (header_id);
CREATE INDEX flap_bid_bid_bid_id_index
    ON maker.flap_bid_bid (bid_id);
CREATE INDEX flap_bid_bid_address_index
    ON maker.flap_bid_bid (address_id);

CREATE TABLE maker.flap_bid_lot
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, bid_id, lot)
);

CREATE INDEX flap_bid_lot_header_id_index
    ON maker.flap_bid_lot (header_id);
CREATE INDEX flap_bid_lot_bid_id_index
    ON maker.flap_bid_lot (bid_id);
CREATE INDEX flap_bid_lot_bid_address_index
    ON maker.flap_bid_lot (address_id);

CREATE TABLE maker.flap_bid_guy
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    guy        TEXT    NOT NULL,
    UNIQUE (diff_id, header_id, address_id, bid_id, guy)
);

CREATE INDEX flap_bid_guy_header_id_index
    ON maker.flap_bid_guy (header_id);
CREATE INDEX flap_bid_guy_bid_id_index
    ON maker.flap_bid_guy (bid_id);
CREATE INDEX flap_bid_guy_bid_address_index
    ON maker.flap_bid_guy (address_id);

CREATE TABLE maker.flap_bid_tic
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    tic        BIGINT  NOT NULL,
    UNIQUE (diff_id, header_id, address_id, bid_id, tic)
);

CREATE INDEX flap_bid_tic_header_id_index
    ON maker.flap_bid_tic (header_id);
CREATE INDEX flap_bid_tic_bid_id_index
    ON maker.flap_bid_tic (bid_id);
CREATE INDEX flap_bid_tic_bid_address_index
    ON maker.flap_bid_tic (address_id);

CREATE TABLE maker.flap_bid_end
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    "end"      BIGINT  NOT NULL,
    UNIQUE (diff_id, header_id, address_id, bid_id, "end")
);

CREATE INDEX flap_bid_end_header_id_index
    ON maker.flap_bid_end (header_id);
CREATE INDEX flap_bid_end_bid_id_index
    ON maker.flap_bid_end (bid_id);
CREATE INDEX flap_bid_end_bid_address_index
    ON maker.flap_bid_end (address_id);

CREATE TABLE maker.flap_vat
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    vat        TEXT    NOT NULL,
    UNIQUE (diff_id, header_id, address_id, vat)
);

CREATE INDEX flap_vat_header_id_index
    ON maker.flap_vat (header_id);
CREATE INDEX flap_vat_address_index
    ON maker.flap_vat (address_id);

CREATE TABLE maker.flap_gem
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    gem        TEXT    NOT NULL,
    UNIQUE (diff_id, header_id, address_id, gem)
);

CREATE INDEX flap_gem_header_id_index
    ON maker.flap_gem (header_id);
CREATE INDEX flap_gem_address_index
    ON maker.flap_gem (address_id);

CREATE TABLE maker.flap_beg
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    beg        NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, beg)
);

CREATE INDEX flap_beg_header_id_index
    ON maker.flap_beg (header_id);
CREATE INDEX flap_beg_address_index
    ON maker.flap_beg (address_id);

CREATE TABLE maker.flap_ttl
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    ttl        INT     NOT NULL,
    UNIQUE (diff_id, header_id, address_id, ttl)
);

CREATE INDEX flap_ttl_header_id_index
    ON maker.flap_ttl (header_id);
CREATE INDEX flap_ttl_address_index
    ON maker.flap_ttl (address_id);

CREATE TABLE maker.flap_tau
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    tau        INT     NOT NULL,
    UNIQUE (diff_id, header_id, address_id, tau)
);

CREATE INDEX flap_tau_header_id_index
    ON maker.flap_tau (header_id);
CREATE INDEX flap_tau_address_index
    ON maker.flap_tau (address_id);

CREATE TABLE maker.flap_kicks
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    kicks      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, kicks)
);

CREATE INDEX flap_kicks_header_id_index
    ON maker.flap_kicks (header_id);
CREATE INDEX flap_kicks_address_index
    ON maker.flap_kicks (address_id);

CREATE TABLE maker.flap_live
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    live       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, live)
);

CREATE INDEX flap_live_header_id_index
    ON maker.flap_live (header_id);
CREATE INDEX flap_live_address_index
    ON maker.flap_live (address_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.flap_live_address_index;
DROP INDEX maker.flap_live_header_id_index;
DROP INDEX maker.flap_kicks_address_index;
DROP INDEX maker.flap_kicks_header_id_index;
DROP INDEX maker.flap_tau_address_index;
DROP INDEX maker.flap_tau_header_id_index;
DROP INDEX maker.flap_ttl_address_index;
DROP INDEX maker.flap_ttl_header_id_index;
DROP INDEX maker.flap_beg_address_index;
DROP INDEX maker.flap_beg_header_id_index;
DROP INDEX maker.flap_gem_address_index;
DROP INDEX maker.flap_gem_header_id_index;
DROP INDEX maker.flap_vat_address_index;
DROP INDEX maker.flap_vat_header_id_index;
DROP INDEX maker.flap_bid_bid_address_index;
DROP INDEX maker.flap_bid_bid_bid_id_index;
DROP INDEX maker.flap_bid_bid_header_id_index;
DROP INDEX maker.flap_bid_lot_bid_address_index;
DROP INDEX maker.flap_bid_lot_bid_id_index;
DROP INDEX maker.flap_bid_lot_header_id_index;
DROP INDEX maker.flap_bid_guy_bid_address_index;
DROP INDEX maker.flap_bid_guy_bid_id_index;
DROP INDEX maker.flap_bid_guy_header_id_index;
DROP INDEX maker.flap_bid_tic_bid_address_index;
DROP INDEX maker.flap_bid_tic_bid_id_index;
DROP INDEX maker.flap_bid_tic_header_id_index;
DROP INDEX maker.flap_bid_end_bid_address_index;
DROP INDEX maker.flap_bid_end_bid_id_index;
DROP INDEX maker.flap_bid_end_header_id_index;

DROP TABLE maker.flap_bid_bid;
DROP TABLE maker.flap_bid_lot;
DROP TABLE maker.flap_bid_guy;
DROP TABLE maker.flap_bid_tic;
DROP TABLE maker.flap_bid_end;
DROP TABLE maker.flap_beg;
DROP TABLE maker.flap_vat;
DROP TABLE maker.flap_gem;
DROP TABLE maker.flap_ttl;
DROP TABLE maker.flap_tau;
DROP TABLE maker.flap_kicks;
DROP TABLE maker.flap_live;