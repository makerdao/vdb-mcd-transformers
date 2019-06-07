-- +goose Up
CREATE TABLE maker.urns
(
    id     SERIAL PRIMARY KEY,
    ilk_id INTEGER NOT NULL REFERENCES maker.ilks (id)ON DELETE CASCADE,
    guy    TEXT,
    UNIQUE (ilk_id, guy)
);

COMMENT ON TABLE maker.urns
    IS E'@name raw_urns';

-- +goose Down
DROP TABLE maker.urns;
