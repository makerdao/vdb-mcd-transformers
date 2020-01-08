-- +goose Up
CREATE TABLE maker.log_value
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    val       NUMERIC,
    UNIQUE (header_id, log_id)
);
COMMENT ON COLUMN maker.log_value.id
    IS E'@omit';

CREATE INDEX log_value_header_index
    ON maker.log_value (header_id);
CREATE INDEX log_value_log_index
    ON maker.log_value (log_id);


-- +goose Down
DROP INDEX maker.log_value_log_index;
DROP INDEX maker.log_value_header_index;

DROP TABLE maker.log_value;
