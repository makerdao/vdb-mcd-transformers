-- +goose Up
CREATE TABLE maker.dent
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC,
    bid        NUMERIC,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

COMMENT ON TABLE maker.dent
    IS E'Note event emitted when dent is called on Flip or Flop contract.';

CREATE INDEX dent_header_index
    ON maker.dent (header_id);
CREATE INDEX dent_log_index
    ON maker.dent (log_id);
CREATE INDEX dent_address_index
    ON maker.dent (address_id);


-- +goose Down
DROP INDEX maker.dent_header_index;
DROP INDEX maker.dent_log_index;
DROP INDEX maker.dent_address_index;

DROP TABLE maker.dent;
