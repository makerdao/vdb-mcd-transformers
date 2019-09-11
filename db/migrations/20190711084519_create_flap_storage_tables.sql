-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.flap_bid_bid
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id       NUMERIC NOT NULL,
    bid          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, address_id, bid_id, bid)
);
CREATE INDEX flap_bid_bid_block_number_index ON maker.flap_bid_bid (block_number);
CREATE INDEX flap_bid_bid_bid_id_index ON maker.flap_bid_bid (bid_id);
CREATE INDEX flap_bid_bid_address_id_index ON maker.flap_bid_bid (address_id);

CREATE TABLE maker.flap_bid_lot
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id       NUMERIC NOT NULL,
    lot          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, address_id, bid_id, lot)
);
CREATE INDEX flap_bid_lot_block_number_index ON maker.flap_bid_lot (block_number);
CREATE INDEX flap_bid_lot_bid_id_index ON maker.flap_bid_lot (bid_id);
CREATE INDEX flap_bid_lot_bid_address_id_index ON maker.flap_bid_lot (address_id);

CREATE TABLE maker.flap_bid_guy
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id       NUMERIC NOT NULL,
    guy          TEXT    NOT NULL,
    UNIQUE (block_number, block_hash, address_id, bid_id, guy)
);
CREATE INDEX flap_bid_guy_block_number_index ON maker.flap_bid_guy (block_number);
CREATE INDEX flap_bid_guy_bid_id_index ON maker.flap_bid_guy (bid_id);
CREATE INDEX flap_bid_guy_bid_address_id_index ON maker.flap_bid_guy (address_id);

CREATE TABLE maker.flap_bid_tic
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id       NUMERIC NOT NULL,
    tic          BIGINT  NOT NULL,
    UNIQUE (block_number, block_hash, address_id, bid_id, tic)
);
CREATE INDEX flap_bid_tic_block_number_index ON maker.flap_bid_tic (block_number);
CREATE INDEX flap_bid_tic_bid_id_index ON maker.flap_bid_tic (bid_id);
CREATE INDEX flap_bid_tic_bid_address_id_index ON maker.flap_bid_tic (address_id);

CREATE TABLE maker.flap_bid_end
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id       NUMERIC NOT NULL,
    "end"        BIGINT  NOT NULL,
    UNIQUE (block_number, block_hash, address_id, bid_id, "end")
);
CREATE INDEX flap_bid_end_block_number_index ON maker.flap_bid_end (block_number);
CREATE INDEX flap_bid_end_bid_id_index ON maker.flap_bid_end (bid_id);
CREATE INDEX flap_bid_end_bid_address_id_index ON maker.flap_bid_end (address_id);

CREATE TABLE maker.flap_vat
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    vat          TEXT    NOT NULL,
    UNIQUE (block_number, block_hash, address_id, vat)
);

CREATE TABLE maker.flap_gem
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    gem          TEXT    NOT NULL,
    UNIQUE (block_number, block_hash, address_id, gem)
);

CREATE TABLE maker.flap_beg
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    beg          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, address_id, beg)
);

CREATE TABLE maker.flap_ttl
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    ttl          INT     NOT NULL,
    UNIQUE (block_number, block_hash, address_id, ttl)
);

CREATE TABLE maker.flap_tau
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    tau          INT     NOT NULL,
    UNIQUE (block_number, block_hash, address_id, tau)
);

CREATE TABLE maker.flap_kicks
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    kicks        NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, address_id, kicks)
);

CREATE INDEX flap_kicks_block_number_index ON maker.flap_kicks (block_number);
CREATE INDEX flap_kicks_kicks_index ON maker.flap_kicks (kicks);
CREATE INDEX flap_kicks_address_id_index ON maker.flap_kicks (address_id);

-- prevent naming conflict with maker.flap_kick in postgraphile
COMMENT ON TABLE maker.flap_kicks IS E'@name flapKicksStorage';

CREATE TABLE maker.flap_live
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    live         NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, address_id, live)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.flap_kicks_address_id_index;
DROP INDEX maker.flap_kicks_kicks_index;
DROP INDEX maker.flap_kicks_block_number_index;
DROP INDEX maker.flap_bid_bid_block_number_index;
DROP INDEX maker.flap_bid_bid_address_id_index;
DROP INDEX maker.flap_bid_bid_bid_id_index;
DROP INDEX maker.flap_bid_lot_block_number_index;
DROP INDEX maker.flap_bid_lot_bid_address_id_index;
DROP INDEX maker.flap_bid_lot_bid_id_index;
DROP INDEX maker.flap_bid_guy_block_number_index;
DROP INDEX maker.flap_bid_guy_bid_address_id_index;
DROP INDEX maker.flap_bid_guy_bid_id_index;
DROP INDEX maker.flap_bid_tic_block_number_index;
DROP INDEX maker.flap_bid_tic_bid_address_id_index;
DROP INDEX maker.flap_bid_tic_bid_id_index;
DROP INDEX maker.flap_bid_end_block_number_index;
DROP INDEX maker.flap_bid_end_bid_address_id_index;
DROP INDEX maker.flap_bid_end_bid_id_index;

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