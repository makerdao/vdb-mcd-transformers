-- +goose Up
CREATE TABLE maker.cat_file_chop_lump
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    tx_idx    INTEGER NOT NULL,
    log_idx   INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX cat_file_chop_lump_header_index
    ON maker.cat_file_chop_lump (header_id);

CREATE INDEX cat_file_chop_lump_ilk_index
    ON maker.cat_file_chop_lump (ilk_id);

CREATE TABLE maker.cat_file_flip
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    what      TEXT,
    flip      TEXT,
    tx_idx    INTEGER NOT NULL,
    log_idx   INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX cat_file_flip_header_index
    ON maker.cat_file_flip (header_id);

CREATE INDEX cat_file_flip_ilk_index
    ON maker.cat_file_flip (ilk_id);

CREATE TABLE maker.cat_file_vow
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    what      TEXT,
    data      TEXT,
    tx_idx    INTEGER NOT NULL,
    log_idx   INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX cat_file_vow_header_index
    ON maker.cat_file_vow (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN cat_file_chop_lump INTEGER NOT NULL DEFAULT 0;

ALTER TABLE public.checked_headers
    ADD COLUMN cat_file_flip INTEGER NOT NULL DEFAULT 0;

ALTER TABLE public.checked_headers
    ADD COLUMN cat_file_vow INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.cat_file_chop_lump_header_index;
DROP INDEX maker.cat_file_chop_lump_ilk_index;
DROP INDEX maker.cat_file_flip_header_index;
DROP INDEX maker.cat_file_flip_ilk_index;
DROP INDEX maker.cat_file_vow_header_index;

DROP TABLE maker.cat_file_chop_lump;
DROP TABLE maker.cat_file_flip;
DROP TABLE maker.cat_file_vow;

ALTER TABLE public.checked_headers
    DROP COLUMN cat_file_chop_lump;

ALTER TABLE public.checked_headers
    DROP COLUMN cat_file_flip;

ALTER TABLE public.checked_headers
    DROP COLUMN cat_file_vow;
