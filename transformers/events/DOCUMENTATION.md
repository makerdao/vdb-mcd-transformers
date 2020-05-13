# Transformers

## Architecture

Transformers fetch logs from Ethereum, convert/decode them into usable data, and then persist them in postgres. A transformer converts the raw chain data into a human friendly representation suitable for consumption in the API.

## Event Types

For Maker there are three main types of log events that we're tracking:
  
1. Custom events that are defined in the contract solidity code.
1. `LogNote` events which utilize the [DSNote library](https://github.com/dapphub/ds-note).
1. `Note` events in the `Vat`

The transformer process for each of these different log types is the same, except for the converting process, as denoted below.

## Creating a Transformer

### Ensure there is a Contract
If the contract isn't already present in the environment you'll need to add it before you can begin transforming its events. To do so:

1. Get the contract address if you don't already have it.
1. Search for the contract on [etherscan](https://etherscan.io/). Bookmark this
   page, it is your new best friend.
1. In each environment (at this time testing, docker, and mcdTransformers) add a
   new contract to the contract section. For example:
   
    ``` toml
    [contract]
    [contrac t.MCD_FLIP_ETH_A]
        address  = "0xd8a04f5412223f513dc55f839574430f5ec15531"
        abi      = '[{"inputs":[{"internalType":"address","name":"vat_","type":"address"},{"internalType":"bytes32","name":"ilk_","type":"bytes32"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"id","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"lot","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"bid","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"tab","type":"uint256"},{"indexed":true,"internalType":"address","name":"usr","type":"address"},{"indexed":true,"internalType":"address","name":"gal","type":"address"}],"name":"Kick","type":"event"},{"anonymous":true,"inputs":[{"indexed":true,"internalType":"bytes4","name":"sig","type":"bytes4"},{"indexed":true,"internalType":"address","name":"usr","type":"address"},{"indexed":true,"internalType":"bytes32","name":"arg1","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"arg2","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"LogNote","type":"event"},{"constant":true,"inputs":[],"name":"beg","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"bids","outputs":[{"internalType":"uint256","name":"bid","type":"uint256"},{"internalType":"uint256","name":"lot","type":"uint256"},{"internalType":"address","name":"guy","type":"address"},{"internalType":"uint48","name":"tic","type":"uint48"},{"internalType":"uint48","name":"end","type":"uint48"},{"internalType":"address","name":"usr","type":"address"},{"internalType":"address","name":"gal","type":"address"},{"internalType":"uint256","name":"tab","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"name":"deal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"id","type":"uint256"},{"internalType":"uint256","name":"lot","type":"uint256"},{"internalType":"uint256","name":"bid","type":"uint256"}],"name":"dent","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"bytes32","name":"what","type":"bytes32"},{"internalType":"uint256","name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"ilk","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"},{"internalType":"address","name":"gal","type":"address"},{"internalType":"uint256","name":"tab","type":"uint256"},{"internalType":"uint256","name":"lot","type":"uint256"},{"internalType":"uint256","name":"bid","type":"uint256"}],"name":"kick","outputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"kicks","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"usr","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"tau","outputs":[{"internalType":"uint48","name":"","type":"uint48"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"id","type":"uint256"},{"internalType":"uint256","name":"lot","type":"uint256"},{"internalType":"uint256","name":"bid","type":"uint256"}],"name":"tend","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"name":"tick","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"ttl","outputs":[{"internalType":"uint48","name":"","type":"uint48"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"vat","outputs":[{"internalType":"contract VatLike","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"wards","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"name":"yank","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]'
        deployed = 8928180
    ```
 
     * The address field is the contract address you searched for, as a string.
    * The abi can be found on the `contracts` tab in Etherscan (scroll down). You'll need the entire thing as a string, so it's best to use single quotes like above.
    * The deployed field is the block number of the transaction where the contract was deployed. That can be found on the Contract's page, by looking for the Contract Creator and clicking it's transaction hash. That will take you to the transaction where the it was deployed on.

### Create the Transformer

To create a custom transformer you will need to create a new `Transformer`
struct with a `ToModels` method on it which converts a `core.EventLog` object
(the raw, untransformed data) to an `event.InsertionModel` (the domain object).

Note that for this step you may _not_ need to create the database
migration, because to write the transformer you don't need to save the transformed
object to the database. You may be saving addresses, which
go into the already existing address table and are referenced by a foreign key. 

This isn't always the case, but if you do need to create the
migration, you'll know, because the tests won't pass. The directions here assume
you do not, yet, as this is most likely.

1. Search for the contract on etherscan using it's signature. 
1. Find the events for the contract in the contract's source.
1. For each event in the contract:
   1. Write a test for the event signature in [`signature_test.go`](../shared/constants/signature_test.go)
      1. To find out the event's signature you can use this [Keccak-256 calculator](https://emn178.github.io/online-tools/keccak_256.html). For example:
         * Search Etherscan for the contract `0x39755357759cE0d7f32dC8dC45414CCa409AE24e`.
         * Open the 'Contract' tab and in the contract source code and search for the `LogItemUpdate` event. It looks like this: `event LogItemUpdate(uint id);`.
         * You can copy that string and paste it into the Keccak-256 calculator (link above). Remove the parameter names, leaving only the types. Use `uint256` for the `uint` and `int256` for `int` as those are the solidity defaults, e.g. `LogItemUpdate(uint256)`. This will output an event signature.
         * To verify you got the signature right you can take that signature (in this case `a2c251311b1a7a475913900a2a73dc9789a21b04bc737e050bbc506dd4eb3488`) and search for it on the events tab of the contract. Make sure you prefix it with `0x`. The results should include the event you are looking for, provided one already exists for it.
      1. Make the test pass by updating the function list in `signature.go` as well as updating `method.go` as needed. Both of these files are in the constants package.
      1. Create test event in [`test_data`](../test_data)
         * Use Data from a real event in Etherscan wherever possible. Particularly when it comes to Topics and the Data entries. You'll thank me later.
         * Update constants as needed. Look at older examples of test data for inspiration.
   1. Create a new directory in `events` named after your new transformer. This will be the package your transformer is stored in.
   1. Create three files - `transformer.go`, `transformer_test.go` and `<package_name>_test.go` in the new directory. See the other transformers for examples. 
   1. `<package_name>_test.go` lets `gingko` run tests in this directory. The file should look something like this (this can be generated with `gingko`):
       
       ```go
       package <package_name>_test

       import (
           "testing"

           . "github.com/onsi/ginkgo"
           . "github.com/onsi/gomega"
       )

       func TestLogMinSell(t *testing.T) {
           RegisterFailHandler(Fail)
           RunSpecs(t, "<PackageName>  Suite")
       }
       ```
   1. Implement the `ToModels` function in `transformer.go` with the appropriate testing in `transformer_test.go` of course. This means converting the raw log into a go struct. For clarity you're implementing this [interface](https://github.com/makerdao/vulcanizedb/blob/staging/libraries/shared/factories/event/converter.go).
      1. For custom events:
         1. You can convert each EventLog entry into an entity using the function `contract.UnpackLog(&entity, "<EventLogName>", log.Log)`. The entity is usually defined (by you) in a file called `entity.go` in the same package as the transformer. You can use the [abigen](https://geth.ethereum.org/docs/install-and-build/installing-geth) tool that comes with `Geth` as well, if you're struggling with getting the fields right, but it's usually simpler to just do it yourself. See the transformer in `transformers/events/log_make` for an example of this method.
         1. After converting the log entry to an entity, look up all of it's foreign keys from it's data and assign them to the entity.
         1. From the entity create an InsertionModel, which you return.
      1. For `LogNote` events you'll need to look at the method signature of the method that is calling them, because LogNote events are a generic structure. For example:
         * The `tend` method is called on the [flip.sol contract](https://github.com/makerdao/dss/blob/master/src/flip.sol#L123), and its method signature looks like this: `tend(uint,uint,uint)`.
         * Only the first four bytes of the Keccak-256 hashed method signature will be located in `topic[0]` on the log, unlike custom events.
         * The message sender will be in `topic[1]`.
         * The first parameter passed to `tend` becomes `topic[2]`.
         * The second parameter passed to `tend` will be `topic[3]`.
         * Any additional parameters will be in the log's data field.
         * More detail is located in the [DSNote repo](https://github.com/dapphub/ds-note).
         * For an example implementation look at the `flip_sol` transformer.
      1. If, while implementing the transformer and its corresponding unit tests you find you need to migrate the database see below. 

### Store the data in the database

1. Use the `make new_migration` task to create a new migration.
   * Each event log has its own table in the database.
   * The specific log event tables are all created in the `maker` schema.
1. The new migration can be run by running `make test` or `make migration NAME=<database_name>`. Note that if you need to modify the migration multiple times you do not need to rollback the new migration, `make test` will drop it and recreate it.
1. To verify that the migration will work create an integration for the shiny
   new transformer in [`integration_tests`](../integration_tests).
   * Unlike the transformer test you wrote earlier integration tests use their own `TransformerConfig` to run the transformer in the same way it will be in production.
   * You don't look at the event transformer directly, but query from the database.
   * These tests hit the internet as well, and will need real block numbers.
   * For best results look at one of the other tests in the `integration_tests` directory, and use real data from etherscan/mainnet.
1. Integration tests can be run with `make integrationtest`.

### Add the transformer to the list of transformers to run

Finally you can add the transformer to the list of transformers by updating the
configuration, and creating an initializer for the transformer.
 
1. In the environments (docker.toml, testing.toml and mcdTransformers.toml) add
   the new package name the list of `transformerNames` in the exporter - alphabetically.
1. Underneath that list add a configuration (alphabetically again) to list of
   configurations. For example:
   
    ``` toml
    [exporter.log_delete]
        path = "transformers/events/log_delete/initializer"
        type = "eth_event"
        repository = "github.com/makerdao/vdb-mcd-transformers"
        migrations = "db/migrations"
        contracts = ["OASIS_MATCHING_MARKET_ONE", "OASIS_MATCHING_MARKET_TWO"]
        rank = "0"
    ```
1. Note the path to the initializer is a directory named initializer in your new event transformer package. Of course you haven't created that yet. Create that directory and inside it create a file named `initializer.go`. 
1. Initializers are boilerplate that tell the system how to create your transformer. They look like this:
   
    ```go
    package initializer

    import (
    "github.com/makerdao/vdb-mcd-transformers/transformers/events/log_delete"
    "github.com/makerdao/vdb-mcd-transformers/transformers/shared"
    "github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
    "github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
    )

    var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
    Config:      shared.GetEventTransformerConfig(constants.LogDeleteTable, constants.LogDeleteSignature()),
    Transformer: log_delete.Transformer{},
    }.NewTransformer
    ```

    Simply replace the constants and package names with your transformer.

1. Finally add your package to the list of transformerExporters in `plugins/transformerExporter.go`. Again alphabetically. This can also be generated using the `./vulcanizedb compose --config=/path/to/config.toml` command.

### Fetching Logs

In the event there are not logs for an event you're looking to transform in etherscan (be it in Mainnet, Kovan or other) you can generate an example raw event by deploying the contract to a local chain and emitting the event manually.

1. Fetch the logs from the chain based on the example event's topic zero:
   * The topic zero is based on the keccak-256 hash of the log event's method signature. These are located in [`pkg/transformers/shared/constants/signature.go`](../shared/constants/signature.go). 
   * Fetching is done in batch from the [`watcher`](https://github.com/vulcanize/maker-vulcanizedb/blob/staging/libraries/shared/watcher/event_watcher.go). 
   * The logs are then chunked up by the [`chunker`](https://github.com/vulcanize/maker-vulcanizedb/blob/staging/libraries/shared/chunker/log_chunker.go) before being delegated to each transformer.

## Useful Documents

[Ethereum Event ABI Specification](https://solidity.readthedocs.io/en/develop/abi-spec.html#events)
