module github.com/makerdao/vdb-mcd-transformers

go 1.12

require (
	github.com/ethereum/go-ethereum v1.9.8
	github.com/jmoiron/sqlx v0.0.0-20181024163419-82935fac6c1a
	github.com/magiconair/properties v1.8.0
	github.com/makerdao/vulcanizedb v0.0.11-0.20191211162058-30c501644255
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.3.2
	golang.org/x/crypto v0.0.0-20191119213627-4f8c1d86b1ba
	gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405
)

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/ethereum/go-ethereum => github.com/vulcanize/go-ethereum v0.0.0-20190731183759-8e20673bd101
