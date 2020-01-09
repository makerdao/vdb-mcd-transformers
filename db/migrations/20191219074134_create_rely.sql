-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE maker.rely
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    address_id INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    usr        INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    UNIQUE (header_id, log_id)
);

CREATE INDEX rely_header_index
    ON maker.rely (header_id);
CREATE INDEX rely_log_index
    ON maker.rely (log_id);
CREATE INDEX rely_address_index
    ON maker.rely (address_id);
CREATE INDEX rely_msg_sender_index
    ON maker.rely (msg_sender);
CREATE INDEX rely_usr_index
    ON maker.rely (usr);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX maker.rely_header_index;
DROP INDEX maker.rely_log_index;
DROP INDEX maker.rely_address_index;
DROP INDEX maker.rely_msg_sender_index;
DROP INDEX maker.rely_usr_index;

DROP TABLE maker.rely;