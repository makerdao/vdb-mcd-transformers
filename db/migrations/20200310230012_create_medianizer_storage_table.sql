-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.median_val
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    val       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, val)
);

CREATE INDEX median_val_header_id_index
    ON maker.median_val (header_id);

CREATE TABLE maker.median_bar
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    bar       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, bar)
);

CREATE INDEX median_bar_header_id_index
    ON maker.median_bar (header_id);

CREATE TABLE maker.median_age
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    age       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, age)
);

CREATE INDEX median_age_header_id_index
    ON maker.median_age (header_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.median_val_header_id_index;
DROP INDEX maker.median_bar_header_id_index;
DROP INDEX maker.median_age_header_id_index;

DROP TABLE maker.median_val;
DROP TABLE maker.median_bar;
DROP TABLE maker.median_age;
