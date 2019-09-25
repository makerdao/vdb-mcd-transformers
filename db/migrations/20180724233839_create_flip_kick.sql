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
    log_id     BIGINT NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

-- prevent naming conflict with maker.flip_kicks in postgraphile
COMMENT ON TABLE maker.flip_kick IS E'@name flipKickEvent';

CREATE INDEX flip_kick_header_index
    ON maker.flip_kick (header_id);
CREATE INDEX flip_kick_bid_id_index
    ON maker.flip_kick (bid_id);
CREATE INDEX flip_kick_address_id_index
    ON maker.flip_kick (address_id);


-- +goose Down
DROP INDEX maker.flip_kick_address_id_index;
DROP INDEX maker.flip_kick_bid_id_index;
DROP INDEX maker.flip_kick_header_index;

DROP TABLE maker.flip_kick;
