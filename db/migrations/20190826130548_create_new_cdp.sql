-- +goose Up
CREATE TABLE maker.new_cdp
(
    id        SERIAL PRIMARY KEY,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    usr       TEXT,
    own       TEXT,
    cdp       NUMERIC,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX new_cdp_log_index
    ON maker.new_cdp (log_id);

-- +goose Down
DROP INDEX maker.new_cdp_log_index;
DROP TABLE maker.new_cdp;
