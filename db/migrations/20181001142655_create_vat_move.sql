-- +goose Up
CREATE TABLE maker.vat_move
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    src       TEXT    NOT NULL,
    dst       TEXT    NOT NULL,
    rad       NUMERIC NOT NULL,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_move_header_index
    ON maker.vat_move (header_id);


-- +goose Down
DROP INDEX maker.vat_move_header_index;
DROP TABLE maker.vat_move;
