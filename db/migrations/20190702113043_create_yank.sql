-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.yank
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX yank_header_index
    ON maker.yank (header_id);

CREATE INDEX yank_bid_id_index
    on maker.yank (bid_id);

ALTER TABLE public.checked_headers
    ADD COLUMN yank INTEGER NOT NULL DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.yank_header_index;
DROP INDEX maker.yank_bid_id_index;

DROP TABLE maker.yank;
