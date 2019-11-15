-- +goose Up
CREATE TABLE public.watched_contracts
(
  contract_id   SERIAL PRIMARY KEY,
  contract_abi  json,
  contract_hash VARCHAR(66) UNIQUE
);

COMMENT ON TABLE public.watched_contracts
    IS E'@omit';

-- +goose Down
DROP TABLE watched_contracts;
