-- +goose Up
CREATE TABLE maker.vat_move
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    src       TEXT    NOT NULL,
    dst       TEXT    NOT NULL,
    rad       NUMERIC NOT NULL,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX vat_move_header_index
    ON maker.vat_move (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN vat_move INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.vat_move_header_index;
DROP TABLE maker.vat_move;
ALTER TABLE public.checked_headers
    DROP COLUMN vat_move;
