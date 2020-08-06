-- +goose Up
CREATE TABLE maker.cat_box
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    box       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, box)
);

CREATE INDEX cat_box_header_id_index
    ON maker.cat_box (header_id);

CREATE TABLE maker.cat_litter
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    litter    NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, litter)
);

CREATE INDEX cat_litter_header_id_index
    ON maker.cat_litter (header_id);


-- +goose Down
DROP TABLE maker.cat_box;
DROP TABLE maker.cat_litter;
