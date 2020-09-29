#!/bin/sh
echo "Starting event backfill to block 10959649"
# Fire up execute
./vulcanizedb backfillEvents --config config.toml -e 10959649 # TODO: update before next run
