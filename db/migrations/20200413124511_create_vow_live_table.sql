-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.vow_live
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    live      NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, live)
);

CREATE INDEX vow_live_header_id_index
    ON maker.vow_live (header_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.vow_live_header_id_index;

DROP TABLE maker.vow_live;