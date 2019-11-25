#!/bin/sh
# Runs the migrations and starts the headerSync and continuousLogSync services

# DEBUG
set -x

if test -z "$VDB_PG_CONNECT"; then
  # Exit if the variable tests fail
  set -e

  # Check the database variables are set
  test $DATABASE_NAME
  test $DATABASE_HOSTNAME
  test $DATABASE_PORT
  test $DATABASE_USER
  test $DATABASE_PASSWORD
  set +e

  # Construct the connection string for postgres
  VDB_PG_CONNECT=postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOSTNAME:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable
fi

# Run the DB migrations
echo "Connecting with: $VDB_PG_CONNECT"
./goose -dir db/migrations postgres "$VDB_PG_CONNECT" up

if [ $? -ne 0 ]; then
  echo "Could not run migrations. Are the database details correct?"
  exit 1
fi


if test -z "$STARTING_BLOCK_NUMBER"
then
    echo "STARTING_BLOCK_NUMBER is required and no value was given"
    exit 1
fi


echo "Starting headerSync and executing the transformers..."
# Fire up the services
if [ $? -eq 0 ]; then
  # Fire up the services
  ./vulcanizedb headerSync --config config.toml -s $STARTING_BLOCK_NUMBER &
  ./vulcanizedb execute --config config.toml &
fi


# Check every 60 seconds to see if either process has excited.
# If grepping for process names finds something, they exit with 0 status. If they are not both 0, then one of the processes has already excited.

while sleep 60; do
  ps aux | grep headerSync | grep -q -v grep
  HEADER_SYNC_STATUS=$?

  ps aux | grep execute | grep -q -v grep
  EXECUTE_STATUS=$?

  if [ $HEADER_SYNC_STATUS -ne 0 -o $EXECUTE_STATUS -ne 0 ]; then
    echo "One of the processes has already exited."
    exit 1
  fi
done
