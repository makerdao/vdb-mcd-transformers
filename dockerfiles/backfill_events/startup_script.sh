#!/bin/sh

if test -z "$ENDING_BLOCK_NUMBER"
then
    echo ENDING_BLOCK_NUMBER is required and no value was given
    exit 1
fi
echo "Starting event backfill from block $ENDING_BLOCK_NUMBER"
# Fire up execute
./vulcanizedb backfillEvents --config config.toml -e $ENDING_BLOCK_NUMBER
