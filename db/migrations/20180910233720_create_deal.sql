-- +goose Up
CREATE TABLE maker.deal
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

COMMENT ON TABLE maker.deal
    IS E'Note event emitted when deal is called on Flap, Flip, or Flop contracts.';

CREATE INDEX deal_header_index
    ON maker.deal (header_id);
CREATE INDEX deal_log_index
    ON maker.deal (log_id);
CREATE INDEX deal_bid_id_index
    ON maker.deal (bid_id);
CREATE INDEX deal_address_index
    ON maker.deal (address_id);


-- +goose Down
DROP INDEX maker.deal_address_index;
DROP INDEX maker.deal_bid_id_index;
DROP INDEX maker.deal_log_index;
DROP INDEX maker.deal_header_index;

DROP TABLE maker.deal;
