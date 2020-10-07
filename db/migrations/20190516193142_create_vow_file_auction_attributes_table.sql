-- +goose Up
CREATE TABLE maker.vow_file_auction_attributes
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    msg_sender BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    what       TEXT,
    data       NUMERIC,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vow_file_auction_attributes_header_index
    ON maker.vow_file_auction_attributes (header_id);
CREATE INDEX vow_file_auction_attributes_log_index
    ON maker.vow_file_auction_attributes (log_id);
CREATE INDEX vow_file_auction_attributes_msg_sender_index
    ON maker.vow_file_auction_attributes (msg_sender);

-- +goose Down
DROP TABLE maker.vow_file_auction_attributes;
