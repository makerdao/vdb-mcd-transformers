-- +goose Up
CREATE TABLE maker.rely
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    usr        BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX rely_header_index
    ON maker.rely (header_id);
CREATE INDEX rely_log_index
    ON maker.rely (log_id);
CREATE INDEX rely_address_index
    ON maker.rely (address_id);
CREATE INDEX rely_msg_sender_index
    ON maker.rely (msg_sender);
CREATE INDEX rely_usr_index
    ON maker.rely (usr);


-- +goose Down
DROP INDEX maker.rely_header_index;
DROP INDEX maker.rely_log_index;
DROP INDEX maker.rely_address_index;
DROP INDEX maker.rely_msg_sender_index;
DROP INDEX maker.rely_usr_index;

DROP TABLE maker.rely;