-- +goose Up
CREATE TABLE maker.vat_move
(
    id        SERIAL PRIMARY KEY,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    src       TEXT    NOT NULL,
    dst       TEXT    NOT NULL,
    rad       NUMERIC NOT NULL,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_move_header_index
    ON maker.vat_move (header_id);
CREATE INDEX vat_move_log_index
    ON maker.vat_move (log_id);


-- +goose Down
DROP INDEX maker.vat_move_log_index;
DROP INDEX maker.vat_move_header_index;
DROP TABLE maker.vat_move;
