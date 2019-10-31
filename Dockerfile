FROM golang:alpine as builder

RUN apk --update --no-cache add make git g++ linux-headers
# DEBUG
RUN apk add busybox-extras

ENV GO111MODULE on

# Build statically linked vDB binary (wonky path because of Dep)
WORKDIR /go/src/github.com/vulcanize/mcd_transformers
ADD . .

WORKDIR /go/src/github.com/vulcanize
RUN git clone https://github.com/vulcanize/vulcanizedb.git
WORKDIR /go/src/github.com/vulcanize/vulcanizedb
RUN git checkout v0.0.9
RUN go build

# build mcd with local vdb
WORKDIR /go/src/github.com/vulcanize/mcd_transformers

RUN go mod edit -replace=github.com/vulcanize/vulcanizedb=/go/src/github.com/vulcanize/vulcanizedb/
RUN make plugin PACKAGE=github.com/vulcanize/mcd_transformers

# Build migration tool
WORKDIR /go
RUN GO111MODULE=auto go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN GO111MODULE=auto go build -a -ldflags '-s' -tags='no_mysql no_sqlite' -o goose


# app container
FROM golang:alpine
# workdir needs to match gopath for building file to correct path
WORKDIR /go/src/github.com/vulcanize/vulcanizedb

# add certificates for node requests via https
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

# add go so we can build the plugin
RUN apk add --update --no-cache git g++ linux-headers

ARG USER
ARG config_file=environments/docker.toml
ARG vdb_connect

# setup environment
ENV GOPATH $HOME/go
ENV GO111MODULE on
ENV VDB_PG_CONNECT="$vdb_connect"

# Direct logs to stdout for docker log driver
RUN ln -sf /dev/stdout /go/src/github.com/vulcanize/vulcanizedb/vulcanizedb.log

RUN adduser -Su 5000 $USER
# container needs to be writable for plugin execution
RUN chown 5000:5000 /go/src/github.com/vulcanize/vulcanizedb

USER $USER

# chown first so dir is writable
# note: using $USER is merged, but not in the stable release yet
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb .
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/mcd_transformers /go/src/github.com/vulcanize/mcd_transformers

# keep binaries immutable
COPY --from=builder /go/src/github.com/vulcanize/mcd_transformers/$config_file config.toml
COPY --from=builder /go/src/github.com/vulcanize/mcd_transformers/dockerfiles/startup_script.sh .
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose goose
COPY --from=builder /go/src/github.com/vulcanize/mcd_transformers/db/migrations/* db/migrations/
COPY --from=builder /go/src/github.com/vulcanize/mcd_transformers/plugins/transformerExporter.so plugins/transformerExporter.so

# DEBUG
COPY --from=builder /usr/bin/telnet /bin/telnet

# need to execute with a shell to access env variables
CMD ["./startup_script.sh"]
