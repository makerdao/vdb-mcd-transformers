module github.com/makerdao/vdb-mcd-transformers

go 1.12

require (
	github.com/OneOfOne/xxhash v1.2.5
	github.com/ethereum/go-ethereum v1.9.11
	github.com/howeyc/fsnotify v0.9.0 // indirect
	github.com/jmoiron/sqlx v0.0.0-20181024163419-82935fac6c1a
	github.com/lib/pq v1.0.0
	github.com/magiconair/properties v1.8.0
	github.com/makerdao/vulcanizedb v0.0.14-rc.1.0.20200522184407-b3c722d0f42a
	github.com/mattn/go-runewidth v0.0.6
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.10.0
	github.com/sirupsen/logrus v1.2.0
	github.com/spaolacci/murmur3 v1.0.1-0.20190317074736-539464a789e9
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20190404164418-38d8ce5564a5
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127
	gopkg.in/yaml.v2 v2.2.2
)

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/ethereum/go-ethereum => github.com/makerdao/go-ethereum v1.9.11-rc2
