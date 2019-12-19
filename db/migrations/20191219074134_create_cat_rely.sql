-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.cat_rely
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    address_id INT,
    UNIQUE (header_id, log_id)
);

CREATE INDEX cat_rely_header_index
    ON maker.cat_rely (header_id);
CREATE INDEX cat_rely_log_index
    ON maker.cat_rely (log_id);
CREATE INDEX cat_rely_address_index
    ON maker.cat_rely (address_id);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.cat_rely_header_index;
DROP INDEX maker.cat_rely_log_index;
DROP INDEX maker.cat_rely_address_index;

DROP TABLE maker.cat_rely;