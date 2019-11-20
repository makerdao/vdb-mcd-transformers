# [MakerDAO](https://makerdao.com) [VulcanizeDB](https://github.com/makerdao/vulcanizedb) Transformers

[![Build Status](https://travis-ci.com/makerdao/vdb-mcd-transformers.svg?token=yVLfQ4hNQCqvt2qjLkus&branch=staging)](https://travis-ci.com/makerdao/vdb-mcd-transformers)

## About

Vulcanize DB is a set of tools that make it easier for developers to write application-specific indexes and caches for dapps built on Ethereum.

## Dependencies
 - Go 1.11+
 - Postgres 11
 - Ethereum Node
   - [Go Ethereum](https://ethereum.github.io/go-ethereum/downloads/) (1.8.23)
   - [Parity 1.8.11+](https://github.com/paritytech/parity/releases)

## Project Setup

Using Vulcanize for the first time requires several steps be done in order to allow use of the software. The following instructions will offer a guide through the steps of the process:

1. Fetching the project
2. Installing dependencies
3. Configuring shell environment
4. Database setup
5. Configuring synced Ethereum node integration
6. Data syncing

## Installation

In order to fetch the project codebase for local use or modification, install it to your `GOPATH` via:

`go get github.com/makerdao/vulcanizedb`
`go get gopkg.in/DataDog/dd-trace-go.v1/ddtrace`

Once fetched, dependencies can be installed via `go get` or (the preferred method) at specific versions via `golang/dep`, the prototype golang pakcage manager. Installation instructions are [here](https://golang.github.io/dep/docs/installation.html).

In order to install packages with `dep`, ensure you are in the project directory now within your `GOPATH` (default location is `~/go/src/github.com/makerdao/vulcanizedb/`) and run:

`dep ensure`

After `dep` finishes, dependencies should be installed within your `GOPATH` at the versions specified in `Gopkg.toml`.

Lastly, ensure that `GOPATH` is defined in your shell. If necessary, `GOPATH` can be set in `~/.bashrc` or `~/.bash_profile`, depending upon your system. It can be additionally helpful to add `$GOPATH/bin` to your shell's `$PATH`.

## Setting up the Database
1. Install Postgres
1. Create a superuser for yourself and make sure `psql --list` works without prompting for a password.
1. `createdb vulcanize_public`
1. `cd $GOPATH/src/github.com/makerdao/vulcanizedb`
1.  Run the migrations: `make migrate HOST_NAME=localhost NAME=vulcanize_public PORT=5432`
    - To rollback a single step: `make rollback NAME=vulcanize_public`
    - To rollback to a certain migration: `make rollback_to MIGRATION=n NAME=vulcanize_public`
    - To see status of migrations: `make migration_status NAME=vulcanize_public`

    * See below for configuring additional environments

### Setting up GraphQL frontend
The database migrations sets up prerequisites for a GraphQL API, which can be generated with Postgraphile.
To present the API as specified [here](https://github.com/makerdao/vulcan.spec/blob/master/mcd.graphql):

1. Install Postgraphile: `npm install -g postgraphile`
1. Invoke `postgraphile --connection postgres://[superuser]:[password]@localhost:5432/vulcanize_public --schema api --disable-default-mutations --no-ignore-rbac --watch --enhance-graphiql`
1. The GraphiQL frontend can be found at `localhost:5000/graphiql`, and the GraphQL endpoint at `/graphql`

## Configuration
- To use a local Ethereum node, copy `environments/public.toml.example` to
  `environments/public.toml` and update the `ipcPath` and `levelDbPath`.
  - `ipcPath` should match the local node's IPC filepath:
      - For Geth:
        - The IPC file is called `geth.ipc`.
        - The geth IPC file path is printed to the console when you start geth.
        - The default location is:
          - Mac: `<full home path>/Library/Ethereum`
          - Linux: `<full home path>/ethereum/geth.ipc`

      - For Parity:
        - The IPC file is called `jsonrpc.ipc`.
        - The default location is:
          - Mac: `<full home path>/Library/Application\ Support/io.parity.ethereum/`
          - Linux: `<full home path>/local/share/io.parity.ethereum/`
          
  - `levelDbPath` should match Geth's chaindata directory path.
      - The geth LevelDB chaindata path is printed to the console when you start geth.
      - The default location is:
          - Mac: `<full home path>/Library/Ethereum/geth/chaindata`
          - Linux: `<full home path>/ethereum/geth/chaindata`
      - `levelDbPath` is irrelevant (and `coldImport` is currently unavailable) if only running parity.

- See `environments/infura.toml` to configure commands to run against infura, if a local node is unavailable.
- Copy `environments/local.toml.example` to `environments/local.toml` to configure commands to run against a local node such as [Ganache](https://truffleframework.com/ganache) or [ganache-cli](https://github.com/trufflesuite/ganache-clihttps://github.com/trufflesuite/ganache-cli).

## Start syncing with postgres
Syncs VulcanizeDB with the configured Ethereum node, populating blocks, transactions, receipts, and logs.
This command is useful when you want to maintain a broad cache of what's happening on the blockchain.
1. Start Ethereum node (**if fast syncing your Ethereum node, wait for initial sync to finish**)
1. In a separate terminal start VulcanizeDB:
    - `./vulcanizedb sync --config <config.toml> --starting-block-number <block-number>`

## Alternatively, sync from Geth's underlying LevelDB
Sync VulcanizeDB from the LevelDB underlying a Geth node.
1. Assure node is not running, and that it has synced to the desired block height.
1. Start vulcanize_db
   - `./vulcanizedb coldImport --config <config.toml>`
1. Optional flags:
    - `--starting-block-number <block number>`/`-s <block number>`: block number to start syncing from
    - `--ending-block-number <block number>`/`-e <block number>`: block number to sync to
    - `--all`/`-a`: sync all missing blocks

## Alternatively, only sync the headers
Syncs VulcanizeDB with the configured Ethereum node, populating only block headers.
This command is useful when you want a minimal baseline from which to track targeted data on the blockchain (e.g. individual smart contract storage values).
1. Start Ethereum node
1. In a separate terminal start VulcanizeDB:
    - `./vulcanizedb headerSync --config <config.toml> --starting-block-number <block-number>`

## Continuously sync Maker event logs from header sync
Continuously syncs Maker event logs from the configured Ethereum node based on the populated block headers.
This includes logs related to auctions, multi-collateral dai, and price feeds.
This command requires that the `headerSync` process is also being run so as to be able to sync in real time.

1. Start Ethereum node (or plan to configure the commands to point to a remote IPC path).
1. In a separate terminal run the headerSync command (see above).
1. In another terminal window run the continuousLogSync command:
  - `./vulcanizedb continuousLogSync --config <config.toml>`
  - An option `--transformers` flag may be passed to the command to specific which transformers to execute, this will default to all transformers if the flag is not passed.
    - `./vulcanizedb continuousLogSync --config environments/private.toml --transformers="priceFeed"`
    - see the `buildTransformerInitializerMap` method in `cmd/continuousLogSync.go` for available transformers

## Backfill Maker event logs from header sync
Backfills Maker event logs from the configured Ethereum node based on the populated block headers.
This includes logs related to auctions, multi-collateral dai, and price feeds.
This command requires that a header sync (see command above) has previously been run.

_Since auction/mcd contracts have not yet been deployed, this command will need to be run a local blockchain at the moment. As such, a new environment file will need to be added. See `environments/local.toml.example`._

1. Start Ethereum node
1. In a separate terminal run the backfill command:
  - `./vulcanizedb backfillMakerLogs --config <config.toml>`
  
## Running the Tests
- `make test` will run the unit tests and skip the integration tests
- `make integrationtest` will run the just the integration tests

## Deploying
1. you will need to make sure you have ssh agent running and your ssh key added to it. instructions [here](https://developer.github.com/v3/guides/using-ssh-agent-forwarding/#your-key-must-be-available-to-ssh-agent)
1. `go get -u github.com/pressly/sup/cmd/sup`
1. `sup staging deploy`
