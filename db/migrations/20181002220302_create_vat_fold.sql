-- +goose Up
CREATE TABLE maker.vat_fold
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    urn_id    INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
    rate      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_fold_header_index
    ON maker.vat_fold (header_id);

CREATE INDEX vat_fold_urn_index
    ON maker.vat_fold (urn_id);

ALTER TABLE public.checked_headers
    ADD COLUMN vat_fold INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.vat_fold_header_index;
DROP INDEX maker.vat_fold_urn_index;

DROP TABLE maker.vat_fold;

ALTER TABLE public.checked_headers
    DROP COLUMN vat_fold;
