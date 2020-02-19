-- +goose Up
CREATE TABLE maker.tend
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

COMMENT ON TABLE maker.tend
    IS E'Note event emitted when tend is called on Flap or Flip contract.';

CREATE INDEX tend_header_index
    ON maker.tend (header_id);
CREATE INDEX tend_log_index
    ON maker.tend (log_id);
CREATE INDEX tend_address_index
    ON maker.tend (address_id);


-- +goose Down
DROP INDEX maker.tend_address_index;
DROP INDEX maker.tend_log_index;
DROP INDEX maker.tend_header_index;

DROP TABLE maker.tend;
