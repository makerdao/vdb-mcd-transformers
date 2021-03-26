-- +goose Up
CREATE TABLE maker.clip_take
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    clip_id    NUMERIC,
    max        NUMERIC,
    price        NUMERIC,
    owe        NUMERIC,
    tab        NUMERIC,
    lot        NUMERIC,
    usr        BIGINT,
    UNIQUE (header_id, log_id)
);

CREATE INDEX clip_take_header_index
    ON maker.clip_take (header_id);
CREATE INDEX clip_take_clip_id_index
    ON maker.clip_take (clip_id);
CREATE INDEX clip_take_address_index
    ON maker.clip_take (address_id);
CREATE INDEX clip_take_log_index
    ON maker.clip_take (log_id);

-- +goose Down
DROP INDEX maker.clip_take_log_index;
DROP INDEX maker.clip_take_address_index;
DROP INDEX maker.clip_take_clip_id_index;
DROP INDEX maker.clip_take_header_index;

DROP TABLE maker.clip_take;
