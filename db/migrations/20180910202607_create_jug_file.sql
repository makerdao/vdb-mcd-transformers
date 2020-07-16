-- +goose Up
CREATE TABLE maker.jug_file_base
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    what       TEXT,
    data       NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX jug_file_base_header_index
    ON maker.jug_file_base (header_id);
CREATE INDEX jug_file_base_log_index
    ON maker.jug_file_base (log_id);

CREATE TABLE maker.jug_file_ilk
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX jug_file_ilk_header_index
    ON maker.jug_file_ilk (header_id);
CREATE INDEX jug_file_ilk_log_index
    ON maker.jug_file_ilk (log_id);
CREATE INDEX jug_file_ilk_ilk_index
    ON maker.jug_file_ilk (ilk_id);

CREATE TABLE maker.jug_file_vow
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    what      TEXT,
    data      TEXT,
    UNIQUE (header_id, log_id)
);

CREATE INDEX jug_file_vow_header_index
    ON maker.jug_file_vow (header_id);
CREATE INDEX jug_file_vow_log_index
    ON maker.jug_file_vow (log_id);


-- +goose Down
DROP TABLE maker.jug_file_ilk;
DROP TABLE maker.jug_file_base;
DROP TABLE maker.jug_file_vow;
