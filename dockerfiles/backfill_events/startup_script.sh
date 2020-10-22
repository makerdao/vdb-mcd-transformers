#!/bin/sh
echo "Starting event backfill to block 11108609"
# Fire up execute
./vulcanizedb backfillEvents --config config.toml -e 11108609 # TODO: update before next run
