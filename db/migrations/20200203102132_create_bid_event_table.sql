-- +goose Up
CREATE TABLE api.bid_event
(
    bid_id           NUMERIC     NOT NULL,
    contract_address TEXT        NOT NULL,
    act              api.bid_act NOT NULL,
    lot              NUMERIC DEFAULT NULL,
    bid_amount       NUMERIC DEFAULT NULL,
    ilk_identifier   TEXT    DEFAULT NULL,
    urn_identifier   TEXT    DEFAULT NULL,
    block_height     BIGINT      NOT NULL
);


CREATE OR REPLACE FUNCTION maker.insert_bid_event(bid_id NUMERIC, address_id INTEGER, header_id INTEGER,
                                                  act api.bid_act, lot NUMERIC, bid_amount NUMERIC) RETURNS VOID AS
$$
INSERT
INTO api.bid_event (bid_id, contract_address, act, lot, bid_amount, ilk_identifier, urn_identifier, block_height)
VALUES (insert_bid_event.bid_id,
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
    PERFORM maker.insert_bid_event(NEW.bid_id, NEW.address_id, NEW.header_id, TG_ARGV[0]::api.bid_act, NEW.lot,
                                   NEW.bid);
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
    PERFORM maker.insert_bid_event(NEW.bid_id, NEW.address_id, NEW.header_id, TG_ARGV[0]::api.bid_act, NULL, NULL);
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
UPDATE api.bid_event
SET ilk_identifier = (SELECT identifier FROM maker.ilks WHERE id = new_diff.ilk_id)
WHERE bid_event.contract_address = (SELECT address FROM public.addresses WHERE id = new_diff.address_id)
$$
    LANGUAGE sql;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_bid_event_ilk() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_bid_event_ilk(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_ilk
    AFTER INSERT
    ON maker.flip_ilk
    FOR EACH ROW
EXECUTE PROCEDURE maker.update_bid_event_ilk();


CREATE OR REPLACE FUNCTION maker.insert_bid_event_urn(new_diff maker.flip_bid_usr) RETURNS VOID AS
$$
UPDATE api.bid_event
SET urn_identifier = new_diff.usr
WHERE bid_event.bid_id = new_diff.bid_id
  AND bid_event.contract_address = (SELECT address FROM public.addresses WHERE id = new_diff.address_id)
$$
    LANGUAGE sql;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION maker.update_bid_event_urn() RETURNS TRIGGER
AS
$$
BEGIN
    PERFORM maker.insert_bid_event_urn(NEW);
    RETURN NULL;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER flip_urn
    AFTER INSERT
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

DROP FUNCTION maker.insert_bid_event_urn(maker.flip_bid_usr);
DROP FUNCTION maker.insert_bid_event_ilk(maker.flip_ilk);
DROP FUNCTION maker.insert_bid_event(NUMERIC, INTEGER, INTEGER, api.bid_act, NUMERIC, NUMERIC);

DROP TABLE api.bid_event CASCADE;
