-- +goose Up
CREATE TABLE maker.vat_frob
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    urn_id    INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
    v         TEXT,
    w         TEXT,
    dink      NUMERIC,
    dart      NUMERIC,
    UNIQUE (header_id, log_id)
);

COMMENT ON TABLE maker.vat_frob
    IS E'Note event emitted when frob is called on Vat contract.';

CREATE INDEX vat_frob_header_index
    ON maker.vat_frob (header_id);
CREATE INDEX vat_frob_log_index
    ON maker.vat_frob (log_id);
CREATE INDEX vat_frob_urn_index
    ON maker.vat_frob (urn_id);


-- +goose Down
DROP INDEX maker.vat_frob_header_index;
DROP INDEX maker.vat_frob_log_index;
DROP INDEX maker.vat_frob_urn_index;

DROP TABLE maker.vat_frob;
