-- +goose Up

CREATE FUNCTION api.urn_bid_events(urn_identifier TEXT, ilk_identifier TEXT) RETURNS SETOF maker.bid_event AS
$$
SELECT *
FROM maker.bid_event
WHERE bid_event.ilk_identifier = urn_bid_events.ilk_identifier
  AND bid_event.urn_identifier = urn_bid_events.urn_identifier
$$
    LANGUAGE sql
    STRICT
    STABLE;

-- +goose Down
DROP FUNCTION api.urn_bid_events(TEXT, TEXT);
