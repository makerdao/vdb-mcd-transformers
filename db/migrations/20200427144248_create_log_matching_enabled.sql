-- +goose Up
CREATE TABLE maker.log_matching_enabled
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    is_enabled BOOLEAN,
    UNIQUE (header_id, log_id)
);

CREATE INDEX log_matching_enabled_header_index
    ON maker.log_matching_enabled (header_id);
CREATE INDEX log_matching_enabled_log_index
    ON maker.log_matching_enabled (log_id);
CREATE INDEX log_matching_enabled_address_index
    ON maker.log_matching_enabled (address_id);

-- +goose Down
DROP TABLE maker.log_matching_enabled;