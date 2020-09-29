## Updating Addresses, ABIs, Deployment Blocks, and Integration Tests for Maker Releases
[MakerDAO Changelog](https://changelog.makerdao.com/)

Visit the changelog for the latest release, ie https://changelog.makerdao.com/

## Overview
### Updating Toml Files
Use the Contract Addresses and ABIs from the changelog to update the toml config file: `mcdTransformers.toml`.

Replace the contract address in the toml files with the updated contract address from the changelog.

Use the new address to update the deployment block in the toml files:
1. Search for the contract address on [Mainnet](https://etherscan.io/)
2. Under "More Info", click the transaction link next to "Contract Creator". This loads a Transaction Details page.
3. Replace the deployment block in the toml files with the block from the Transaction Details page.

Replace the ABI in the toml files with the updated ABI from the changelog.

### Referencing the Updated Contract
To find the source of the deployed contract:
https://github.com/makerdao/dss/releases

### Updating Tests
Since we've updated the addresses and deployment blocks, we need to make sure our tests are up to date.
1. Go to etherscan and make sure you're on the page for the contract address that you updated, and click Events. The url should look like https://etherscan.io/address/{0xAddress}#events.
2. Visit the `/transformers/integration_tests` directory, and find tests for the contract that you're updating. For example, if you are updating `MCD_JUG`, you'll see tests in with filenames `jug_[something]`. Choose a file to update, for example, `jug_init`.
3. Find the signature for the given contract function from `/shared/constants/signature_test.go`
4. On etherscan, we want to filter the events by the signature, which is `topic0`. To do this, paste the signature in the search bar in the events pane.
5. If there are no results, then we don't have to update any tests. If there are results, scroll down to the bottom of the page, and copy the block number. Use this block number in the relevant test.
6. Run the integration tests and fix discrepancies. You can often validate discrepancies by converting the `topic1`, `topic2`, and `topic3` hex values from the block to strings, but it depends on how the expectation is written.
