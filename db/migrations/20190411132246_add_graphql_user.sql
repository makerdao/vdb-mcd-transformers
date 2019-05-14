-- +goose Up
-- Selectively grant execute to the graphql user to limit the API
REVOKE EXECUTE ON ALL FUNCTIONS IN SCHEMA maker FROM public;
ALTER DEFAULT PRIVILEGES IN SCHEMA maker REVOKE EXECUTE ON FUNCTIONS FROM public;

-- +goose StatementBegin
DO
  $do$
    BEGIN
      IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'graphql')
      THEN
        CREATE USER graphql WITH PASSWORD 'graphql';
      END IF;
    END
  $do$;
-- +goose StatementEnd
GRANT USAGE ON SCHEMA maker TO graphql;
GRANT EXECUTE ON FUNCTION maker.all_bites(text) TO graphql;
GRANT EXECUTE ON FUNCTION maker.all_frobs(text) TO graphql;
GRANT EXECUTE ON FUNCTION maker.all_urns(bigint) TO graphql;
GRANT EXECUTE ON FUNCTION maker.all_ilks(bigint) TO graphql;
GRANT EXECUTE ON FUNCTION maker.get_ilk(bigint, integer) TO graphql;
GRANT EXECUTE ON FUNCTION maker.all_ilk_states(bigint, integer) TO graphql;
GRANT EXECUTE ON FUNCTION maker.all_urn_states(text, text, bigint) TO graphql;
GRANT EXECUTE ON FUNCTION maker.get_urn(text, text, bigint) TO graphql;

-- +goose Down
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA maker FROM graphql;
REVOKE ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA maker FROM graphql;
REVOKE ALL PRIVILEGES ON SCHEMA maker FROM graphql;
