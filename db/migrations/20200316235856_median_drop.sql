-- +goose Up
CREATE TABLE maker.median_drop
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    a0        INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    a1        INTEGER DEFAULT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    a2        INTEGER DEFAULT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    a3        INTEGER DEFAULT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
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
CREATE INDEX median_drop_a0_index
    ON maker.median_drop (a0);
CREATE INDEX median_drop_a1_index
    ON maker.median_drop (a1);
CREATE INDEX median_drop_a2_index
    ON maker.median_drop (a2);
CREATE INDEX median_drop_a3_index
    ON maker.median_drop (a3);
-- +goose Down

DROP TABLE maker.median_drop;
