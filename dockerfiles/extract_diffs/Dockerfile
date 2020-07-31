FROM golang:alpine as builder

RUN apk --update --no-cache add g++ git linux-headers

# Build migration tool
WORKDIR /go
RUN GO111MODULE=auto go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN GO111MODULE=auto go build -a -ldflags '-s' -tags='no_mysql no_sqlite' -o goose

ARG VDB_VERSION=staging
ENV GO111MODULE on

# Get and build VulcanizeDB
WORKDIR /go/src/github.com/makerdao
RUN git clone https://github.com/makerdao/vulcanizedb.git
WORKDIR /go/src/github.com/makerdao/vulcanizedb
RUN git checkout $VDB_VERSION
RUN go build

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
        && update-ca-certificates 2>/dev/null || true

# Direct logs to stdout for docker log driver
RUN ln -sf /dev/stdout /go/src/github.com/makerdao/vulcanizedb/vulcanizedb.log

ARG CONFIG_FILE=environments/docker.toml

# Add required files
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb .
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb/dockerfiles/extract_diffs/startup_script.sh .
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb/db/migrations/* db/migrations/
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose goose
# needed for waiting until postgres is ready before starting from docker-compose
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb/dockerfiles/wait-for-it.sh .

HEALTHCHECK CMD test -f /tmp/connection

# need to execute with a shell to access env variables
CMD ["./startup_script.sh"]