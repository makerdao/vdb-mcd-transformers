-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE maker.urns
(
    id         SERIAL PRIMARY KEY,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    identifier CITEXT,
    UNIQUE (ilk_id, identifier)
);

COMMENT ON TABLE maker.urns
    IS E'@name raw_urns\nCDP.';

CREATE INDEX urn_ilk_index
    ON maker.urns (ilk_id);

-- +goose Down
DROP INDEX maker.urn_ilk_index;
DROP TABLE maker.urns;
