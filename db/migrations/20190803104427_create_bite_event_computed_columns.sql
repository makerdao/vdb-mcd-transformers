-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Extend type bite_event with ilk field
CREATE FUNCTION api.bite_event_ilk(event api.bite_event)
    RETURNS api.ilk_state AS
$$
SELECT *
FROM api.get_ilk(event.ilk_identifier, event.block_height)
$$
    LANGUAGE sql
    STABLE;


-- Extend type bite_event with urn field
CREATE FUNCTION api.bite_event_urn(event api.bite_event)
    RETURNS api.urn_state AS
$$
SELECT *
FROM api.get_urn(event.ilk_identifier, event.urn_identifier, event.block_height)
$$
    LANGUAGE sql
    STABLE;


-- Extend type bite_event with bid field
CREATE FUNCTION api.bite_event_bid(event api.bite_event)
    RETURNS api.flip_state AS
$$
SELECT *
FROM api.get_flip(event.bid_id, event.ilk_identifier, event.block_height)
$$
    LANGUAGE sql
    STABLE;


-- Extend type bite_event with txs field
CREATE FUNCTION api.bite_event_tx(event api.bite_event)
    RETURNS api.tx AS
$$
SELECT txs.hash, txs.tx_index, headers.block_number, headers.hash, tx_from, tx_to
FROM public.header_sync_transactions txs
         LEFT JOIN headers ON txs.header_id = headers.id
WHERE block_number <= event.block_height
  AND txs.tx_index = event.tx_idx
ORDER BY block_number DESC
$$
    LANGUAGE sql
    STABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION api.bite_event_tx(api.bite_event);
DROP FUNCTION api.bite_event_bid(api.bite_event);
DROP FUNCTION api.bite_event_urn(api.bite_event);
DROP FUNCTION api.bite_event_ilk(api.bite_event);
