-- +goose Up
CREATE TABLE maker.vat_fold
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    u         TEXT    NOT NULL,
    rate      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vat_fold_header_index
    ON maker.vat_fold (header_id);
CREATE INDEX vat_fold_log_index
    ON maker.vat_fold (log_id);
CREATE INDEX vat_fold_ilk_index
    ON maker.vat_fold (ilk_id);


-- +goose Down
DROP INDEX maker.vat_fold_header_index;
DROP INDEX maker.vat_fold_log_index;
DROP INDEX maker.vat_fold_ilk_index;

DROP TABLE maker.vat_fold;
