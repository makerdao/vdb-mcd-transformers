-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.spot_poke
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    value     NUMERIC,
    spot      NUMERIC,
    log_idx   INTEGER NOT NULL,
    tx_idx    INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

CREATE INDEX spot_poke_header_index
    ON maker.spot_poke (header_id);

CREATE INDEX spot_poke_ilk_index
    ON maker.spot_poke (ilk_id);

ALTER TABLE public.checked_headers
    ADD COLUMN spot_poke INTEGER NOT NULL DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.spot_poke_header_index;
DROP INDEX maker.spot_poke_ilk_index;

DROP TABLE maker.spot_poke;

ALTER TABLE public.checked_headers
    DROP COLUMN spot_poke;
