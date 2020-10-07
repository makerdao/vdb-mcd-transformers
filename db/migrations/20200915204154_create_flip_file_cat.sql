-- +goose Up
CREATE TABLE maker.flip_file_cat
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    what       TEXT,
    data       INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX flip_file_cat_header_index
    ON maker.flip_file_cat (header_id);
CREATE INDEX flip_file_cat_log_index
    ON maker.flip_file_cat (log_id);
CREATE INDEX flip_file_cat_address_id_index
    ON maker.flip_file_cat (address_id);
CREATE INDEX flip_file_cat_msg_sender_index
    ON maker.flip_file_cat (msg_sender);
-- +goose Down
DROP TABLE maker.flip_file_cat;
