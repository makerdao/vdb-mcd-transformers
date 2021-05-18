-- +goose Up
CREATE TABLE maker.clip_ilk
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (diff_id, header_id, address_id, ilk_id)
);

CREATE INDEX clip_ilk_header_id_index
    ON maker.clip_ilk (header_id);
CREATE INDEX clip_ilk_ilk_id_index
    ON maker.clip_ilk (ilk_id);
CREATE INDEX clip_ilk_address_index
    ON maker.clip_ilk (address_id);

-- +goose Down
DROP TABLE maker.clip_ilk;
