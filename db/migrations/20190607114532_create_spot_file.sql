-- +goose Up
CREATE TABLE maker.spot_file_mat
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

CREATE TABLE maker.spot_file_pip
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    pip       TEXT,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
    ADD COLUMN spot_file_mat_checked INTEGER NOT NULL DEFAULT 0;

ALTER TABLE public.checked_headers
    ADD COLUMN spot_file_pip_checked INTEGER NOT NULL DEFAULT 0;


-- +goose Down
DROP TABLE maker.spot_file_mat;
DROP TABLE maker.spot_file_pip;
ALTER TABLE public.checked_headers
    DROP COLUMN spot_file_mat_checked;
ALTER TABLE public.checked_headers
    DROP COLUMN spot_file_pip_checked;

