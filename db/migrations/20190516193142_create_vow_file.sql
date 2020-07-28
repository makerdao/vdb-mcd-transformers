-- +goose Up
CREATE TABLE maker.vow_file_auction_attributes
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vow_file_header_index
    ON maker.vow_file_auction_attributes (header_id);
CREATE INDEX vow_file_log_index
    ON maker.vow_file_auction_attributes (log_id);


-- +goose Down
DROP TABLE maker.vow_file_auction_attributes;
