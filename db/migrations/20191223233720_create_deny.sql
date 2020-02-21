-- +goose Up
CREATE TABLE maker.deny
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    usr        INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX deny_header_index
    ON maker.deny (header_id);
CREATE INDEX deny_log_index
    ON maker.deny (log_id);
CREATE INDEX deny_address_index
    ON maker.deny (address_id);
CREATE INDEX deny_msg_sender_index
    ON maker.deny (msg_sender);
CREATE INDEX deny_usr_index
    ON maker.deny (usr);


-- +goose Down
DROP INDEX maker.deny_usr_index;
DROP INDEX maker.deny_address_index;
DROP INDEX maker.deny_msg_sender_index;
DROP INDEX maker.deny_log_index;
DROP INDEX maker.deny_header_index;

DROP TABLE maker.deny;
