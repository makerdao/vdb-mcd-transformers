FROM makerdao/vdb-builder:latest as builder

ARG VDB_VERSION

WORKDIR /go/src/github.com/makerdao/vdb-mcd-transformers
COPY . .

WORKDIR /go/src/github.com/makerdao
RUN git clone https://github.com/makerdao/vulcanizedb.git
WORKDIR /go/src/github.com/makerdao/vulcanizedb
RUN git checkout $VDB_VERSION
RUN go build

# build mcd with local vdb
WORKDIR /go/src/github.com/makerdao/vdb-mcd-transformers
RUN go mod edit -replace=github.com/makerdao/vulcanizedb=/go/src/github.com/makerdao/vulcanizedb/
RUN make plugin PACKAGE=github.com/makerdao/vdb-mcd-transformers \
    TARGET_LOCATION=$GOPATH/src/github.com/makerdao/vdb-mcd-transformers/plugins/execute/transformerExporter.go \
    OUTPUT_LOCATION=$GOPATH/src/github.com/makerdao/vdb-mcd-transformers/plugins/transformerExporter.so

# app container
FROM makerdao/vdb-runner:latest as runner
WORKDIR /go/src/github.com/makerdao/vulcanizedb

ARG CONFIG_FILE=environments/mcdTransformers.toml

# Direct logs to stdout for docker log driver
RUN ln -sf /dev/stdout /go/src/github.com/makerdao/vulcanizedb/vulcanizedb.log

# keep binaries immutable
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb .
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/$CONFIG_FILE config.toml
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/dockerfiles/execute/startup_script.sh .
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/db/migrations db/vdb-mcd-transformers/migrations/
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/plugins/transformerExporter.so plugins/transformerExporter.so
COPY --from=builder /go/bin/goose goose
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/dockerfiles/wait-for-it.sh .
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/dockerfiles/execute/run_migrations.sh .

# NOTE: You may need to update this in your cloud provider config in addition to changes here
HEALTHCHECK CMD grep -q "event watcher starting" /tmp/execute_health_check \
    && grep -q "storage watcher for new diffs starting" /tmp/execute_health_check \
    && grep -q "storage watcher for unrecognized diffs starting" /tmp/execute_health_check

# need to execute with a shell to access env variables
CMD ["./startup_script.sh"]
