#!/bin/bash
# Starts the getStorageValue command

# Verify required args present
MISSING_VAR_MESSAGE=" is required and no value was given"

function testDatabaseVariables() {
  for a in DATABASE_NAME DATABASE_HOSTNAME DATABASE_PORT DATABASE_USER DATABASE_PASSWORD BACKFILL_START_BLOCK BACKFILL_END_BLOCK
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

# Run backfillStorage
echo "Running backfillStorage from block $BACKFILL_START_BLOCK to $BACKFILL_END_BLOCK"
./vulcanizedb backfillStorage -s=$BACKFILL_START_BLOCK -e=$BACKFILL_END_BLOCK --config config.toml
