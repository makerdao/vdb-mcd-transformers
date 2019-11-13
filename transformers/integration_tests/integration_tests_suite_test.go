package integration_tests

import (
	"errors"
	"github.com/sirupsen/logrus"
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"io/ioutil"
)

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTests Suite")
}

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
		testConfig.BindEnv("url", "INFURA_URL")
		ipc = testConfig.GetString("url")
	}
	if ipc == "" {
		logrus.Fatal(errors.New("infura.toml IPC path or $INFURA_URL env variable need to be set"))
	}
	// Set log to discard logs emitted by dependencies
	log.SetOutput(ioutil.Discard)
	// Set logrus to discard logs emitted by mcd_transformers
	logrus.SetOutput(ioutil.Discard)
})
