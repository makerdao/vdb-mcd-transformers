-- +goose Up
CREATE TABLE maker.spot_file_mat
(
    id          SERIAL PRIMARY KEY,
    header_id   INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id      BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    ilk_id      INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    msg_sender  INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    what        TEXT,
    data        NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX spot_file_mat_header_index
    ON maker.spot_file_mat (header_id);
CREATE INDEX spot_file_mat_log_index
    ON maker.spot_file_mat (log_id);
CREATE INDEX spot_file_mat_ilk_index
    ON maker.spot_file_mat (ilk_id);
CREATE INDEX spot_file_mat_msg_sender_index
    ON maker.spot_file_mat (msg_sender);

CREATE TABLE maker.spot_file_par
(
    id          SERIAL PRIMARY KEY,
    header_id   INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id      BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    msg_sender  INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    what        TEXT,
    data        NUMERIC,
    UNIQUE (header_id, log_id)
);

CREATE INDEX spot_file_par_header_index
    ON maker.spot_file_par (header_id);
CREATE INDEX spot_file_par_log_index
    ON maker.spot_file_par (log_id);
CREATE INDEX spot_file_par_msg_sender_index
    ON maker.spot_file_par (msg_sender);

CREATE TABLE maker.spot_file_pip
(
    id          SERIAL PRIMARY KEY,
    header_id   INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    log_id      BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
    ilk_id      INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
    msg_sender  INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    what        TEXT,
    pip         TEXT,
    UNIQUE (header_id, log_id)
);

CREATE INDEX spot_file_pip_ilk_index
    ON maker.spot_file_pip (ilk_id);
CREATE INDEX spot_file_pip_log_index
    ON maker.spot_file_pip (log_id);
CREATE INDEX spot_file_pip_header_index
    ON maker.spot_file_pip (header_id);
CREATE INDEX spot_file_pip_msg_sender_index
    ON maker.spot_file_pip (msg_sender);


-- +goose Down
DROP TABLE maker.spot_file_mat;
DROP TABLE maker.spot_file_par;
DROP TABLE maker.spot_file_pip;
