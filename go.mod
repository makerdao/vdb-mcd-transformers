module github.com/vulcanize/mcd_transformers

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/btcsuite/btcd v0.0.0-20190115013929-ed77733ec07d // indirect
	github.com/cespare/cp v1.1.1 // indirect
	github.com/elastic/gosigar v0.10.4
	github.com/ethereum/go-ethereum v1.9.5
	github.com/jmoiron/sqlx v0.0.0-20181024163419-82935fac6c1a
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.3.2
	github.com/vulcanize/vulcanizedb v0.0.8
	golang.org/x/crypto v0.0.0-20190926114937-fa1a29108794
	golang.org/x/tools v0.0.0-20190606124116-d0a3d012864b // indirect
	google.golang.org/appengine v1.6.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
)

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/ethereum/go-ethereum => github.com/vulcanize/go-ethereum v0.0.0-20190731183759-8e20673bd101
