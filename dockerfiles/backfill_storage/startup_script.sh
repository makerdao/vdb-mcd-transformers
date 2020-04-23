#!/bin/sh
# Starts the getStorageValue command

# Verify required args present
MISSING_VAR_MESSAGE=" is required and no value was given"

if test -z "$STARTING_BLOCK_NUMBER"
then
    echo STARTING_BLOCK_NUMBER $MISSING_VAR_MESSAGE
    exit 1
fi

if test -z "$ENDING_BLOCK_NUMBER"
then
    echo ENDING_BLOCK_NUMBER $MISSING_VAR_MESSAGE
    exit 1
fi

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

# Run getStorageValue
echo "Running backfillStorage from block $STARTING_BLOCK_NUMBER to $ENDING_BLOCK_NUMBER"
./vulcanizedb backfillStorage -s=$STARTING_BLOCK_NUMBER -e=$ENDING_BLOCK_NUMBER --config config.toml
