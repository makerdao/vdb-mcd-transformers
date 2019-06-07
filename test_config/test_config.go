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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var TestConfig *viper.Viper
var DBConfig config.Database
var TestClient config.Client
var Infura *viper.Viper
var InfuraClient config.Client
var ABIFilePath string
var wipeTableQueries []string

func init() {
	setTestConfig()
	setABIPath()
}

func setTestConfig() {
	TestConfig = viper.New()
	TestConfig.SetConfigName("testing")
	TestConfig.AddConfigPath("$GOPATH/src/github.com/vulcanize/mcd_transformers/environments/")
	err := TestConfig.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	ipc := TestConfig.GetString("client.ipcPath")
	hn := TestConfig.GetString("database.hostname")
	port := TestConfig.GetInt("database.port")
	name := TestConfig.GetString("database.name")
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
	ABIFilePath = gp + "/src/github.com/vulcanize/vulcanizedb/pkg/geth/testing/"
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
