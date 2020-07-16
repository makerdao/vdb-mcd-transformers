#!/bin/sh
source run_migrations.sh

DIFF_STARTING_BLOCK="-1"

if test -n "$DIFF_BLOCK_FROM_HEAD_OF_CHAIN"; then
    DIFF_STARTING_BLOCK=$DIFF_BLOCK_FROM_HEAD_OF_CHAIN
fi

echo "Starting transformer execution..."
# Fire up execute
./vulcanizedb execute --config config.toml -d $DIFF_STARTING_BLOCK
