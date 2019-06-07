S
`Dockerfile` will build an alpine image containing:
- An $GOPATH with vulcanizedb, mcd_transformers, and goose
- An app directory with the vulcanizedb binary, startup_script.sh, and a (configurable) config.toml
Build with (e.g. from the project directory) `docker build ./ -t mcd:0.0.1 --build-arg USER`


## To use the container:
1. Setup a postgres database matching your config (e.g. `vulcanize_public`)
1. Set the config's (default: `environments/example.toml`) ipc path to a node endpoint
1. Run with e.g. `docker run mcd:0.0.1` [this triggers `headerSync` + `composeAndExecute` with the specified config]

NOTE: this file is written for execution on OS X, making use of `host.docker.internal` to access Postgres from the host.
For execution on linux, replace instances of `host.docker.internal` with `localhost` and run with `--network="host"`.

