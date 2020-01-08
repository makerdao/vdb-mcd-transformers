package integration_tests

import (
	"errors"
	"io/ioutil"
	"log"
	"testing"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTests Suite")
}

var (
	db         *postgres.DB
	blockChain core.BlockChain
)

var _ = BeforeSuite(func() {
	testConfig := viper.New()
	testConfig.SetConfigName("testing")
	testConfig.AddConfigPath("$GOPATH/src/github.com/makerdao/vdb-mcd-transformers/environments/")
	err := testConfig.ReadInConfig()
	ipc = testConfig.GetString("client.ipcPath")
	if err != nil {
		logrus.Fatal(err)
	}
	// If we don't have an ipc path in the config file, check the env variable
	if ipc == "" {
		configErr := testConfig.BindEnv("url", "INFURA_URL")
		if configErr != nil {
			logrus.Fatalf("Unable to bind url to INFURA_URL env var")
		}
		ipc = testConfig.GetString("url")
	}
	if ipc == "" {
		logrus.Fatal(errors.New("infura.toml IPC path or $INFURA_URL env variable need to be set"))
	}

	rpcClient, ethClient, clientErr := getClients(ipc)
	Expect(clientErr).NotTo(HaveOccurred())
	var blockChainErr error
	blockChain, blockChainErr = getBlockChain(rpcClient, ethClient)
	Expect(blockChainErr).NotTo(HaveOccurred())

	db = test_config.NewTestDB(blockChain.Node())
	test_config.CleanTestDB(db)

	// Set log to discard logs emitted by dependencies
	log.SetOutput(ioutil.Discard)
	// Set logrus to discard logs emitted by mcd_transformers
	logrus.SetOutput(ioutil.Discard)
})
