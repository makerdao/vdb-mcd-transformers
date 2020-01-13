-- +goose Up
CREATE TABLE maker.flap_kick
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC NOT NULL,
    bid        NUMERIC NOT NULL,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

-- prevent naming conflict with maker.flap_kicks in postgraphile
COMMENT ON TABLE maker.flap_kick IS E'@name flapKickEvent';

CREATE INDEX flap_kick_header_index
    ON maker.flap_kick (header_id);
CREATE INDEX flap_kick_log_index
    ON maker.flap_kick (log_id);
CREATE INDEX flap_kick_address_index
    ON maker.flap_kick (address_id);


-- +goose Down
DROP INDEX maker.flap_kick_address_index;
DROP INDEX maker.flap_kick_log_index;
DROP INDEX maker.flap_kick_header_index;
DROP TABLE maker.flap_kick;
