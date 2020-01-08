-- +goose Up
CREATE TABLE maker.wards
(
    id         SERIAL PRIMARY KEY,
    diff_id    BIGINT  NOT NULL REFERENCES storage_diff (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    usr        INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    wards      INTEGER NOT NULL,
    UNIQUE (diff_id, header_id, address_id, usr, wards)
);

CREATE INDEX wards_header_id_index
    ON maker.wards (header_id);
CREATE INDEX wards_address_index
    ON maker.wards (address_id);
CREATE INDEX wards_usr_index
    ON maker.wards (usr);


-- +goose Down
DROP INDEX maker.wards_usr_index;
DROP INDEX maker.wards_address_index;
DROP INDEX maker.wards_header_id_index;

DROP TABLE maker.wards;
