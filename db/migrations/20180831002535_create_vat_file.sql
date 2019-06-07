-- +goose Up
CREATE TABLE maker.vat_file_ilk
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE TABLE maker.vat_file_debt_ceiling
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
    ADD COLUMN vat_file_debt_ceiling_checked BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE public.checked_headers
    ADD COLUMN vat_file_ilk_checked BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
DROP TABLE maker.vat_file_ilk;
DROP TABLE maker.vat_file_debt_ceiling;

ALTER TABLE public.checked_headers
    DROP COLUMN vat_file_debt_ceiling_checked;

ALTER TABLE public.checked_headers
    DROP COLUMN vat_file_ilk_checked;
