-- +goose Up
CREATE TABLE maker.cat_file_chop_dunk
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    msg_sender BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    what       TEXT,
    data       NUMERIC,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX cat_file_chop_dunk_header_index
    ON maker.cat_file_chop_dunk (header_id);
CREATE INDEX cat_file_chop_dunk_log_index
    ON maker.cat_file_chop_dunk (log_id);
CREATE INDEX cat_file_chop_dunk_ilk_index
    ON maker.cat_file_chop_dunk (ilk_id);

-- +goose Down
DROP TABLE maker.cat_file_chop_dunk;
