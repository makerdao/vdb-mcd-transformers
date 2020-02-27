-- +goose Up
CREATE TABLE maker.val_void
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX val_void_header_index
    ON maker.val_void (header_id);
CREATE INDEX val_void_log_index
    ON maker.val_void (log_id);
CREATE INDEX val_void_address_index
    ON maker.val_void (address_id);
CREATE INDEX val_void_msg_sender_index
    ON maker.val_void (msg_sender);


-- +goose Down
DROP TABLE maker.val_void;
