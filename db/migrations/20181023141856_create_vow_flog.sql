-- +goose Up
CREATE TABLE maker.vow_flog
(
    id          SERIAL PRIMARY KEY,
    header_id   INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    msg_sender  INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    log_id      BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    era         INTEGER NOT NULL,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vow_flog_era_index
    ON maker.vow_flog (era);
CREATE INDEX vow_flog_log_index
    ON maker.vow_flog (log_id);
CREATE INDEX vow_flog_header_index
    ON maker.vow_flog (header_id);
CREATE INDEX vow_flog_msg_sender_index
    ON maker.vow_flog (msg_sender);


-- +goose Down
DROP INDEX maker.vow_flog_era_index;
DROP INDEX maker.vow_flog_log_index;
DROP INDEX maker.vow_flog_header_index;

DROP TABLE maker.vow_flog;
