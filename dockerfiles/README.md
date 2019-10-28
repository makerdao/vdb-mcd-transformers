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
1. Run with required environment variables: `docker run -e CLIENT_IPCPATH="https://kovan.infura.io/v3/token" -e DATABASE_NAME="vulcanize_public" -e DATABASE_HOSTNAME="host.docker.internal" -e DATABASE_PORT="5432" -e DATABASE_USER="vulcanize" -e DATABASE_PASSWORD="vulcanize" -e FILESYSTEM_STORAGEDIFFSPATH="/path/to/diffs" vulcanize_mcd_transformers:0.0.1`.
    - This triggers `headerSync` + `composeAndExecute`.
    - NOTE: contract addresses are currently configured in `environments/docker.toml` to point at the given release's Kovan deployment.
       You can optionally replace any address with an environment variable, e.g. `-e CONTRACT_CONTRACT_MCD_FLIP_ETH_A_ADDRESS=0x1234"`.
    - To use a config file other than the default (`environments/docker.toml`), pass the following flag when building the image `--build-arg config_file=path/to/your/config/file`

NOTE: this file is written for execution on OS X, making use of `host.docker.internal` to access Postgres from the host.
For execution on linux, replace instances of `host.docker.internal` with `localhost` and run with `--network="host"`.

Running with a local geth client emitting statediffs:
- Use the v1.10-alpha.0 release of the vulcanize go-ethereum patch.
    - this patch allows for statediffs to be emitted
    - start geth with the following flags `--statediff`, `--ws`
- Make sure to set CLIENT_IPCPATH to use either the websocket or ipc interface with geth (notifications required for subscriptions are not supported with http).
    e.g. `docker run -e CLIENT_IPCPATH="ws://host.docker.internal:8546" -e DATABASE_NAME="vulcanize_public" -e DATABASE_HOSTNAME="host.docker.internal" -e DATABASE_PORT="5432" -e DATABASE_USER="vulcanize" -e DATABASE_PASSWORD="vulcanize" -e STORAGEDIFFS_SOURCE="geth" vulcanize_mcd_transformers:0.0.1`.
