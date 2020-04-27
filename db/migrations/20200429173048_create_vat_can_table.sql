-- +goose Up
CREATE TABLE maker.vat_can
(
    id        SERIAL PRIMARY KEY,
    diff_id   BIGINT  NOT NULL REFERENCES public.storage_diff (id) ON DELETE CASCADE,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    bit       INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    usr       INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    can       NUMERIC NOT NULL,
    UNIQUE (diff_id, header_id, bit, usr, can)
);

CREATE INDEX vat_can_header_id_index
    ON maker.vat_can (header_id);
CREATE INDEX vat_can_bit_index
    ON maker.vat_can (bit);
CREATE INDEX vat_can_usr_index
    ON maker.vat_can (usr);

-- +goose Down
DROP INDEX maker.vat_can_header_id_index;
DROP INDEX maker.vat_can_bit_index;
DROP INDEX maker.vat_can_usr_index;

DROP TABLE maker.vat_can;