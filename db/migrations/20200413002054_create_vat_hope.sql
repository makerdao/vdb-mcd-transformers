-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.vat_hope
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    usr       INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_hope_header_index
    ON maker.vat_hope (header_id);
CREATE INDEX vat_hope_log_index
    ON maker.vat_hope (log_id);
CREATE INDEX vat_hope_usr_index
    ON maker.vat_hope (usr);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.vat_hope_header_index;
DROP INDEX maker.vat_hope_log_index;
DROP INDEX maker.vat_hope_usr_index;

DROP TABLE maker.vat_hope;