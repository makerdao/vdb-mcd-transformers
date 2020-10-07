FROM golang:1.15.2-buster as builder

RUN apt-get update && \
apt-get install -y make git g++ linux-headers-amd64

ARG VDB_VERSION=prod

ENV GO111MODULE on

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
FROM golang:1.15.2-buster
WORKDIR /go/src/github.com/makerdao/vulcanizedb

# add certificates for node requests via https
# bash for wait-for-it
RUN apt-get update \
        && apt-get upgrade -y \
        && apt-get install -y \
        ca-certificates \
        bash \
        postgresql-client \
        && update-ca-certificates 2>/dev/null || true

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
