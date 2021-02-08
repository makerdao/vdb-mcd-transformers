-- +goose NO TRANSACTION
-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS vat_dai_guy_index
  ON maker.vat_dai (guy);

-- +goose Down
DROP INDEX CONCURRENTLY maker.vat_dai_guy_index;
