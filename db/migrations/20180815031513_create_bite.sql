-- +goose Up
CREATE TABLE maker.bite
(
    id              SERIAL PRIMARY KEY,
    header_id       INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    urn_id          INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
    ink             NUMERIC,
    art             NUMERIC,
    tab             NUMERIC,
    flip            TEXT,
    bite_identifier NUMERIC,
    tx_idx          INTEGER NOT NULL,
    log_idx         INTEGER NOT NULL,
    raw_log         JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

COMMENT ON TABLE maker.bite
    IS E'@name raw_bites';
COMMENT ON COLUMN maker.bite.bite_identifier
    IS E'@name id';
COMMENT ON COLUMN maker.bite.id
    IS E'@omit';

CREATE INDEX bite_header_index
    ON maker.bite (header_id);

CREATE INDEX bite_urn_index
    ON maker.bite (urn_id);

ALTER TABLE public.checked_headers
    ADD COLUMN bite INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.bite_header_index;
DROP INDEX maker.bite_urn_index;

DROP TABLE maker.bite;

ALTER TABLE public.checked_headers
    DROP COLUMN bite;
