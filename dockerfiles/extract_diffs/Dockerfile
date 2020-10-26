FROM makerdao/vdb-builder:latest as builder

ARG VDB_VERSION=prod

WORKDIR /go/src/github.com/makerdao/vdb-mcd-transformers
COPY . .

WORKDIR /go/src/github.com/makerdao
RUN git clone https://github.com/makerdao/vulcanizedb.git
WORKDIR /go/src/github.com/makerdao/vulcanizedb
RUN git checkout $VDB_VERSION
RUN go build

# App container
FROM makerdao/vdb-runner:latest
WORKDIR /go/src/github.com/makerdao/vulcanizedb

# Direct logs to stdout for docker log driver
RUN ln -sf /dev/stdout /go/src/github.com/makerdao/vulcanizedb/vulcanizedb.log

ARG CONFIG_FILE=environments/extractDiffs.toml

# Add required files
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb .
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/$CONFIG_FILE config.toml
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/dockerfiles/extract_diffs/startup_script.sh .
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/db/migrations db/vdb-mcd-transformers/migrations/
COPY --from=builder /go/bin/goose goose
# needed for waiting until postgres is ready before starting from docker-compose
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb/dockerfiles/wait-for-it.sh .

HEALTHCHECK CMD test -f /tmp/connection

# need to execute with a shell to access env variables
CMD ["./startup_script.sh"]
