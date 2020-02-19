-- +goose Up
CREATE TABLE maker.pot_file_dsr
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    what      TEXT,
    data      NUMERIC,
    UNIQUE (header_id, log_id)
);

COMMENT ON TABLE maker.pot_file_dsr
    IS E'Note event emitted when file(bytes32,uint256) is called on Pot contract.';

CREATE INDEX pot_file_dsr_header_index
    ON maker.pot_file_dsr (header_id);
CREATE INDEX pot_file_dsr_log_index
    ON maker.pot_file_dsr (log_id);

CREATE TABLE maker.pot_file_vow
(
    id        SERIAL PRIMARY KEY,
    header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    what      TEXT,
    data      TEXT,
    UNIQUE (header_id, log_id)
);

COMMENT ON TABLE maker.pot_file_vow
    IS E'Note event emitted when file(bytes32,address) is called on Pot contract.';

CREATE INDEX pot_file_vow_header_index
    ON maker.pot_file_vow (header_id);
CREATE INDEX pot_file_vow_log_index
    ON maker.pot_file_vow (log_id);


-- +goose Down
DROP INDEX maker.pot_file_dsr_header_index;
DROP INDEX maker.pot_file_dsr_log_index;
DROP INDEX maker.pot_file_vow_header_index;
DROP INDEX maker.pot_file_vow_log_index;

DROP TABLE maker.pot_file_dsr;
DROP TABLE maker.pot_file_vow;
