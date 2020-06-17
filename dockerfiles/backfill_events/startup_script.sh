#!/bin/sh

echo "Starting event backfill..."
# Fire up execute
./vulcanizedb backfillEvents --config config.toml
