-- +goose Up
CREATE TABLE maker.log_median_price
    (
        id serial primary key,
        header_id integer NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
        address_id integer not null references public.addresses (id) on delete cascade,
        log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
        val numeric,
        age numeric,
        UNIQUE (header_id, log_id)
);

CREATE INDEX log_median_price_header_index
    ON maker.log_median_price (header_id);
CREATE INDEX log_median_price_log_index
    ON maker.log_median_price (log_id);
CREATE INDEX log_median_price_address_index
    ON maker.log_median_price (address_id);

-- +goose Down
DROP TABLE maker.log_median_price;