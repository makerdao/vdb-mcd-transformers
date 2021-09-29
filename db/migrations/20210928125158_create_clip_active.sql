-- +goose Up
CREATE TABLE maker.clip_active
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    active     NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, active)
);

CREATE TABLE maker.clip_active_sales
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    sale_id    NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, address_id, sale_id)
);

-- +goose Down
DROP TABLE maker.clip_active;
DROP TABLE maker.clip_active_sales;
