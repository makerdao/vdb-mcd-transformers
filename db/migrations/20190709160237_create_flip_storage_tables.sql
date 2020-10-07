-- +goose Up
CREATE TABLE maker.flip_bid_bid
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    bid        NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, bid_id, address_id, bid)
);

CREATE INDEX flip_bid_bid_header_id_index
    ON maker.flip_bid_bid (header_id);
CREATE INDEX flip_bid_bid_bid_id_index
    ON maker.flip_bid_bid (bid_id);
CREATE INDEX flip_bid_bid_address_index
    ON maker.flip_bid_bid (address_id);

CREATE TABLE maker.flip_bid_lot
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, bid_id, address_id, lot)
);

CREATE INDEX flip_bid_lot_header_id_index
    ON maker.flip_bid_lot (header_id);
CREATE INDEX flip_bid_lot_bid_id_index
    ON maker.flip_bid_lot (bid_id);
CREATE INDEX flip_bid_lot_address_index
    ON maker.flip_bid_lot (address_id);

CREATE TABLE maker.flip_bid_guy
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    guy        TEXT,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, bid_id, address_id, guy)
);

CREATE INDEX flip_bid_guy_header_id_index
    ON maker.flip_bid_guy (header_id);
CREATE INDEX flip_bid_guy_bid_id_index
    ON maker.flip_bid_guy (bid_id);
CREATE INDEX flip_bid_guy_address_index
    ON maker.flip_bid_guy (address_id);

CREATE TABLE maker.flip_bid_tic
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    tic        BIGINT  NOT NULL,
    bid_id     NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, bid_id, address_id, tic)
);

CREATE INDEX flip_bid_tic_header_id_index
    ON maker.flip_bid_tic (header_id);
CREATE INDEX flip_bid_tic_bid_id_index
    ON maker.flip_bid_tic (bid_id);
CREATE INDEX flip_bid_tic_address_index
    ON maker.flip_bid_tic (address_id);

CREATE TABLE maker.flip_bid_end
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    "end"      BIGINT  NOT NULL,
    bid_id     NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, bid_id, address_id, "end")
);

CREATE INDEX flip_bid_end_header_id_index
    ON maker.flip_bid_end (header_id);
CREATE INDEX flip_bid_end_bid_id_index
    ON maker.flip_bid_end (bid_id);
CREATE INDEX flip_bid_end_address_index
    ON maker.flip_bid_end (address_id);

CREATE TABLE maker.flip_bid_usr
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    usr        TEXT,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, bid_id, address_id, usr)
);

CREATE INDEX flip_bid_usr_header_id_index
    ON maker.flip_bid_usr (header_id);
CREATE INDEX flip_bid_usr_bid_id_index
    ON maker.flip_bid_usr (bid_id);
CREATE INDEX flip_bid_usr_address_index
    ON maker.flip_bid_usr (address_id);

CREATE TABLE maker.flip_bid_gal
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    gal        TEXT,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, bid_id, address_id, gal)
);

CREATE INDEX flip_bid_gal_header_id_index
    ON maker.flip_bid_gal (header_id);
CREATE INDEX flip_bid_gal_bid_id_index
    ON maker.flip_bid_gal (bid_id);
CREATE INDEX flip_bid_gal_address_index
    ON maker.flip_bid_gal (address_id);

CREATE TABLE maker.flip_bid_tab
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    tab        NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, bid_id, address_id, tab)
);

CREATE INDEX flip_bid_tab_header_id_index
    ON maker.flip_bid_tab (header_id);
CREATE INDEX flip_bid_tab_bid_id_index
    ON maker.flip_bid_tab (bid_id);
CREATE INDEX flip_bid_tab_address_index
    ON maker.flip_bid_tab (address_id);

CREATE TABLE maker.flip_vat
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    vat        TEXT,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, address_id, vat)
);

CREATE INDEX flip_vat_header_id_index
    ON maker.flip_vat (header_id);
CREATE INDEX flip_vat_address_index
    ON maker.flip_vat (address_id);

CREATE TABLE maker.flip_ilk
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, address_id, ilk_id)
);

CREATE INDEX flip_ilk_header_id_index
    ON maker.flip_ilk (header_id);
CREATE INDEX flip_ilk_ilk_id_index
    ON maker.flip_ilk (ilk_id);
CREATE INDEX flip_ilk_address_index
    ON maker.flip_ilk (address_id);

CREATE TABLE maker.flip_beg
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    beg        NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, address_id, beg)
);

CREATE INDEX flip_beg_header_id_index
    ON maker.flip_beg (header_id);
CREATE INDEX flip_beg_address_index
    ON maker.flip_beg (address_id);

CREATE TABLE maker.flip_ttl
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    ttl        NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, address_id, ttl)
);

CREATE INDEX flip_ttl_header_id_index
    ON maker.flip_ttl (header_id);
CREATE INDEX flip_ttl_address_index
    ON maker.flip_ttl (address_id);

CREATE TABLE maker.flip_tau
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    tau        NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, address_id, tau)
);

CREATE INDEX flip_tau_header_id_index
    ON maker.flip_tau (header_id);
CREATE INDEX flip_tau_address_index
    ON maker.flip_tau (address_id);

CREATE TABLE maker.flip_kicks
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    kicks      NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, address_id, kicks)
);

CREATE INDEX flip_kicks_header_id_index
    ON maker.flip_kicks (header_id);
CREATE INDEX flip_kicks_address_index
    ON maker.flip_kicks (address_id);

CREATE TABLE maker.flip_cat
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    cat        INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE ,
    UNIQUE (diff_id, header_id, address_id, cat)
);

CREATE INDEX flip_cat_header_id_index
    ON maker.flip_cat (header_id);
CREATE INDEX flip_cat_address_index
    ON maker.flip_cat (address_id);
CREATE INDEX flip_cat_cat_index
    ON maker.flip_cat (cat);

-- +goose Down
DROP TABLE maker.flip_kicks;
DROP TABLE maker.flip_tau;
DROP TABLE maker.flip_ttl;
DROP TABLE maker.flip_beg;
DROP TABLE maker.flip_ilk;
DROP TABLE maker.flip_vat;
DROP TABLE maker.flip_bid_tab;
DROP TABLE maker.flip_bid_gal;
DROP TABLE maker.flip_bid_usr;
DROP TABLE maker.flip_bid_end;
DROP TABLE maker.flip_bid_tic;
DROP TABLE maker.flip_bid_guy;
DROP TABLE maker.flip_bid_lot;
DROP TABLE maker.flip_bid_bid;
DROP TABLE maker.flip_cat;
