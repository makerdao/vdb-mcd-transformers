#!/bin/sh
echo "Starting event backfill from block 10875761"
# Fire up execute
./vulcanizedb backfillEvents --config config.toml -e 10875761
