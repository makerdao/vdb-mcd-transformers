-- +goose Up
CREATE TABLE maker.osm_change
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    src        INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX osm_change_log_index
    ON maker.osm_change (log_id);
CREATE INDEX osm_change_header_index
    ON maker.osm_change (header_id);
CREATE INDEX osm_change_address_index
    ON maker.osm_change (address_id);
CREATE INDEX osm_change_msg_sender_index
    ON maker.osm_change (msg_sender);
CREATE INDEX osm_change_src_index
    ON maker.osm_change (src);


-- +goose Down
DROP TABLE maker.osm_change;
