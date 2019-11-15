-- +goose Up
CREATE TABLE public.full_sync_logs
(
    id           SERIAL PRIMARY KEY,
    block_number BIGINT,
    address      VARCHAR(66),
    tx_hash      VARCHAR(66),
    index        BIGINT,
    topic0       VARCHAR(66),
    topic1       VARCHAR(66),
    topic2       VARCHAR(66),
    topic3       VARCHAR(66),
    data         TEXT
);

COMMENT ON TABLE public.full_sync_logs
    IS E'@omit';

-- +goose Down
DROP TABLE full_sync_logs;
