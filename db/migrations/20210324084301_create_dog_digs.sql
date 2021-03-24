-- +goose Up
CREATE TABLE maker.dog_digs
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    rad        NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX dog_digs_log_index
    ON maker.dog_digs (log_id);
CREATE INDEX dog_digs_header_index
    ON maker.dog_digs (header_id);
CREATE INDEX dog_digs_address_index
    ON maker.dog_digs (address_id);
CREATE INDEX dog_digs_ilk_index
    ON maker.dog_digs (ilk_id);

-- +goose Down
DROP TABLE maker.dog_digs;