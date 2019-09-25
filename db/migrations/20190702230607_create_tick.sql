-- +goose Up
CREATE TABLE maker.tick
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX tick_header_index
    ON maker.tick (header_id);


CREATE INDEX tick_bid_id_index
    ON maker.tick (bid_id);

-- +goose Down
DROP INDEX maker.tick_header_index;
DROP INDEX maker.tick_bid_id_index;
DROP TABLE maker.tick;
