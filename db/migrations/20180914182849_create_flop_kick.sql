-- +goose Up
CREATE TABLE maker.flop_kick
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC NOT NULL,
    bid        NUMERIC NOT NULL,
    gal        TEXT,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

-- prevent naming conflict with maker.flop_kicks in postgraphile
COMMENT ON TABLE maker.flop_kick IS E'@name flopKickEvent\nKick event emitted by Flop contract.';

CREATE INDEX flop_kick_header_index
    ON maker.flop_kick (header_id);
CREATE INDEX flop_kick_log_index
    ON maker.flop_kick (log_id);
CREATE INDEX flop_kick_address_index
    ON maker.flop_kick (address_id);


-- +goose Down
DROP INDEX maker.flop_kick_address_index;
DROP INDEX maker.flop_kick_log_index;
DROP INDEX maker.flop_kick_header_index;
DROP TABLE maker.flop_kick;
