#!/bin/sh
echo "Starting event backfill to block <NEEDS UPDATE>"
# Fire up execute
./vulcanizedb backfillEvents --config config.toml -e 10875761 # TODO: update before next run
