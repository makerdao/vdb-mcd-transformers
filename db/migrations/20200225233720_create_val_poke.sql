-- +goose Up
CREATE TABLE maker.val_poke
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    wut        TEXT,
    UNIQUE (header_id, log_id)
);

CREATE INDEX val_poke_header_index
    ON maker.val_poke (header_id);
CREATE INDEX val_poke_log_index
    ON maker.val_poke (log_id);
CREATE INDEX val_poke_address_index
    ON maker.val_poke (address_id);
CREATE INDEX val_poke_msg_sender_index
    ON maker.val_poke (msg_sender);


-- +goose Down
DROP TABLE maker.val_poke;
