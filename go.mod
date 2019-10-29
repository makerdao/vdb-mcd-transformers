module github.com/vulcanize/mcd_transformers

go 1.12

require (
	github.com/apilayer/freegeoip v3.5.0+incompatible // indirect
	github.com/beorn7/perks v1.0.0
	github.com/btcsuite/btcd v0.0.0-20190115013929-ed77733ec07d // indirect
	github.com/btcsuite/btcutil v0.0.0-20180706230648-ab6388e0c60a
	github.com/cespare/cp v1.1.1 // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/ethereum/go-ethereum v1.9.6
	github.com/go-logfmt/logfmt v0.3.0
	github.com/gogo/protobuf v1.1.1
	github.com/graph-gophers/graphql-go v0.0.0-20191024035216-0a9cfbec35a1 // indirect
	github.com/influxdata/influxdb v1.7.9 // indirect
	github.com/jmoiron/sqlx v0.0.0-20181024163419-82935fac6c1a
	github.com/json-iterator/go v1.1.6
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/oschwald/maxminddb-golang v1.5.0 // indirect
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/common v0.4.1
	github.com/prometheus/procfs v0.0.2
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.3.2
	github.com/vulcanize/vulcanizedb v0.0.9
	golang.org/x/crypto v0.0.0-20190926114937-fa1a29108794
	google.golang.org/appengine v1.6.0 // indirect
	gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/ethereum/go-ethereum => github.com/vulcanize/go-ethereum v0.0.0-20190731183759-8e20673bd101
