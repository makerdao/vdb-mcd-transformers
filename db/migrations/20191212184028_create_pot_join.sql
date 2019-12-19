-- +goose Up
CREATE TABLE maker.pot_join
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    log_id     BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
    msg_sender INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    wad        NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX pot_join_header_index
    ON maker.pot_join (header_id);
CREATE INDEX pot_join_log_index
    ON maker.pot_join (log_id);
CREATE INDEX pot_join_msg_sender_index
    ON maker.pot_join (msg_sender);


-- +goose Down
DROP INDEX maker.pot_join_msg_sender_index;
DROP INDEX maker.pot_join_log_index;
DROP INDEX maker.pot_join_header_index;
DROP TABLE maker.pot_join;
