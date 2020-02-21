-- +goose Up
CREATE TABLE maker.bite
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    urn_id    INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
    ink       NUMERIC,
    art       NUMERIC,
    tab       NUMERIC,
    flip      TEXT,
    bid_id    NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX bite_header_index
    ON maker.bite (header_id);
CREATE INDEX bite_log_index
    ON maker.bite (log_id);
CREATE INDEX bite_urn_index
    ON maker.bite (urn_id);


-- +goose Down
DROP INDEX maker.bite_header_index;
DROP INDEX maker.bite_log_index;
DROP INDEX maker.bite_urn_index;

DROP TABLE maker.bite;
