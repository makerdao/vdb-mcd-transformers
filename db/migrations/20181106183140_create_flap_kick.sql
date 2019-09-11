-- +goose Up
CREATE TABLE maker.flap_kick
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC NOT NULL,
    bid        NUMERIC NOT NULL,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    tx_idx     INTEGER NOT NULL,
    log_idx    INTEGER NOT NULL,
    raw_log    JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

-- prevent naming conflict with maker.flap_kicks in postgraphile
COMMENT ON TABLE maker.flap_kick IS E'@name flapKickEvent';

CREATE INDEX flap_kick_header_index
    ON maker.flap_kick (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN flap_kick INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.flap_kick_header_index;
DROP TABLE maker.flap_kick;
ALTER TABLE public.checked_headers
    DROP COLUMN flap_kick;
