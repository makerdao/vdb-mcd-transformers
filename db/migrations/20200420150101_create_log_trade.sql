-- +goose Up
CREATE TABLE maker.log_trade
(
    id         serial primary key,
    log_id     bigint  not null references public.event_logs (id) on delete cascade,
    header_id  integer not null references public.headers (id) on delete cascade,
    address_id integer not null references public.addresses (id) on delete cascade,
    pay_gem    integer not null references public.addresses (id) on delete cascade,
    buy_gem    integer not null references public.addresses (id) on delete cascade,
    pay_amt    numeric,
    buy_amt    numeric,
    UNIQUE (header_id, log_id)
);

CREATE INDEX log_trade_header_index
    ON maker.log_trade (header_id);
CREATE INDEX log_trade_log_index
    ON maker.log_trade (log_id);
CREATE INDEX log_trade_address_index
    ON maker.log_trade (address_id);
CREATE INDEX log_trade_pay_gem_index
    ON maker.log_trade (pay_gem);
CREATE INDEX log_trade_buy_gem_index
    ON maker.log_trade (buy_gem);

-- +goose Down
DROP TABLE maker.log_trade;