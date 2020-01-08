FROM golang:alpine as builder

RUN apk --update --no-cache add make git g++ linux-headers

# Build migration tool
WORKDIR /go
RUN GO111MODULE=auto go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN GO111MODULE=auto go build -a -ldflags '-s' -tags='no_mysql no_sqlite' -o goose

ARG VDB_VERSION=staging
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
RUN make plugin PACKAGE=github.com/makerdao/vdb-mcd-transformers


# app container
FROM golang:alpine
WORKDIR /go/src/github.com/makerdao/vulcanizedb

# add certificates for node requests via https
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

# add go so we can build the plugin
RUN apk add --update --no-cache git g++ linux-headers

ARG CONFIG_FILE=environments/docker.toml

# Direct logs to stdout for docker log driver
RUN ln -sf /dev/stdout /go/src/github.com/makerdao/vulcanizedb/vulcanizedb.log

# keep binaries immutable
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb .
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/$CONFIG_FILE config.toml
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/dockerfiles/startup_script.sh .
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose goose
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/db/migrations/* db/migrations/
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/plugins/transformerExporter.so plugins/transformerExporter.so
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb/dockerfiles/wait-for-it.sh .

# need to execute with a shell to access env variables
CMD ["./startup_script.sh"]
