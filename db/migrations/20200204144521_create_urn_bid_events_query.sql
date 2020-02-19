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

COMMENT ON FUNCTION api.urn_bid_events(urn_identifier TEXT, ilk_identifier TEXT)
    IS E'Get bid events related to an auction associated with a given Urn. urnIdentifier (e.g. "0xC93C178EC17B06bddBa0CC798546161aF9D25e8A") and ilkIdentifier (e.g. "ETH-A") are required.';

-- +goose Down
DROP FUNCTION api.urn_bid_events(TEXT, TEXT);
