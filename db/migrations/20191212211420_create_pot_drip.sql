-- +goose Up
CREATE TABLE maker.pot_drip
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

COMMENT ON TABLE maker.pot_drip
    IS E'Note event emitted when drip is called on Pot contract.';

CREATE INDEX pot_drip_header_index
    ON maker.pot_drip (header_id);
CREATE INDEX pot_drip_log_index
    ON maker.pot_drip (log_id);
CREATE INDEX pot_drip_msg_sender_index
    ON maker.pot_drip (msg_sender);

-- +goose Down
DROP INDEX maker.pot_drip_header_index;
DROP INDEX maker.pot_drip_log_index;
DROP INDEX maker.pot_drip_msg_sender_index;

DROP TABLE maker.pot_drip;