-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.set_min_sell(
    id         serial primary key,
    log_id     bigint  not null references public.event_logs (id) on delete cascade,
    header_id  integer not null references public.headers (id) on delete cascade,
    address_id integer not null references public.addresses (id) on delete cascade,
    pay_gem    integer not null references public.addresses (id) on delete cascade,
    msg_sender integer not null references public.addresses (id) on delete cascade,
    dust       numeric,
    UNIQUE (header_id, log_id)
);

CREATE INDEX set_min_sell_header_index
    ON maker.set_min_sell (header_id);
CREATE INDEX set_min_sell_log_index
    ON maker.set_min_sell (log_id);
CREATE INDEX set_min_sell_address_index
    ON maker.set_min_sell (address_id);
CREATE INDEX set_min_sell_pay_gem_index
    ON maker.set_min_sell (pay_gem);
CREATE INDEX set_min_sell_msg_sender
    ON maker.set_min_sell (msg_sender);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE maker.set_min_sell;
