When a new collateral is added to the MCD system, it also needs to be added to the `vdb-mcd-transformers` repository in order to track changes to the collateral's Flip, Median and OSM contracts.

## 1. Tracking New Collateral Events and Storage Changes
- New collateral(s) transformers need to be added to the `execute` command in order to track events and storage changes for the new collateral's Flip, Median and OSM contract. 
- Use the collateral generator to add new collateral config to vdb-mcd-transformers:
    1. After building the vdb-mcd-transformers project, run `./vdb-mcd-transformers addNewCollateral` from the `vdb-mcd-transformers` directory.
    1. The generator will prompt you to enter the following information about the new collateral:
        - Collateral Name: this can be found at [https://changelog.makerdao.com/](https://changelog.makerdao.com/).
        - Collateral Version: this can be found at [https://changelog.makerdao.com/](https://changelog.makerdao.com/), and can be formatted with underscore or period delimiters (e.g. 1_2_3 or 1.2.3).
        - Flip Address: this can be found at [https://changelog.makerdao.com/](https://changelog.makerdao.com/).
        - Flip Address ABI: this can be found at [https://changelog.makerdao.com/](https://changelog.makerdao.com/) or on [Etherscan](https://etherscan.io/).
        - Flip Address Deployment Block: this be can be found by looking at the creation transaction on [Etherscan](https://etherscan.io/).
    1. It will then ask if Median and OSM contracts are required, and prompt for the address, ABI and deployment block for each of those contract types if required.
        - There will be a new Flip contract for each new collateral, but the price feed contract(s) to be watched depends on the collateral type.
        - **Standard Crypto Currencies**: Standard crypto currency collaterals will have an OSM and a Median contract to track. The [changelog](https://changelog.makerdao.com/) should have a PIP contract address, which will be the OSM contract. To get the associated Median contract, look at the `src` value on the OSM contract, this will be the Median contract's address. An easy place to find this value is to look at the [Read Contract tab on Etherscan.](https://etherscan.io/address/0x8Df8f06DC2dE0434db40dcBb32a82A104218754c#readContract)
        - **Stable Coins**: In general, stable coin collateral prices are assumed to be $1, so there isn't an Oracle contract to track - so, OSM and Median contracts are unnecessary for stable coin collaterals.
        - **Liquidity Provider Tokens**: Oracle contracts for LP tokens are still in development, and vdb-mcd-transformers does not yet have event or storage transformers for these contracts.
     1. Once completed, the generator will have added the following components to the repo:
        -  To `mcdTransformers.toml` config file:
            - contract details
            - added the new contract to the appropriate event transformer exporters
            - storage transformer exporter
            - new storage transformer name
        - To ` plugins/execute/transformerExporter.go`: the new Flip StorageTransformerInitializer.
        - To `transformers/storage/flip/initializers`: a new initializer directory for the new collateral, and an initializer.go file.
    1. A couple of manual changes also need to be made.
        1. Add the new collateral's contract name to the [FlipV110ABI() method](https://github.com/makerdao/vdb-mcd-transformers/blob/prod/transformers/shared/constants/method.go#L57) - this will ensure that the new Flip contract has the same ABI as the existing contracts. If the ABI has changed, we will need to make sure to handle the new version. The `FlipV110ABI` method is for all Flip contracts from v1.1.0+.
        1. Add an integration test for the [deny](https://github.com/makerdao/vdb-mcd-transformers/blob/prod/transformers/integration_tests/deny_test.go) and [rely](https://github.com/makerdao/vdb-mcd-transformers/blob/prod/transformers/integration_tests/rely_test.go) events to ensure that the transformers are wired up correctly.
        1. The additions to the integration tests will also require a few changes in [transformers/test_data/config_values.go](https://github.com/makerdao/vdb-mcd-transformers/blob/prod/transformers/test_data/config_values.go).
            - Add the new contract to the [FlipV110Addresses() method](https://github.com/makerdao/vdb-mcd-transformers/blob/prod/transformers/test_data/config_values.go#L43)..
            - Add a helper method to get the new Flip contract's address, e.g. [https://github.com/makerdao/vdb-mcd-transformers/blob/prod/transformers/test_data/config_values.go#L83](https://github.com/makerdao/vdb-mcd-transformers/blob/prod/transformers/test_data/config_values.go#L83).
         
## 2. Backfill Collateral Events
- There will likely be a gap between when the new collateral's Flip and Price Feed (OSM, Median, LPOracle) contracts are deployed and when `vdb-mcd-transformers` begins watching for their new events. So, in order to ensure that no events are missed, a `backfillEvents` process needs to be started.
- The `backfillEvents` process needs to be started **after** the new collaterals have been added to the `execute` process, step 1.
- Once the backfill code changes have been merged in, the backfill process will be deployed and started automatically.
- Example Pull Request for configuring `backfillEvents`: [https://github.com/makerdao/vdb-mcd-transformers/pull/453/files](https://github.com/makerdao/vdb-mcd-transformers/pull/453/files).
- Once the `backfillEvents` process has completed, make sure to clear out the changes to `./environments/backfillEvents.toml`, `./plugins/backfill/transformerExporter.go` and
`./dockerfiles/backfill_events/startup_script.sh` so that the process isn't restarted when new changes are merged to the repository.

### Changes Required for `backfillEvents` command:
1. Add the following components the `backfillEvents.toml` config file. These components can all be copied and pasted from the `mcdTransformers.toml` file which was updated in the previous step with the collateral generator script.
    1. Contract Details
        - The Flip contract details are required. The Price Feed contracts details are required as well if they're relevant for the given collateral. See the note above to determine if an OSM or Median contracts are relevant for the new collateral.
        - contract details format:
            ```toml
            [contract.CONTRACT-NAME]
                  address = "contract address"
                  abi = "contract-abi"
                  deployed = 0 # contract deployment block
            ```
    1. Event Transformer Exporters
        - event transformer exporter format:
            ```toml
            [exporter.transformer_name]
                contracts = ["CONTRACT-NAME"]
                migrations = "db/migrations"
                path = "transformers/events/transformer_name/initializer"
                rank = "0"
                repository = "github.com/makerdao/vdb-mcd-transformers"
                type = "eth_event"
            ```
        - Make sure to fill out the contracts array with only the new collateral contracts that correspond with the given transformer.
        - [Example event transformer exporters for all collateral contracts (Flip and Price Feed contracts).](#all-exporter-example)
        - [Example event transformer exporters required for Flip contracts, there are currently 8 of them.](#flip-exporter-example)
        - [Example event transformer exporters required for Median contracts, there are currently 7 of them.](#median-exporter-example)
        - [Example event transformer exporters are required for OSM contracts, there are currently 2 of them.](#osm-exporter-example)
    1. Transformer Names
        - The transformer name for each of the transformer exporters added in the previous step need to be added to the `transformer_names` collection.
            ```toml
            transformerNames = [
                "auction_file", "deal", "dent", "deny", "flip_file_cat", "flip_kick",
                "log_median_price", "log_value", "median_diss_batch",
                "median_diss_single", "median_drop", "median_kiss_batch",
                "median_kiss_single", "median_lift", "osm_change", "rely",
                "tend", "tick", "yank"
            ]
            ```
        - Note: The Median and OSM transformer names should be omitted if those transformers are not included/necessary for this collateral.
1. Add transformer initializers to the backfill event transformerExporter file: `./plugins/backfill/transformerExporter.go`.
    - A transformer initializer needs to be added for each of the transformer exporters that were added to `backfillEvents.toml`.
    - The initializers can be found in the `execute` transformerExporter file: `./plugins/execute/transformerExporter.go`. It is important to make sure to only include the transformers for the new collateral. Also, the following example can be used:
        ```go
            []event.TransformerInitializer{
                    auction_file.EventTransformerInitializer,
                    deal.EventTransformerInitializer,
                    dent.EventTransformerInitializer,
                    deny.EventTransformerInitializer,
                    flip_file_cat.EventTransformerInitializer,
                    flip_kick.EventTransformerInitializer,
                    log_median_price.EventTransformerInitializer,
                    log_value.EventTransformerInitializer,
                    median_diss_batch.EventTransformerInitializer,
                    median_diss_single.EventTransformerInitializer,
                    median_drop.EventTransformerInitializer,
                    median_kiss_batch.EventTransformerInitializer,
                    median_kiss_single.EventTransformerInitializer,
                    median_lift.EventTransformerInitializer,
                    osm_change.EventTransformerInitializer,
                    rely.EventTransformerInitializer,
                    tend.EventTransformerInitializer,
                    tick.EventTransformerInitializer,
                    yank.EventTransformerInitializer,
                },
        ```
1. Update `./dockerfiles/backfill_events/startup_script.sh`:
    - The `backfillEvents` command allows for an ending block configuration which is set in the dockerfile startup script. The command will start at the earliest deployment block from all contracts that are configured in `backfillEvents.toml`.
    - The ending block number should be close to the block when the new collaterals have been added to the `execute` command in step 1.

## 3. Extract Storage Diffs From Collateral Contracts
- Once the `backfillEvents` process finishes the `extractDiffs` process can be restarted to include diffs from the new collateral contracts.
- It's preferred to wait until `backfillEvents` completes so that we have all of the events necessary to generate hashed storage keys in order to be able to transform storage diffs. If `extractDiffs` is started before `backfillEvents` completes, the new diff's keys would likely not be recognized, because their corresponding events (which allow us to decode the storage keys) may not be backfilled yet. In this case, the diffs would be marked as `unrecognized`  and would end up getting transformed later.
- Once the code changes are merged, the `extractDiffs` process is restarted and any contract addresses in the `extractDiffs.toml` file will be used in creating the geth diff subscription.
### Changes required for `extractDiffs` command:
Update `extractDiffs.toml` with the new collateral's contract details. E.g.:

    ```toml
    [contract.CONTRACT-NAME]
          address = "contract address"
          abi = "contract-abi"
          deployed = 0 # contract deployment block
    ```

## 4. Backfill Collateral Storage
- Like the `extractDiffs` storage process above, it's preferred to wait until `backfillEvents` completes so that we have all of the events necessary to generate hashed storage keys in order to be able to transform storage diffs. 
- Similar to how `backfillEvents` is updated, several components need to be added to the `backfillStorage.toml` config file as well. These components can all be copied and pasted from the `mcdTransformers.toml` file which was updated with the collateral generator script.
- Example Pull Request configuring `backfillStorage`: [https://github.com/makerdao/vdb-mcd-transformers/pull/455](https://github.com/makerdao/vdb-mcd-transformers/pull/455).
- Once the `backfillStorage` process has completed, make sure to clear out the changes to `./environments/backfillStorage.toml`, `./plugins/execute/storage/transformerExporter.go` and
`./dockerfiles/backfill_storage/startup_script.sh` so that the process isn't restarted when new changes are merged to the repository.

### Changes Required for `backfillStorage` command:

1. Add the following components the `backfillStorage.toml` config file.
    1. Contract Details
        - The Flip contract details are required. The Price Feed contracts details are required as well if they're relevant for the given collateral. See the note above to determine if an OSM or Median contracts are relevant for the new collateral.
        - contract details format:
            ```toml
            [contract.CONTRACT-NAME]
                  address = "contract address"
                  abi = "contract-abi"
                  deployed = 0 # contract deployment block
            ```
    1. Storage Transformer Exporters
        - There will be a storage transformer exporter for each collateral's Flip contract, and Median contract (if applicable). Note: vdb-mcd-transformers is not currently watching OSM contract storage slots.
        - Example storage transform exporters for the AAVE collateral:
            ```toml
            [exporter.flip_aave_a_v1_2_2]
                migrations = "db/migrations"
                path = "transformers/storage/flip/initializers/aave_a/v1_2_2"
                rank = "0"
                repository = "github.com/makerdao/vdb-mcd-transformers"
                type = "eth_storage"
            [exporter.median_aave_v1_2_2]
                migrations = "db/migrations"
                path = "transformers/storage/median/initializers/median_aave/v1_2_2"
                rank = "0"
                repository = "github.com/makerdao/vdb-mcd-transformers"
                type = "eth_storage"
            ```
    1. Transformer Names
        - The transformer name for each of the transformer exporters added in the previous step need to be added to the `transformer_names` collection.
            ```toml
            transformerNames = [
              "flip_aave_a_v1_2_2",
              "median_aave_v1_2_2",
            ]
            ```
        - Note: The median transformer name should be omitted if it is not included/necessary for this collateral.
1. Add transformer initializers to the backfill storage transformerExporter file: `./plugins/backfill/storage/transformerExporter.go`.
    - A transformer initializer needs to be added for each of the transformer exporters that were added to `backfillStorage.toml`.
    - The initializers can be found in the `execute` transformerExporter file: `./plugins/execute/transformerExporter.go`. It is important to make sure to only include the transformers for the new collateral.
        ```go
            []storage.TransformerInitializer{
                    flip_univ2daiusdc_a_v1_2_5.StorageTransformerInitializer,
                },
        ```
1. Update `./dockerfiles/backfill_storage/startup_script.sh`:
    - The `backfillStorage` command allows for configuring starting and ending blocks which are set in the dockerfile startup script.
    - The starting block number should be the earliest deployment block of the transformers that are being backfilled.
    - The ending block number should be close to the block when the new collaterals were added to `extractDiffs`.

## <a name="restarting-backfill"></a>Restarting A Backfill Process
If a backfill process stops either due to a failure, or an intentional pause, it may need to be restarted. 
1. Take a look at the logs from the last container that was running the backfill - this should give you information about the last block that was being processed when the process stopped.
1. Create a PR re-configuring the backfill including the startup script, the transformerExporter.go and the toml config file.
1. For the backfillStorage process, the starting block number will need to be updated in the startup script. It can now be set to the last block that the previous process was working on - to be sure that no blocks are missed, it may be a good idea to add in a small buffer of a couple of blocks before the block number from step 1.

## Gotchas:
- The `backfillStorage` process uses `getStorageAt` to get the storage values for all storage keys for the contract we're watching over the range of blocks. As it's inserting the backfilled diffs, it also checks that they are in fact new/updated values before inserting them. This process is not very performant, and we have seen it take about 24 hours to process 50,000 diffs - but please note that this is very dependent on how many contracts are being backfilled, and how many storage value changes have occurred on the given contracts. This particular benchmark is from a process backfilling seven contracts.
- `backfillEvents` and `backfillStorage` processes are both auto-deployed on merge to the vdb-mcd-transformers and vdb-oasis-transformers repositories. The implication of this is that every time a PR is merged into either of these repositories, it's associated backfill processes are restarted with the current configuration in the startup script, transformerExporter and config files. A few things to do/keep in mind:
    - You will not be able to merge another PR until the current backfill process is finished, otherwise it will be restarted.
    - Once the current backfill process successfully completes, it's a good idea clean up the backfill configuration so that a duplicate backfill process isn't kicked off on the next merge to the repository.. This means removing:
        - the starting and ending block number from the startup script
        - the transformer initializers from the transformerExporter.go file
        - the transformerNames and transformerExporters from the toml config file
    - The backfillStorage process is rather time consuming, and there is a chance that you may need to pause and restart the process, if another PR needs to be merged. To do this:
        - clean up the existing config (as mention in the previous bullet point) and merge it in to stop the backfill process
        - [restart the backfill process](#restarting-backfill) when all of the pressing PRs are merged in

---

## Event Transformer Exporter Examples:
### <a name="all-exporter-example"></a> Example event transformer exporters for all collateral contracts (Flip and Price Feed contracts).
```toml
[exporter.deny]
  contracts = [] # add flip, median and osm contracts here
  migrations = "db/migrations"
  path = "transformers/events/auth/deny_initializer"
  rank = "0"
  repository = "github.com/makerdao/vdb-mcd-transformers"
  type = "eth_event"

[exporter.rely]
  contracts = [] # add flip, median and osm contracts here
  migrations = "db/migrations"
  path = "transformers/events/auth/rely_initializer"
  rank = "0"
  repository = "github.com/makerdao/vdb-mcd-transformers"
  type = "eth_event"
```

### <a name="flip-exporter-example"></a> Example event transformer exporters required for Flip contracts, there are currently 8 of them.
```toml
[exporter.auction_file]
  contracts = [] # add flip contract(s) here
  migrations = "db/migrations"
  path = "transformers/events/auction_file/initializer"
  rank = "0"
  repository = "github.com/makerdao/vdb-mcd-transformers"
  type = "eth_event"

[exporter.deal]
  contracts = [] # add flip contract(s) here
  migrations = "db/migrations"
  path = "transformers/events/deal/initializer"
  rank = "0"
  repository = "github.com/makerdao/vdb-mcd-transformers"
  type = "eth_event"

[exporter.dent]
  contracts = [] # add flip contract(s) here
  migrations = "db/migrations"
  path = "transformers/events/dent/initializer"
  rank = "0"
  repository = "github.com/makerdao/vdb-mcd-transformers"
  type = "eth_event"

[exporter.flip_file_cat]
  contracts = [] # add flip contract(s) here
  migrations = "db/migrations"
  path = "transformers/events/flip_file/cat/initializer"
  rank = "0"
  repository = "github.com/makerdao/vdb-mcd-transformers"
  type = "eth_event"

[exporter.flip_kick]
  contracts = [] # add flip contract(s) here
  migrations = "db/migrations"
  path = "transformers/events/flip_kick/initializer"
  rank = "0"
  repository = "github.com/makerdao/vdb-mcd-transformers"
  type = "eth_event"

[exporter.tend]
  contracts = [] # add flip contract(s) here
  migrations = "db/migrations"
  path = "transformers/events/tend/initializer"
  rank = "0"
  repository = "github.com/makerdao/vdb-mcd-transformers"
  type = "eth_event"

[exporter.tick]
  contracts = [] # add flip contract(s) here
  migrations = "db/migrations"
  path = "transformers/events/tick/initializer"
  rank = "0"
  repository = "github.com/makerdao/vdb-mcd-transformers"
  type = "eth_event"

[exporter.yank]
  contracts = [] # add flip contract(s) here
  migrations = "db/migrations"
  path = "transformers/events/yank/initializer"
  rank = "0"
  repository = "github.com/makerdao/vdb-mcd-transformers"
  type = "eth_event"
```

### <a name="median-exporter-example"></a> Example event transformer exporters required for Median contracts, there are currently 7 of them.
```toml
[exporter.log_median_price]
    contracts = [] # add median contracts here
    migrations = "db/migrations"
    path = "transformers/events/log_median_price/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"

  [exporter.median_diss_batch]
    contracts = [] # add median contracts here
    migrations = "db/migrations"
    path = "transformers/events/median_diss/batch/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"

  [exporter.median_diss_single]
    contracts = [] # add median contracts here
    migrations = "db/migrations"
    path = "transformers/events/median_diss/single/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"

  [exporter.median_drop]
    contracts = [] # add median contracts here
    migrations = "db/migrations"
    path = "transformers/events/median_drop/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"

  [exporter.median_kiss_batch]
    contracts = [] # add median contracts here
    migrations = "db/migrations"
    path = "transformers/events/median_kiss/batch/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"

  [exporter.median_kiss_single]
    contracts = [] # add median contracts here
    migrations = "db/migrations"
    path = "transformers/events/median_kiss/single/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"

  [exporter.median_lift]
    contracts = [] # add median contracts here
    migrations = "db/migrations"
    path = "transformers/events/median_lift/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"
```

### <a name="osm-exporter-example"></a> Example event transformer exporters are required for OSM contracts, there are currently 2 of them.
```toml
[exporter.log_value]
    contracts = [] # add osm contracts here
    migrations = "db/migrations"
    path = "transformers/events/log_value/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"

[exporter.osm_change]
    contracts = [] # add osm contracts here
    migrations = "db/migrations"
    path = "transformers/events/osm_change/initializer"
    rank = "0"
    repository = "github.com/makerdao/vdb-mcd-transformers"
    type = "eth_event"
```
