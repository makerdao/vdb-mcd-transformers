-- +goose Up
CREATE TABLE maker.tick
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX tick_header_index
    ON maker.tick (header_id);
CREATE INDEX tick_log_index
    ON maker.tick (log_id);
CREATE INDEX tick_bid_id_index
    ON maker.tick (bid_id);
CREATE INDEX tick_address_index
    ON maker.tick (address_id);
CREATE INDEX tick_msg_sender_index
    ON maker.tick (msg_sender);

-- +goose Down
DROP TABLE maker.tick;
