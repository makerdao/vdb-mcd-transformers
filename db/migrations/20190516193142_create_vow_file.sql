-- +goose Up
CREATE TABLE maker.vow_file
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    UNIQUE (header_id, log_id)
);

COMMENT ON TABLE maker.vow_file
    IS E'Note event emitted when file(bytes32,uint256) is called on Vow contract.';

CREATE INDEX vow_file_header_index
    ON maker.vow_file (header_id);
CREATE INDEX vow_file_log_index
    ON maker.vow_file (log_id);


-- +goose Down
DROP INDEX maker.vow_file_log_index;
DROP INDEX maker.vow_file_header_index;

DROP TABLE maker.vow_file;
