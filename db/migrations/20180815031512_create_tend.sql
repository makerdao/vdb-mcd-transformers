-- +goose Up
CREATE TABLE maker.tend
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC,
    bid        NUMERIC,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX tend_header_index
    ON maker.tend (header_id);
CREATE INDEX tend_log_index
    ON maker.tend (log_id);
CREATE INDEX tend_address_index
    ON maker.tend (address_id);
CREATE INDEX tend_msg_sender_index
    ON maker.tend (msg_sender);


-- +goose Down
DROP TABLE maker.tend;
