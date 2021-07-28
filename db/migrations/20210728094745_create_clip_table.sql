-- +goose Up
CREATE TABLE maker.clip
(
    address_id   BIGINT    NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    block_number BIGINT    NOT NULL,
    sale_id      NUMERIC   NOT NULL,
    pos          NUMERIC   DEFAULT NULL,
    tab          NUMERIC   DEFAULT NULL,
    lot          NUMERIC   DEFAULT NULL,
    usr          TEXT      DEFAULT NULL,
    tic          NUMERIC   DEFAULT NULL,
    "top"          NUMERIC   DEFAULT NULL,
    created      TIMESTAMP DEFAULT NULL,
    updated      TIMESTAMP NOT NULL,
    PRIMARY KEY (address_id, sale_id, block_number)
);

CREATE INDEX clip_address_index
    ON maker.clip (address_id);

CREATE FUNCTION clip_sale_pos_before_block(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT pos
FROM maker.clip_sale_pos
         LEFT JOIN public.headers ON clip_sale_pos.header_id = headers.id
WHERE clip_sale_pos.sale_id = clip_sale_pos_before_block.sale_id
  AND clip_sale_pos.address_id = clip_sale_pos_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = clip_sale_pos_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_pos_before_block
    IS E'@omit';

CREATE FUNCTION clip_sale_tab_before_block(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT tab
FROM maker.clip_sale_tab
         LEFT JOIN public.headers ON clip_sale_tab.header_id = headers.id
WHERE clip_sale_tab.sale_id = clip_sale_tab_before_block.sale_id
  AND clip_sale_tab.address_id = clip_sale_tab_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = clip_sale_tab_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_tab_before_block
    IS E'@omit';

CREATE FUNCTION clip_sale_lot_before_block(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT lot
FROM maker.clip_sale_lot
         LEFT JOIN public.headers ON clip_sale_lot.header_id = headers.id
WHERE clip_sale_lot.sale_id = clip_sale_lot_before_block.sale_id
  AND clip_sale_lot.address_id = clip_sale_lot_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = clip_sale_lot_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_lot_before_block
    IS E'@omit';

CREATE FUNCTION clip_sale_usr_before_block(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS TEXT AS
$$
SELECT usr
FROM maker.clip_sale_usr
         LEFT JOIN public.headers ON clip_sale_usr.header_id = headers.id
WHERE clip_sale_usr.sale_id = clip_sale_usr_before_block.sale_id
  AND clip_sale_usr.address_id = clip_sale_usr_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = clip_sale_usr_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_usr_before_block
    IS E'@omit';

CREATE FUNCTION clip_sale_tic_before_block(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT tic
FROM maker.clip_sale_tic
         LEFT JOIN public.headers ON clip_sale_tic.header_id = headers.id
WHERE clip_sale_tic.sale_id = clip_sale_tic_before_block.sale_id
  AND clip_sale_tic.address_id = clip_sale_tic_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = clip_sale_tic_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_tic_before_block
    IS E'@omit';

CREATE FUNCTION clip_sale_top_before_block(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS NUMERIC AS
$$
SELECT top
FROM maker.clip_sale_top
         LEFT JOIN public.headers ON clip_sale_top.header_id = headers.id
WHERE clip_sale_top.sale_id = clip_sale_top_before_block.sale_id
  AND clip_sale_top.address_id = clip_sale_top_before_block.address_id
  AND headers.block_number < (SELECT block_number FROM public.headers WHERE id = clip_sale_top_before_block.header_id)
ORDER BY block_number DESC
LIMIT 1
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_top_before_block
    IS E'@omit';

CREATE FUNCTION clip_sale_time_created(address_id BIGINT, sale_id NUMERIC) RETURNS TIMESTAMP AS
$$
SELECT api.epoch_to_datetime(MIN(block_timestamp))
FROM public.headers
         LEFT JOIN maker.clip_kick ON clip_kick.header_id = headers.id
WHERE clip_kick.address_id = clip_sale_time_created.address_id
  AND clip_kick.sale_id = clip_sale_time_created.sale_id
$$
    LANGUAGE sql;

COMMENT ON FUNCTION clip_sale_time_created
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.delete_obsolete_clip(sale_id NUMERIC, address_id BIGINT, header_id INTEGER) RETURNS VOID AS
$$
DECLARE
    clip_block      BIGINT     := (
        SELECT block_number
        FROM public.headers
        WHERE id = header_id);
    clip_state      maker.clip := (
        SELECT (clip.address_id, block_number, clip.sale_id, pos, tab, lot, usr, tic, "top", created, updated)
        FROM maker.clip
        WHERE clip.sale_id = delete_obsolete_clip.sale_id
          AND clip.address_id = delete_obsolete_clip.address_id
          AND clip.block_number = clip_block);
    prev_clip_state maker.clip := (
        SELECT (clip.address_id, block_number, clip.sale_id, pos, tab, lot, usr, tic, "top", created, updated)
        FROM maker.clip
        WHERE clip.sale_id = delete_obsolete_clip.sale_id
          AND clip.address_id = delete_obsolete_clip.address_id
          AND clip.block_number < clip_block
        ORDER BY clip.block_number DESC
        LIMIT 1);
BEGIN
    DELETE
    FROM maker.clip
    WHERE clip.sale_id = delete_obsolete_clip.sale_id
      AND clip.address_id = delete_obsolete_clip.address_id
      AND clip.block_number = clip_block
      AND (clip_state.pos IS NULL OR clip_state.pos = prev_clip_state.pos)
      AND (clip_state.tab IS NULL OR clip_state.tab = prev_clip_state.tab)
      AND (clip_state.lot IS NULL OR clip_state.lot = prev_clip_state.lot)
      AND (clip_state.usr IS NULL OR clip_state.usr = prev_clip_state.usr)
      AND (clip_state.tic IS NULL OR clip_state.tic = prev_clip_state.tic)
      AND (clip_state."top" IS NULL OR clip_state."top" = prev_clip_state."top");
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.delete_obsolete_clip
    IS E'@omit';

CREATE OR REPLACE FUNCTION maker.insert_clip_created(new_event maker.clip_kick) RETURNS VOID
AS
$$
UPDATE maker.clip
SET created = api.epoch_to_datetime(headers.block_timestamp)
FROM public.headers
WHERE headers.id = new_event.header_id
  AND clip.address_id = new_event.address_id
  AND clip.sale_id = new_event.sale_id
  AND clip.created IS NULL
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.insert_clip_created
    IS E'@omit';

CREATE OR REPLACE FUNCTION maker.clear_clip_created(old_event maker.clip_kick) RETURNS VOID
AS
$$
UPDATE maker.clip
SET created = clip_sale_time_created(old_event.address_id, old_event.sale_id)
WHERE clip.address_id = old_event.address_id
  AND clip.sale_id = old_event.sale_id
$$
    LANGUAGE sql;

COMMENT ON FUNCTION maker.clear_clip_created
    IS E'@omit';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_clip_created() RETURNS TRIGGER
AS
$$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        PERFORM maker.insert_clip_created(NEW);
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM maker.clear_clip_created(OLD);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

COMMENT ON FUNCTION maker.update_clip_created
    IS E'@omit';

CREATE TRIGGER clip_created_trigger
    AFTER INSERT OR DELETE
    ON maker.clip_kick
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_clip_created();

-- +goose Down
DROP TRIGGER clip_created_trigger ON maker.clip_kick;

DROP FUNCTION maker.insert_clip_created(maker.clip_kick);
DROP FUNCTION maker.clear_clip_created(maker.clip_kick);
DROP FUNCTION maker.update_clip_created();
DROP FUNCTION clip_sale_pos_before_block(NUMERIC, BIGINT, INTEGER);
DROP FUNCTION clip_sale_tab_before_block(NUMERIC, BIGINT, INTEGER);
DROP FUNCTION clip_sale_lot_before_block(NUMERIC, BIGINT, INTEGER);
DROP FUNCTION clip_sale_usr_before_block(NUMERIC, BIGINT, INTEGER);
DROP FUNCTION clip_sale_tic_before_block(NUMERIC, BIGINT, INTEGER);
DROP FUNCTION clip_sale_top_before_block(NUMERIC, BIGINT, INTEGER);

DROP FUNCTION clip_sale_time_created(BIGINT, NUMERIC);

DROP INDEX maker.clip_address_index;
DROP TABLE maker.clip;