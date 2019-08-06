-- +goose Up
CREATE TABLE maker.flop_kick
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC NOT NULL,
    bid        NUMERIC NOT NULL,
    gal        TEXT,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
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
