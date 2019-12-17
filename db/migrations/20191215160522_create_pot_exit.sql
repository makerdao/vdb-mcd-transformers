-- +goose Up
CREATE TABLE maker.pot_exit
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    wad       NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX pot_exit_header_index
    ON maker.pot_exit (header_id);
CREATE INDEX pot_exit_log_index
    ON maker.pot_exit (log_id);


-- +goose Down
DROP INDEX maker.pot_exit_log_index;
DROP INDEX maker.pot_exit_header_index;
DROP TABLE maker.pot_exit;
