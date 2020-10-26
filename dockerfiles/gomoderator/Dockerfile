FROM makerdao/vdb-builder:latest as builder

ARG VDB_VERSION=prod

WORKDIR /go/src/github.com/makerdao
RUN git clone https://github.com/makerdao/vulcanizedb.git
WORKDIR /go/src/github.com/makerdao/vulcanizedb
RUN git checkout $VDB_VERSION
RUN go build
WORKDIR /go/src/github.com/makerdao/vdb-mcd-transformers
COPY . .

# app container
FROM makerdao/vdb-runner:latest
WORKDIR /go/src/github.com/makerdao/vulcanizedb

# get access to go.mod files in directories
COPY --from=builder /go/src/github.com/makerdao/vulcanizedb .
COPY --from=builder /go/src/github.com/makerdao/vdb-mcd-transformers/go.mod ./mcd-transformers/go.mod

RUN python3 scripts/gomoderator.py ./ mcd-transformers/
