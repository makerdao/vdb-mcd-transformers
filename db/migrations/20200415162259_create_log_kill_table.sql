-- +goose Up
CREATE TABLE maker.log_kill
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    offer_id   numeric,
    pair       character varying(66),
    maker      INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    pay_gem    INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    buy_gem    INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    pay_amt    NUMERIC,
    buy_amt    NUMERIC,
    timestamp  BIGINT,
    UNIQUE (header_id, log_id)
);

CREATE INDEX log_kill_index
    ON maker.log_kill (log_id);
CREATE INDEX log_kill_header_index
    ON maker.log_kill (header_id);
CREATE INDEX log_kill_address_index
    ON maker.log_kill (address_id);
CREATE INDEX log_kill_maker_index
    ON maker.log_kill (maker);
CREATE INDEX log_kill_pay_gem_index
    ON maker.log_kill (pay_gem);
CREATE INDEX log_kill_buy_gem_index
    ON maker.log_kill (buy_gem);

-- +goose Down
DROP TABLE maker.log_kill;