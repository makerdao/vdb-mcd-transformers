module github.com/makerdao/vdb-mcd-transformers

go 1.12

require (
	github.com/beorn7/perks v1.0.0
	github.com/ethereum/go-ethereum v1.9.7
	github.com/go-logfmt/logfmt v0.3.0
	github.com/gogo/protobuf v1.1.1
	github.com/jmoiron/sqlx v0.0.0-20181024163419-82935fac6c1a
	github.com/json-iterator/go v1.1.6
	github.com/magiconair/properties v1.8.0
	github.com/makerdao/vulcanizedb v0.0.10
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/common v0.4.1
	github.com/prometheus/procfs v0.0.2
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.3.2
	golang.org/x/crypto v0.0.0-20191119213627-4f8c1d86b1ba
	golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e
	gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405
)

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/ethereum/go-ethereum => github.com/vulcanize/go-ethereum v0.0.0-20190731183759-8e20673bd101
