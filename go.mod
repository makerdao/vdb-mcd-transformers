module github.com/makerdao/vdb-mcd-transformers

go 1.12

require (
	github.com/ethereum/go-ethereum v1.9.9
	github.com/gogo/protobuf v1.1.1
	github.com/howeyc/fsnotify v0.9.0 // indirect
	github.com/jmoiron/sqlx v0.0.0-20181024163419-82935fac6c1a
	github.com/lib/pq v1.0.0
	github.com/magiconair/properties v1.8.0
	github.com/makerdao/vulcanizedb v0.0.14-rc.1.0.20200401050751-582ceb4d07d8
	github.com/mattn/go-runewidth v0.0.6
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.3.2
	golang.org/x/crypto v0.0.0-20191119213627-4f8c1d86b1ba
	gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405
	gopkg.in/ini.v1 v1.51.1 // indirect
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
)

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/ethereum/go-ethereum => github.com/vulcanize/go-ethereum v0.0.0-20190731183759-8e20673bd101
