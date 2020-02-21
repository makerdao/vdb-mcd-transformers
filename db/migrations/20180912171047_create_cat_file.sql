-- +goose Up
CREATE TABLE maker.cat_file_chop_lump
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX cat_file_chop_lump_header_index
    ON maker.cat_file_chop_lump (header_id);
CREATE INDEX cat_file_chop_lump_log_index
    ON maker.cat_file_chop_lump (log_id);
CREATE INDEX cat_file_chop_lump_ilk_index
    ON maker.cat_file_chop_lump (ilk_id);

CREATE TABLE maker.cat_file_flip
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    what      TEXT,
    flip      TEXT,
    UNIQUE (header_id, log_id)
);

CREATE INDEX cat_file_flip_header_index
    ON maker.cat_file_flip (header_id);
CREATE INDEX cat_file_flip_log_index
    ON maker.cat_file_flip (log_id);
CREATE INDEX cat_file_flip_ilk_index
    ON maker.cat_file_flip (ilk_id);

CREATE TABLE maker.cat_file_vow
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    what      TEXT,
    data      TEXT,
    UNIQUE (header_id, log_id)
);

CREATE INDEX cat_file_vow_header_index
    ON maker.cat_file_vow (header_id);
CREATE INDEX cat_file_vow_log_index
    ON maker.cat_file_vow (log_id);


-- +goose Down
DROP INDEX maker.cat_file_chop_lump_header_index;
DROP INDEX maker.cat_file_chop_lump_log_index;
DROP INDEX maker.cat_file_chop_lump_ilk_index;
DROP INDEX maker.cat_file_flip_header_index;
DROP INDEX maker.cat_file_flip_log_index;
DROP INDEX maker.cat_file_flip_ilk_index;
DROP INDEX maker.cat_file_vow_header_index;
DROP INDEX maker.cat_file_vow_log_index;

DROP TABLE maker.cat_file_chop_lump;
DROP TABLE maker.cat_file_flip;
DROP TABLE maker.cat_file_vow;
