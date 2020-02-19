-- +goose Up
CREATE TABLE maker.pot_exit
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    wad        NUMERIC,
    UNIQUE (header_id, log_id)
);

COMMENT ON TABLE maker.pot_exit
    IS E'Note event emitted when exit is called on Pot contract.';

CREATE INDEX pot_exit_header_index
    ON maker.pot_exit (header_id);
CREATE INDEX pot_exit_log_index
    ON maker.pot_exit (log_id);
CREATE INDEX pot_exit_msg_sender_index
    ON maker.pot_exit (msg_sender);


-- +goose Down
DROP INDEX maker.pot_exit_msg_sender_index;
DROP INDEX maker.pot_exit_log_index;
DROP INDEX maker.pot_exit_header_index;
DROP TABLE maker.pot_exit;
