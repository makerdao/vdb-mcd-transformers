-- +goose Up
CREATE TABLE public.addresses
(
    id             SERIAL PRIMARY KEY,
    address        character varying(42),
    hashed_address character varying(66),
    UNIQUE (address)
);

-- +goose Down
DROP TABLE public.addresses;