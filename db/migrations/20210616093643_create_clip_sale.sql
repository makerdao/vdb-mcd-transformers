-- +goose Up
CREATE TABLE maker.clip_sale_pos
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    sale_id    NUMERIC NOT NULL,
    pos        NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, sale_id, address_id, pos)
);

CREATE TABLE maker.clip_sale_tab
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    sale_id    NUMERIC NOT NULL,
    tab        NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, sale_id, address_id, tab)
);

CREATE TABLE maker.clip_sale_lot
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    sale_id    NUMERIC NOT NULL,
    lot        NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, sale_id, address_id, lot)
);

CREATE TABLE maker.clip_sale_usr
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    sale_id    NUMERIC NOT NULL,
    usr        TEXT,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, sale_id, address_id, usr)
);

CREATE TABLE maker.clip_sale_tic
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    sale_id    NUMERIC NOT NULL,
    tic        NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, sale_id, address_id, tic)
);

CREATE TABLE maker.clip_sale_top
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    sale_id    NUMERIC NOT NULL,
    top        NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, sale_id, address_id, top)
);

-- +goose Down
DROP TABLE maker.clip_sale_pos;
DROP TABLE maker.clip_sale_tab;
DROP TABLE maker.clip_sale_lot;
DROP TABLE maker.clip_sale_usr;
DROP TABLE maker.clip_sale_tic;
DROP TABLE maker.clip_sale_top;
