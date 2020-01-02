-- +goose Up
CREATE TABLE maker.deny
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    usr        INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX deny_header_index
    ON maker.deny (header_id);
CREATE INDEX deny_log_index
    ON maker.deny (log_id);
CREATE INDEX deny_address_index
    ON maker.deny (address_id);
CREATE INDEX deny_usr_index
    ON maker.deny (usr);


-- +goose Down
DROP INDEX maker.deny_usr_index;
DROP INDEX maker.deny_address_index;
DROP INDEX maker.deny_log_index;
DROP INDEX maker.deny_header_index;

DROP TABLE maker.deny;
