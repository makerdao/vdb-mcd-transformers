-- +goose Up
CREATE TABLE maker.cat_claw
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    msg_sender BIGINT  NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    rad        NUMERIC,
    UNIQUE (header_id, log_id)
);
CREATE INDEX cat_claw_header_index
    ON maker.cat_claw (header_id);
CREATE INDEX cat_claw_address_index
    ON maker.cat_claw (address_id);
CREATE INDEX cat_claw_msg_sender_index
    ON maker.cat_claw (msg_sender);
CREATE INDEX cat_claw_log_index
    ON maker.cat_claw (log_id);
-- +goose Down
DROP TABLE maker.cat_claw;
