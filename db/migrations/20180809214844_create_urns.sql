-- +goose Up
CREATE TABLE maker.urns (
  id    SERIAL PRIMARY KEY,
  ilk   INTEGER NOT NULL REFERENCES maker.ilks (id),
  guy   TEXT,
  UNIQUE (ilk, guy)
);

-- +goose Down
DROP TABLE maker.urns;
