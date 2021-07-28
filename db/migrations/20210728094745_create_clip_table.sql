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
    top          NUMERIC   DEFAULT NULL,
    created      TIMESTAMP DEFAULT NULL,
    updated      TIMESTAMP NOT NULL,
    PRIMARY KEY (address_id, sale_id, block_number)
);

CREATE INDEX clip_address_index
    ON maker.clip (address_id);

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
DROP FUNCTION clip_sale_time_created(BIGINT, NUMERIC);

DROP INDEX maker.clip_address_index;
DROP TABLE maker.clip;