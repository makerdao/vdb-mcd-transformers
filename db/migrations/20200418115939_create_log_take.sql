-- +goose Up
CREATE TABLE maker.log_take
(
    id         serial primary key,
    log_id     bigint  not null references public.event_logs (id) on delete cascade,
    address_id bigint  not null references public.addresses (id) on delete cascade,
    maker      bigint  not null references public.addresses (id) on delete cascade,
    pay_gem    bigint  not null references public.addresses (id) on delete cascade,
    buy_gem    bigint  not null references public.addresses (id) on delete cascade,
    taker      bigint  not null references public.addresses (id) on delete cascade,
    take_amt   numeric,
    give_amt   numeric,
    offer_id   numeric,
    pair       character varying(66),
    header_id  integer not null references public.headers (id) on delete cascade,
    timestamp  integer,
    UNIQUE (header_id, log_id)
);

CREATE INDEX log_take_header_index
    ON maker.log_take (header_id);
CREATE INDEX log_take_log_index
    ON maker.log_take (log_id);
CREATE INDEX log_take_address_index
    ON maker.log_take (address_id);
CREATE INDEX log_take_maker_index
    ON maker.log_take (maker);
CREATE INDEX log_take_pay_gem_index
    ON maker.log_take (pay_gem);
CREATE INDEX log_take_buy_gem_index
    ON maker.log_take (buy_gem);
CREATE INDEX log_take_taker_index
    ON maker.log_take (taker);

-- +goose Down
DROP TABLE maker.log_take;