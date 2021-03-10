-- +goose Up
CREATE TABLE maker.dog_file_ilk_unit
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    what       TEXT,
    data       NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX dog_file_ilk_unit_header_index
    ON maker.dog_file_ilk_unit (header_id);
CREATE INDEX dog_file_ilk_unit_log_index
    ON maker.dog_file_ilk_unit (log_id);
CREATE INDEX dog_file_ilk_unit_address_index
    ON maker.dog_file_ilk_unit (address_id);
CREATE INDEX dog_file_ilk_unit_ilk_index
    ON maker.dog_file_ilk_unit (ilk_id);

-- +goose Down
DROP INDEX maker.dog_file_ilk_unit_header_index;
DROP INDEX maker.dog_file_ilk_unit_log_index;
DROP INDEX maker.dog_file_ilk_unit_address_index;
DROP INDEX maker.dog_file_ilk_unit_ilk_index;

DROP TABLE maker.dog_file_ilk_unit;