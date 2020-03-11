-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.medianizer_val
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    val       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, val)
);

CREATE INDEX medianizer_val_header_id_index
    ON maker.medianizer_val (header_id);

CREATE TABLE maker.medianizer_bar
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    bar       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, bar)
);

CREATE INDEX medianizer_bar_header_id_index
    ON maker.medianizer_bar (header_id);

CREATE TABLE maker.medianizer_age
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    age       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, age)
);

CREATE INDEX medianizer_age_header_id_index
    ON maker.medianizer_age (header_id);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.medianizer_val_header_id_index;
DROP INDEX maker.medianizer_bar_header_id_index;
DROP INDEX maker.medianizer_age_header_id_index;

DROP TABLE maker.medianizer_val;
DROP TABLE maker.medianizer_bar;
DROP TABLE maker.medianizer_age;
