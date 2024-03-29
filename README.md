# [MakerDAO](https://makerdao.com) [VulcanizeDB](https://github.com/makerdao/vulcanizedb) Transformers

[![Build Status](https://travis-ci.com/makerdao/vdb-mcd-transformers.svg?token=yVLfQ4hNQCqvt2qjLkus&branch=prod)](https://travis-ci.com/makerdao/vdb-mcd-transformers)


## Table of Contents
1. [Background](#background)
1. [Install](#install)
1. [Usage](#usage)
1. [Contributing](#contributing)
1. [License](#license)

## Background
This repository is a collection of transformers to be used along with VDB as a plugin to fetch, transform and persist log events and storage slots of specified MCD contracts.

## Install
### Docker Setup

A full VDB environment requires several concurrent processes to read and transform raw data. The [vdb-transformer-utilities](https://github.com/makerdao/vdb-transformer-utilities) repository provides a docker-compose setup to manage the multiple containerized processes for you. Please see the [README](https://github.com/makerdao/vdb-transformer-utilities/dockerfiles/README.md) for setup instructions. 

### Manual Setup

#### Dependencies
 - Ethereum Node
    - VDB requires a patched version of [Go Ethereum](https://github.com/makerdao/go-ethereum) (1.10.8+) in order to access contract storage.
 - [VulcanizeDB](https://github.com/makerdao/vulcanizedb)
 - Go 1.15+
 - Postgres 11.2
 
#### Updating VulcanizeDB

Updating to a new version requires updating the dependency in go.mod and updating the vulcanizedb schema copy for testing, to keep them in sync. The simplest thing to do is run the `make update_vulcanize BRANCH=<branch name>` task which will handle both.
    
### Getting the project
Download the transformer codebase to your local `GOPATH` via: `go get github.com/makerdao/vdb-mcd-transformers`

## Usage

### Plugin

This repository is a plugin for VulcanizeDB.  As such, it cannot be run as a standalone executable, but instead is intended to be included as part of a VulcanizeDB core command. There are two VulcanizeDB core commands that are required for events and storage slots to be transformed and persisted to the Postgres database:

1. `headerSync` fetches raw Ethereum data and syncs it into VulcanizeDB's Postgres database where then it can be used for transformations. More information about the `headerSync` command can be found in the [VulcanizeDB repository](https://github.com/makerdao/vulcanizedb/blob/prod/documentation/data-syncing.md#headersync).

1. `execute` uses the raw Ethereum data that has been synced into Postgres and applies transformations to configured MCD contract data via [event](./transformers/events) and [storage](./transformers/storage) transformers. The VulcanizeDB repository includes [more information](https://github.com/makerdao/vulcanizedb/blob/prod/documentation/custom-transformers.md) on transformers.

The core commands can also be run via images or built and run via the command line interface. In either method, a postgres database will first need to be created:
1. Install Postgres
   1. Create a user for yourself that is able run migrations and add extensions.
1. `createdb vulcanize_public`
1. Migrate the database using the `make migrate` task in this repository.

### Backfill
This project contains one executable, `backfillUrns`.
This command enables you to get storage data for Urns by performing lookups at blocks where we detected events indicating an Urn state change.
To run it, you first need to `go build`.
After that, you can run something like `./vdb-mcd-transformers backfillUrns --starting-block 10000000` - which would backfill Urn storage since block 10,000,000 through the head of the chain.
Note that, as with other commands, executing this requires either a config file or env vars to specify the database and ethereum node connection.

### Running With Docker

#### Database Initialization

**NOTE** The database must be migrated once before running the `headerSync` command, otherwise the database will not be able to properly create its schema. Assuming you are not using the `docker-compose` setup above you can migrate the database once using:

```
docker run -e DATABASE_USER=<user> -e DATABASE_PASSWORD=<pw> -e DATABASE_HOSTNAME=<host> -e DATABASE_PORT=<port> -e DATABASE_NAME=<name> -e CLIENT_IPCPATH=<path> makerdao/vdb-mcd-execute:latest ./run_migrations.sh
```

#### Running `headerSync`
`headerSync` Docker images are located in the [MakerDao Dockerhub organization](https://hub.docker.com/repository/docker/makerdao/vdb-headersync).

Start `headerSync`:
```
docker run -e DATABASE_USER=<user> -e DATABASE_PASSWORD=<pw> -e DATABASE_HOSTNAME=<host> -e DATABASE_PORT=<port> -e DATABASE_NAME=<name> -e STARTING_BLOCK_NUMBER=<starting block number> -e CLIENT_IPCPATH=<path> makerdao/vdb-headersync:latest
```
- `STARTING_BLOCK_NUMBER` indicates where to start syncing raw headers. If you don't care about headers before the contracts of interest were deployed, it is recommended to use the earliest deployment block of those contracts. Otherwise, the command will sync all headers starting at the genesis block.
- To allow the docker container to connect to a local database and a local Ethereum node: 
  - when running on Linux include `--network=host` 
  - when running on MacOSX use `host.docker.internal` as the `DATABASE_HOSTNAME` and as the host in the `CLIENT_IPCPATH`

#### Running `execute`
`execute` Docker images are located in the [MakerDao Dockerhub organization](https://hub.docker.com/repository/docker/makerdao/vdb-mcd-execute). See the [Docker README](./dockerfiles/README.md) for further information.

### With the CLI

1. Move to the project directory:
```cd $GOPATH/src/github.com/makerdao/vulcanizedb```

1. Be sure you have enabled Go Modules (`export GO111MODULE=on`), and build the executable with:
```make build```

#### Running `headerSync`
```./vulcanizedb headerSync --config <config.toml> --starting-block-number <block-number>```

For more information, see the [VulcanizeDB data-syncing documentation](https://github.com/makerdao/vulcanizedb/blob/prod/documentation/data-syncing.md).

#### Running `execute`
Instead of running `execute`, you will also need to `compose` the transformer initializer plugin prior to execution. (This command builds and links configured transformers). This step is not explicitly required when using Docker because it is included in the `vdb-execute` container.

There is a convenience command called `composeAndExecute` in `vulcanizedb` which encapsulates both composing the plugin, and then
executing it. 

```
./vulcanizedb composeAndExecute --config=$GOPATH/makerdao/vdb-mcd-transformers/environments/mcdTransformers.toml \
    DATABASE_NAME=vulcanize_public \
    DATABASE_PASSWORD=vulcanize \
    DATABASE_HOST=localhost \
    DATABASE_PORT=5432 \
    DATABASE_USER=vulcanize
```
   
Notes:
- Make sure that `vulcanizedb` and `vdb-mcd-transformers` versions are compatible. `vulcanizedb` will default to grabbing
the [most recent vdb-mcd-transformers release](https://github.com/makerdao/vdb-mcd-transformers/releases). You can check
the `vdb-mcd-transformers` [go.mod](./go.mod) file to see what `vulcanizedb` version is expected.
- Make sure that the transformers in the config file you're using match up with the ones included in the release.
- The dependencies between the two repositories need to be in sync, otherwise the plugins will not be able to be composed properly.
There is a [script](https://github.com/makerdao/vulcanizedb/blob/prod/scripts/gomoderator.py) in the VulcanizeDB
repository to take care of this. This mismatch of dependencies versions should not happen if two compatible releases are
used, but is possible in development.
- If you need to use a different dependency than what is currently defined in `go.mod` in either repository, it may be
helpful to look into [the replace directive](https://github.com/golang/go/wiki/Modules#when-should-i-use-the-replace-directive).
This instruction enables you to point at a fork or the local filesystem for dependency resolution.
- If you are running into issues, ensure that `GOPATH` is defined in your shell. If necessary, `GOPATH` can be set in 
`~/.bashrc` or `~/.bash_profile`, depending upon your system. It can be additionally helpful to add `$GOPATH/bin` to your
shell's `$PATH`.

### Exposing the data
[Postgraphile](https://www.graphile.org/postgraphile/) is used to expose GraphQL endpoints for our database schemas. Please see the more detailed documentation [here](https://github.com/makerdao/vulcanizedb/blob/prod/documentation/postgraphile.md).

### Running tests
- Set the ipc path to a node by setting the CLIENT_IPCPATH environment variable.
- `make test` will run the unit tests and skip the integration tests
- `make integrationtest` will run just the integration tests
- `make test` and `make integrationtest` set up a clean `vulcanize_testing` db

## Contributing
Contributions are welcome!

VulcanizeDB follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/1/4/code-of-conduct).

For more information on contributing, please see [here](https://github.com/makerdao/vulcanizedb/blob/prod/documentation/contributing.md).

## License
[AGPL-3.0](LICENSE) © Vulcanize Inc
