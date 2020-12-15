#!/bin/bash
source run_migrations.sh

NEW_DIFF_STARTING_BLOCK="-1"
UNRECOGNIZED_DIFF_STARTING_BLOCK="-1"

if test -n "$NEW_DIFF_BLOCK_FROM_HEAD_OF_CHAIN"; then
    NEW_DIFF_STARTING_BLOCK=$NEW_DIFF_BLOCK_FROM_HEAD_OF_CHAIN
fi

if test -n "$UNRECOGNIZED_DIFF_BLOCK_FROM_HEAD_OF_CHAIN"; then
    UNRECOGNIZED_DIFF_STARTING_BLOCK=$UNRECOGNIZED_DIFF_BLOCK_FROM_HEAD_OF_CHAIN
fi
echo "Starting transformer execution..."
# Fire up execute
./vulcanizedb execute --config config.toml --recheck-headers -d $NEW_DIFF_STARTING_BLOCK -u $UNRECOGNIZED_DIFF_STARTING_BLOCK
