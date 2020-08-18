-- +goose Up
CREATE TABLE maker.deal
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    msg_sender BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    bid_id     NUMERIC NOT NULL,
    UNIQUE (header_id, log_id)
);

CREATE INDEX deal_header_index
    ON maker.deal (header_id);
CREATE INDEX deal_log_index
    ON maker.deal (log_id);
CREATE INDEX deal_bid_id_index
    ON maker.deal (bid_id);
CREATE INDEX deal_address_index
    ON maker.deal (address_id);
CREATE INDEX deal_msg_sender_index
    ON maker.deal (msg_sender);


-- +goose Down
DROP TABLE maker.deal;
