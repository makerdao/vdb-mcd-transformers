-- +goose Up
CREATE TABLE maker.vat_file_ilk
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_file_ilk_header_index
    ON maker.vat_file_ilk (header_id);
CREATE INDEX vat_file_ilk_log_index
    ON maker.vat_file_ilk (log_id);
CREATE INDEX vat_file_ilk_ilk_index
    ON maker.vat_file_ilk (ilk_id);

CREATE TABLE maker.vat_file_debt_ceiling
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_file_debt_ceiling_header_index
    ON maker.vat_file_debt_ceiling (header_id);
CREATE INDEX vat_file_debt_ceiling_log_index
    ON maker.vat_file_debt_ceiling (log_id);


-- +goose Down
DROP INDEX maker.vat_file_ilk_header_index;
DROP INDEX maker.vat_file_ilk_log_index;
DROP INDEX maker.vat_file_ilk_ilk_index;
DROP INDEX maker.vat_file_debt_ceiling_header_index;
DROP INDEX maker.vat_file_debt_ceiling_log_index;

DROP TABLE maker.vat_file_ilk;
DROP TABLE maker.vat_file_debt_ceiling;
