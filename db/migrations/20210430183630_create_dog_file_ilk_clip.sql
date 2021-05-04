-- +goose Up
CREATE TABLE maker.dog_file_ilk_clip
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    what       TEXT,
    clip_id    BIGINT,
    UNIQUE (header_id, log_id)
);

CREATE INDEX dog_file_ilk_clip_header_index
    ON maker.dog_file_ilk_clip (header_id);
CREATE INDEX dog_file_ilk_clip_log_index
    ON maker.dog_file_ilk_clip (log_id);
CREATE INDEX dog_file_ilk_clip_address_index
    ON maker.dog_file_ilk_clip (address_id);
CREATE INDEX dog_file_ilk_clip_ilk_index
    ON maker.dog_file_ilk_clip (ilk_id);
CREATE INDEX dog_file_ilk_clip_what_index
    ON maker.dog_file_ilk_clip (what);
CREATE INDEX dog_file_ilk_clip_clip_index
    ON maker.dog_file_ilk_clip (clip_id);

-- +goose Down
DROP INDEX maker.dog_file_ilk_clip_header_index;
DROP INDEX maker.dog_file_ilk_clip_log_index;
DROP INDEX maker.dog_file_ilk_clip_address_index;
DROP INDEX maker.dog_file_ilk_clip_ilk_index;
DROP INDEX maker.dog_file_ilk_clip_what_index;
DROP INDEX maker.dog_file_ilk_clip_clip_index;

DROP TABLE maker.dog_file_ilk_clip;
