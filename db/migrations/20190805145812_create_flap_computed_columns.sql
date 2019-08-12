-- +goose Up
-- +goose StatementBegin

-- Extend flap with bid events
CREATE FUNCTION api.flap_bid_events(flap api.flap)
RETURNS SETOF api.flap_bid_event AS
    $$
    SELECT *
    FROM api.all_flap_bid_events()
    WHERE bid_id = flap.bid_id
    $$
LANGUAGE sql
STABLE;

-- +goose StatementEnd
-- +goose Down
DROP FUNCTION api.flap_bid_events(api.flap);


