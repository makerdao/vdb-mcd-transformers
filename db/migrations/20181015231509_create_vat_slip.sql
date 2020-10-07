-- +goose Up
CREATE TABLE maker.vat_slip
(
    id        SERIAL PRIMARY KEY,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    usr       TEXT,
    wad       NUMERIC,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_slip_header_index
    ON maker.vat_slip (header_id);
CREATE INDEX vat_slip_log_index
    ON maker.vat_slip (log_id);
CREATE INDEX vat_slip_ilk_index
    ON maker.vat_slip (ilk_id);


-- +goose Down
DROP INDEX maker.vat_slip_header_index;
DROP INDEX maker.vat_slip_log_index;
DROP INDEX maker.vat_slip_ilk_index;

DROP TABLE maker.vat_slip;
