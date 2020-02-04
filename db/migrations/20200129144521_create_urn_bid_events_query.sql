-- +goose Up

CREATE FUNCTION api.urn_bid_events(urn_identifier TEXT, ilk_identifier TEXT) RETURNS SETOF api.flip_bid_event AS
$$
SELECT bid_events.*
FROM api.all_flip_bid_events() bid_events
         JOIN public.addresses ON bid_events.contract_address = addresses.address
         JOIN maker.flip_ilk ON addresses.id = flip_ilk.address_id
         JOIN maker.ilks ON flip_ilk.ilk_id = ilks.id
         JOIN maker.flip_bid_usr ON bid_events.bid_id = flip_bid_usr.bid_id AND addresses.id = flip_bid_usr.address_id
WHERE ilks.identifier = ilk_identifier
  AND flip_bid_usr.usr = urn_identifier
$$
    LANGUAGE sql
    STRICT
    STABLE;

-- +goose Down
DROP FUNCTION api.urn_bid_events(TEXT, TEXT);
