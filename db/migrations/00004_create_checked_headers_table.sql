-- +goose Up
CREATE TABLE public.checked_headers
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER UNIQUE NOT NULL REFERENCES headers (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE public.checked_headers;
