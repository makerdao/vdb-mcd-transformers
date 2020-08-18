-- +goose Up
CREATE TABLE maker.log_insert
(
    id         serial primary key,
    log_id     bigint  not null references public.event_logs (id) on delete cascade,
    address_id bigint  not null references public.addresses (id) on delete cascade,
    keeper     bigint  not null references public.addresses (id) on delete cascade,
    offer_id   numeric,
    header_id  integer not null references public.headers (id) on delete cascade,
    UNIQUE (header_id, log_id)
);

CREATE INDEX log_insert_header_index
    ON maker.log_insert (header_id);
CREATE INDEX log_insert_log_index
    ON maker.log_insert (log_id);
CREATE INDEX log_insert_address_index
    ON maker.log_insert (address_id);
CREATE INDEX log_insert_keeper_index
    ON maker.log_insert (keeper);

-- +goose Down
DROP TABLE maker.log_insert;