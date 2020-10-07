-- +goose Up
CREATE TABLE maker.jug_drip
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    msg_sender BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX jug_drip_header_index
    ON maker.jug_drip (header_id);
CREATE INDEX jug_drip_log_index
    ON maker.jug_drip (log_id);
CREATE INDEX jog_drip_msg_sender
    ON maker.jug_drip (msg_sender);
CREATE INDEX jug_drip_ilk_index
    ON maker.jug_drip (ilk_id);


-- +goose Down
DROP TABLE maker.jug_drip;
