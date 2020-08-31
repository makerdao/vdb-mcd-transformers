-- +goose Up
CREATE TABLE maker.cat_file_box
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

CREATE INDEX cat_file_box_log_index
    ON maker.cat_file_box (log_id);
CREATE INDEX cat_file_box_address_index
    ON maker.cat_file_box (address_id);
CREATE INDEX cat_file_box_msg_sender
    ON maker.cat_file_box (msg_sender);
CREATE INDEX cat_file_box_header_index
    ON maker.cat_file_box (header_id);


-- +goose Down
DROP TABLE maker.cat_file_box;
