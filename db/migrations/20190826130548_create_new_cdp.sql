-- +goose Up
CREATE TABLE maker.new_cdp
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    usr       TEXT,
    own       TEXT,
    cdp       NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX new_cdp_log_index
    ON maker.new_cdp (log_id);

ALTER TABLE public.checked_headers
    ADD COLUMN new_cdp INTEGER NOT NULL DEFAULT 0;

-- +goose Down
DROP INDEX maker.new_cdp_log_index;
DROP TABLE maker.new_cdp;
ALTER TABLE public.checked_headers
    DROP COLUMN new_cdp;
