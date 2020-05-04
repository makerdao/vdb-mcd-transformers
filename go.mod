module github.com/makerdao/vdb-mcd-transformers

go 1.12

require (
	github.com/ethereum/go-ethereum v1.9.9
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/graph-gophers/graphql-go v0.0.0-20191024035216-0a9cfbec35a1 // indirect
	github.com/howeyc/fsnotify v0.9.0 // indirect
	github.com/jmoiron/sqlx v0.0.0-20181024163419-82935fac6c1a
	github.com/lib/pq v1.0.0
	github.com/magiconair/properties v1.8.0
	github.com/makerdao/vulcanizedb v0.0.14-rc.1.0.20200512153458-6d857133e389
	github.com/mattn/go-runewidth v0.0.6
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/pkg/errors v0.8.1 // indirect
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.3.2
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570 // indirect
	github.com/steakknife/hamming v0.0.0-20180906055917-c99c65617cd3 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/wsddn/go-ecdh v0.0.0-20161211032359-48726bab9208 // indirect
	golang.org/x/crypto v0.0.0-20191119213627-4f8c1d86b1ba
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20190709231704-1e4459ed25ff // indirect
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
)

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/ethereum/go-ethereum => github.com/vulcanize/go-ethereum v0.0.0-20190731183759-8e20673bd101
