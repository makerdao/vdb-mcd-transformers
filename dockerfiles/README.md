`Dockerfile` will build an alpine image containing:
- An $GOPATH with vulcanizedb, mcd_transformers, and goose
- An app directory with the vulcanizedb binary, startup_script.sh, and a (configurable) config.toml
Build with (e.g. from the project directory) `docker build ./ -t vulcanize_mcd_transformers:0.0.1 --build-arg USER`


## To use the container:
1. Setup a postgres database matching your config (e.g. `vulcanize_public`)
1. Determine values for the following _required_ environment variables:
    - `CLIENT_IPCPATH`
    - `DATABASE_NAME`
    - `DATABASE_HOSTNAME`
    - `DATABASE_PORT`
    - `DATABASE_USER`
    - `DATABASE_PASSWORD`
    - `STARTING_BLOCK_NUMBER` which is the block that the `headerSync` process starts at
    - `FILESYSTEM_STORAGEDIFFSPATH` or `STORAGEDIFFS_SOURCE` - this depends on where you're getting storage diffs from, see below.
        - Use `FILESYSTEM_STORAGEDIFFSPATH` only when getting storage diffs from a CSV file - this will be the path to the CSV file.
        - Set `STORAGEDIFFS_SOURCE` to `geth` when getting storage diffs from a subscription to a geth client. The default `STORAGEDIFFS_SOURCE` is `csv`.
1. Run with required environment variables: `docker run -e CLIENT_IPCPATH="https://kovan.infura.io/v3/token" -e DATABASE_NAME="vulcanize_public" -e DATABASE_HOSTNAME="host.docker.internal" -e DATABASE_PORT="5432" -e DATABASE_USER="vulcanize" -e DATABASE_PASSWORD="vulcanize" -e STARTING_BLOCK_NUMBER=0 -e FILESYSTEM_STORAGEDIFFSPATH="/path/to/diffs" vulcanize_mcd_transformers:0.0.1`.
    - This triggers `headerSync` + `composeAndExecute`.
    - NOTE: contract addresses are currently configured in `environments/docker.toml` to point at the given release's Kovan deployment.
       You can optionally replace any address with an environment variable, e.g. `-e CONTRACT_CONTRACT_MCD_FLIP_ETH_A_ADDRESS=0x1234"`.
    - To use a config file other than the default (`environments/docker.toml`), pass the following flag when building the image `--build-arg config_file=path/to/your/config/file`

NOTE: this file is written for execution on OS X, making use of `host.docker.internal` to access Postgres from the host.
For execution on linux, replace instances of `host.docker.internal` with `localhost` and run with `--network="host"`.

#### Running with a local geth client emitting statediffs:
- Use the v1.10-alpha.0 release of the vulcanize go-ethereum patch.
    - this patch allows for statediffs to be emitted
    - start geth with the following flags `--statediff`, `--ws`
- Make sure to set CLIENT_IPCPATH to use either the websocket or ipc interface with geth (notifications required for subscriptions are not supported with http).
    e.g. `docker run -e CLIENT_IPCPATH="ws://host.docker.internal:8546" -e DATABASE_NAME="vulcanize_public" -e DATABASE_HOSTNAME="host.docker.internal" -e DATABASE_PORT="5432" -e DATABASE_USER="vulcanize" -e DATABASE_PASSWORD="vulcanize" -e STORAGEDIFFS_SOURCE="geth" -e STARTING_BLOCK_NUMBER=0 vulcanize_mcd_transformers:0.0.1`.

#### Running with a geth docker container:
- On the `v.10-alpha.0` geth release, run `docker build ./ -t geth-statediffing`
- Run the previously built docker container with the following flags:
    - Public the following ports from the geth container to the host:
          - `-p 8545:8545` //used in the rpc calls
          - `-p 8546:8546` //used in web socket subscription
          - `-p 30303:30303`
    - To have the geth process (in a container) use chaindata on your host machine, create a shared volume: `-v <path to shared volumen on host>:/root/.ethereum`.
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
        docker run -v /Users/elizabethengelman/Library/Ethereum:/root/.ethereum -p 8545:8545 -p 8546:8546 -p 30303:30303 statediffing-geth
          --rpc --rpcaddr "0.0.0.0" --ws --wsaddr "0.0.0.0" --statediff --syncmode full --statediff.watchedaddresses <contract address>
          --statediff.watchedaddresses <contract address>
    ```
