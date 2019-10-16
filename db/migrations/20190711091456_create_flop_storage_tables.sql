-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.flop_bid_bid
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id       NUMERIC NOT NULL,
    bid          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, bid_id, address_id, bid)
);

CREATE INDEX flop_bid_bid_block_number_index
    ON maker.flop_bid_bid (block_number);
CREATE INDEX flop_bid_bid_bid_id_index
    ON maker.flop_bid_bid (bid_id);
CREATE INDEX flop_bid_bid_address_id_index
    ON maker.flop_bid_bid (address_id);

CREATE TABLE maker.flop_bid_lot
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id       NUMERIC NOT NULL,
    lot          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, bid_id, address_id, lot)
);

CREATE INDEX flop_bid_lot_block_number_index
    ON maker.flop_bid_lot (block_number);
CREATE INDEX flop_bid_lot_bid_id_index
    ON maker.flop_bid_lot (bid_id);
CREATE INDEX flop_bid_lot_bid_address_id_index
    ON maker.flop_bid_lot (address_id);

CREATE TABLE maker.flop_bid_guy
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id       NUMERIC NOT NULL,
    guy          TEXT,
    UNIQUE (block_number, block_hash, bid_id, address_id, guy)
);

CREATE INDEX flop_bid_guy_block_number_index
    ON maker.flop_bid_guy (block_number);
CREATE INDEX flop_bid_guy_bid_id_index
    ON maker.flop_bid_guy (bid_id);
CREATE INDEX flop_bid_guy_bid_address_id_index
    ON maker.flop_bid_guy (address_id);

CREATE TABLE maker.flop_bid_tic
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id       NUMERIC NOT NULL,
    tic          BIGINT  NOT NULL,
    UNIQUE (block_number, block_hash, bid_id, address_id, tic)
);

CREATE INDEX flop_bid_tic_block_number_index
    ON maker.flop_bid_tic (block_number);
CREATE INDEX flop_bid_tic_bid_id_index
    ON maker.flop_bid_tic (bid_id);
CREATE INDEX flop_bid_tic_bid_address_id_index
    ON maker.flop_bid_tic (address_id);

CREATE TABLE maker.flop_bid_end
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id       NUMERIC NOT NULL,
    "end"        BIGINT  NOT NULL,
    UNIQUE (block_number, block_hash, bid_id, address_id, "end")
);

CREATE INDEX flop_bid_end_block_number_index
    ON maker.flop_bid_end (block_number);
CREATE INDEX flop_bid_end_bid_id_index
    ON maker.flop_bid_end (bid_id);
CREATE INDEX flop_bid_end_bid_address_id_index
    ON maker.flop_bid_end (address_id);

CREATE TABLE maker.flop_vat
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    vat          TEXT,
    UNIQUE (block_number, block_hash, address_id, vat)
);

CREATE TABLE maker.flop_gem
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    gem          TEXT,
    UNIQUE (block_number, block_hash, address_id, gem)
);

CREATE TABLE maker.flop_beg
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    beg          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, address_id, beg)
);

CREATE TABLE maker.flop_pad
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    pad          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, address_id, pad)
);

CREATE TABLE maker.flop_ttl
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    ttl          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, address_id, ttl)
);

CREATE TABLE maker.flop_tau
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    tau          NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, address_id, tau)
);

CREATE TABLE maker.flop_kicks
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash   TEXT,
    address_id   INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    kicks        NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, address_id, kicks)
);

CREATE INDEX flop_kicks_block_number_index
    ON maker.flop_kicks (block_number);
CREATE INDEX flop_kicks_kicks_index
    ON maker.flop_kicks (kicks);
CREATE INDEX flop_kicks_address_id_index
    ON maker.flop_kicks (address_id);

-- prevent naming conflict with maker.flop_kick in postgraphile
COMMENT ON TABLE maker.flop_kicks IS E'@name flopKicksStorage';

CREATE TABLE maker.flop_live
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
DROP INDEX maker.flop_kicks_address_id_index;
DROP INDEX maker.flop_kicks_kicks_index;
DROP INDEX maker.flop_kicks_block_number_index;
DROP INDEX maker.flop_bid_end_block_number_index;
DROP INDEX maker.flop_bid_end_bid_address_id_index;
DROP INDEX maker.flop_bid_end_bid_id_index;
DROP INDEX maker.flop_bid_tic_block_number_index;
DROP INDEX maker.flop_bid_tic_bid_address_id_index;
DROP INDEX maker.flop_bid_tic_bid_id_index;
DROP INDEX maker.flop_bid_guy_block_number_index;
DROP INDEX maker.flop_bid_guy_bid_address_id_index;
DROP INDEX maker.flop_bid_guy_bid_id_index;
DROP INDEX maker.flop_bid_lot_block_number_index;
DROP INDEX maker.flop_bid_lot_bid_address_id_index;
DROP INDEX maker.flop_bid_lot_bid_id_index;
DROP INDEX maker.flop_bid_bid_block_number_index;
DROP INDEX maker.flop_bid_bid_address_id_index;
DROP INDEX maker.flop_bid_bid_bid_id_index;

DROP TABLE maker.flop_live;
DROP TABLE maker.flop_kicks;
DROP TABLE maker.flop_tau;
DROP TABLE maker.flop_ttl;
DROP TABLE maker.flop_beg;
DROP TABLE maker.flop_pad;
DROP TABLE maker.flop_gem;
DROP TABLE maker.flop_vat;
DROP TABLE maker.flop_bid_end;
DROP TABLE maker.flop_bid_tic;
DROP TABLE maker.flop_bid_guy;
DROP TABLE maker.flop_bid_lot;
DROP TABLE maker.flop_bid_bid;
