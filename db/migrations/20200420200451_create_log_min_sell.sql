-- +goose Up
CREATE TABLE maker.log_min_sell
(
    id         serial primary key,
    log_id     bigint  not null references public.event_logs (id) on delete cascade,
    address_id bigint  not null references public.addresses (id) on delete cascade,
    pay_gem    bigint  not null references public.addresses (id) on delete cascade,
    min_amount numeric,
    header_id  integer not null references public.headers (id) on delete cascade,
    UNIQUE (header_id, log_id)
);

CREATE INDEX log_min_sell_header_index
    ON maker.log_min_sell (header_id);
CREATE INDEX log_min_sell_log_index
    ON maker.log_min_sell (log_id);
CREATE INDEX log_min_sell_address_index
    ON maker.log_min_sell (address_id);
CREATE INDEX log_min_sell_pay_gem_index
    ON maker.log_min_sell (pay_gem);

-- +goose Down
DROP TABLE maker.log_min_sell;
