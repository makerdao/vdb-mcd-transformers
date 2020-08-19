-- +goose Up
CREATE TABLE maker.pot_cage
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX pot_cage_header_index
    ON maker.pot_cage (header_id);
CREATE INDEX pot_cage_log_index
    ON maker.pot_cage (log_id);
CREATE INDEX pot_cage_msg_sender_index
    ON maker.pot_cage (msg_sender);


-- +goose Down
DROP TABLE maker.pot_cage;
