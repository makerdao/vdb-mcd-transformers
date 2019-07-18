-- +goose Up
CREATE TABLE maker.jug_init
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX jug_init_header_index
    ON maker.jug_init (header_id);

CREATE INDEX jug_init_ilk_index
    ON maker.jug_init (ilk_id);

ALTER TABLE public.checked_headers
    ADD COLUMN jug_init INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.jug_init_header_index;
DROP INDEX maker.jug_init_ilk_index;

DROP TABLE maker.jug_init;

ALTER TABLE public.checked_headers
    DROP COLUMN jug_init;
