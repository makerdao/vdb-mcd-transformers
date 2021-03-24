-- +goose Up
CREATE TABLE maker.dog_file_ilk_chop_hole
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

CREATE INDEX dog_file_ilk_chop_hole_header_index
    ON maker.dog_file_ilk_chop_hole (header_id);
CREATE INDEX dog_file_ilk_chop_hole_log_index
    ON maker.dog_file_ilk_chop_hole (log_id);
CREATE INDEX dog_file_ilk_chop_hole_address_index
    ON maker.dog_file_ilk_chop_hole (address_id);
CREATE INDEX dog_file_ilk_chop_hole_ilk_index
    ON maker.dog_file_ilk_chop_hole (ilk_id);

-- +goose Down
DROP INDEX maker.dog_file_ilk_chop_hole_header_index;
DROP INDEX maker.dog_file_ilk_chop_hole_log_index;
DROP INDEX maker.dog_file_ilk_chop_hole_address_index;
DROP INDEX maker.dog_file_ilk_chop_hole_ilk_index;

DROP TABLE maker.dog_file_ilk_chop_hole;