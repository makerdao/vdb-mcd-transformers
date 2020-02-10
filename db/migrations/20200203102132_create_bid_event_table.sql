-- +goose Up
CREATE TABLE maker.bid_event
(
    log_id           BIGINT PRIMARY KEY REFERENCES event_logs (id) ON DELETE CASCADE,
    bid_id           NUMERIC     NOT NULL,
    contract_address TEXT        NOT NULL,
    act              api.bid_act NOT NULL,
    lot              NUMERIC DEFAULT NULL,
    bid_amount       NUMERIC DEFAULT NULL,
    ilk_identifier   TEXT    DEFAULT NULL,
    urn_identifier   TEXT    DEFAULT NULL,
    block_height     BIGINT      NOT NULL
);
COMMENT ON COLUMN maker.bid_event.log_id IS E'@omit';

CREATE INDEX bid_event_index ON maker.bid_event (contract_address, bid_id);
CREATE INDEX bid_event_urn_index ON maker.bid_event (ilk_identifier, urn_identifier);


CREATE OR REPLACE FUNCTION maker.insert_bid_event(log_id BIGINT, bid_id NUMERIC, address_id INTEGER, header_id INTEGER,
                                                  act api.bid_act, lot NUMERIC, bid_amount NUMERIC) RETURNS VOID AS
$$
INSERT
INTO maker.bid_event (log_id, bid_id, contract_address, act, lot, bid_amount, ilk_identifier, urn_identifier,
                      block_height)
VALUES (insert_bid_event.log_id,
        insert_bid_event.bid_id,
        (SELECT address FROM public.addresses WHERE addresses.id = insert_bid_event.address_id),
        insert_bid_event.act,
        insert_bid_event.lot,
        insert_bid_event.bid_amount,
        (SELECT ilks.identifier
         FROM maker.flip_ilk
                  JOIN maker.ilks ON flip_ilk.ilk_id = ilks.id
                  JOIN public.headers ON flip_ilk.header_id = headers.id
         WHERE flip_ilk.address_id = insert_bid_event.address_id
         ORDER BY headers.block_number DESC
         LIMIT 1),
        (SELECT usr
         FROM maker.flip_bid_usr
                  JOIN public.headers ON flip_bid_usr.header_id = headers.id
         WHERE flip_bid_usr.bid_id = insert_bid_event.bid_id
           AND flip_bid_usr.address_id = insert_bid_event.address_id
         ORDER BY headers.block_number DESC
         LIMIT 1),
        (SELECT block_number FROM public.headers WHERE id = insert_bid_event.header_id))
$$
    LANGUAGE sql;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_bid_kick_tend_dent_event() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_bid_event(NEW.log_id, NEW.bid_id, NEW.address_id, NEW.header_id, TG_ARGV[0]::api.bid_act,
                                   NEW.lot, NEW.bid);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flap_kick
    AFTER INSERT
    ON maker.flap_kick
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_bid_kick_tend_dent_event('kick');

CREATE TRIGGER flip_kick
    AFTER INSERT
    ON maker.flip_kick
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_bid_kick_tend_dent_event('kick');

CREATE TRIGGER flop_kick
    AFTER INSERT
    ON maker.flop_kick
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_bid_kick_tend_dent_event('kick');

CREATE TRIGGER tend
    AFTER INSERT
    ON maker.tend
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_bid_kick_tend_dent_event('tend');

CREATE TRIGGER dent
    AFTER INSERT
    ON maker.dent
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_bid_kick_tend_dent_event('dent');

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_bid_tick_deal_yank_event() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_bid_event(NEW.log_id, NEW.bid_id, NEW.address_id, NEW.header_id, TG_ARGV[0]::api.bid_act, NULL,
                                   NULL);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER tick
    AFTER INSERT
    ON maker.tick
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_bid_tick_deal_yank_event('tick');

CREATE TRIGGER deal
    AFTER INSERT
    ON maker.deal
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_bid_tick_deal_yank_event('deal');

CREATE TRIGGER yank
    AFTER INSERT
    ON maker.yank
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_bid_tick_deal_yank_event('yank');

CREATE OR REPLACE FUNCTION maker.insert_bid_event_ilk(new_diff maker.flip_ilk) RETURNS VOID AS
$$
UPDATE maker.bid_event
SET ilk_identifier = (SELECT identifier FROM maker.ilks WHERE id = new_diff.ilk_id)
WHERE bid_event.contract_address = (SELECT address FROM public.addresses WHERE id = new_diff.address_id)
$$
    LANGUAGE sql;

CREATE OR REPLACE FUNCTION maker.clear_bid_event_ilk(old_diff maker.flip_ilk) RETURNS VOID AS
$$
UPDATE maker.bid_event
SET ilk_identifier = NULL
WHERE bid_event.contract_address = (SELECT address FROM public.addresses WHERE id = old_diff.address_id)
$$
    LANGUAGE sql;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_bid_event_ilk() RETURNS TRIGGER
AS
$$
BEGIN
    IF TG_OP = 'INSERT' THEN
        PERFORM maker.insert_bid_event_ilk(NEW);
    ELSIF TG_OP = 'DELETE' THEN
        PERFORM maker.clear_bid_event_ilk(OLD);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_ilk
    AFTER INSERT OR DELETE
    ON maker.flip_ilk
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_bid_event_ilk();


CREATE OR REPLACE FUNCTION maker.insert_bid_event_urn(diff maker.flip_bid_usr, new_usr TEXT) RETURNS VOID AS
$$
UPDATE maker.bid_event
SET urn_identifier = new_usr
WHERE bid_event.bid_id = diff.bid_id
  AND bid_event.contract_address = (SELECT address FROM public.addresses WHERE id = diff.address_id)
$$
    LANGUAGE sql;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_bid_event_urn() RETURNS TRIGGER
AS
$$
BEGIN
    IF TG_OP = 'INSERT' THEN
        PERFORM maker.insert_bid_event_urn(NEW, NEW.usr);
    ELSIF TG_OP = 'DELETE' THEN
        PERFORM maker.insert_bid_event_urn(OLD, NULL);
    END IF;
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_urn
    AFTER INSERT OR DELETE
    ON maker.flip_bid_usr
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_bid_event_urn();

-- +goose Down
DROP TRIGGER flip_urn ON maker.flip_bid_usr;
DROP TRIGGER flip_ilk ON maker.flip_ilk;
DROP TRIGGER yank ON maker.yank;
DROP TRIGGER deal ON maker.deal;
DROP TRIGGER tick ON maker.tick;
DROP TRIGGER dent ON maker.dent;
DROP TRIGGER tend ON maker.tend;
DROP TRIGGER flop_kick ON maker.flop_kick;
DROP TRIGGER flip_kick ON maker.flip_kick;
DROP TRIGGER flap_kick ON maker.flap_kick;

DROP FUNCTION maker.update_bid_event_urn();
DROP FUNCTION maker.update_bid_event_ilk();
DROP FUNCTION maker.update_bid_tick_deal_yank_event();
DROP FUNCTION maker.update_bid_kick_tend_dent_event();

DROP FUNCTION maker.insert_bid_event_urn(maker.flip_bid_usr, TEXT);
DROP FUNCTION maker.clear_bid_event_ilk(maker.flip_ilk);
DROP FUNCTION maker.insert_bid_event_ilk(maker.flip_ilk);
DROP FUNCTION maker.insert_bid_event(BIGINT, NUMERIC, INTEGER, INTEGER, api.bid_act, NUMERIC, NUMERIC);

DROP TABLE maker.bid_event CASCADE;
