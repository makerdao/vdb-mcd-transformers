-- +goose Up
CREATE TABLE maker.vow_file_auction_address
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    what      TEXT,
    data      INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vow_file_auction_address_header_index
    ON maker.vow_file_auction_address (header_id);
CREATE INDEX vvow_file_auction_address_log_index
    ON maker.vow_file_auction_address (log_id);
CREATE INDEX vow_file_auction_address_msg_sender_index
    ON maker.vow_file_auction_address (msg_sender);
CREATE INDEX vow_file_auction_address_data_index
    ON maker.vow_file_auction_address (data);

-- +goose Down
DROP TABLE maker.vow_file_auction_address;
