FROM golang:alpine as builder

RUN apk --update --no-cache add make git g++
# DEBUG
RUN apk add busybox-extras

# Build statically linked vDB binary (wonky path because of Dep)
WORKDIR /go/src/github.com/vulcanize/mcd_transformers
ADD . .

# Build migration tool
RUN go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN go build -a -ldflags '-s' -tags='no_mysql no_sqlite' -o goose

RUN go get -u -d github.com/vulcanize/vulcanizedb
WORKDIR /go/src/github.com/vulcanize/vulcanizedb
RUN go build -a -ldflags '-s' .

# app container
FROM golang:alpine
WORKDIR /app

# add certificates for node requests via https
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

# add go so we can build the plugin
RUN apk add --update --no-cache g++

ARG USER
ARG config_file=environments/example.toml
ARG vdb_pg_host="host.docker.internal"
ARG vdb_pg_port="5432"
ARG vdb_dbname="vulcanize_public"
ARG vdb_pg_connect="postgres://$USER@$vdb_pg_host:$vdb_pg_port/$vdb_dbname?sslmode=disable"

# setup environment
ENV VDB_PG_CONNECT="$vdb_pg_connect"
ENV GOPATH $HOME/go

RUN adduser -Su 5000 $USER
USER $USER

# chown first so dir is writable
# note: using $USER is merged, but not in the stable release yet
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb /go/src/github.com/vulcanize/vulcanizedb
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/mcd_transformers /go/src/github.com/vulcanize/mcd_transformers

# keep binaries immutable
COPY --from=builder /go/src/github.com/vulcanize/mcd_transformers/$config_file config.toml
COPY --from=builder /go/src/github.com/vulcanize/mcd_transformers/dockerfiles/startup_script.sh .
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose goose
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/db/migrations migrations/vulcanizedb
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/vulcanizedb vulcanizedb

# DEBUG
COPY --from=builder /usr/bin/telnet /bin/telnet

# need to execute with a shell to access env variables
CMD ["./startup_script.sh"]