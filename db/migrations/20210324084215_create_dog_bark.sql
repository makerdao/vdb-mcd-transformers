-- +goose Up
CREATE TABLE maker.dog_bark
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    ilk_id     INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    urn_id     INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
    clip       BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    ink        NUMERIC,
    art        NUMERIC,
    due        NUMERIC,
    sale_id    NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX dog_bark_header_index
    ON maker.dog_bark (header_id);
CREATE INDEX dog_bark_log_index
    ON maker.dog_bark (log_id);
CREATE INDEX dog_bark_ilk_index
    ON maker.dog_bark (ilk_id);
CREATE INDEX dog_bark_urn_index
    ON maker.dog_bark (urn_id);
CREATE INDEX dog_bark_address_index
    ON maker.dog_bark (address_id);
CREATE INDEX dog_bark_clip_index
    ON maker.dog_bark (clip);

-- +goose Down
DROP TABLE maker.dog_bark;
