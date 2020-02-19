-- +goose Up
CREATE TABLE maker.pot_cage
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

COMMENT ON TABLE maker.pot_cage
    IS E'Note event emitted when cage is called on Pot contract.';

CREATE INDEX pot_cage_header_index
    ON maker.pot_cage (header_id);
CREATE INDEX pot_cage_log_index
    ON maker.pot_cage (log_id);


-- +goose Down
DROP INDEX maker.pot_cage_log_index;
DROP INDEX maker.pot_cage_header_index;

DROP TABLE maker.pot_cage;
