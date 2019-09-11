-- +goose Up
CREATE TABLE maker.flop_kick
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC NOT NULL,
    bid        NUMERIC NOT NULL,
    gal        TEXT,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    tx_idx     INTEGER NOT NULL,
    log_idx    INTEGER NOT NULL,
    raw_log    JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

-- prevent naming conflict with maker.flop_kicks in postgraphile
COMMENT ON TABLE maker.flop_kick IS E'@name flopKickEvent';

CREATE INDEX flop_kick_header_index
    ON maker.flop_kick (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN flop_kick INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.flop_kick_header_index;
DROP TABLE maker.flop_kick;
ALTER TABLE public.checked_headers
    DROP COLUMN flop_kick;
