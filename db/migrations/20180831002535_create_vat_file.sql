-- +goose Up
CREATE TABLE maker.vat_file_ilk
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_file_ilk_header_index
    ON maker.vat_file_ilk (header_id);

CREATE INDEX vat_file_ilk_ilk_index
    ON maker.vat_file_ilk (ilk_id);

CREATE TABLE maker.vat_file_debt_ceiling
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_file_debt_ceiling_header_index
    ON maker.vat_file_debt_ceiling (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN vat_file_debt_ceiling INTEGER NOT NULL DEFAULT 0;

ALTER TABLE public.checked_headers
    ADD COLUMN vat_file_ilk INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.vat_file_ilk_header_index;
DROP INDEX maker.vat_file_ilk_ilk_index;
DROP INDEX maker.vat_file_debt_ceiling_header_index;

DROP TABLE maker.vat_file_ilk;
DROP TABLE maker.vat_file_debt_ceiling;

ALTER TABLE public.checked_headers
    DROP COLUMN vat_file_debt_ceiling;

ALTER TABLE public.checked_headers
    DROP COLUMN vat_file_ilk;
