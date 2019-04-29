-- +goose Up
CREATE TABLE maker.ilks (
  id        SERIAL PRIMARY KEY,
  ilk       TEXT UNIQUE NOT NULL,
  name      TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE maker.ilks;
