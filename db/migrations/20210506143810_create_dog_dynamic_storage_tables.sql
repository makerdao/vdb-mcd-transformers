-- +goose Up
CREATE TABLE maker.dog_ilk_clip
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    clip       TEXT,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, ilk_id, clip)
);

CREATE INDEX dog_ilk_clip_header_id_index
    ON maker.dog_ilk_clip (header_id);
CREATE INDEX dog_ilk_clip_ilk_index
    ON maker.dog_ilk_clip (ilk_id);
CREATE INDEX dog_ilk_clip_address_id_index
    ON maker.dog_ilk_clip (address_id);

CREATE TABLE maker.dog_ilk_chop
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    chop       NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, ilk_id, chop)
);

CREATE INDEX dog_ilk_chop_header_id_index
    ON maker.dog_ilk_chop (header_id);
CREATE INDEX dog_ilk_chop_ilk_index
    ON maker.dog_ilk_chop (ilk_id);
CREATE INDEX dog_ilk_chop_address_id_index
    ON maker.dog_ilk_chop (address_id);

CREATE TABLE maker.dog_ilk_hole
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    hole       NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, ilk_id, hole)
);

CREATE INDEX dog_ilk_hole_header_id_index
    ON maker.dog_ilk_hole (header_id);
CREATE INDEX dog_ilk_hole_ilk_index
    ON maker.dog_ilk_hole (ilk_id);
CREATE INDEX dog_ilk_hole_address_id_index
    ON maker.dog_ilk_hole (address_id);

CREATE TABLE maker.dog_ilk_dirt
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    dirt       NUMERIC NOT NULL,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, ilk_id, dirt)
);

CREATE INDEX dog_ilk_dirt_header_id_index
    ON maker.dog_ilk_dirt (header_id);
CREATE INDEX dog_ilk_dirt_ilk_index
    ON maker.dog_ilk_dirt (ilk_id);
CREATE INDEX dog_ilk_dirt_address_id_index
    ON maker.dog_ilk_dirt (address_id);

-- +goose Down

DROP TABLE maker.dog_ilk_clip;
DROP TABLE maker.dog_ilk_chop;
DROP TABLE maker.dog_ilk_hole;
DROP TABLE maker.dog_ilk_dirt;
