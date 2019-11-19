module github.com/makerdao/vdb-mcd-transformers

go 1.12

require (
	github.com/ethereum/go-ethereum v1.9.6
	github.com/jmoiron/sqlx v0.0.0-20181024163419-82935fac6c1a
	github.com/makerdao/vulcanizedb v0.0.9
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.3.2
	golang.org/x/crypto v0.0.0-20190926114937-fa1a29108794
)

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/ethereum/go-ethereum => github.com/vulcanize/go-ethereum v0.0.0-20190731183759-8e20673bd101
