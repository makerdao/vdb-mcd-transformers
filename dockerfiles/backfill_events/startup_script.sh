#!/bin/sh
echo "Starting event backfill to block 11376200"
# Fire up execute
./vulcanizedb backfillEvents --config config.toml -e 11376200 # TODO: update before next run
