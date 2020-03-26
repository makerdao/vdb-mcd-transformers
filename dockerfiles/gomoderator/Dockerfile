FROM golang:alpine as builder

RUN apk --update --no-cache add make git g++ linux-headers python3

ARG VDB_VERSION=staging
ENV GO111MODULE on

WORKDIR /go/src/github.com/makerdao
RUN git clone https://github.com/makerdao/vulcanizedb.git
WORKDIR /go/src/github.com/makerdao/vulcanizedb
RUN git checkout $VDB_VERSION
RUN go build
WORKDIR /go/src/github.com/makerdao/vdb-mcd-transformers
COPY . .
RUN go build

# app container
FROM golang:alpine
WORKDIR /go/src/github.com/makerdao/vulcanizedb

# add go so we can build the plugin
RUN apk add --update --no-cache git g++ linux-headers python3

# get access to go.mod files in directories
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb .
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/go.mod ./mcd-transformers/go.mod

RUN python3 scripts/gomoderator.py ./ mcd-transformers/
