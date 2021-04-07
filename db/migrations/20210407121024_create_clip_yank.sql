-- +goose Up
CREATE TABLE maker.clip_yank
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    sale_id    NUMERIC NOT NULL,
    UNIQUE (header_id, log_id)
);

CREATE INDEX clip_yank_header_index
    ON maker.clip_yank (header_id);
CREATE INDEX clip_yank_sale_id_index
    ON maker.clip_yank (sale_id);
CREATE INDEX clip_yank_address_index
    ON maker.clip_yank (address_id);
CREATE INDEX clip_yank_log_index
    ON maker.clip_yank (log_id);

-- +goose Down
DROP INDEX maker.clip_yank_log_index;
DROP INDEX maker.clip_yank_address_index;
DROP INDEX maker.clip_yank_sale_id_index;
DROP INDEX maker.clip_yank_header_index;

DROP TABLE maker.clip_yank;
