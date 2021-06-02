-- +goose Up
CREATE TABLE maker.clip_dog
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    dog        BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, dog)
);

CREATE TABLE maker.clip_vow
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    vow        BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, vow)
);

-- +goose Down
DROP TABLE maker.clip_dog;
DROP TABLE maker.clip_vow;
