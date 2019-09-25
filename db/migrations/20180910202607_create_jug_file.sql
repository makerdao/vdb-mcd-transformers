-- +goose Up
CREATE TABLE maker.jug_file_base
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX jug_file_base_header_index
    ON maker.jug_file_base (header_id);

CREATE TABLE maker.jug_file_ilk
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX jug_file_ilk_header_index
    ON maker.jug_file_ilk (header_id);

CREATE INDEX jug_file_ilk_ilk_index
    ON maker.jug_file_ilk (ilk_id);

CREATE TABLE maker.jug_file_vow
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    what      TEXT,
    data      TEXT,
    UNIQUE (header_id, log_id)
);

CREATE INDEX jug_file_vow_header_index
    ON maker.jug_file_vow (header_id);


-- +goose Down
DROP INDEX maker.jug_file_base_header_index;
DROP INDEX maker.jug_file_ilk_header_index;
DROP INDEX maker.jug_file_ilk_ilk_index;
DROP INDEX maker.jug_file_vow_header_index;

DROP TABLE maker.jug_file_ilk;
DROP TABLE maker.jug_file_base;
DROP TABLE maker.jug_file_vow;
