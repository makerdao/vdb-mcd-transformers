-- +goose Up
CREATE TABLE maker.vat_heal
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    rad       NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_heal_header_index
    ON maker.vat_heal (header_id);
CREATE INDEX vat_heal_log_index
    ON maker.vat_heal (log_id);


-- +goose Down
DROP INDEX maker.vat_heal_log_index;
DROP INDEX maker.vat_heal_header_index;
DROP TABLE maker.vat_heal;
