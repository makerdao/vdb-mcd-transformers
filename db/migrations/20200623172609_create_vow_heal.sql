-- +goose Up
CREATE TABLE maker.vow_heal
(
    id         serial primary key,
    log_id     bigint  not null references public.event_logs (id) on delete cascade,
    msg_sender bigint  not null references public.addresses (id) on delete cascade,
    rad        numeric,
    header_id  integer not null references public.headers (id) on delete cascade,
    UNIQUE (header_id, log_id)
);

CREATE INDEX vow_heal_header_index
    ON maker.vow_heal (header_id);
CREATE INDEX vow_heal_log_index
    ON maker.vow_heal (log_id);
CREATE INDEX vow_heal_msg_sender
    ON maker.vow_heal (msg_sender);

-- +goose Down
DROP TABLE maker.vow_heal;
