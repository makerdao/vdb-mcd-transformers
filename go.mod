module github.com/makerdao/vdb-mcd-transformers

go 1.12

require (
	github.com/ClickHouse/clickhouse-go v1.4.3 // indirect
	github.com/apilayer/freegeoip v3.5.0+incompatible // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20200620013148-b91950f658ec // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/elastic/gosigar v0.10.5 // indirect
	github.com/ethereum/go-ethereum v1.9.16
	github.com/influxdata/influxdb v1.7.9 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/karalabe/usb v0.0.0-20191104083709-911d15fe12a9 // indirect
	github.com/lib/pq v1.0.0
	github.com/makerdao/vulcanizedb v0.0.15-rc.1.0.20200909153029-7aa1ed76f791
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.10.0
	github.com/oschwald/maxminddb-golang v1.5.0 // indirect
	github.com/sirupsen/logrus v1.2.0
	github.com/spaolacci/murmur3 v0.0.0-20180118202830-f09979ecbc72
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.2
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20200603215123-a4a8cb9d2cbc
)

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/ethereum/go-ethereum => github.com/makerdao/go-ethereum v1.9.15-statechange-filter
