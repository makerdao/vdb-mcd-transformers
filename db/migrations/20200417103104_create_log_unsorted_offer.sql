-- +goose Up
CREATE TABLE maker.log_unsorted_offer
(
    id         SERIAL PRIMARY KEY,
    log_id     BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    header_id  INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    offer_id   NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX log_unsorted_offer_header_index
    ON maker.log_unsorted_offer (header_id);
CREATE INDEX log_unsorted_offer_log_index
    ON maker.log_unsorted_offer (log_id);
CREATE INDEX log_unsorted_offer_address_index
    ON maker.log_unsorted_offer (address_id);

-- +goose Down
DROP TABLE maker.log_unsorted_offer;
