-- +goose Up
CREATE OR REPLACE FUNCTION api.flip_bid_event_bid(event api.flip_bid_event)
    RETURNS api.flip_bid_snapshot AS
$$
SELECT * FROM api.get_flip_with_address(event.bid_id, event.contract_address)
$$
    LANGUAGE sql
    STABLE;

-- +goose Down

-- put flip_bid_event_bid back as it was
CREATE OR REPLACE FUNCTION api.flip_bid_event_bid(event api.flip_bid_event)
    RETURNS api.flip_bid_snapshot AS
$$
WITH ilk AS (
    SELECT ilks.identifier
    FROM maker.flip_ilk
             LEFT JOIN maker.ilks ON ilks.id = flip_ilk.ilk_id
    WHERE flip_ilk.address_id = (SELECT id FROM addresses WHERE address = event.contract_address)
    LIMIT 1
)
SELECT *
FROM api.get_flip(event.bid_id, (SELECT identifier FROM ilk))
$$
    LANGUAGE sql
    STABLE;
