-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION api.clip_sale_event_kick_event(event api.clip_sale_event)
    RETURNS maker.clip_kick
    LANGUAGE sql STABLE
AS $$
SELECT * FROM maker.clip_kick WHERE log_id = event.log_id
$$;

CREATE FUNCTION api.clip_sale_event_take_event(event api.clip_sale_event)
    RETURNS maker.clip_take
    LANGUAGE sql STABLE
AS $$
SELECT * FROM maker.clip_take WHERE log_id = event.log_id
$$;

CREATE FUNCTION api.clip_sale_event_redo_event(event api.clip_sale_event)
    RETURNS maker.clip_redo
    LANGUAGE sql STABLE
AS $$
SELECT * FROM maker.clip_redo WHERE log_id = event.log_id
$$;

CREATE FUNCTION api.clip_sale_event_yank_event(event api.clip_sale_event)
    RETURNS maker.clip_yank
    LANGUAGE sql STABLE
AS $$
SELECT * FROM maker.clip_yank WHERE log_id = event.log_id
$$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION api.clip_sale_event_kick_event(api.clip_sale_event);
DROP FUNCTION api.clip_sale_event_take_event(api.clip_sale_event);
DROP FUNCTION api.clip_sale_event_redo_event(api.clip_sale_event);
DROP FUNCTION api.clip_sale_event_yank_event(api.clip_sale_event);
