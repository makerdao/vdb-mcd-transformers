-- +goose Up
CREATE TABLE maker.vat_fork
(
    id        SERIAL PRIMARY KEY,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    src       TEXT,
    dst       TEXT,
    dink      NUMERIC,
    dart      NUMERIC,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_fork_header_index
    ON maker.vat_fork (header_id);
CREATE INDEX vat_fork_log_index
    ON maker.vat_fork (log_id);
CREATE INDEX vat_fork_ilk_index
    ON maker.vat_fork (ilk_id);


-- +goose Down
DROP INDEX maker.vat_fork_header_index;
DROP INDEX maker.vat_fork_log_index;
DROP INDEX maker.vat_fork_ilk_index;

DROP TABLE maker.vat_fork;
