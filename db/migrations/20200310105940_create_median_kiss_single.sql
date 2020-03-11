-- +goose Up
CREATE TABLE maker.median_kiss_single
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    a          INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX median_kiss_single_log_index
    ON maker.median_kiss_single (log_id);
CREATE INDEX median_kiss_single_header_index
    ON maker.median_kiss_single (header_id);
CREATE INDEX median_kiss_single_address_index
    ON maker.median_kiss_single (address_id);
CREATE INDEX median_kiss_single_msg_sender_index
    ON maker.median_kiss_single (msg_sender);
CREATE INDEX median_kiss_single_a_index
    ON maker.median_kiss_single (a);


-- +goose Down
DROP TABLE maker.median_kiss_single;
