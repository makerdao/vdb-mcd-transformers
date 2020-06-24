module github.com/makerdao/vdb-mcd-transformers

go 1.12

require (
	github.com/ethereum/go-ethereum v1.9.11
	github.com/jmoiron/sqlx v0.0.0-20181024163419-82935fac6c1a
	github.com/lib/pq v1.0.0
	github.com/makerdao/vulcanizedb v0.0.14-rc.1.0.20200526232336-ac90e4306d86
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.10.0
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.2
	golang.org/x/crypto v0.0.0-20190404164418-38d8ce5564a5
)

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/ethereum/go-ethereum => github.com/makerdao/go-ethereum v1.9.11-rc2
