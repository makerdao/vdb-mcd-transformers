-- +goose Up
DROP TABLE maker.vow_woe;

-- +goose Down

CREATE TABLE maker.vow_woe (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  woe           numeric
);
