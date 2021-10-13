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
    TARGET_LOCATION=$GOPATH/src/github.com/makerdao/vdb-mcd-transformers/plugins/backfill/storage/transformerExporter.go \
    OUTPUT_LOCATION=$GOPATH/src/github.com/makerdao/vdb-mcd-transformers/plugins/transformerExporter.so

# app container
FROM makerdao/vdb-runner:latest
WORKDIR /go/src/github.com/makerdao/vulcanizedb

ARG CONFIG_FILE=environments/backfillStorage.toml

# Direct logs to stdout for docker log driver
RUN ln -sf /dev/stdout /go/src/github.com/makerdao/vulcanizedb/vulcanizedb.log

# keep binaries immutable
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb .
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/$CONFIG_FILE config.toml
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/dockerfiles/backfill_storage/startup_script.sh .
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/plugins/transformerExporter.so plugins/transformerExporter.so

# need to execute with a shell to access env variables
CMD ["./startup_script.sh"]
