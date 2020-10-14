# Dockerfile

## execute
Builds a buster image for running the vulcanizedb `execute` command against transformers in this repo.

### Build
Build from the project root directory with: `docker build -f dockerfiles/execute/Dockerfile . -t vdb-mcd-execute:latest`.
The following options are available at build time:
- `VDB_VERSION` - target a specific vulcanizedb branch/release to generate the binary (ex: `docker build -f dockerfiles/execute/Dockerfile --build-arg VDB_VERSION=v0.0.14-rc.1 . -t vdb-mcd-execute:latest`).
- `CONFIG_FILE` - path to desired config file for this container (ex: `docker build -f dockerfiles/execute/Dockerfile --build-arg CONFIG_FILE=path . -t vdb-mcd-execute:latest`).

There is a `make` task as well named `dockerbuild` which can quickly build the image with `make dockerbuild IMAGE=execute`, but that is not suitable for a production build.

### Run
Running the container requires an existing DB with which the container can interact.

The following arguments are required at runtime:

- `DATABASE_NAME`
- `DATABASE_HOSTNAME`
- `DATABASE_PORT`
- `DATABASE_USER`
- `DATABASE_PASSWORD`
- `CLIENT_IPCPATH`

#### Example

With arguments correctly populated, the following command will run the container on OS X:

```
docker run -e DATABASE_NAME=vulcanize_public -e DATABASE_HOSTNAME=host.docker.internal -e DATABASE_PORT=5432 -e DATABASE_USER=user -e DATABASE_PASSWORD=pw -e CLIENT_IPCPATH=https://mainnet.infura.io/v3/token -it vdb-mcd-execute:latest
```

#### Explanation

With the above command, we assume the host is exposing the database `vulcanize_public` on `localhost:5432` and user `user` with password `pw` has write access to that db.
We expect that we can successfully make calls against the [Ethereum JSON RPC API](https://github.com/ethereum/wiki/wiki/JSON-RPC) at `https://kovan.infura.io/v3/token`.

Note that on OS X, we use `host.docker.internal` to access `localhost`.
For execution on linux, replace instances of `host.docker.internal` with `localhost` and run with `--network="host"`.

#### Running with a geth docker container:
- On the `latest` geth release, run `docker build ./ -t geth-statediffing`
- Run the previously built docker container with the following flags:
    - Public the following ports from the geth container to the host:
          - `-p 8545:8545` //used in the rpc calls
          - `-p 8546:8546` //used in web socket subscription
          - `-p 30303:30303`
    - To have the geth process (in a container) use chaindata on your host machine, create a shared volume: `-v <path to shared volume on host>:/root/.ethereum`.
        - The host path of this volume could be one of the following if you're hoping to use chaindata from an exisiting node:
           - `~/Library/Ethereum/` on macOS
           - `~/.ethereum/` on Linux
     - Include the following geth flags:
        - `--rpc`
        - `--rpcaddr "0.0.0.0"`
        - `--ws`
        - `--wsaddr "0.0.0.0"`
        - `--statediff`
        - `--syncmode full`
        - `--statediff.watchedaddresses <contract address>`
        - `--statediff.watchedaddresses <contract address>`
    - e.g.
    ```shell script
        docker run -v </path/to/user/geth/directory>:/root/.ethereum -p 8545:8545 -p 8546:8546 -p 30303:30303 geth-statediffing
          --rpc --rpcaddr "0.0.0.0" --ws --wsaddr "0.0.0.0" --statediff --syncmode full --statediff.watchedaddresses <contract address>
          --statediff.watchedaddresses <contract address>
    ```
  
## extractDiffs
Dockerfile for populating storage diffs in db to be transformed by `execute`.

## Build
From project root directory:
```
docker build -f dockerfiles/extract_diffs/Dockerfile . -t extract_diffs:latest
```

### Run
Against statediffing Geth pubsub:
```
docker run -e DATABASE_USER=user -e DATABASE_PASSWORD=password -e DATABASE_HOSTNAME=host -e DATABASE_PORT=port -e DATABASE_NAME=name -e CLIENT_IPCPATH=path -e STORAGEDIFFS_SOURCE=geth -it extract_diffs:latest
```

## backfillEvents
Dockerfile for getting events from previously checked headers when adding a new event transformer.


### Build
`docker build -f dockerfiles/backfill_events/Dockerfile . -t backfill_events:latest`

## Run
```
docker run -e DATABASE_NAME=vulcanize_public -e DATABASE_HOSTNAME=host.docker.internal -e DATABASE_PORT=5432 -e DATABASE_USER=user -e DATABASE_PASSWORD=pw -e CLIENT_IPCPATH=https://mainnet.infura.io/v3/token -e ENDING_BLOCK_NUMBER=1 -it backfill_events:latest
```

## backfillStorage
Dockerfile for getting storage for all configured transformers in a range of blocks, and persisting them into the
`public.storage_diff` table. This is useful in case storage diffs were missed. Please note that the storage value for
every storage key for every transformer will be fetched, regardless of if the value actually changed at the given block.

### Build
From project root directory:
```
docker build -f dockerfiles/backfill_storage/Dockerfile . -t backfill_storage:latest
```

### Run
```
docker run -e CLIENT_IPCPATH=ipc_path -e DATABASE_USER=user -e DATABASE_PASSWORD=password -e DATABASE_HOSTNAME=host -e DATABASE_PORT=port -e DATABASE_NAME=name -e STARTING_BLOCK_NUMBER=block-number -e ENDING_BLOCK_NUMBER=block-number -it backfill_storage:latest
```
Notes:
- `STARTING_BLOCK_NUMBER` and `ENDING_BLOCK_NUMBER` variables are required
- `CLIENT_IPCPATH` should point at an archive node to enable accessing historical state
