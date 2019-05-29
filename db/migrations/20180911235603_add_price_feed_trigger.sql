-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION notify_pip_log_value()
  RETURNS trigger AS $$
BEGIN
  PERFORM pg_notify(
      CAST('postgraphile:pip_log_value' AS text),
      json_build_object('__node__', json_build_array('pip_log_value', NEW.id)) :: text
  );
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER notify_pip_log_value
  AFTER INSERT
  ON maker.pip_log_value
  FOR EACH ROW
EXECUTE PROCEDURE notify_pip_log_value();

-- +goose Down
DROP TRIGGER notify_pip_log_value
ON maker.pip_log_value;
