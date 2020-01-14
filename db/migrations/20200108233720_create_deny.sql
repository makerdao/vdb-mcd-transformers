-- +goose Up
CREATE TABLE maker.vat_deny
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    usr        INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_deny_header_index
    ON maker.vat_deny (header_id);
CREATE INDEX vat_deny_log_index
    ON maker.vat_deny (log_id);
CREATE INDEX vat_deny_address_index
    ON maker.vat_deny (address_id);
CREATE INDEX vat_deny_usr_index
    ON maker.vat_deny (usr);


-- +goose Down
DROP INDEX maker.vat_deny_usr_index;
DROP INDEX maker.vat_deny_address_index;
DROP INDEX maker.vat_deny_log_index;
DROP INDEX maker.vat_deny_header_index;

DROP TABLE maker.vat_deny;
