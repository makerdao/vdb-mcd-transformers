-- +goose Up
CREATE TABLE maker.log_median_price
(
    id         SERIAL PRIMARY KEY,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    val        NUMERIC,
    age        NUMERIC,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
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