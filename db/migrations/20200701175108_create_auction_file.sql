-- +goose Up
CREATE TABLE maker.auction_file
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    what       TEXT,
    data       NUMERIC,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX auction_file_header_index
    ON maker.auction_file (header_id);
CREATE INDEX auction_file_log_index
    ON maker.auction_file (log_id);
CREATE INDEX auction_file_address_id_index
    ON maker.auction_file (address_id);
CREATE INDEX auction_file_msg_sender_index
    ON maker.auction_file (msg_sender);

-- +goose Down
DROP INDEX maker.auction_file_header_index;
DROP INDEX maker.auction_file_log_index;
DROP INDEX maker.auction_file_address_id_index;
DROP INDEX maker.auction_file_msg_sender_index;

DROP TABLE maker.auction_file;
