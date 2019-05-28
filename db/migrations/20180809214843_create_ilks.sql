-- +goose Up
CREATE TABLE maker.ilks (
  id         SERIAL PRIMARY KEY,
  ilk        TEXT UNIQUE NOT NULL,
  identifier TEXT UNIQUE NOT NULL
);

COMMENT ON TABLE maker.ilks IS E'@name raw_ilks';

-- +goose Down
DROP TABLE maker.ilks;
