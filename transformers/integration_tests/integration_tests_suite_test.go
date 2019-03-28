package integration_tests

import (
	"errors"
	log "github.com/sirupsen/logrus"
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
	testConfig.SetConfigName("infura")
	testConfig.AddConfigPath("$GOPATH/src/github.com/vulcanize/mcd_transformers/environments/")
	err := testConfig.ReadInConfig()
	ipc = testConfig.GetString("client.ipcPath")
	if err != nil {
		log.Fatal(err)
	}
	// If we don't have an ipc path in the config file, check the env variable
	if ipc == "" {
		testConfig.BindEnv("url", "INFURA_URL")
		ipc = testConfig.GetString("url")
	}
	if ipc == "" {
		log.Fatal(errors.New("infura.toml IPC path or $INFURA_URL env variable need to be set"))
	}
	log.SetOutput(ioutil.Discard)
})
