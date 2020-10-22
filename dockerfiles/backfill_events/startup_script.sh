#!/bin/sh
echo "Starting event backfill to block UPDATE"
# Fire up execute
./vulcanizedb backfillEvents --config config.toml -e UPDATE # TODO: update before next run
