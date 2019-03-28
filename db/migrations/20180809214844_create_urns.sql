-- +goose Up
CREATE TABLE maker.urns (
  id      SERIAL PRIMARY KEY,
  ilk_id  INTEGER NOT NULL REFERENCES maker.ilks (id),
  guy     TEXT,
  UNIQUE (ilk_id, guy)
);

-- +goose Down
DROP TABLE maker.urns;
