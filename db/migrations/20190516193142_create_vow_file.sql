-- +goose Up
CREATE TABLE maker.vow_file
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vow_file_header_index
    ON maker.vow_file (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN vow_file INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.vow_file_header_index;

DROP TABLE maker.vow_file;

ALTER TABLE public.checked_headers
    DROP COLUMN vow_file;

