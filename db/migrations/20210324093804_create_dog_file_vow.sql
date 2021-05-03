-- +goose Up
CREATE TABLE maker.dog_file_vow
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    what       TEXT,
    data       NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX dog_file_vow_header_index
    ON maker.dog_file_vow (header_id);
CREATE INDEX dog_file_vow_log_index
    ON maker.dog_file_vow (log_id);
CREATE INDEX dog_file_vow_address_index
    ON maker.dog_file_vow (address_id);

-- +goose Down
DROP INDEX maker.dog_file_vow_header_index;
DROP INDEX maker.dog_file_vow_log_index;
DROP INDEX maker.dog_file_vow_address_index;

DROP TABLE maker.dog_file_vow;
