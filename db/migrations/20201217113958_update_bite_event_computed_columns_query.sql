-- +goose Up
--add flip_address attribute to bite_event type
ALTER TYPE api.bite_event ADD ATTRIBUTE flip_address text;

-- update all_bites to return a bite_event with a flip_address attribute
CREATE OR REPLACE FUNCTION api.all_bites(ilk_identifier TEXT, max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier)

SELECT ilk_identifier,
       identifier AS urn_identifier,
       bid_id,
       ink,
       art,
       tab,
       block_number,
       log_id,
       flip
FROM maker.bite
         LEFT JOIN maker.urns ON bite.urn_id = urns.id
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE urns.ilk_id = (SELECT id FROM ilk)
ORDER BY urn_identifier, block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
    OFFSET
    all_bites.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;

-- update urn_bites to return a bite_event with a flip_address attribute
CREATE OR REPLACE FUNCTION api.urn_bites(ilk_identifier TEXT, urn_identifier TEXT, max_results INTEGER DEFAULT -1,
                              result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier),
     urn AS (SELECT id
             FROM maker.urns
             WHERE ilk_id = (SELECT id FROM ilk)
               AND identifier = urn_bites.urn_identifier)

SELECT ilk_identifier,
       urn_bites.urn_identifier,
       bid_id,
       ink,
       art,
       tab,
       block_number,
       log_id,
       flip
FROM maker.bite
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE bite.urn_id = (SELECT id FROM urn)
ORDER BY block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
    OFFSET
    urn_bites.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;

-- update bite_event_bid function to make sure to return the flip from the correct flip contract
CREATE OR REPLACE FUNCTION api.bite_event_bid(event api.bite_event)
    RETURNS api.flip_bid_snapshot AS
$$
SELECT *
FROM api.get_flip_with_address(event.bid_id, event.flip_address, event.ilk_identifier, event.block_height)
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- put bite event back to where it was before
ALTER TYPE api.bite_event DROP ATTRIBUTE flip_address;

-- put all_bites query back to where it was before this migration
CREATE OR REPLACE FUNCTION api.all_bites(ilk_identifier TEXT, max_results INTEGER DEFAULT -1, result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier)

SELECT ilk_identifier,
       identifier AS urn_identifier,
       bid_id,
       ink,
       art,
       tab,
       block_number,
       log_id
FROM maker.bite
         LEFT JOIN maker.urns ON bite.urn_id = urns.id
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE urns.ilk_id = (SELECT id FROM ilk)
ORDER BY urn_identifier, block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
    OFFSET
    all_bites.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;

-- put urn_bites query back to where it was before this migration
CREATE OR REPLACE FUNCTION api.urn_bites(ilk_identifier TEXT, urn_identifier TEXT, max_results INTEGER DEFAULT -1,
                              result_offset INTEGER DEFAULT 0)
    RETURNS SETOF api.bite_event AS
$$
WITH ilk AS (SELECT id FROM maker.ilks WHERE ilks.identifier = ilk_identifier),
     urn AS (SELECT id
             FROM maker.urns
             WHERE ilk_id = (SELECT id FROM ilk)
               AND identifier = urn_bites.urn_identifier)

SELECT ilk_identifier,
       urn_bites.urn_identifier,
       bid_id,
       ink,
       art,
       tab,
       block_number,
       log_id
FROM maker.bite
         LEFT JOIN headers ON bite.header_id = headers.id
WHERE bite.urn_id = (SELECT id FROM urn)
ORDER BY block_number DESC
LIMIT CASE WHEN max_results = -1 THEN NULL ELSE max_results END
    OFFSET
    urn_bites.result_offset
$$
    LANGUAGE sql
    STRICT --necessary for postgraphile queries with required arguments
    STABLE;



-- put bite_event_bit query back to where it was before this migration
CREATE OR REPLACE FUNCTION api.bite_event_bid(event api.bite_event)
    RETURNS api.flip_bid_snapshot AS
$$
SELECT *
FROM api.get_flip(event.bid_id, event.ilk_identifier, event.block_height)
$$
    LANGUAGE sql
    STABLE;

