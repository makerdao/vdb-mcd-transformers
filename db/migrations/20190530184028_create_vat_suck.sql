-- +goose Up
CREATE TABLE maker.vat_suck
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    u         TEXT,
    v         TEXT,
    rad       NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_suck_header_index
    ON maker.vat_suck (header_id);
CREATE INDEX vat_suck_log_index
    ON maker.vat_suck (log_id);


-- +goose Down
DROP INDEX maker.vat_suck_log_index;
DROP INDEX maker.vat_suck_header_index;
DROP TABLE maker.vat_suck;
