-- +goose Up
ALTER TABLE maker.bid_event
    ALTER CONSTRAINT bid_event_log_id_fkey DEFERRABLE INITIALLY DEFERRED;

-- +goose Down
ALTER TABLE maker.bid_event
    ALTER CONSTRAINT bid_event_log_id_fkey NOT DEFERRABLE;
