When a new collateral is added to the MCD system, it also needs to be added to the `vdb-mcd-transformers` repository in order to track changes to the collateral's Flip, Median and OSM contracts.

## 1. Tracking New Collateral Events and Storage Changes
- New collateral(s) transformers need to be added to the `execute` command in order to track events and storage changes for the new collateral's Flip, Median and OSM contract. 
- Use the collateral generator by running `./vdb-mcd-transformers addNewCollateral` from the `vdb-mcd-transformers` directory.
- The collateral generator adds the following components to the `mcdTransformers.toml` config file for each new collateral:
    - contract details
    - event transformer exporters
    - storage transformer exporters
    - transformer names
- Note: There will be a new Flip contract for each new collateral, but the price feed contract(s) to be watched depends on the collateral type.
    - **Standard Crypto Currencies**: Standard crypto currency collaterals will have an OSM and a Median contract to track. The [changelog](https://changelog.makerdao.com/) should have a PIP contract address, which will be the OSM contract. To get the associated Median contract, look at the `src` value on the OSM contract, this will be the Median contract's address. An easy place to find this value is to look at the [Read Contract tab on Etherscan.](https://etherscan.io/address/0x8Df8f06DC2dE0434db40dcBb32a82A104218754c#readContract)
    - **Stable Coins**: In general, stable coin collateral prices are assumed to be $1, so there isn't an Oracle contract to track - so, OSM and Median contracts are unnecessary for stable coin collaterals.
    - **Liquidity Provider Tokens**: Oracle contracts for LP tokens are still in development, and vdb-mcd-transformers does not yet have event or storage transformers for these contracts.
    
## 2. Backfill Collateral Events
- There will likely be a gap between when the new collateral's Flip and Price Feed (OSM, Median, LPOracle) contracts are deployed and when `vdb-mcd-transformers` begins watching for their new events. So, in order to ensure that no events are missed, a `backfillEvents` process needs to be started.
- The `backfillEvents` process needs to be started **after** the new collaterals have been added to the `execute` process, step 1.
- Several components need to be added to the `backfillEvents.toml` config file. These components can all be copied and pasted from the `mcdTransformers.toml` file which was updated in the previous step with the collateral generator script.
    1. Contract Details
        - The Flip contract details are required. The Price Feed contract details are required as well if they're relevant for the given collateral. See the note above to determine if an OSM or Median contracts are relevant for the new collateral.
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
- Add transformer initializers to the backfill event transformerExporter file: `./plugins/backfill/transformerExporter.go`.
    - A transformer initializer needs to be added for each of the transformer exporters that was added to `backfillEvents.toml`.
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
- Update `./dockerfiles/backfill_events/startup_script.sh`:
    - The `backfillEvents` command allows for an ending block configuration which is set in the dockerfile startup script. The command will start at the earliest deployment block from all contracts that are configured in `backfillEvents.toml`.
    - The ending block number should be close to the block when the new collaterals have been added to the `execute` command in step 1.
- Once the backfill code changes have been merged in, the backfill process will be deployed and started automatically.

## 3. Extract Storage Diffs From Collateral Contracts
- Once the `backfillEvents` process finishes the `extractDiffs` process can be restarted to include diffs from the new collateral contracts.
- It's preferred to wait until `backfillEvents` completes so that we have all of the events necessary to generate hashed storage keys in order to be able to transform storage diffs. If `extractDiffs` is started before `backfillEvents` completes, the new diffs would be marked as `pending` if we didn't receive it's corresponding event, and would end up getting transformed later.
- `extractDiffs.toml` will need to be updated with the new collateral's contract details, and when the process is restarted any contract addresses in that file will be used in creating the geth diff subscription.
    
## 4. Backfill Collateral Storage
- Once storage diffs for the new contracts are being tracked going forward, the `backfillStorage` process can be started.
- Similar to how `backfillEvents` is updated, several components need to be added to the `backfillStorage.toml` config file as well. These components can all be copied and pasted from the `mcdTransformers.toml` file which was updated in the previous step with the collateral generator script.
    1. Contract Details
        - The Flip contract details are required. The Price Feed contract details are required as well if they're relevant for the given collateral. See the note above to determine if an OSM or Median contracts are relevant for the new collateral.
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
