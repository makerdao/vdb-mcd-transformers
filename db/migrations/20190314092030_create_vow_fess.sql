-- +goose Up
CREATE TABLE maker.vow_fess
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    tab       NUMERIC NOT NULL,
    UNIQUE (header_id, log_id)
);

COMMENT ON TABLE maker.vow_fess
    IS E'Note event emitted when fess is called on Vow contract.';

CREATE INDEX vow_fess_header_index
    ON maker.vow_fess (header_id);
CREATE INDEX vow_fess_log_index
    ON maker.vow_fess (log_id);


-- +goose Down
DROP INDEX maker.vow_fess_log_index;
DROP INDEX maker.vow_fess_header_index;
DROP TABLE maker.vow_fess;
