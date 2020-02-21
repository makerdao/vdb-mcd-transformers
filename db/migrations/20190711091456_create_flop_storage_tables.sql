-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.flop_bid_bid
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    bid        NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, bid_id, address_id, bid)
);

CREATE INDEX flop_bid_bid_header_id_index
    ON maker.flop_bid_bid (header_id);
CREATE INDEX flop_bid_bid_bid_id_index
    ON maker.flop_bid_bid (bid_id);
CREATE INDEX flop_bid_bid_address_index
    ON maker.flop_bid_bid (address_id);

COMMENT ON TABLE maker.flop_bid_bid
    IS E'Value of a Bid\'s bid field on the Flop contract as of a block header.';

CREATE TABLE maker.flop_bid_lot
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, bid_id, address_id, lot)
);

CREATE INDEX flop_bid_lot_header_id_index
    ON maker.flop_bid_lot (header_id);
CREATE INDEX flop_bid_lot_bid_id_index
    ON maker.flop_bid_lot (bid_id);
CREATE INDEX flop_bid_lot_bid_address_index
    ON maker.flop_bid_lot (address_id);

COMMENT ON TABLE maker.flop_bid_lot
    IS E'Value of a Bid\'s lot field on the Flop contract as of a block header.';

CREATE TABLE maker.flop_bid_guy
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    guy        TEXT,
    UNIQUE (diff_id, header_id, bid_id, address_id, guy)
);

CREATE INDEX flop_bid_guy_header_id_index
    ON maker.flop_bid_guy (header_id);
CREATE INDEX flop_bid_guy_bid_id_index
    ON maker.flop_bid_guy (bid_id);
CREATE INDEX flop_bid_guy_bid_address_index
    ON maker.flop_bid_guy (address_id);

COMMENT ON TABLE maker.flop_bid_guy
    IS E'Value of a Bid\'s guy field on the Flop contract as of a block header.';

CREATE TABLE maker.flop_bid_tic
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    tic        BIGINT  NOT NULL,
    UNIQUE (diff_id, header_id, bid_id, address_id, tic)
);

CREATE INDEX flop_bid_tic_header_id_index
    ON maker.flop_bid_tic (header_id);
CREATE INDEX flop_bid_tic_bid_id_index
    ON maker.flop_bid_tic (bid_id);
CREATE INDEX flop_bid_tic_bid_address_index
    ON maker.flop_bid_tic (address_id);

COMMENT ON TABLE maker.flop_bid_tic
    IS E'Value of a Bid\'s tic field on the Flop contract as of a block header.';

CREATE TABLE maker.flop_bid_end
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    "end"      BIGINT  NOT NULL,
    UNIQUE (diff_id, header_id, bid_id, address_id, "end")
);

CREATE INDEX flop_bid_end_header_id_index
    ON maker.flop_bid_end (header_id);
CREATE INDEX flop_bid_end_bid_id_index
    ON maker.flop_bid_end (bid_id);
CREATE INDEX flop_bid_end_bid_address_index
    ON maker.flop_bid_end (address_id);

COMMENT ON TABLE maker.flop_bid_end
    IS E'Value of a Bid\'s end field on the Flop contract as of a block header.';

CREATE TABLE maker.flop_vat
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    vat        TEXT,
    UNIQUE (diff_id, header_id, address_id, vat)
);

CREATE INDEX flop_vat_header_id_index
    ON maker.flop_vat (header_id);
CREATE INDEX flop_vat_address_index
    ON maker.flop_vat (address_id);

COMMENT ON TABLE maker.flop_vat
    IS E'Value of the Flop contract\'s vat variable as of a block header.';

CREATE TABLE maker.flop_gem
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    gem        TEXT,
    UNIQUE (diff_id, header_id, address_id, gem)
);

CREATE INDEX flop_gem_header_id_index
    ON maker.flop_gem (header_id);
CREATE INDEX flop_gem_address_index
    ON maker.flop_gem (address_id);

COMMENT ON TABLE maker.flop_gem
    IS E'Value of the Flop contract\'s gem variable as of a block header.';

CREATE TABLE maker.flop_beg
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    beg        NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, beg)
);

CREATE INDEX flop_beg_header_id_index
    ON maker.flop_beg (header_id);
CREATE INDEX flop_beg_address_index
    ON maker.flop_beg (address_id);

COMMENT ON TABLE maker.flop_beg
    IS E'Value of the Flop contract\'s beg variable as of a block header.';

CREATE TABLE maker.flop_pad
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    pad        NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, pad)
);

CREATE INDEX flop_pad_header_id_index
    ON maker.flop_pad (header_id);
CREATE INDEX flop_pad_address_index
    ON maker.flop_pad (address_id);

COMMENT ON TABLE maker.flop_pad
    IS E'Value of the Flop contract\'s pad variable as of a block header.';

CREATE TABLE maker.flop_ttl
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    ttl        NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, ttl)
);

CREATE INDEX flop_ttl_header_id_index
    ON maker.flop_ttl (header_id);
CREATE INDEX flop_ttl_address_index
    ON maker.flop_ttl (address_id);

COMMENT ON TABLE maker.flop_ttl
    IS E'Value of the Flop contract\'s ttl variable as of a block header.';

CREATE TABLE maker.flop_tau
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    tau        NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, tau)
);

CREATE INDEX flop_tau_header_id_index
    ON maker.flop_tau (header_id);
CREATE INDEX flop_tau_address_index
    ON maker.flop_tau (address_id);

COMMENT ON TABLE maker.flop_tau
    IS E'Value of the Flop contract\'s tau variable as of a block header.';

CREATE TABLE maker.flop_kicks
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    kicks      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, kicks)
);

CREATE INDEX flop_kicks_header_id_index
    ON maker.flop_kicks (header_id);
CREATE INDEX flop_kicks_address_index
    ON maker.flop_kicks (address_id);

-- prevent naming conflict with maker.flop_kick in postgraphile
COMMENT ON TABLE maker.flop_kicks
    IS E'@name flopKicksStorage\nValue of the Flop contract\'s kicks variable as of a block header.';

CREATE TABLE maker.flop_live
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    live       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, live)
);

CREATE INDEX flop_live_header_id_index
    ON maker.flop_live (header_id);
CREATE INDEX flop_live_address_index
    ON maker.flop_live (address_id);

COMMENT ON TABLE maker.flop_live
    IS E'Value of the Flop contract\'s live variable as of a block header.';

CREATE TABLE maker.flop_vow
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    vow        TEXT,
    UNIQUE (diff_id, header_id, address_id, vow)
);

CREATE INDEX flop_vow_header_id_index
    ON maker.flop_vow (header_id);
CREATE INDEX flop_vow_address_index
    ON maker.flop_vow (address_id);

COMMENT ON TABLE maker.flop_vow
    IS E'Value of the Flop contract\'s vow variable as of a block header.';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.flop_vow_address_index;
DROP INDEX maker.flop_vow_header_id_index;
DROP INDEX maker.flop_live_address_index;
DROP INDEX maker.flop_live_header_id_index;
DROP INDEX maker.flop_kicks_address_index;
DROP INDEX maker.flop_kicks_header_id_index;
DROP INDEX maker.flop_tau_address_index;
DROP INDEX maker.flop_tau_header_id_index;
DROP INDEX maker.flop_ttl_address_index;
DROP INDEX maker.flop_ttl_header_id_index;
DROP INDEX maker.flop_pad_address_index;
DROP INDEX maker.flop_pad_header_id_index;
DROP INDEX maker.flop_beg_address_index;
DROP INDEX maker.flop_beg_header_id_index;
DROP INDEX maker.flop_gem_address_index;
DROP INDEX maker.flop_gem_header_id_index;
DROP INDEX maker.flop_vat_address_index;
DROP INDEX maker.flop_vat_header_id_index;
DROP INDEX maker.flop_bid_end_bid_address_index;
DROP INDEX maker.flop_bid_end_bid_id_index;
DROP INDEX maker.flop_bid_end_header_id_index;
DROP INDEX maker.flop_bid_tic_bid_address_index;
DROP INDEX maker.flop_bid_tic_bid_id_index;
DROP INDEX maker.flop_bid_tic_header_id_index;
DROP INDEX maker.flop_bid_guy_bid_address_index;
DROP INDEX maker.flop_bid_guy_bid_id_index;
DROP INDEX maker.flop_bid_guy_header_id_index;
DROP INDEX maker.flop_bid_lot_bid_address_index;
DROP INDEX maker.flop_bid_lot_bid_id_index;
DROP INDEX maker.flop_bid_lot_header_id_index;
DROP INDEX maker.flop_bid_bid_address_index;
DROP INDEX maker.flop_bid_bid_bid_id_index;
DROP INDEX maker.flop_bid_bid_header_id_index;

DROP TABLE maker.flop_vow;
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
