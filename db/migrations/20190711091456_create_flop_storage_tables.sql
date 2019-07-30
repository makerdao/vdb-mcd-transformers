-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.flop_bid_bid
(
    id SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash TEXT,
    contract_address TEXT,
    bid_id NUMERIC NOT NULL,
    bid NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, bid_id, contract_address, bid)
);

CREATE INDEX flop_bid_bid_block_number_index
    ON maker.flop_bid_bid (block_number);

CREATE TABLE maker.flop_bid_lot
(
    id SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash TEXT,
    contract_address TEXT,
    bid_id NUMERIC NOT NULL,
    lot NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, bid_id, contract_address, lot)
);

CREATE INDEX flop_bid_lot_block_number_index
    ON maker.flop_bid_lot (block_number);

CREATE TABLE maker.flop_bid_guy
(
    id SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash TEXT,
    contract_address TEXT,
    bid_id NUMERIC NOT NULL,
    guy TEXT,
    UNIQUE (block_number, block_hash, bid_id, contract_address, guy)
);

CREATE INDEX flop_bid_guy_block_number_index
    ON maker.flop_bid_guy (block_number);

CREATE TABLE maker.flop_bid_tic
(
    id SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash TEXT,
    contract_address TEXT,
    bid_id NUMERIC NOT NULL,
    tic BIGINT NOT NULL,
    UNIQUE (block_number, block_hash, bid_id, contract_address, tic)
);

CREATE INDEX flop_bid_tic_block_number_index
    ON maker.flop_bid_tic (block_number);

CREATE TABLE maker.flop_bid_end
(
    id SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash TEXT,
    contract_address TEXT,
    bid_id NUMERIC NOT NULL,
    "end" BIGINT NOT NULL,
    UNIQUE (block_number, block_hash, bid_id, contract_address, "end")
);

CREATE INDEX flop_bid_end_block_number_index
    ON maker.flop_bid_end (block_number);

CREATE TABLE maker.flop_vat
(
    id SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash TEXT,
    contract_address TEXT,
    vat TEXT,
    UNIQUE (block_number, block_hash, contract_address, vat)
);

CREATE TABLE maker.flop_gem
(
    id SERIAL PRIMARY KEY,
    block_number BIGINT,
    block_hash TEXT,
    contract_address TEXT,
    gem TEXT,
    UNIQUE (block_number, block_hash, contract_address, gem)
);

CREATE TABLE  maker.flop_beg
(
    id SERIAL PRIMARY KEY ,
    block_number BIGINT,
    block_hash TEXT,
    contract_address TEXT,
    beg NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, contract_address, beg)
);

CREATE TABLE  maker.flop_ttl
(
    id SERIAL PRIMARY KEY ,
    block_number BIGINT,
    block_hash TEXT,
    contract_address TEXT,
    ttl NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, contract_address, ttl)
);

CREATE TABLE  maker.flop_tau
(
    id SERIAL PRIMARY KEY ,
    block_number BIGINT,
    block_hash TEXT,
    contract_address TEXT,
    tau NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, contract_address, tau)
);

CREATE TABLE  maker.flop_kicks
(
    id SERIAL PRIMARY KEY ,
    block_number BIGINT,
    block_hash TEXT,
    contract_address TEXT,
    kicks NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, contract_address, kicks)
);

-- prevent naming conflict with maker.flop_kick in postgraphile
COMMENT ON TABLE maker.flop_kicks IS E'@name flopKicksStorage';

CREATE TABLE  maker.flop_live
(
    id SERIAL PRIMARY KEY ,
    block_number BIGINT,
    block_hash TEXT,
    contract_address TEXT,
    live NUMERIC NOT NULL,
    UNIQUE (block_number, block_hash, contract_address, live)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.flop_bid_end_block_number_index;
DROP INDEX maker.flop_bid_tic_block_number_index;
DROP INDEX maker.flop_bid_guy_block_number_index;
DROP INDEX maker.flop_bid_lot_block_number_index;
DROP INDEX maker.flop_bid_bid_block_number_index;

DROP TABLE maker.flop_live;
DROP TABLE maker.flop_kicks;
DROP TABLE maker.flop_tau;
DROP TABLE maker.flop_ttl;
DROP TABLE maker.flop_beg;
DROP TABLE maker.flop_gem;
DROP TABLE maker.flop_vat;
DROP TABLE maker.flop_bid_end;
DROP TABLE maker.flop_bid_tic;
DROP TABLE maker.flop_bid_guy;
DROP TABLE maker.flop_bid_lot;
DROP TABLE maker.flop_bid_bid;