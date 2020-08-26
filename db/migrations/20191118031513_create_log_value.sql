-- +goose Up
CREATE TABLE maker.log_value
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    val        NUMERIC,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX log_value_header_index
    ON maker.log_value (header_id);
CREATE INDEX log_value_log_index
    ON maker.log_value (log_id);
CREATE INDEX log_value_address_index
    ON maker.log_value (address_id);


-- +goose Down
DROP TABLE maker.log_value;
