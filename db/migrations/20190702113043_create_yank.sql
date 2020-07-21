-- +goose Up
CREATE TABLE maker.yank
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    UNIQUE (header_id, log_id)
);

CREATE INDEX yank_header_index
    ON maker.yank (header_id);
CREATE INDEX yank_log_index
    ON maker.yank (log_id);
CREATE INDEX yank_address_index
    ON maker.yank (address_id);
CREATE INDEX yank_msg_sender
    ON maker.yank (msg_sender);
CREATE INDEX yank_bid_id_index
    ON maker.yank (bid_id);

-- +goose Down
DROP TABLE maker.yank;
