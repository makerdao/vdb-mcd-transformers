-- +goose Up
CREATE TABLE maker.bite
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    urn_id    INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
    ink       NUMERIC,
    art       NUMERIC,
    tab       NUMERIC,
    flip      TEXT,
    tx_idx    INTEGER NOT NULL,
    log_idx   INTEGER NOT NULL,
    raw_log   JSONB,
    UNIQUE (header_id, tx_idx, log_idx)
);

COMMENT ON TABLE maker.bite
    IS E'@name raw_bites';

ALTER TABLE public.checked_headers
    ADD COLUMN bite_checked BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
DROP TABLE maker.bite;

ALTER TABLE public.checked_headers
    DROP COLUMN bite_checked;
