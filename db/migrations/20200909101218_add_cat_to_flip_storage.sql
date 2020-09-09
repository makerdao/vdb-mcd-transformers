-- +goose Up
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
DROP TABLE maker.flip_cat;
