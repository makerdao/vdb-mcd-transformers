-- +goose Up
CREATE TABLE maker.jug_init
(
    id        SERIAL PRIMARY KEY,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX jug_init_header_index
    ON maker.jug_init (header_id);

CREATE INDEX jug_init_ilk_index
    ON maker.jug_init (ilk_id);


-- +goose Down
DROP INDEX maker.jug_init_header_index;
DROP INDEX maker.jug_init_ilk_index;

DROP TABLE maker.jug_init;
