FROM golang:alpine as builder

RUN apk --update --no-cache add g++ git linux-headers

# Build migration tool
ARG VDB_VERSION=staging
ENV GO111MODULE on

WORKDIR /go/src/github.com/makerdao/vdb-mcd-transformers
COPY . .

# Get and build VulcanizeDB
WORKDIR /go/src/github.com/makerdao
RUN git clone https://github.com/makerdao/vulcanizedb.git
WORKDIR /go/src/github.com/makerdao/vulcanizedb
RUN git checkout $VDB_VERSION
RUN go build
# Build migration tool
RUN go install -tags='no_mysql no_sqlite3 no_mssql no_redshift' github.com/pressly/goose/cmd/goose

# App container
FROM golang:alpine
WORKDIR /go/src/github.com/makerdao/vulcanizedb

# Add certificates for node requests via https
# bash for wait-for-it.sh
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        bash \
        postgresql-client \
        && update-ca-certificates 2>/dev/null || true

ARG CONFIG_FILE=environments/extractDiffs.toml

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
