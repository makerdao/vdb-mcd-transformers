-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION api.clip_sale_snapshot_active(clip_sale_snapshot api.clip_sale_snapshot) RETURNS boolean
    LANGUAGE sql
    STABLE
AS
$$
SELECT EXISTS(SELECT 1
              FROM maker.clip_active_sales
                       JOIN headers h ON clip_active_sales.header_id = h.id
              WHERE clip_sale_snapshot.sale_id = sale_id
                AND clip_sale_snapshot.block_height = h.block_number)
$$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.clip_sale_snapshot_active(api.clip_sale_snapshot);