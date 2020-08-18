-- +goose Up
CREATE TABLE maker.dent
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC,
    bid        NUMERIC,
    msg_sender BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX dent_header_index
    ON maker.dent (header_id);
CREATE INDEX dent_log_index
    ON maker.dent (log_id);
CREATE INDEX dent_msg_sender_index
    ON maker.dent (msg_sender);
CREATE INDEX dent_address_index
    ON maker.dent (address_id);


-- +goose Down
DROP TABLE maker.dent;
