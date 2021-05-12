-- +goose Up
CREATE TABLE maker.dog_deny
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    usr        BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX dog_deny_header_index
    ON maker.dog_deny (header_id);
CREATE INDEX dog_deny_log_index
    ON maker.dog_deny (log_id);
CREATE INDEX dog_deny_address_index
    ON maker.dog_deny (address_id);
CREATE INDEX dog_deny_usr_index
    ON maker.dog_deny (usr);

-- +goose Down
DROP INDEX maker.dog_deny_header_index;
DROP INDEX maker.dog_deny_log_index;
DROP INDEX maker.dog_deny_address_index;
DROP INDEX maker.dog_deny_usr_index;

DROP TABLE maker.dog_deny;
