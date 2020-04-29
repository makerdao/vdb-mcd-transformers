-- +goose Up
CREATE TABLE maker.log_buy_enabled
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    is_enabled BOOLEAN,
    UNIQUE (header_id, log_id)
);

CREATE INDEX log_buy_enabled_header_index
    ON maker.log_buy_enabled (header_id);
CREATE INDEX log_buy_enabled_log_index
    ON maker.log_buy_enabled (log_id);
CREATE INDEX log_buy_enabled_address_index
    ON maker.log_buy_enabled (address_id);

-- +goose Down
DROP TABLE maker.log_buy_enabled;
