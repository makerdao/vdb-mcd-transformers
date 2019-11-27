-- +goose Up
CREATE TABLE maker.pot_cage
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX pot_cage_header_index
    ON maker.pot_cage (header_id);


-- +goose Down
DROP INDEX maker.pot_cage_header_index;

DROP TABLE maker.pot_cage;
