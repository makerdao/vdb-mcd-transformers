// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package test_config

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/makerdao/vulcanizedb/pkg/eth"
	"github.com/makerdao/vulcanizedb/pkg/eth/client"
	"github.com/makerdao/vulcanizedb/pkg/eth/converters"
	"github.com/makerdao/vulcanizedb/pkg/eth/node"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

var DBConfig config.Database
var TestClient config.Client
var ABIFilePath string
var wipeTableQueries []string

func init() {
	SetTestConfig()
	setABIPath()
}

func SetTestConfig() {
	viper.AddConfigPath("$GOPATH/src/github.com/makerdao/vdb-mcd-transformers/environments/")

	viper.SetConfigName("testDatabase")
	readConfigErr := viper.ReadInConfig()
	if readConfigErr != nil {
		log.Fatal(readConfigErr)
	}

	viper.SetConfigName("mcdTransformers")
	mergeConfigErr := viper.MergeInConfig()
	if mergeConfigErr != nil {
		log.Fatal(mergeConfigErr)
	}

	ipc := viper.GetString("client.ipcPath")
	hn := viper.GetString("database.hostname")
	port := viper.GetInt("database.port")
	name := viper.GetString("database.name")

	DBConfig = config.Database{
		Hostname: hn,
		Name:     name,
		Port:     port,
	}
	TestClient = config.Client{
		IPCPath: ipc,
	}
}

func setABIPath() {
	gp := os.Getenv("GOPATH")
	ABIFilePath = gp + "/src/github.com/makerdao/vulcanizedb/pkg/eth/testing/"
}

func NewTestBlockchain() (core.BlockChain, error) {
	ipc := TestClient.IPCPath
	// If we don't have an ipc path in the config file, check the env variable
	if ipc == "" {
		configErr := viper.BindEnv("url", "CLIENT_IPCPATH")
		if configErr != nil {
			return nil, fmt.Errorf("unable to bind url to CLIENT_IPCPATH env var %w", configErr)
		}
		ipc = viper.GetString("url")
	}

	rpcClient, ethClient, clientErr := getClients(ipc)
	if clientErr != nil {
		return nil, fmt.Errorf("failed to get test clients: %w", clientErr)
	}

	return getBlockChain(rpcClient, ethClient), nil
}

func getClients(ipc string) (client.RpcClient, *ethclient.Client, error) {
	raw, err := rpc.Dial(ipc)
	if err != nil {
		return client.RpcClient{}, &ethclient.Client{}, err
	}
	return client.NewRpcClient(raw, ipc), ethclient.NewClient(raw), nil
}

func getBlockChain(rpcClient client.RpcClient, ethClient *ethclient.Client) core.BlockChain {
	testClient := client.NewEthClient(ethClient)
	testNode := node.MakeNode(rpcClient)
	transactionConverter := converters.NewTransactionConverter(testClient)
	return eth.NewBlockChain(testClient, rpcClient, testNode, transactionConverter)
}

func NewTestDB(node core.Node) *postgres.DB {
	db, err := postgres.NewDB(DBConfig, node)
	if err != nil {
		panic(fmt.Sprintf("Could not create new test db: %v", err))
	}
	return db
}

// Cleans all tables in the DB. Note that this requires cascade constraints to be in place,
// so deletion can be run in any order.
func CleanTestDB(db *postgres.DB) {
	if len(wipeTableQueries) == 0 {
		// The generated queries delete from all tables in the public and maker schemas,
		// except eth_nodes and goose_db_version.
		err := db.Select(&wipeTableQueries,
			`SELECT 'DELETE FROM ' || schemaname || '.' || relname || ';'
			FROM pg_stat_user_tables
			WHERE schemaname IN ('public', 'maker', 'api')
			AND relname NOT IN ('eth_nodes', 'goose_db_version');`)
		if err != nil {
			panic("Failed to generate DB cleaning query: " + err.Error())
		}
	}

	for _, query := range wipeTableQueries {
		db.MustExec(query)
	}
}

// Returns a new test node, with the same ID
func NewTestNode() core.Node {
	return core.Node{
		GenesisBlock: "GENESIS",
		NetworkID:    1,
		ID:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
		ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
	}
}
