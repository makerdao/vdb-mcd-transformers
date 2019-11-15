-- +goose Up
CREATE TABLE public.nodes (
  id            SERIAL PRIMARY KEY,
  client_name   VARCHAR,
  genesis_block VARCHAR(66),
  network_id    NUMERIC,
  node_id       VARCHAR(128),
  CONSTRAINT node_uc UNIQUE (genesis_block, network_id, node_id)
);

COMMENT ON TABLE public.nodes
    IS E'@omit';

-- +goose Down
DROP TABLE nodes;
