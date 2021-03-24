-- +goose Up
CREATE TABLE maker.clip_kick
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    top        NUMERIC,
    tab        NUMERIC,
    lot        NUMERIC,
    usr        BIGINT,
    kpr        BIGINT,
    coin       NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX clip_kick_header_index
    ON maker.clip_kick (header_id);
CREATE INDEX clip_kick_bid_id_index
    ON maker.clip_kick (bid_id);
CREATE INDEX clip_kick_address_index
    ON maker.clip_kick (address_id);
CREATE INDEX clip_kick_log_index
    ON maker.clip_kick (log_id);

-- +goose Down
DROP INDEX maker.clip_kick_log_index;
DROP INDEX maker.clip_kick_address_index;
DROP INDEX maker.clip_kick_bid_id_index;
DROP INDEX maker.clip_kick_header_index;

DROP TABLE maker.clip_kick;
