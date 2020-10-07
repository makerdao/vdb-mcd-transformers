-- +goose Up
CREATE TABLE maker.vat_rely
(
    id        SERIAL PRIMARY KEY,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    usr       BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_rely_header_index
    ON maker.vat_rely (header_id);
CREATE INDEX vat_rely_log_index
    ON maker.vat_rely (log_id);
CREATE INDEX vat_rely_usr_index
    ON maker.vat_rely (usr);


-- +goose Down
DROP INDEX maker.vat_rely_header_index;
DROP INDEX maker.vat_rely_log_index;
DROP INDEX maker.vat_rely_usr_index;

DROP TABLE maker.vat_rely;