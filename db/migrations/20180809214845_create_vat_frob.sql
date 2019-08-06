-- +goose Up
CREATE TABLE maker.vat_frob
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    urn_id    INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
    v         TEXT,
    w         TEXT,
    dink      NUMERIC,
    dart      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_frob_header_index
    ON maker.vat_frob (header_id);

CREATE INDEX vat_frob_urn_index
    ON maker.vat_frob (urn_id);

ALTER TABLE public.checked_headers
    ADD COLUMN vat_frob INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.vat_frob_header_index;
DROP INDEX maker.vat_frob_urn_index;

DROP TABLE maker.vat_frob;

ALTER TABLE public.checked_headers
    DROP COLUMN vat_frob;
