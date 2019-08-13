package cdp_manager_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestCdpManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CDP Manager Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
