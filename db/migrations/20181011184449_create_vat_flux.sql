-- +goose Up
CREATE TABLE maker.vat_flux
(
    id        SERIAL PRIMARY KEY,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    src       TEXT,
    dst       TEXT,
    wad       NUMERIC,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_flux_header_index
    ON maker.vat_flux (header_id);
CREATE INDEX vat_flux_log_index
    ON maker.vat_flux (log_id);
CREATE INDEX vat_flux_ilk_index
    ON maker.vat_flux (ilk_id);

-- +goose Down
DROP INDEX maker.vat_flux_header_index;
DROP INDEX maker.vat_flux_log_index;
DROP INDEX maker.vat_flux_ilk_index;

DROP TABLE maker.vat_flux;
