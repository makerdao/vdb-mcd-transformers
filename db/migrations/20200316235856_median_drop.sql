-- +goose Up
CREATE TABLE maker.median_drop
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    a_length   INTEGER NOT NULL,
    a          TEXT ARRAY NOT NULL,
    UNIQUE (header_id, log_id)
);

CREATE INDEX median_drop_log_index
    ON maker.median_drop (log_id);
CREATE INDEX median_drop_header_index
    ON maker.median_drop (header_id);
CREATE INDEX median_drop_address_index
    ON maker.median_drop (address_id);
CREATE INDEX median_drop_msg_sender_index
    ON maker.median_drop (msg_sender);
-- +goose Down

DROP TABLE maker.median_drop;
