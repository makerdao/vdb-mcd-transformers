-- +goose Up
CREATE TABLE maker.checked_headers
(
    id          SERIAL PRIMARY KEY,
    check_count INTEGER,
    header_id   INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    UNIQUE (header_id)
);

CREATE INDEX checked_headers_header_index
    ON maker.checked_headers (header_id);
CREATE INDEX checked_headers_check_count
    ON maker.checked_headers (check_count);

-- +goose Down
DROP TABLE maker.checked_headers;
