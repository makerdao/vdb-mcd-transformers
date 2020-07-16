-- +goose Up
CREATE TABLE maker.log_delete
(
    id         serial primary key,
    log_id     bigint  not null references public.event_logs (id) on delete cascade,
    header_id  integer not null references public.headers (id) on delete cascade,
    address_id integer not null references public.addresses (id) on delete cascade,
    keeper    integer not null references public.addresses (id) on delete cascade,
    offer_id numeric,
    UNIQUE (header_id, log_id)
);

CREATE INDEX log_delete_header_index
    ON maker.log_delete (header_id);
CREATE INDEX log_delete_log_index
    ON maker.log_delete (log_id);
CREATE INDEX log_delete_address_index
    ON maker.log_delete (address_id);
CREATE INDEX log_delete_keeper_index
    ON maker.log_delete (keeper);

-- +goose Down
DROP TABLE maker.log_delete;