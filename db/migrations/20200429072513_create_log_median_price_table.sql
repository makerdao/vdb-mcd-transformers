-- +goose Up
CREATE TABLE maker.log_median_price
    (
        id serial primary key,
        header_id integer NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
        log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
        val numeric,
        age numeric,
        UNIQUE (header_id, log_id)
);

-- +goose Down
DROP TABLE maker.log_median_price;