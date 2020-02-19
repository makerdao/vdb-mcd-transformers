-- +goose Up
CREATE TABLE maker.jug_drip
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

COMMENT ON TABLE maker.jug_drip
    IS E'Note event emitted when drip is called on Jug contract.';

CREATE INDEX jug_drip_header_index
    ON maker.jug_drip (header_id);
CREATE INDEX jug_drip_log_index
    ON maker.jug_drip (log_id);
CREATE INDEX jug_drip_ilk_index
    ON maker.jug_drip (ilk_id);


-- +goose Down
DROP INDEX maker.jug_drip_header_index;
DROP INDEX maker.jug_drip_log_index;
DROP INDEX maker.jug_drip_ilk_index;

DROP TABLE maker.jug_drip;
