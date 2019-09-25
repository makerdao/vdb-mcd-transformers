-- +goose Up
CREATE TABLE maker.jug_drip
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX jug_drip_header_index
    ON maker.jug_drip (header_id);

CREATE INDEX jug_drip_ilk_index
    ON maker.jug_drip (ilk_id);


-- +goose Down
DROP INDEX maker.jug_drip_header_index;
DROP INDEX maker.jug_drip_ilk_index;

DROP TABLE maker.jug_drip;
