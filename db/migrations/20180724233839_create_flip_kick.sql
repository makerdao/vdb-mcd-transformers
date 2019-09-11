-- +goose Up
CREATE TABLE maker.flip_kick
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC,
    bid        NUMERIC,
    tab        NUMERIC,
    usr        TEXT,
    gal        TEXT,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    tx_idx     INTEGER NOT NULL,
    log_idx    INTEGER NOT NULL,
    raw_log    JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

-- prevent naming conflict with maker.flip_kicks in postgraphile
COMMENT ON TABLE maker.flip_kick IS E'@name flipKickEvent';

CREATE INDEX flip_kick_header_index
    ON maker.flip_kick (header_id);
CREATE INDEX flip_kick_bid_id_index
    ON maker.flip_kick (bid_id);
CREATE INDEX flip_kick_address_id_index
    ON maker.flip_kick (address_id);

ALTER TABLE public.checked_headers
    ADD COLUMN flip_kick INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.flip_kick_address_id_index;
DROP INDEX maker.flip_kick_bid_id_index;
DROP INDEX maker.flip_kick_header_index;

DROP TABLE maker.flip_kick;

ALTER TABLE public.checked_headers
    DROP COLUMN flip_kick;
