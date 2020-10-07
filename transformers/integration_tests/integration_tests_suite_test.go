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
	test_config.SetTestConfig()
	var blockChainErr error
	blockChain, blockChainErr = test_config.NewTestBlockchain()
	Expect(blockChainErr).NotTo(HaveOccurred())

	db = test_config.NewTestDB(blockChain.Node())
	test_config.CleanTestDB(db)

	// Set log to discard logs emitted by dependencies
	log.SetOutput(ioutil.Discard)
})
