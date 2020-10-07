#!/bin/sh
# Runs the migrations and executes transformers

MISSING_VAR_MESSAGE=" is required and no value was given"

function testDatabaseVariables() {
  for a in DATABASE_NAME DATABASE_HOSTNAME DATABASE_PORT DATABASE_USER DATABASE_PASSWORD
  do
    eval arg="$"$a
    test $arg
    if [ $? -ne 0 ]; then
      echo $a $MISSING_VAR_MESSAGE
      exit 1
    fi
  done
}

if test -z "$VDB_PG_CONNECT"; then
  # Exits if the variable tests fail
  testDatabaseVariables
  if [ $? -ne 0 ]; then
    exit 1
  fi

  # Construct the connection string for postgres
  VDB_PG_CONNECT=postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOSTNAME:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable
fi

# Run the DB migrations
echo "Connecting with: $VDB_PG_CONNECT"
./goose -dir db/migrations postgres "$VDB_PG_CONNECT" up
psql $VDB_PG_CONNECT -c 'CREATE SCHEMA IF NOT EXISTS maker;'
./goose -table maker.goose_db_version -dir db/vdb-mcd-transformers/migrations postgres "$VDB_PG_CONNECT" up

if [ $? -ne 0 ]; then
  echo "Could not run migrations. Are the database details correct?"
  exit 1
fi
