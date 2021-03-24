-- +goose Up
CREATE TABLE maker.dog_vat
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    vat        BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, vat)
);

CREATE INDEX dog_vat_header_index
    ON maker.dog_vat(header_id);
CREATE INDEX dog_vat_address_index
    ON maker.dog_vat(address_id);

CREATE TABLE maker.dog_vow
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    vow        BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, vow)
);

CREATE INDEX dog_vow_header_index
    ON maker.dog_vow(header_id);
CREATE INDEX dog_vow_address_index
    ON maker.dog_vow(address_id);

CREATE TABLE maker.dog_live
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    live       NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, live)
);

CREATE INDEX dog_live_header_index
    ON maker.dog_live(header_id);
CREATE INDEX dog_live_address_index
    ON maker.dog_live(address_id);

CREATE TABLE maker.dog_hole
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    hole       NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, hole)
);

CREATE INDEX dog_hole_header_index
    ON maker.dog_hole(header_id);
CREATE INDEX dog_hole_address_index
    ON maker.dog_hole(address_id);

CREATE TABLE maker.dog_dirt
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    dirt       NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, dirt)
);

CREATE INDEX dog_dirt_header_index
    ON maker.dog_dirt(header_id);
CREATE INDEX dog_dirt_address_index
    ON maker.dog_dirt(address_id);

-- +goose Down
DROP TABLE maker.dog_vat;
DROP TABLE maker.dog_vow;
DROP TABLE maker.dog_live;
DROP TABLE maker.dog_hole;
DROP TABLE maker.dog_dirt;
