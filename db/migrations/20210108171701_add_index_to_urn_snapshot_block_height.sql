-- +goose NO TRANSACTION
-- +goose Up
CREATE INDEX CONCURRENTLY urn_snapshot_block_height_index
  ON api.urn_snapshot (block_height);

-- +goose Down
DROP INDEX CONCURRENTLY api.urn_snapshot_block_height_index;
