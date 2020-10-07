-- +goose Up
CREATE TABLE maker.jug_init
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    msg_sender BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX jug_init_log_index
    ON maker.jug_init (log_id);
CREATE INDEX jug_init_header_index
    ON maker.jug_init (header_id);
CREATE INDEX jug_init_msg_sender_index
    ON maker.jug_init (msg_sender);
CREATE INDEX jug_init_ilk_index
    ON maker.jug_init (ilk_id);


-- +goose Down
DROP TABLE maker.jug_init;
