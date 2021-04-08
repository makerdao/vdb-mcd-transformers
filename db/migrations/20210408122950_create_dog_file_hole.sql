-- +goose Up
CREATE TABLE maker.dog_file_hole
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    what       TEXT,
    data       BIGINT,
    UNIQUE (header_id, log_id)
);

CREATE INDEX dog_file_hole_header_index
    ON maker.dog_file_hole (header_id);
CREATE INDEX dog_file_hole_log_index
    ON maker.dog_file_hole (log_id);
CREATE INDEX dog_file_hole_address_index
    ON maker.dog_file_hole (address_id);
CREATE INDEX dog_file_hole_what_index
    ON maker.dog_file_hole (what);

-- +goose Down
DROP INDEX maker.dog_file_hole_header_index;
DROP INDEX maker.dog_file_hole_log_index;
DROP INDEX maker.dog_file_hole_address_index;
DROP INDEX maker.dog_file_hole_what_index;

DROP TABLE maker.dog_file_hole;
