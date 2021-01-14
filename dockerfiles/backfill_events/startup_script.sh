#!/bin/bash

# The ENDING_BLOCK_NUMBER needs to be set to around the block number where VDB began tracking the collateral(s) events.
# In other words, just after the block number where the new collateral(s) was added to exeucte.
if test -z "$ENDING_BLOCK_NUMBER"
then
    echo ENDING_BLOCK_NUMBER is required and no value was given
    exit 1
fi
echo "Starting event backfill from block 11602562"
# Fire up execute
./vulcanizedb backfillEvents --config config.toml -e 11602562
