module github.com/makerdao/vdb-mcd-transformers

go 1.12

require (
	github.com/OneOfOne/xxhash v1.2.2
	github.com/VictoriaMetrics/fastcache v1.5.7
	github.com/dop251/goja v0.0.0-20200219165308-d1232e640a87
	github.com/ethereum/go-ethereum v1.9.16
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4
	github.com/howeyc/fsnotify v0.9.0 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.0.0
	github.com/makerdao/vulcanizedb v0.0.15-rc.1.0.20200826164038-125a904f02eb
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.10.0
	github.com/pressly/goose v2.7.0-rc5+incompatible // indirect
	github.com/sirupsen/logrus v1.2.0
	github.com/spaolacci/murmur3 v0.0.0-20180118202830-f09979ecbc72
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.2
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20200603215123-a4a8cb9d2cbc
)

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/ethereum/go-ethereum => github.com/makerdao/go-ethereum v1.9.15-statechange-filter
