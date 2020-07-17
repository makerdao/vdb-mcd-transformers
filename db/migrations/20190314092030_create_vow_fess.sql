-- +goose Up
CREATE TABLE maker.vow_fess
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    tab        NUMERIC NOT NULL,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vow_fess_header_index
    ON maker.vow_fess (header_id);
CREATE INDEX vow_fess_log_index
    ON maker.vow_fess (log_id);
CREATE INDEX vow_fess_msg_sender_index
    ON maker.vow_fess (msg_sender);

-- +goose Down
DROP TABLE maker.vow_fess;
