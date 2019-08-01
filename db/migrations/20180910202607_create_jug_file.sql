-- +goose Up
CREATE TABLE maker.jug_file_base
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX jug_file_base_header_index
    ON maker.jug_file_base (header_id);

CREATE TABLE maker.jug_file_ilk
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX jug_file_ilk_header_index
    ON maker.jug_file_ilk (header_id);

CREATE INDEX jug_file_ilk_ilk_index
    ON maker.jug_file_ilk (ilk_id);

CREATE TABLE maker.jug_file_vow
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    what      TEXT,
    data      TEXT,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX jug_file_vow_header_index
    ON maker.jug_file_vow (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN jug_file_base INTEGER NOT NULL DEFAULT 0;

ALTER TABLE public.checked_headers
    ADD COLUMN jug_file_ilk INTEGER NOT NULL DEFAULT 0;

ALTER TABLE public.checked_headers
    ADD COLUMN jug_file_vow INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.jug_file_base_header_index;
DROP INDEX maker.jug_file_ilk_header_index;
DROP INDEX maker.jug_file_ilk_ilk_index;
DROP INDEX maker.jug_file_vow_header_index;

DROP TABLE maker.jug_file_ilk;
DROP TABLE maker.jug_file_base;
DROP TABLE maker.jug_file_vow;

ALTER TABLE public.checked_headers
    DROP COLUMN jug_file_base;

ALTER TABLE public.checked_headers
    DROP COLUMN jug_file_ilk;

ALTER TABLE public.checked_headers
    DROP COLUMN jug_file_vow;
