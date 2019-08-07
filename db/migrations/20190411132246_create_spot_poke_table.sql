-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.spot_poke
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    value     NUMERIC,
    spot      NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX spot_poke_header_index
    ON maker.spot_poke (header_id);

CREATE INDEX spot_poke_ilk_index
    ON maker.spot_poke (ilk_id);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.spot_poke_header_index;
DROP INDEX maker.spot_poke_ilk_index;

DROP TABLE maker.spot_poke;
