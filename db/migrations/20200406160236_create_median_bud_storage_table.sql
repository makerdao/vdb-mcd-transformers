-- +goose Up
CREATE TABLE maker.median_bud
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    a          INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    bud        INTEGER NOT NULL,
    UNIQUE (diff_id, header_id, address_id, a, bud)
);

CREATE INDEX median_bud_header_id_index
    ON maker.median_bud (header_id);
CREATE INDEX median_bud_address_index
    ON maker.median_bud (address_id);
CREATE INDEX median_bud_a_index
    ON maker.median_bud (a);


-- +goose Down
DROP TABLE maker.median_bud;
