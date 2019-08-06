-- +goose Up
CREATE TABLE maker.vat_heal
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    rad       NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_heal_header_index
    ON maker.vat_heal (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN vat_heal INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.vat_heal_header_index;
DROP TABLE maker.vat_heal;
ALTER TABLE public.checked_headers
    DROP COLUMN vat_heal;
