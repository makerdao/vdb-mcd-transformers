-- +goose Up
CREATE TABLE maker.vow_fess
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    tab       NUMERIC NOT NULL,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vow_fess_header_index
    ON maker.vow_fess (header_id);


-- +goose Down
DROP INDEX maker.vow_fess_header_index;
DROP TABLE maker.vow_fess;
