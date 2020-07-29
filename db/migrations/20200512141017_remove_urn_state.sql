-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE FUNCTION api.bite_event_urn(event api.bite_event) RETURNS api.urn_snapshot
    LANGUAGE sql STABLE
    AS $$
SELECT *
FROM api.get_urn(event.ilk_identifier, event.urn_identifier, event.block_height)
$$;

CREATE FUNCTION api.flip_bid_snapshot_urn(flip api.flip_bid_snapshot) RETURNS api.urn_snapshot
       LANGUAGE sql STABLE
       AS $$
SELECT *
FROM api.get_urn(
     (SELECT identifier FROM maker.ilks WHERE ilks.id = flip.ilk_id),
     (SELECT identifier FROM maker.urns WHERE urns.id = flip.urn_id),
     flip.block_height)
$$;

CREATE FUNCTION api.frob_event_urn(event api.frob_event) RETURNS api.urn_snapshot
       LANGUAGE sql STABLE
       AS $$
SELECT *
FROM api.get_urn(event.ilk_identifier, event.urn_identifier, event.block_height)
$$;

CREATE FUNCTION api.managed_cdp_urn(cdp api.managed_cdp) RETURNS api.urn_snapshot
       LANGUAGE sql STABLE
       AS $$
SELECT *
FROM api.get_urn(cdp.ilk_identifier, cdp.urn_identifier, api.max_block())
$$;

-- +goose Down
DROP FUNCTION api.bite_event_urn;
DROP FUNCTION api.flip_bid_snapshot_urn;
DROP FUNCTION api.frob_event_urn;
DROP FUNCTION api.managed_cdp_urn;
