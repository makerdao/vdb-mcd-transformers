-- +goose Up
CREATE TABLE maker.log_make
(
    id         serial primary key,
    log_id     bigint  not null references public.event_logs (id) on delete cascade,
    address_id bigint  not null references public.addresses (id) on delete cascade,
    maker      bigint  not null references public.addresses (id) on delete cascade,
    pay_gem    bigint  not null references public.addresses (id) on delete cascade,
    buy_gem    bigint  not null references public.addresses (id) on delete cascade,
    pay_amt    numeric,
    buy_amt    numeric,
    offer_id   numeric,
    pair       character varying(66),
    header_id  integer not null references public.headers (id) on delete cascade,
    timestamp  integer,
    UNIQUE (header_id, log_id)
);

CREATE INDEX log_make_header_index
    ON maker.log_make (header_id);
CREATE INDEX log_make_log_index
    ON maker.log_make (log_id);
CREATE INDEX log_make_address_index
    ON maker.log_make (address_id);
CREATE INDEX log_make_maker_index
    ON maker.log_make (maker);
CREATE INDEX log_make_pay_gem_index
    ON maker.log_make (pay_gem);
CREATE INDEX log_make_buy_gem_index
    ON maker.log_make (buy_gem);

-- +goose Down
DROP TABLE maker.log_make;