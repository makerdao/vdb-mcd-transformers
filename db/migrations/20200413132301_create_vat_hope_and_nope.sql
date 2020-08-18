-- +goose Up
CREATE TABLE maker.vat_hope
(
    id        SERIAL PRIMARY KEY,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    usr       BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_hope_header_index
    ON maker.vat_hope (header_id);
CREATE INDEX vat_hope_log_index
    ON maker.vat_hope (log_id);
CREATE INDEX vat_hope_usr_index
    ON maker.vat_hope (usr);

CREATE TABLE maker.vat_nope
(
    id        SERIAL PRIMARY KEY,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    usr       BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_nope_header_index
    ON maker.vat_nope (header_id);
CREATE INDEX vat_nope_log_index
    ON maker.vat_nope (log_id);
CREATE INDEX vat_nope_usr_index
    ON maker.vat_nope (usr);

-- +goose Down
DROP INDEX maker.vat_hope_header_index;
DROP INDEX maker.vat_hope_log_index;
DROP INDEX maker.vat_hope_usr_index;
DROP INDEX maker.vat_nope_header_index;
DROP INDEX maker.vat_nope_log_index;
DROP INDEX maker.vat_nope_usr_index;

DROP TABLE maker.vat_hope;
DROP TABLE maker.vat_nope;
