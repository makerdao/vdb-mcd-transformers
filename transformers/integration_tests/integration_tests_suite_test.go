package integration_tests

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	Expect(err).To(BeNil())
	ipc = testConfig.GetString("client.ipcPath")
	// If we don't have an ipc path in the config file, check the env variable
	if ipc == "" {
		configErr := testConfig.BindEnv("url", "CLIENT_IPCPATH")
		Expect(configErr).To(BeNil(), "Unable to bind url to CLIENT_IPCPATH env var")
		ipc = testConfig.GetString("url")
	}
	Expect(ipc).NotTo(BeEmpty(), "testing.toml IPC path or $CLIENT_IPCPATH env variable need to be set")

	rpcClient, ethClient, clientErr := getClients(ipc)
	Expect(clientErr).NotTo(HaveOccurred())
	var blockChainErr error
	blockChain, blockChainErr = getBlockChain(rpcClient, ethClient)
	Expect(blockChainErr).NotTo(HaveOccurred())

	db = test_config.NewTestDB(blockChain.Node())
	test_config.CleanTestDB(db)

	// Set log to discard logs emitted by dependencies
	log.SetOutput(ioutil.Discard)
})
