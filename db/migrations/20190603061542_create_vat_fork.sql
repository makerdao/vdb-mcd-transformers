-- +goose Up
CREATE TABLE maker.vat_fork
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    src       TEXT,
    dst       TEXT,
    dink      NUMERIC,
    dart      NUMERIC,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX vat_fork_header_index
    ON maker.vat_fork (header_id);

CREATE INDEX vat_fork_ilk_index
    ON maker.vat_fork (ilk_id);

ALTER TABLE public.checked_headers
    ADD COLUMN vat_fork INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.vat_fork_header_index;
DROP INDEX maker.vat_fork_ilk_index;

DROP TABLE maker.vat_fork;

ALTER TABLE public.checked_headers
    DROP COLUMN vat_fork;
