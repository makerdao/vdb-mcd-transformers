#!/bin/bash
source run_migrations.sh

echo "Starting transformer execution from 200000 blocks back"
# Fire up execute
# rescanning diffs in last 200,000 to cover missed diffs between blocks 11234894 and 11315955
./vulcanizedb execute --config config.toml -d 200000
