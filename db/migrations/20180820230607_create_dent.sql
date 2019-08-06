-- +goose Up
CREATE TABLE maker.dent
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    lot        NUMERIC,
    bid        NUMERIC,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX dent_header_index
    ON maker.dent (header_id);

ALTER TABLE public.checked_headers
    ADD COLUMN dent INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.dent_header_index;

DROP TABLE maker.dent;

ALTER TABLE public.checked_headers
    DROP COLUMN dent;