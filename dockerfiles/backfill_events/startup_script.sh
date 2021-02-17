#!/bin/bash

# The ending block number needs to be set to around the block number where VDB began tracking the collateral(s) events.
# In other words, just after the block number where the new collateral(s) was added to execute.
echo "Starting event backfill to block [11862497]"
# Fire up execute
./vulcanizedb backfillEvents --config config.toml -e 11862497 
