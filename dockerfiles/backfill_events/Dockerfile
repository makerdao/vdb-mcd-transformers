FROM golang:1.15.5-buster as builder

RUN apt-get update && \
    apt-get install -y \
    make git g++ linux-headers-amd64

# Build migration tool
WORKDIR /go
RUN GO111MODULE=auto go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN GO111MODULE=auto go build -a -ldflags '-s' -tags='no_mysql no_sqlite' -o goose

ARG VDB_VERSION=beta
ENV GO111MODULE on

WORKDIR /go/src/github.com/makerdao/vdb-mcd-transformers
ADD . .

WORKDIR /go/src/github.com/makerdao
RUN git clone https://github.com/makerdao/vulcanizedb.git
WORKDIR /go/src/github.com/makerdao/vulcanizedb
RUN git checkout $VDB_VERSION
RUN go build

# build mcd with local vdb
WORKDIR /go/src/github.com/makerdao/vdb-mcd-transformers
RUN go mod edit -replace=github.com/makerdao/vulcanizedb=/go/src/github.com/makerdao/vulcanizedb/
RUN make plugin PACKAGE=github.com/makerdao/vdb-mcd-transformers \
    TARGET_LOCATION=$GOPATH/src/github.com/makerdao/vdb-mcd-transformers/plugins/backfill/transformerExporter.go \
    OUTPUT_LOCATION=$GOPATH/src/github.com/makerdao/vdb-mcd-transformers/plugins/transformerExporter.so

# app container
FROM golang:1.15.5-buster

RUN apt-get update \
        && apt-get install -y \
        ca-certificates \
        bash \
        busybox \
        postgresql-client \
        && update-ca-certificates 2>/dev/null || true

WORKDIR /go/src/github.com/makerdao/vulcanizedb

ARG CONFIG_FILE=environments/backfillEvents.toml

# Direct logs to stdout for docker log driver
RUN ln -sf /dev/stdout /go/src/github.com/makerdao/vulcanizedb/vulcanizedb.log

# keep binaries immutable
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb .
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/$CONFIG_FILE config.toml
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/dockerfiles/backfill_events/startup_script.sh .
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose goose
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/db/migrations/* db/migrations/
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/plugins/transformerExporter.so plugins/transformerExporter.so
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/dockerfiles/wait-for-it.sh .

# need to execute with a shell to access env variables
CMD ["./startup_script.sh"]
